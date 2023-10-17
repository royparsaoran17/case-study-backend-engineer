package dto

import (
	"{{.ModuleName}}/internal/consts"
	"{{.ModuleName}}/internal/entity"
	"{{.ModuleName}}/internal/presentations"
)

func {{.ObjectName}}ToResponse(src entity.{{.ObjectName}}) presentations.{{.ObjectName}}Detail {
	x := presentations.{{.ObjectName}}Detail{
        {{- range $key, $value := .Column }}
            {{ $value.Name }}: src.{{ $value.Name }},
        {{- end }}
	}

	if !src.CreatedAt.IsZero() {
		x.CreatedAt = src.CreatedAt.Format(consts.LayoutDateTimeFormat)
	}
	//
	//if !src.UpdatedAt.IsZero() {
	//	x.UpdatedAt = src.UpdatedAt.Format(consts.LayoutDateTimeFormat)
	//}
	//
	//if !src.DeletedAt.IsZero() {
	//	x.DeletedAt = src.DeletedAt.Format(consts.LayoutDateTimeFormat)
	//}

	return x
}

func {{.ObjectName}}sToResponse(inputs []entity.{{.ObjectName}}) []presentations.{{.ObjectName}}Detail {
	var (
		result = []presentations.{{.ObjectName}}Detail{}
	)

	for i := 0; i {{.LessThenSign}} len(inputs); i++ {
		result = append(result, {{.ObjectName}}ToResponse(inputs[i]))
	}

	return result
}