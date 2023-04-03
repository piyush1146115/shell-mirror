package api

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type CommandOutput struct {
	Output string `json:"output"`
}

type CommandInput struct {
	Command string `json:"command"`
}

var MyRouter = mux.NewRouter().StrictSlash(true)
var wait time.Duration
var server *http.Server

func CreateServer(port string) {
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "The duration for which our server gracefully wait for existing connections to finish")
	flag.Parse()

	server = &http.Server{
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      MyRouter,
	}
}

func StartServer() {
	log.Println("Starting server on Address: ", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	GracefulShutDown()
}

func GracefulShutDown() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	server.Shutdown(ctx)
	log.Println("Shutting down!!")
	os.Exit(0)
}

func init() {
	handleRequests()
}
