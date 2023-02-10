package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_bitsbRouter "github.com/sainak/bitsb/bitsb/delivery/http/router"
	_bitsbRepo "github.com/sainak/bitsb/bitsb/repo/postgres"
	_bitsbService "github.com/sainak/bitsb/bitsb/service"
	_rootRouter "github.com/sainak/bitsb/root/delivery/http/router"
	_userRouter "github.com/sainak/bitsb/users/delivery/http/router"
	_userRepo "github.com/sainak/bitsb/users/repo/postgres"
	_userService "github.com/sainak/bitsb/users/service"
	"github.com/sainak/bitsb/utils/jwt"
	middl "github.com/sainak/bitsb/utils/middleware"
)

var (
	version     = "nil"
	environment = ""
)

func init() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Print(err)
	}

	environment = viper.GetString("ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}
	environment = strings.ToLower(environment)

	if environment != "local" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
		logrus.SetReportCaller(true)
	}

	if viper.GetBool("SERVER_DEBUG") {
		logrus.Println("SERVER running in debug mode")
	}
}

func main() {
	logrus.Println("Version: ", version)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("SENTRY_DSN"),
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
		ServerName:       "bitsb",
		Release:          "bitsb@" + version, //-ldflags='-X main.release=VALUE'
		Dist:             "",
		Environment:      environment,
	})
	if err != nil {
		logrus.Errorf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	dsn := viper.GetString("DB_DSN")
	logrus.Info("DB_DSN: ", dsn)
	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	timeout := time.Duration(viper.GetInt("SERVER_TIMEOUT")) * time.Second

	jwtInstance := jwt.New(
		viper.GetString("JWT_SECRET"),
		viper.GetString("JWT_EXPIRY"),
		viper.GetString("JWT_REFRESH_EXPIRY"),
	)

	r := chi.NewRouter()

	r.Use(
		middleware.Maybe(middleware.CleanPath, func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/debug/")
		}),
		middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/debug/")
		}),
	)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Important: Chi has a middleware stack and thus it is important to put the
	// Sentry handler on the appropriate place. If using middleware.Recoverer,
	// the Sentry middleware must come afterwards (and configure it with
	// Repanic: true).
	r.Use(sentryMiddleware.Handle)

	_rootRouter.RegisterRoutes(r)

	userRepo := _userRepo.NewUserRepository(dbConn)
	locationRepo := _bitsbRepo.NewLocationRepository(dbConn)
	busRouteRepo := _bitsbRepo.NewBusRouteRepository(dbConn)

	userService := _userService.NewUserService(userRepo, jwtInstance, timeout)
	locationService := _bitsbService.NewLocationService(locationRepo)
	busRouteService := _bitsbService.NewBusRouteService(busRouteRepo, locationRepo)

	jwtMiddleware := middl.JWTAuth(jwtInstance, userRepo)

	_userRouter.RegisterRoutes(r, userService, jwtMiddleware)
	_bitsbRouter.RegisterLocationRoutes(r, locationService, jwtMiddleware)
	_bitsbRouter.RegisterBusRouteRoutes(r, busRouteService, jwtMiddleware)

	if viper.GetBool("SERVER_DEBUG") {
		r.Mount("/debug", middleware.Profiler())
	}

	server := &http.Server{
		Addr:              ":" + viper.GetString("WEBSITE_PORT"),
		Handler:           r,
		ReadHeaderTimeout: timeout,
	}
	logrus.Println("Listening on: http://0.0.0.0" + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		logrus.Println(err)
	}
}
