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

	_rootRouter "github.com/sainak/bitsb/root/delivery/http/router"
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

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Recoverer)
	// Important: Chi has a middleware stack and thus it is important to put the
	// Sentry handler on the appropriate place. If using middleware.Recoverer,
	// the Sentry middleware must come afterwards (and configure it with
	// Repanic: true).
	r.Use(sentryMiddleware.Handle)

	_rootRouter.RegisterRoutes(r)

	if viper.GetBool("SERVER_DEBUG") {
		r.Mount("/debug", middleware.Profiler())
	}

	timeout := viper.GetInt("SERVER_TIMEOUT")

	server := &http.Server{
		Addr:              ":" + viper.GetString("WEBSITE_PORT"),
		Handler:           r,
		ReadHeaderTimeout: time.Duration(timeout) * time.Second,
	}
	logrus.Println("Listening on: http://0.0.0.0" + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		logrus.Println(err)
	}
}
