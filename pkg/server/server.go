package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/normanjaeckel/Jamora/pkg/handler"
	"github.com/normanjaeckel/Jamora/pkg/model"
)

func registerHandler(mux *http.ServeMux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, handler)
}

func Run(ctx context.Context) error {
	addr := ":8080"
	dbFile := "database.sqlite"
	mux := http.NewServeMux()

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	createTablesQuery := model.CampaignTableQuery + model.GroupTableQuery
	if _, err := db.ExecContext(ctx, createTablesQuery); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	registerHandler(mux, "GET /{$}", handler.MainPage)
	registerHandler(mux, "GET /assets/htmx.min.js", handler.Htmx)

	campaignHandler := handler.NewCampaignHandler(ctx, db)

	// if _, ok := req.Header[http.CanonicalHeaderKey("HX-Request")]; !ok {
	// 	MainPage(w, req)
	// 	return
	// }

	registerHandler(mux, "GET /campaign", campaignHandler.List)
	registerHandler(mux, "GET /campaign/create-form", campaignHandler.CreateForm)
	registerHandler(mux, "POST /campaign", campaignHandler.Create)
	registerHandler(mux, "GET /campaign/{id}", campaignHandler.Detail)
	registerHandler(mux, "GET /campaign/{id}/update", http.NotFound)
	registerHandler(mux, "POST /campaign/{id}", http.NotFound)
	registerHandler(mux, "GET /campaign/{id}/delete", http.NotFound)
	registerHandler(mux, "DELETE /campaign/{id}", http.NotFound)

	// mux.HandleFunc("GET /class/{id}", http.NotFound)
	// mux.HandleFunc("GET /class/new", http.NotFound)
	// mux.HandleFunc("POST /class", http.NotFound)
	// mux.HandleFunc("GET /class/{id}/update", http.NotFound)
	// mux.HandleFunc("POST /class/{id}", http.NotFound)
	// mux.HandleFunc("GET /class/{id}/delete", http.NotFound)
	// mux.HandleFunc("DELETE /class/{id}", http.NotFound)

	// mux.HandleFunc("GET /pupil/{id}", http.NotFound)
	// mux.HandleFunc("GET /pupil/new", http.NotFound)
	// mux.HandleFunc("POST /pupil", http.NotFound)
	// mux.HandleFunc("GET /pupil/{id}/update", http.NotFound)
	// mux.HandleFunc("POST /pupil/{id}", http.NotFound)
	// mux.HandleFunc("GET /pupil/{id}/delete", http.NotFound)
	// mux.HandleFunc("DELETE /pupil/{id}", http.NotFound)

	// mux.HandleFunc("GET /group/{id}", http.NotFound)
	// mux.HandleFunc("GET /group/new", http.NotFound)
	// mux.HandleFunc("POST /group", http.NotFound)
	// mux.HandleFunc("GET /group/{id}/update", http.NotFound)
	// mux.HandleFunc("POST /group/{id}", http.NotFound)
	// mux.HandleFunc("GET /group/{id}/delete", http.NotFound)
	// mux.HandleFunc("DELETE /group/{id}", http.NotFound)

	// mux.HandleFunc("GET /rule/{id}", http.NotFound)
	// mux.HandleFunc("GET /rule/new", http.NotFound)
	// mux.HandleFunc("POST /rule", http.NotFound)
	// mux.HandleFunc("GET /rule/{id}/update", http.NotFound)
	// mux.HandleFunc("POST /rule/{id}", http.NotFound)
	// mux.HandleFunc("GET /rule/{id}/delete", http.NotFound)
	// mux.HandleFunc("DELETE /rule/{id}", http.NotFound)

	// mux.HandleFunc("GET /assign", http.NotFound)
	// mux.HandleFunc("POST /assign", http.NotFound)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
		//BaseContext: func(net.Listener) context.Context { return ctx },
	}

	// Shutdown logic in a separate goroutine.
	shutdownResult := make(chan error)
	go func() {
		// Wait for the context to be closed.
		<-ctx.Done()

		// Shutdown server
		log.Print("Shutdown server")
		if err := srv.Shutdown(context.WithoutCancel(ctx)); err != nil {
			shutdownResult <- fmt.Errorf("server shutdown: %w", err)
			return
		}
		shutdownResult <- nil
	}()

	// Start server
	log.Printf("Webserver is listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return <-shutdownResult
}
