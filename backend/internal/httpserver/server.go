package httpserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sendtolinux/internal/dbussvc"
)

type Server struct {
	svc *dbussvc.Service
}

func Start(svc *dbussvc.Service) (*http.Server, error) {
	bind := getenvDefault("STL_BIND", "0.0.0.0")
	port := getenvInt("STL_PORT", 8000)
	addr := net.JoinHostPort(bind, strconv.Itoa(port))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	host := hostnameOrLocal()
	url := fmt.Sprintf("http://%s:%d/", host, actualPort)
	svc.SetStatus(url, uint32(actualPort), true)

	h := &Server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/text", h.handleText)

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("http server error: %v", err)
		}
		svc.SetStatus(url, uint32(actualPort), false)
	}()

	log.Printf("HTTP server listening on %s", addr)
	return server, nil
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!doctype html>
<html>
<head><meta charset="utf-8"><title>Send to Linux</title></head>
<body>
  <h1>Send to Linux</h1>
  <form action="/text" method="post">
    <textarea name="text" rows="12" cols="60" placeholder="Paste text"></textarea><br>
    <button type="submit">Send</button>
  </form>
</body>
</html>`)
}

func (s *Server) handleText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxMB := getenvInt("STL_MAX_UPLOAD_MB", 100)
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxMB)*1024*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	if text == "" {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			text = strings.TrimSpace(string(body))
		}
	}
	if text == "" {
		http.Error(w, "empty text", http.StatusBadRequest)
		return
	}

	saveDir, err := resolveSaveDir()
	if err != nil {
		http.Error(w, "cannot resolve save dir", http.StatusInternalServerError)
		return
	}

	filename := uniqueTextFilename(saveDir, time.Now())
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		http.Error(w, "cannot create save dir", http.StatusInternalServerError)
		return
	}
	if err := os.WriteFile(filename, []byte(text), 0o644); err != nil {
		http.Error(w, "cannot save text", http.StatusInternalServerError)
		return
	}

	item := dbussvc.RecentItem{
		ID:    strconv.FormatInt(time.Now().UnixNano(), 10),
		Type:  "text",
		Value: text,
		Size:  uint32(len([]byte(text))),
	}
	s.svc.AddRecent(item)
	s.svc.EmitItemReceived(item)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "OK")
}

func resolveSaveDir() (string, error) {
	if dir := os.Getenv("STL_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Downloads", "SendToLinux"), nil
}

func uniqueTextFilename(dir string, now time.Time) string {
	base := fmt.Sprintf("text-%s.txt", now.Format("20060102-150405"))
	path := filepath.Join(dir, base)
	if _, err := os.Stat(path); err != nil {
		return path
	}
	for i := 1; i < 1000; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("text-%s-%d.txt", now.Format("20060102-150405"), i))
		if _, err := os.Stat(candidate); err != nil {
			return candidate
		}
	}
	return path
}

func getenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func hostnameOrLocal() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "localhost"
	}
	return host
}
