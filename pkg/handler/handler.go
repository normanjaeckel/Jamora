package handler

import (
	"embed"
	"html/template"
	"maps"
	"net/http"

	"github.com/normanjaeckel/Jamora/pkg/model"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed index.html
var index []byte

//go:embed assets/htmx.min.js
var htmx []byte

func MainPage(w http.ResponseWriter, req *http.Request) {
	if _, err := w.Write(index); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Htmx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	if _, err := w.Write(htmx); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type CampaignHandler struct {
	m *model.Model
	t *template.Template
}

func NewCampaignHandler(m *model.Model) CampaignHandler {
	t := template.Must(template.ParseFS(templateFS, "templates/campaign.html"))
	return CampaignHandler{m, t}
}

func (h *CampaignHandler) List(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := h.t.ExecuteTemplate(w, "list", h.m); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *CampaignHandler) CreateForm(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := h.t.ExecuteTemplate(w, "create form", nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *CampaignHandler) Create(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
	}
	if !req.PostForm.Has("title") {
		http.Error(w, "Missing title", http.StatusBadRequest)
		return
	}
	var newId int64
	for k := range maps.Keys(*h.m) {
		newId = max(k, newId)
	}
	x := *h.m
	x[newId+1] = model.Campaign{Title: req.PostForm.Get("title"), Description: req.PostForm.Get("description")}
	h.m = &x

	h.List(w, req)
}
