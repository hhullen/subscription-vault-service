package api

import (
	"context"
	"net/http"
	"time"

	ds "subscription-vault-service/internal/datastruct"
	"subscription-vault-service/internal/supports"

	_ "subscription-vault-service/internal/docs"

	"github.com/gorilla/schema"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:generate mockgen -source=api.go -destination=api_mock.go -package=api ISubscriptionService,ISecretProvider,IServer,IRouter,ILogger
//go:generate mockgen -destination=http_mock.go -package=api net/http Handler

const (
	readTimeout  = time.Second * 5
	writeTimeout = time.Second * 5

	globalRateLimit      = 500
	globalLimiterKey     = "global"
	overlimitMessage     = "service exhausted"
	userIdLimitPerSecond = 100
	limiterUserIdPrefix  = "user_id"

	apiPrefix     = "/api/v1"
	swaggerPrefix = "/swagger/"
)

var schemaDecoder = schema.NewDecoder()

func init() {
	schemaDecoder.RegisterConverter(ds.DateType{}, ds.ParseSchemaDateType)
}

type IWithStatus interface {
	GetStatus() string
}

type ISubscriptionService interface {
	AddSubscription(*ds.CreateSubscriptionRequest) *ds.CreateSubscriptionResponse
	GetSubscription(*ds.GetSubscriptionRequest) *ds.GetSubscriptionResponse
	UpdateSubscription(*ds.UpdateSubscriptionRequest) *ds.UpdateSubscriptionResponse
	DeleteSubscription(*ds.DeleteSubscriptionRequest) *ds.DeleteSubscriptionResponse
	ListSubscriptions(*ds.ListSubscriptionsRequest) *ds.ListSubscriptionsResponse
	CalculateSubscriptionsPrice(*ds.CalculateSubscriptionsPriceRequest) *ds.CalculateSubscriptionsPriceResponse
}

type ISecretProvider interface {
	ReadSecret(key string) (string, error)
}

type IServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type IRouter interface {
	Handle(pattern string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ILogger interface {
	InfoKV(message string, argsKV ...any)
	WarnKV(message string, argsKV ...any)
	ErrorKV(message string, argsKV ...any)
	FatalKV(message string, argsKV ...any)
}

type ResponseWriterInterceptor struct {
	rw   http.ResponseWriter
	code int
}

func (rwi *ResponseWriterInterceptor) Header() http.Header {
	return rwi.rw.Header()
}

func (rwi *ResponseWriterInterceptor) Write(data []byte) (int, error) {
	return rwi.rw.Write(data)
}

func (rwi *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.code = statusCode
	rwi.rw.WriteHeader(statusCode)
}

type API struct {
	ctx         context.Context
	logger      ILogger
	subsService ISubscriptionService
	secret      ISecretProvider
	server      IServer
	router      IRouter
}

func NewAPI(ctx context.Context, address string,
	subs ISubscriptionService,
	sec ISecretProvider,
	log ILogger) (*API, error) {
	router := http.NewServeMux()

	router.Handle(swaggerPrefix, httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         address,
		Handler:      mainMiddleware(router, log),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return buildAPI(ctx, subs, log, sec, server, router), nil
}

func buildAPI(ctx context.Context,
	subs ISubscriptionService,
	log ILogger,
	sec ISecretProvider,
	srv IServer,
	rou IRouter) *API {
	api := &API{
		ctx:         ctx,
		logger:      log,
		subsService: subs,
		secret:      sec,
		server:      srv,
		router:      rou,
	}

	api.setupSubscriptionHandlers()

	return api
}

func (a *API) StartListening() error {
	a.logger.InfoKV("Server is listening")
	return a.server.ListenAndServe()
}

func (a *API) Stop() error {
	a.logger.InfoKV("Server Shutdown")
	return a.server.Shutdown(a.ctx)
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func mainMiddleware(next http.Handler, log ILogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		rwi := &ResponseWriterInterceptor{rw: w}
		next.ServeHTTP(rwi, r)

		te := time.Since(ts)

		loggerFunc := log.InfoKV
		if rwi.code >= 500 {
			loggerFunc = log.ErrorKV
		} else if rwi.code >= 400 {
			loggerFunc = log.WarnKV
		}
		loggerFunc("Request", "method", r.Method, "url", r.URL.String(), "duration(ms)", te.Milliseconds(), "status_code", rwi.code)
	})
}

func pattern(method, prefixPath string) string {
	return supports.Concat(method, " ", prefixPath)
}
