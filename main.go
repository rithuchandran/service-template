package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"hotels-service-template/hotel"
	"hotels-service-template/hotel_handler"
	"hotels-service-template/route"

	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repo := hotel.NewRepository(getDb())

	expediaClient := hotel.NewClient("https://test.ean.com/2.2")
	regionService := hotel.NewRegionService(repo, expediaClient)
	regionHandler := hotel_handler.NewRegionHandler(regionService)
	router := route.New(mux.NewRouter())
	router.Configure(regionHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Wrap(route.SetContentTypeHeader),
	}
	start(server)
}

func getDb() *sql.DB {
	db, err := sql.Open("postgres", "dbname=hotels-poc user=postgres password='password' host=localhost port=5432 sslmode=disable")
	if err != nil {
		fmt.Println("db open error", err)
		panic(err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Ping failed", err)
		panic(err)
	}
	return db
}

func start(server *http.Server) {
	go func() {
		fmt.Printf("Starting server on Port: %v", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	gracefulStop(server)
}

//listens for quit, terminate and interrupt signals and shuts the server gracefully without interrupting any active connections
func gracefulStop(server *http.Server) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stop

	fmt.Printf("Shutting the server down...")
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Server stopped")
	}
}
