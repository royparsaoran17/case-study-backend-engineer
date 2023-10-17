// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"


    validation "github.com/go-ozzo/ozzo-validation/v4"

	"{{.ModuleName}}/internal/appctx"
	"{{.ModuleName}}/internal/consts"
	"{{.ModuleName}}/internal/presentations"
	"{{.ModuleName}}/internal/repositories"
	"{{.ModuleName}}/pkg/logger"
	"{{.ModuleName}}/pkg/tracer"

	ucase "{{.ModuleName}}/internal/ucase/contract"
)

type {{.StructName}} struct {
	repo repositories.{{.RepoContractName}}
}

// New{{.EntityName}} new instance
func New{{.EntityName}}(repo repositories.{{.RepoContractName}}) ucase.UseCase {
	return &{{.StructName}}{repo: repo}
}

// Serve store {{.StructName}} data
func (u *{{.StructName}}) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Param
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.create")
		lf    = logger.NewFields(
			logger.EventName("{{.EntityName}}"),
		)
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

	af, err := u.repo.Store(ctx, param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.WarnWithContext(ctx, fmt.Sprintf("store data to database error: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	lf.Append(logger.Any("affected_rows", af))

	logger.InfoWithContext(ctx, fmt.Sprintf("success store data to database"), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
