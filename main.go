package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/signal"

	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"os/exec"
	"time"
)

type CommandOutput struct {
	Output string `json:"output"`
}

type CommandInput struct {
	Command string `json:"command"`
}

func main() {
	CreateServer("8088")
	StartServer()
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

var MyRouter = mux.NewRouter().StrictSlash(true)
var wait time.Duration
var server *http.Server

func CreateServer(port string) {
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "The duration for which our server gracefully wait for existing connections to finish")
	flag.Parse()

	port = "127.0.0.1:" + port
	server = &http.Server{
		Addr:         port,
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

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Homepage!")
}

func handleRequests() {
	MyRouter.HandleFunc("/", homePage)

	MyRouter.HandleFunc("/execute", execute).Methods("POST")
}

func execute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input CommandInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println(fmt.Errorf("failed to decode json request: %w", err))
		http.Error(w, fmt.Sprintf("failed to decode json request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Println("Got command:", input.Command)

	output, err := executeCommand(input.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := CommandOutput{Output: output}

	log.Println("Output: ", response)

	json.NewEncoder(w).Encode(response)
}

func init() {
	handleRequests()
}
