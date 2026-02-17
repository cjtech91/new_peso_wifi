package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjtech-nads/new_peso_wifi/internal/hardware"
)

var detectedBoard hardware.BoardConfig
var clientTemplate *template.Template
var adminTemplate *template.Template

func main() {
	board, err := hardware.DetectBoard()
	if err != nil {
		log.Fatalf("detect board: %v", err)
	}
	detectedBoard = board

	tmplPath := filepath.Join(workDir(), "web", "client_portal.html")
	if t, err := template.ParseFiles(tmplPath); err != nil {
		log.Printf("client template: %v", err)
	} else {
		clientTemplate = t
	}
	adminPath := filepath.Join(workDir(), "web", "admin.html")
	if t, err := template.ParseFiles(adminPath); err != nil {
		log.Printf("admin template: %v", err)
	} else {
		adminTemplate = t
	}

	fmt.Printf("Detected board: %s (%s)\n", board.Name, board.ID)
	if !board.HasGPIO {
		fmt.Println("GPIO not available on this platform, running in simulation mode")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", clientPortalHandler)
	mux.HandleFunc("/voucher", voucherHandler)
	mux.HandleFunc("/insert-coin", insertCoinHandler)
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

	if clientTemplate != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = clientTemplate.Execute(w, nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no'><title>Piso WiFi Portal</title><style>body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;background:#f5f6fa;color:#2d3436} .card{max-width:360px;margin:16px auto;background:#fff;border-radius:12px;box-shadow:0 4px 6px rgba(0,0,0,.1);overflow:hidden} .head{padding:12px 16px;border-bottom:1px solid #eee;font-weight:700} .content{padding:16px} .row{margin:8px 0} .btn{display:block;width:100%;padding:12px;border:none;border-radius:8px;background:#2ecc71;color:#fff;font-weight:700} .form{display:flex;gap:8px} .input{flex:1;padding:10px;border:1px solid #ddd;border-radius:8px} .submit{padding:10px 16px;border:none;border-radius:8px;background:#3498db;color:#fff;font-weight:700}</style></head><body>")
	fmt.Fprintf(w, "<div class='card'><div class='head'>Piso WiFi Client Portal</div><div class='content'>")
	fmt.Fprintf(w, "<div class='row'>Board: %s (%s)</div>", detectedBoard.Name, detectedBoard.ID)
	fmt.Fprintf(w, "<div class='row'><form method='POST' action='/insert-coin'><button class='btn' type='submit'>Insert Coin</button></form></div>")
	fmt.Fprintf(w, "<div class='row form'><form method='POST' action='/voucher' style='display:flex;gap:8px;width:100%%'><input class='input' name='code' placeholder='Enter Voucher Code...' autofocus><button class='submit' type='submit'>Submit</button></form></div>")
	fmt.Fprintf(w, "</div></div></body></html>")
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

func insertCoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Insert Coin</title></head><body>")
	fmt.Fprintf(w, "<h1>Coin Registered</h1>")
	fmt.Fprintf(w, "<p>A simulated coin insert was received.</p>")
	fmt.Fprintf(w, "<a href='/'>Back to portal</a>")
	fmt.Fprintf(w, "</body></html>")
}

func adminPortalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if adminTemplate != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = adminTemplate.Execute(w, nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=1.0'><title>Piso WiFi Admin</title><style>body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;background:#f5f6fa;color:#111827} .card{max-width:720px;margin:16px auto;background:#fff;border-radius:12px;box-shadow:0 4px 6px rgba(0,0,0,.1);overflow:hidden} .head{padding:12px 16px;border-bottom:1px solid #eee;font-weight:700} .content{padding:16px} .row{margin:8px 0} .btn{display:inline-block;padding:10px 14px;border:none;border-radius:10px;background:#2563eb;color:#fff;font-weight:700;text-decoration:none}</style></head><body>")
	fmt.Fprintf(w, "<div class='card'><div class='head'>Piso WiFi Admin</div><div class='content'>")
	fmt.Fprintf(w, "<div class='row'>Board: %s (%s)</div>", detectedBoard.Name, detectedBoard.ID)
	fmt.Fprintf(w, "<div class='row'><a class='btn' href='/admin/status'>View JSON Status</a> <a class='btn' href='/'>Client Portal</a></div>")
	fmt.Fprintf(w, "</div></div></body></html>")
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

func workDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}
