package handler

import (
	"context"
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/normanjaeckel/Jamora/pkg/model"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed index.html
var index []byte

//go:embed assets/htmx/htmx.min.js
var htmx []byte

func MainPage(w http.ResponseWriter, req *http.Request) {
	if _, err := w.Write(index); err != nil {
		log.Printf("Error: Writing index file to response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Htmx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	if _, err := w.Write(htmx); err != nil {
		log.Printf("Error: Writing htmx file to response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type CampaignHandler struct {
	ctx context.Context
	db  *sql.DB
	t   *template.Template
}

func NewCampaignHandler(ctx context.Context, db *sql.DB) CampaignHandler {
	t := template.Must(template.ParseFS(templateFS, "templates/campaign.html"))
	return CampaignHandler{ctx, db, t}
}

func (h *CampaignHandler) List(w http.ResponseWriter, req *http.Request) {
	rows, err := h.db.QueryContext(h.ctx, "SELECT id, title, description FROM campaigns")
	if err != nil {
		log.Printf("Error: Select campaings from database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var campaigns []model.Campaign
	for rows.Next() {
		var c model.Campaign
		if err := rows.Scan(&c.Id, &c.Title, &c.Description); err != nil {
			log.Printf("Error: Scan campaigns from database query response: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		campaigns = append(campaigns, c)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.t.ExecuteTemplate(w, "list", campaigns); err != nil {
		log.Printf("Error: Execute template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *CampaignHandler) CreateForm(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := h.t.ExecuteTemplate(w, "create form", nil); err != nil {
		log.Printf("Error: Execute template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *CampaignHandler) Create(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	if !req.PostForm.Has("title") {
		http.Error(w, "Missing title", http.StatusBadRequest)
		return
	}
	if _, err := h.db.ExecContext(
		h.ctx,
		"INSERT INTO campaigns (title, description) VALUES ($1, $2)",
		req.PostForm.Get("title"),
		req.PostForm.Get("description"),
	); err != nil {
		log.Printf("Error: Write new campaign into database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.List(w, req)
}

func (h *CampaignHandler) Detail(w http.ResponseWriter, req *http.Request) {

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid campaign id", http.StatusBadRequest)
		return
	}

	row := h.db.QueryRowContext(h.ctx, "SELECT id, title, description FROM campaigns WHERE id=$1", id)

	var c model.Campaign
	if err := row.Scan(&c.Id, &c.Title, &c.Description); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unknown campaign id", http.StatusBadRequest)
			return
		}
		log.Printf("Error: Scan campaigns from database query response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.t.ExecuteTemplate(w, "detail", c); err != nil {
		log.Printf("Error: Execute template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
