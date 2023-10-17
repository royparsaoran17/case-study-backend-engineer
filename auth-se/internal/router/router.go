// Package router
package router

import (
	"auth-se/internal/presentations"
	"auth-se/internal/repositories"
	"auth-se/internal/service/auth"
	"auth-se/internal/service/customer"
	"auth-se/internal/service/role"
	"auth-se/pkg/tracer"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"

	"auth-se/internal/appctx"
	"auth-se/internal/bootstrap"
	"auth-se/internal/consts"
	"auth-se/internal/handler"
	"auth-se/internal/middleware"
	"auth-se/internal/ucase"
	"auth-se/pkg/logger"
	"auth-se/pkg/msgx"
	"auth-se/pkg/routerkit"

	ucaseContract "auth-se/internal/ucase/contract"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rid := r.Header.Get(consts.HeaderXRequestID)
		if rid == "" {
			rid = uuid.NewString()
			r.Header.Set(consts.HeaderXRequestID, rid)
		}

		// create the initial state
		state := presentations.RequestState{
			ID:        rid,
			CreatedAt: time.Now().Local(),
		}

		// Set an initial state value for each request context.
		ctx := context.WithValue(r.Context(), consts.CtxRequestState, state)

		// Re-usable response body for logging
		requestBody, _ := io.ReadAll(r.Body)
		r.Body.Close() // must close
		r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		tracer.AddSpanTag(r.Context(),
			tracer.NewSpanTag("http.request.headers.*", r.Header),
			tracer.NewSpanTag("http.request.body", string(requestBody)),
			tracer.NewSpanTag("http.request.query_params", r.URL.Query()),
			tracer.NewSpanTag("http.x_request_id", rid),
		)

		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msgx.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("error %v", err))
				json.NewEncoder(w).Encode(res.Byte())

				return
			}
		}()

		ctx = context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)
		resp := appctx.Response{}

		// validate middleware
		if err := middleware.FilterFunc(w, req, rtr.config, mdws); err != nil {
			logger.Error(errors.Wrap(err, "error on middleware"))

			switch e := err.(type) {
			case middleware.Error:
				resp = e.Response

			default:
				resp = *appctx.NewResponse().WithContext(ctx).
					WithCode(http.StatusInternalServerError).
					WithMessage(http.StatusText(http.StatusInternalServerError))
			}

			rtr.response(w, resp)
			return
		}

		resp = hfn(req, svc, rtr.config)
		resp.WithLang(lang)
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	return
}

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	//rtr.router.NotFoundHandler = http.HandlerFunc(middleware.NotFound)
	root := rtr.router.PathPrefix("/").Subrouter()
	//in := root.PathPrefix("/internal/").Subrouter()
	liveness := root.PathPrefix("/").Subrouter()
	//_ := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)
	db := bootstrap.RegistryPostgres(rtr.config.WriteDB)

	// repositories
	repo := repositories.NewRepository(db)

	// init redis
	//rdb := bootstrap.RegistryRedisNative(rtr.config)

	// initiate services
	var (
		roleService     = role.NewService(repo)
		authService     = auth.NewService(repo, *rtr.config.App)
		customerService = customer.NewService(repo)
	)

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		ucase.NewHealthCheck(),
	)).Methods(http.MethodGet)

	rtr.mountRoles(roleService)
	rtr.mountCustomer(customerService)
	rtr.mountAuth(authService)

	return rtr.router

}
