// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"{{.ModuleName}}/internal/appctx"
	"{{.ModuleName}}/internal/common"
	"{{.ModuleName}}/internal/consts"
	"{{.ModuleName}}/internal/dto"
	"{{.ModuleName}}/internal/presentations"
	"{{.ModuleName}}/internal/repositories"
	"{{.ModuleName}}/pkg/logger"
	"{{.ModuleName}}/pkg/tracer"
	"{{.ModuleName}}/internal/validator"

	ucase "{{.ModuleName}}/internal/ucase/contract"
)

type {{.StructName}}List struct {
	repo repositories.{{.RepoContractName}}
}

func New{{.EntityName}}List(repo repositories.{{.RepoContractName}}) ucase.UseCase {
	return &{{.StructName}}List{repo: repo}
}

// Serve {{.EntityName}} list data
func (u *{{.StructName}}List) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Query
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.{{.FileName}}_list")
		lf    = logger.NewFields(
			logger.EventName("{{.StructName}}List"),
		)
	)
    defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	err = validation.ValidateStruct(&param,
		validation.Field(&param.Page, validation.Min(int64(1))),
		validation.Field(&param.Limit, validation.Min(int64(1))),
		validation.Field(&param.StartDate, validation.Required.When(param.EndDate != ""), validator.ValidDateTime()),
		validation.Field(&param.EndDate, validation.Required.When(param.StartDate != ""), validator.ValidDateTime()),
	)

	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError).WithError(err)
	}

	param.Limit = common.LimitDefaultValue(param.Limit)
	param.Page = common.PageDefaultValue(param.Page)

	dr, count, err := u.repo.FindWithCount(ctx, param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success fetch {{.TableName}} to database"), lf...)
	return *appctx.NewResponse().
            WithMsgKey(consts.RespSuccess).
            WithData(dto.{{.EntityName}}sToResponse(dr)).
            WithMeta(appctx.MetaData{
                    Page:       param.Page,
                    Limit:      param.Limit,
                    TotalCount: count,
                    TotalPage:  common.PageCalculate(count, param.Limit),
            })
}