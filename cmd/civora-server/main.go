package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type workflowRunEvent struct {
	Action      string `json:"action"`
	Repository  struct { FullName string `json:"full_name"` } `json:"repository"`
	WorkflowRun struct { ID int64 `json:"id"`; Name string `json:"name"`; Conclusion string `json:"conclusion"`; Status string `json:"status"` } `json:"workflow_run"`
}

func main() {
	addr := getenv("CIVORA_ADDR", ":8080")
	secret := os.Getenv("CIVORA_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("CIVORA_WEBHOOK_SECRET is required")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		body, err := io.ReadAll(io.LimitReader(r.Body, 2<<20))
		if err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		if !validSignature(body, r.Header.Get("X-Hub-Signature-256"), secret) {
			http.Error(w, "invalid signature", http.StatusUnauthorized); return
		}
		event := r.Header.Get("X-GitHub-Event")
		if event != "workflow_run" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		var payload workflowRunEvent
		if err := json.Unmarshal(body, &payload); err != nil { http.Error(w, "invalid JSON", http.StatusBadRequest); return }
		log.Printf("workflow_run action=%s repo=%s run_id=%d workflow=%q status=%s conclusion=%s", payload.Action, payload.Repository.FullName, payload.WorkflowRun.ID, payload.WorkflowRun.Name, payload.WorkflowRun.Status, payload.WorkflowRun.Conclusion)
		w.WriteHeader(http.StatusAccepted)
	})
	log.Printf("Civora webhook server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func validSignature(body []byte, got, secret string) bool {
	if len(got) != len("sha256=")+64 || got[:7] != "sha256=" { return false }
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	want := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(got), []byte(want))
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" { return v }
	return fallback
}
