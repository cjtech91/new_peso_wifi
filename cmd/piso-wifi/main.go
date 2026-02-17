package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/cjtech-nads/new_peso_wifi/internal/hardware"
)

var detectedBoard hardware.BoardConfig
var clientTemplate *template.Template
var adminTemplate *template.Template
var adminSessions = map[string]bool{}
var adminSessionsMu sync.Mutex

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
	mux.HandleFunc("/admin/login", adminLoginHandler)
	mux.HandleFunc("/admin/logout", adminLogoutHandler)
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

	if !isAdminAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusFound)
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

	if !isAdminAuthenticated(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if t, err := template.ParseFiles(filepath.Join(workDir(), "web", "admin_login.html")); err == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = t.Execute(w, nil)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<!DOCTYPE html><html><head><meta name='viewport' content='width=device-width, initial-scale=1.0'><title>Admin Login</title><style>body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;background:#f5f6fa;color:#111827} .card{max-width:360px;margin:16px auto;background:#fff;border-radius:12px;box-shadow:0 4px 6px rgba(0,0,0,.1)} .head{padding:12px 16px;border-bottom:1px solid #eee;font-weight:700} .content{padding:16px} .field{margin:8px 0} .input{width:100%%;padding:10px;border:1px solid #ddd;border-radius:8px} .btn{width:100%%;padding:12px;border:none;border-radius:10px;background:#2563eb;color:#fff;font-weight:700;margin-top:8px}</style></head><body>")
		fmt.Fprintf(w, "<div class='card'><div class='head'>Admin Login</div><div class='content'><form method='POST' action='/admin/login'><div class='field'><input class='input' name='username' placeholder='Username' autofocus></div><div class='field'><input class='input' name='password' type='password' placeholder='Password'></div><button class='btn' type='submit'>Sign In</button></form></div></div></body></html>")
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		u := r.Form.Get("username")
		p := r.Form.Get("password")
		if u == "admin" && p == "admin" {
			tok := newToken()
			adminSessionsMu.Lock()
			adminSessions[tok] = true
			adminSessionsMu.Unlock()
			http.SetCookie(w, &http.Cookie{
				Name:     "admin_session",
				Value:    tok,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("admin_session"); err == nil {
		adminSessionsMu.Lock()
		delete(adminSessions, c.Value)
		adminSessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func isAdminAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}
	adminSessionsMu.Lock()
	_, ok := adminSessions[c.Value]
	adminSessionsMu.Unlock()
	return ok
}

func newToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
