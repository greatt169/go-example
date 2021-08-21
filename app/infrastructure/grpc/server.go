package grpc

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	cache "github.com/AeroAgency/golang-bigcache-lib"
	helpers "github.com/AeroAgency/golang-helpers-lib"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"news-ms/application/usecases/news"
	tagrepository "news-ms/domain/tag/repository"
	"news-ms/infrastructure/env"
	"news-ms/interfaces/grpc/handler"
	access_control_v1 "news-ms/interfaces/grpc/proto/v1/access_control"
	content_v1 "news-ms/interfaces/grpc/proto/v1/news"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

var (
	// Create a metrics registry.
	reg = prometheus.NewRegistry()
	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()
	// access_control
	accessControlGrpcPort = os.Getenv("ACCESS_CONTROL_GRPC_PORT")
	accessControlGrpcHost = os.Getenv("ACCESS_CONTROL_GRPC_HOST")
	healthCheckUrl        = "/eapi/news-ms/manage/health"
)

func init() {
	// Регистрация Метрик Prometheus.
	reg.MustRegister(helpers.RequestCount, helpers.RequestDuration)
}

// Сервер контента
type ContentServer struct {
	Server               *grpc.Server
	Logger               log.FieldLogger
	serverCert           string
	serverKey            string
	certificateAuthority string
	middleware           *helpers.Middleware
}

// Запуск сервиса
func (c *ContentServer) Run() error {
	return errors.New("error")
}

// Конструктор
func NewGrpcServer(
	logger log.FieldLogger,
	newsService *news.NewsService,
	cache cache.CacheInterface,
	tagRepository tagrepository.TagRepository,
) *ContentServer {
	// Инициализация нового grpc-сервера
	var grpcServerUnaryInterceptor grpc.UnaryServerInterceptor
	logMode, _ := strconv.Atoi(os.Getenv("GRPC_LOGMODE"))
	if logMode == 1 {
		grpcServerUnaryInterceptor = helpers.GrpcServerUnaryInterceptor
	}
	gserver := grpc.NewServer(
		grpc.MaxRecvMsgSize(10010241024), // снимаем ограничение на получение файлов
		grpc.MaxSendMsgSize(10010241024), // снимаем ограничение на отправку файлов
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcServerUnaryInterceptor),
	)
	// Инициализация Метрик Prometheus.
	grpcMetrics.InitializeMetrics(gserver)
	// Регистрация хэндлеров
	content_v1.RegisterNewsServiceServer(gserver, handler.NewNewsHandler(newsService))
	content_v1.RegisterTagServiceServer(gserver, handler.NewTagHandler(tagRepository, cache))
	content_v1.RegisterContentCheckServiceServer(gserver, handler.NewCheckHandler())

	return &ContentServer{
		Server:               gserver,
		Logger:               logger,
		serverCert:           "",
		serverKey:            "",
		certificateAuthority: "",
		middleware:           &helpers.Middleware{},
	}
}

// Запуск grpc сервера
func (c *ContentServer) RunGrpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	// Слушаем grpc-порт из настроек
	lis, err := net.Listen("tcp", env.GrpcPort)
	if err != nil {
		log.Println(err)
	}

	// Создаем HTTP сервер для prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%s", env.PrometheusPort)}
	// Запуск http сервера для prometheus.
	go func() {
		log.Printf("start prometheus server port: %s", env.PrometheusPort)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()
	go func() {
		log.Printf("start grpc server port%s", env.GrpcPort)
		c.Server.Serve(lis)
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping grpc server...")
	c.Server.GracefulStop()
}

// Запуск proxy-сервера rest -> grpc
func (c *ContentServer) RunProxyServer(wg *sync.WaitGroup) {
	defer wg.Done()
	var grpcServerEndpoint = flag.String("grpc-server-endpoint", env.GrpcPort, "gRPC server endpoint")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Устанавливаем кастомный обработчик ошибок
	runtime.HTTPError = c.middleware.CustomHTTPError
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}),
		runtime.WithIncomingHeaderMatcher(c.middleware.CustomMatcher),
		runtime.WithForwardResponseOption(c.middleware.HttpResponseModifier),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := content_v1.RegisterNewsServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	err = content_v1.RegisterPromoServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	err = content_v1.RegisterContentCheckServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	err = content_v1.RegisterTagServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatal("error reg endpoint", err)
	}
	go func() {
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Println("start Http server on port: " + env.Port)
		log.Fatal(http.ListenAndServe(env.Port, c.MiddlewaresMidHandler(mux)))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping proxy grpc server...")
}

// Выполняет запрос к МС Access Control, возвращает привилегии по токену
func getPrivileges(token string) (string, error) {
	accessControlServiceConnection, _ := grpc.Dial(accessControlGrpcHost+":"+accessControlGrpcPort, grpc.WithInsecure())
	accessControlServiceClient := access_control_v1.NewPrivilegesServiceClient(accessControlServiceConnection)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("access-token", token, "basic-auth", ""),
	)
	privileges, err := accessControlServiceClient.GetUserPrivileges(ctx, &access_control_v1.EmptyRequest{})
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(privileges)
	if err != nil {
		return "", err
	}
	privilegesString := b64.StdEncoding.EncodeToString(data)
	return privilegesString, nil
}

// Middleware для обработки привилегий
func PrivilegesMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == healthCheckUrl {
			h.ServeHTTP(w, r)
		} else {
			token := r.Header.Get("Access-Token")
			privileges, err := getPrivileges(token)
			r.Header.Set("Privileges", privileges)
			w.Header().Add("Privileges", privileges)
			if err == nil {
				h.ServeHTTP(w, r)
			} else {
				handlePrivilegesError(err, w)
			}
		}

	})
}

// Обработка ошибок привилегий (создание корректного REST Response)
func handlePrivilegesError(err error, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnauthorized)
	errorObject := map[string]string{
		"error":   err.Error(),
		"message": helpers.ErrSettings[codes.Unauthenticated].Message,
	}
	json.NewEncoder(w).Encode(errorObject)
}

// Перехватывает запрос, добавляет в headers необъодимые данные для микросервисов (ID Запроса, токен, привилегии и т.д.)
func (c *ContentServer) MiddlewaresMidHandler(h http.Handler) http.Handler {
	m := &helpers.Middleware{}
	h = m.LoggerMiddleware(h)
	h = PrivilegesMiddleware(h)
	return m.MiddlewaresHandler(h)
}
