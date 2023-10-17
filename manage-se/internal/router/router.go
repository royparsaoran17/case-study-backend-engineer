// Package router
package router

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/presentations"
	"manage-se/internal/provider"
	"manage-se/internal/service/auth"
	"manage-se/internal/service/customer"
	"manage-se/internal/service/role"
	"manage-se/internal/service/wallet"
	"manage-se/pkg/httpclientx"
	"manage-se/pkg/tracer"
	"net/http"
	"time"

	"manage-se/internal/appctx"
	"manage-se/internal/bootstrap"
	"manage-se/internal/consts"
	"manage-se/internal/handler"
	"manage-se/internal/middleware"
	"manage-se/internal/ucase"
	"manage-se/pkg/logger"
	"manage-se/pkg/msgx"
	"manage-se/pkg/routerkit"

	ucaseContract "manage-se/internal/ucase/contract"
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
	//db := bootstrap.RegistryPostgres(rtr.config.WriteDB)

	// repositories
	//repo := repositories.NewRepository(db)

	provider := provider.NewProviders(
		&rtr.config.Provider,
		provider.WithHttpClient(httpclientx.NewHttpClientx()),
	)

	// init redis
	rdb := bootstrap.RegistryRedisNative(rtr.config)

	// initiate services
	var (
		authService     = auth.NewService(provider, rdb)
		roleService     = role.NewService(provider)
		customerService = customer.NewService(provider)
		walletService   = wallet.NewService(provider)
	)

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		ucase.NewHealthCheck(),
	)).Methods(http.MethodGet)

	rtr.mountAuth(authService)
	rtr.mountRole(roleService)
	rtr.mountWallet(walletService, authService)
	rtr.mountCustomer(customerService, authService)

	return rtr.router

}
