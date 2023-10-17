package middleware

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"manage-se/internal/appctx"
	"manage-se/internal/consts"
	"manage-se/internal/entity"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/providererrors"
	"manage-se/internal/service/auth"
	"manage-se/pkg/logger"
	"manage-se/pkg/tracer"
	"net/http"
)

func Authorize(authService auth.Auth, roleNames ...string) MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, conf *appctx.Config) error {
		ctx := tracer.SpanStart(r.Context(), "middleware.auth")

		bearerToken := entity.BearerToken(r.Header.Get("Authorization"))

		responder := appctx.NewResponse().WithState("authenticate").WithContext(ctx)

		fl := logger.Fields{
			logger.EventName("middleware.auth"),
			logger.Field{
				Key:   "bearer-token",
				Value: bearerToken.String(),
			},
		}

		defer tracer.SpanFinish(ctx)

		if bearerToken.TokenEmpty() {
			return Error{
				Response: *responder.WithContext(ctx).
					WithMessage(http.StatusText(http.StatusUnauthorized)).
					WithCode(http.StatusUnauthorized),
				err: nil,
			}
		}

		customer, err := authService.VerifyToken(ctx, presentations.Verify{Token: bearerToken.GetToken()})
		if err != nil {
			fmt.Println(err)
			logger.Error(err, fl...)

			errCause := errors.Cause(err)
			switch errCause {

			case consts.ErrForbidden:
				return Error{
					Response: *responder.WithContext(ctx).
						WithMessage(http.StatusText(http.StatusForbidden)).
						WithCode(http.StatusForbidden),
					err: nil,
				}

			case consts.ErrUnauthorized, consts.ErrBearerTokenNotProvided:
				return Error{
					Response: *responder.WithContext(ctx).
						WithMessage(http.StatusText(http.StatusUnauthorized)).
						WithCode(http.StatusUnauthorized),
					err: nil,
				}

			default:
				switch causer := errCause.(type) {

				case consts.Error:
					return Error{
						Response: *responder.
							WithContext(ctx).
							WithMessage(errCause.Error()).
							WithCode(http.StatusBadRequest),
						err: nil,
					}

				case providererrors.Error:
					return Error{
						Response: *responder.
							WithContext(ctx).
							WithMessage(causer.Message).
							WithError(causer.Errors).
							WithCode(causer.Code),
						err: nil,
					}

				default:
					return errors.Wrap(err, "Authorize Middleware")
				}
			}
		}

		ctx = context.WithValue(ctx, consts.CtxUserAuth, entity.CustomerContext{
			ID:     customer.ID.String(),
			Name:   customer.Name,
			Phone:  customer.Phone,
			RoleID: customer.RoleID.String(),
			Role:   customer.Role,
		})

		if len(roleNames) == 0 {
			return nil
		}

		req := r.WithContext(ctx)
		*r = *req

		for _, roleName := range roleNames {

			// by pass if middleware validation role is all roles
			if roleName == consts.AllRoles {
				return nil
			}

			if roleName == customer.Role.Name {
				return nil
			}
		}

		return Error{
			Response: *responder.WithContext(ctx).
				WithMessage(http.StatusText(http.StatusForbidden)).
				WithCode(http.StatusForbidden),
			err: nil,
		}

	}
}
