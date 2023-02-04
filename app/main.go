package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_rootRouter "github.com/sainak/bitsb/root/delivery/http/router"
)

func init() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(err)
	}

	if viper.GetBool("SERVER_DEBUG") {
		log.Println("SERVER running in debug mode")
	}
}

func main() {
	dsn := viper.GetString("DB_DSN")
	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
