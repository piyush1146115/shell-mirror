package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handleRequests() {
	MyRouter.HandleFunc("/", homePage)

	MyRouter.HandleFunc("/execute", execute).Methods("POST")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Homepage!")
}

func execute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input CommandInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println(fmt.Errorf("failed to decode json request: %w", err))
		http.Error(w, fmt.Sprintf("failed to decode json request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if !validator(input.Command) {
		log.Println("Invalid command")
		http.Error(w, fmt.Sprintf("Invalid command: %s", input.Command), http.StatusForbidden)
		return
	}

	log.Println("Got command:", input.Command)

	output, err := executeCommand(input.Command)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute command: %s. Got errror: %s", input.Command, err.Error()), http.StatusInternalServerError)
		return
	}
	response := CommandOutput{Output: output}

	log.Println("Output: ", response)

	json.NewEncoder(w).Encode(response)
}

func validator(cmd string) bool {
	if strings.Contains(cmd, "sudo") {
		return false
	}

	return true
}