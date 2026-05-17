package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FileWriteEvent struct {
	Tool          string `json:"tool"`
	FilePath      string `json:"file_path"`
	Timestamp     string `json:"timestamp"`
	ContentLength int    `json:"content_length"`
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var event FileWriteEvent
	if err := json.Unmarshal(body, &event); err != nil {
		// Print raw payload when it doesn't match the expected shape
		fmt.Printf("[%s] RAW PAYLOAD: %s\n", time.Now().Format(time.RFC3339), string(body))
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Printf(
		"[%s] FILE WRITE  tool=%-6s  path=%s  bytes=%d\n",
		time.Now().Format(time.RFC3339),
		event.Tool,
		event.FilePath,
		event.ContentLength,
	)

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/claude/create", handleCreate)

	addr := ":8080"
	log.Printf("listening on %s — waiting for file-write events\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
