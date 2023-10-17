// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"

    validation "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/gorilla/mux"

	"{{.ModuleName}}/internal/appctx"
	"{{.ModuleName}}/internal/consts"
	"{{.ModuleName}}/internal/presentations"
	"{{.ModuleName}}/internal/repositories"
	"{{.ModuleName}}/internal/entity"
	"{{.ModuleName}}/pkg/logger"
	"{{.ModuleName}}/pkg/tracer"

	ucase "{{.ModuleName}}/internal/ucase/contract"
)

type {{.StructName}}Update struct {
	repo repositories.{{.RepoContractName}}
}

// New{{.EntityName}} new instance
func New{{.EntityName}}Update(repo repositories.{{.RepoContractName}}) ucase.UseCase {
	return &{{.StructName}}Update{repo: repo}
}

// Serve store {{.StructName}} data
func (u *{{.StructName}}Update) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Param
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.update")
		lf    = logger.NewFields(
			logger.EventName("{{.EntityName}}"),
		)

		id = mux.Vars(dctx.Request)["id"]
	)

	defer tracer.SpanFinish(ctx)

    err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	err = validation.ValidateStruct(&param,

    	)

    if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	df, err := u.repo.FindOne(ctx, entity.{{.EntityName}}{})
    if err != nil {
        tracer.SpanError(ctx, err)
        logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
        return *appctx.NewResponse().WithMsgKey(consts.RespError)
    }

    if df == nil {
        logger.WarnWithContext(ctx, fmt.Sprintf("find data %s in database not found", id), lf...)
        return *appctx.NewResponse().WithMsgKey(consts.RespDataNotFound)
    }

    af, err := u.repo.Update(ctx, param, entity.{{.EntityName}}{})
    if err != nil {
        tracer.SpanError(ctx, err)
        logger.ErrorWithContext(ctx, fmt.Sprintf("error update data to database: %v", err), lf...)
        return *appctx.NewResponse().WithMsgKey(consts.RespError)
    }

    logger.InfoWithContext(ctx, fmt.Sprintf("success update data to database with affected rows %d", af), lf...)
    return *appctx.NewResponse().
        WithMsgKey(consts.RespSuccess)
}
