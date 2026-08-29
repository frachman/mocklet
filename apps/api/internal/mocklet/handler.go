package mocklet

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var pathPattern = regexp.MustCompile(`^/[A-Za-z0-9._~!$&'()*+,;=:@%{}\-/]+$`)
var methods = map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true}

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) http.Handler {
	h := &Handler{repo}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/api/v1/mocks", h.create)
	mux.HandleFunc("/api/v1/mocks/", h.manage)
	mux.HandleFunc("/m/", h.runtime)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

type createRequest struct {
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Body        string            `json:"body"`
	ContentType string            `json:"content_type"`
	StatusCode  int               `json:"status_code"`
	DelayMS     int               `json:"delay_ms"`
	Headers     map[string]string `json:"headers"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var in createRequest
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&in) != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	if in.Name == "" {
		in.Name = "Untitled mock"
	}
	if err := validateEndpoint(&in); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	token, _ := randomToken()
	key, _ := randomToken()
	expires := time.Now().Add(24 * time.Hour)
	// The repository generates its own public key; key is only used to keep token generation uniform.
	_ = key
	e := Endpoint{Method: strings.ToUpper(in.Method), Path: in.Path, StatusCode: in.StatusCode, Body: in.Body, ContentType: in.ContentType, DelayMS: in.DelayMS, Headers: in.Headers}
	m, err := h.repo.Create(r.Context(), in.Name, hashToken(token), expires, e)
	if err != nil {
		http.Error(w, "could not create mock", 500)
		return
	}
	writeJSON(w, 201, map[string]any{"public_key": m.PublicKey, "management_token": token, "name": m.Name, "expires_at": m.ExpiresAt, "endpoint": e})
}

func (h *Handler) manage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "api" || parts[1] != "v1" || parts[2] != "mocks" {
		http.NotFound(w, r)
		return
	}
	ok, err := h.repo.Authenticate(r.Context(), parts[3], hashToken(tokenFromRequest(r)))
	if err != nil || !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	m, err := h.repo.FindByPublicKey(r.Context(), parts[3])
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	writeJSON(w, 200, m)
}

func (h *Handler) runtime(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/m/"), "/", 2)
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	m, err := h.repo.FindByPublicKey(r.Context(), parts[0])
	if err != nil || time.Now().After(m.ExpiresAt) {
		http.NotFound(w, r)
		return
	}
	if m.Endpoint.Method != r.Method || m.Endpoint.Path != "/"+parts[1] {
		http.NotFound(w, r)
		return
	}
	if m.Endpoint.DelayMS > 0 {
		time.Sleep(time.Duration(m.Endpoint.DelayMS) * time.Millisecond)
	}
	for k, v := range m.Endpoint.Headers {
		w.Header().Set(k, v)
	}
	if m.Endpoint.ContentType != "" {
		w.Header().Set("Content-Type", m.Endpoint.ContentType)
	}
	w.WriteHeader(m.Endpoint.StatusCode)
	w.Write([]byte(m.Endpoint.Body))
}

func validateEndpoint(in *createRequest) error {
	in.Method = strings.ToUpper(strings.TrimSpace(in.Method))
	if !methods[in.Method] {
		return &clientError{"method must be GET, POST, PUT, PATCH, or DELETE"}
	}
	if !pathPattern.MatchString(in.Path) || !strings.HasPrefix(in.Path, "/") {
		return &clientError{"path must be an absolute route path"}
	}
	if in.StatusCode == 0 {
		in.StatusCode = 200
	}
	if in.StatusCode < 100 || in.StatusCode > 599 {
		return &clientError{"status_code must be between 100 and 599"}
	}
	if in.DelayMS < 0 || in.DelayMS > 10000 {
		return &clientError{"delay_ms must be between 0 and 10000"}
	}
	if in.ContentType == "" {
		in.ContentType = "application/json"
	}
	if len(in.Body) > 1<<20 {
		return &clientError{"body is limited to 1 MiB"}
	}
	return nil
}

type clientError struct{ message string }

func (e *clientError) Error() string { return e.message }
func randomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func tokenFromRequest(r *http.Request) string {
	if v := r.Header.Get("X-Management-Token"); v != "" {
		return v
	}
	v, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return v
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
