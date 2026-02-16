package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/cjtech-nads/new_peso_wifi/internal/hardware"
)

var detectedBoard hardware.BoardConfig
var clientTemplate *template.Template

func main() {
	board, err := hardware.DetectBoard()
	if err != nil {
		log.Fatalf("detect board: %v", err)
	}
	detectedBoard = board

	t, err := template.ParseFiles("web/client_portal.html")
	if err != nil {
		log.Fatalf("parse client template: %v", err)
	}
	clientTemplate = t

	fmt.Printf("Detected board: %s (%s)\n", board.Name, board.ID)
	if !board.HasGPIO {
		fmt.Println("GPIO not available on this platform, running in simulation mode")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", clientPortalHandler)
	mux.HandleFunc("/voucher", voucherHandler)
	mux.HandleFunc("/admin", adminPortalHandler)
	mux.HandleFunc("/admin/status", adminStatusHandler)

	addr := ":8080"
	if v := os.Getenv("PISO_HTTP_ADDR"); v != "" {
		addr = v
	}

	log.Printf("starting piso-wifi HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

func clientPortalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := clientTemplate.Execute(w, nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func voucherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	safeCode := html.EscapeString(code)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Voucher</title></head><body>")
	fmt.Fprintf(w, "<h1>Voucher Submitted</h1>")
	fmt.Fprintf(w, "<p>Code: %s</p>", safeCode)
	fmt.Fprintf(w, "<p>Voucher validation and time credit logic is not implemented yet.</p>")
	fmt.Fprintf(w, "<a href='/'>Back to portal</a>")
	fmt.Fprintf(w, "</body></html>")
}

func adminPortalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Piso WiFi Admin</title></head><body>")
	fmt.Fprintf(w, "<h1>Piso WiFi Admin Portal</h1>")
	fmt.Fprintf(w, "<p>Board: %s (%s)</p>", detectedBoard.Name, detectedBoard.ID)
	fmt.Fprintf(w, "<p><a href='/admin/status'>View JSON status</a></p>")
	fmt.Fprintf(w, "</body></html>")
}

func adminStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(struct {
		Board hardware.BoardConfig `json:"board"`
	}{
		Board: detectedBoard,
	})
}
