package httpserver

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sendtolinux/internal/config"
	"sendtolinux/internal/dbussvc"
	"sendtolinux/internal/version"
)

type Server struct {
	svc      *dbussvc.Service
	template *template.Template
	assetDir string
	assetFS  fs.FS
	version  string
	config   config.Config
}

type pageData struct {
	Message string
	Version string
}

func Start(svc *dbussvc.Service, cfg config.Config) (*http.Server, error) {
	bind := cfg.Bind
	port := cfg.Port
	addr := net.JoinHostPort(bind, strconv.Itoa(port))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port
	host := hostnameOrLocal()
	url := fmt.Sprintf("http://%s:%d/", host, actualPort)
	svc.SetStatus(url, uint32(actualPort), true)
	if qrPath, err := generateQrPng(url); err != nil {
		log.Printf("qr png: %v", err)
	} else {
		svc.SetQrPath(qrPath)
	}

	assetDir := os.Getenv("STL_ASSET_DIR")
	tmpl, assetFS, err := loadTemplateAndFS(assetDir)
	if err != nil {
		return nil, err
	}

	h := &Server{
		svc:      svc,
		template: tmpl,
		assetDir: assetDir,
		assetFS:  assetFS,
		version:  version.Version,
		config:   cfg,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/text", h.handleText)
	mux.HandleFunc("/file", h.handleFile)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assetFS))))

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("http server error: %v", err)
		}
		svc.SetStatus(url, uint32(actualPort), false)
		svc.SetQrPath("")
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
	if err := s.template.Execute(w, pageData{Version: s.version}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (s *Server) handleText(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxMB := s.config.MaxUploadMB
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

	saveDir, err := resolveSaveDir(s.config)
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
		Size:  uint64(len([]byte(text))),
	}
	s.svc.AddRecent(item)
	s.svc.EmitItemReceived(item)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.template.Execute(w, pageData{Message: "Text sent successfully.", Version: s.version}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxMB := s.config.MaxUploadMB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxMB)*1024*1024)
	if err := r.ParseMultipartForm(int64(maxMB) * 1024 * 1024); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	saveDir, err := resolveSaveDir(s.config)
	if err != nil {
		http.Error(w, "cannot resolve save dir", http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		http.Error(w, "cannot create save dir", http.StatusInternalServerError)
		return
	}

	fileHeaders := multipartHeaders(r.MultipartForm)
	if len(fileHeaders) == 0 {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}

	if len(fileHeaders) == 1 {
		header := fileHeaders[0]
		filename := filepath.Base(header.Filename)
		if filename == "" || filename == "." {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		item, err := s.saveUploadedFile(saveDir, header, filename, 0)
		if err != nil {
			http.Error(w, "cannot save file", http.StatusInternalServerError)
			return
		}
		s.svc.AddRecent(item)
		s.svc.EmitItemReceived(item)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := s.template.Execute(w, pageData{Message: "File sent successfully.", Version: s.version}); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
		return
	}

	var firstItem dbussvc.RecentItem
	firstItemSet := false
	for i, header := range fileHeaders {
		filename := filepath.Base(header.Filename)
		if filename == "" || filename == "." {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		item, err := s.saveUploadedFile(saveDir, header, filename, i)
		if err != nil {
			http.Error(w, "cannot save file", http.StatusInternalServerError)
			return
		}
		s.svc.AddRecent(item)
		if !firstItemSet {
			firstItem = item
			firstItemSet = true
		}
	}

	if firstItemSet {
		s.svc.EmitItemReceived(firstItem)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	message := "File sent successfully."
	if len(fileHeaders) > 1 {
		message = fmt.Sprintf("%d files sent successfully.", len(fileHeaders))
	}
	if err := s.template.Execute(w, pageData{Message: message, Version: s.version}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func multipartHeaders(form *multipart.Form) []*multipart.FileHeader {
	if form == nil {
		return nil
	}

	total := 0
	for _, group := range form.File {
		total += len(group)
	}
	if total == 0 {
		return nil
	}

	headers := make([]*multipart.FileHeader, 0, total)
	for _, group := range form.File {
		headers = append(headers, group...)
	}
	return headers
}

func (s *Server) saveUploadedFile(saveDir string, header *multipart.FileHeader, filename string, index int) (dbussvc.RecentItem, error) {
	file, err := header.Open()
	if err != nil {
		return dbussvc.RecentItem{}, err
	}
	defer file.Close()

	targetPath := uniqueFilePath(saveDir, filename)
	out, err := os.Create(targetPath)
	if err != nil {
		return dbussvc.RecentItem{}, err
	}
	defer out.Close()

	written, err := io.Copy(out, file)
	if err != nil {
		return dbussvc.RecentItem{}, err
	}

	item := dbussvc.RecentItem{
		ID:    fmt.Sprintf("%d-%d", time.Now().UnixNano(), index),
		Type:  "file",
		Value: targetPath,
		Size:  uint64(written),
	}
	return item, nil
}

func resolveSaveDir(cfg config.Config) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	downloadsDir := filepath.Join(home, "Downloads")
	if cfg.Dir == "" {
		return filepath.Join(downloadsDir, "SendToLinux"), nil
	}
	folder := filepath.Base(strings.TrimSpace(cfg.Dir))
	if folder == "" || folder == "." || folder == ".." || folder == string(os.PathSeparator) {
		folder = "SendToLinux"
	}
	return filepath.Join(downloadsDir, folder), nil
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

func uniqueFilePath(dir, filename string) string {
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); err != nil {
		return path
	}

	ext := filepath.Ext(filename)
	stem := strings.TrimSuffix(filename, ext)
	for i := 1; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d%s", stem, i, ext)
		path = filepath.Join(dir, candidate)
		if _, err := os.Stat(path); err != nil {
			return path
		}
	}
	return path
}

func hostnameOrLocal() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "localhost"
	}
	if strings.Contains(host, ".") {
		return host
	}
	return host + ".local"
}

func generateQrPng(url string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil || cacheDir == "" {
		cacheDir = os.TempDir()
	}
	dir := filepath.Join(cacheDir, "SendToLinux")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "qr.png")
	cmd := exec.Command("qrencode", "-o", path, "-s", "6", "-m", "1", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message != "" {
			return "", fmt.Errorf("qrencode failed: %w: %s", err, message)
		}
		return "", fmt.Errorf("qrencode failed: %w", err)
	}
	return path, nil
}

//go:embed assets/index.html assets/static/*
var embeddedAssets embed.FS

func loadTemplateAndFS(assetDir string) (*template.Template, fs.FS, error) {
	if assetDir == "" {
		tmpl, err := template.ParseFS(embeddedAssets, "assets/index.html")
		if err != nil {
			return nil, nil, err
		}
		staticFS, err := fs.Sub(embeddedAssets, "assets/static")
		if err != nil {
			return nil, nil, err
		}
		return tmpl, staticFS, nil
	}

	tmpl, err := template.ParseFiles(filepath.Join(assetDir, "index.html"))
	if err != nil {
		return nil, nil, err
	}
	staticDir := filepath.Join(assetDir, "static")
	return tmpl, os.DirFS(staticDir), nil
}
