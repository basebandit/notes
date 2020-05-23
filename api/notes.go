package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/parish/notes/store"
)

//NoteResource defines the requirements for using a note resource.
type NoteResource struct {
	Store  store.Store
	Logger *log.Logger
}

//NewNoteResource returns a configured note resource.
func NewNoteResource(store store.Store, logger *log.Logger) *NoteResource {
	return &NoteResource{
		Store: store,
	}
}

//Router provides necessary routes for note resources.
func (n *NoteResource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", n.handleGetNotes)

	return r
}

func (n NoteResource) handleGetNotes(w http.ResponseWriter, r *http.Request) {

	var err error

	offset := 0
	offsetParam := r.URL.Query().Get("offset")
	if offsetParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			n.Logger.Printf("error: BadRequest reason: invalid offset %d", offset)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	limit := 10
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 1 || limit > 1000 {
			n.Logger.Printf("error: %s; status: %s;  reason: invalid limit %d", err, http.StatusText(http.StatusBadRequest), limit)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	var notes []*store.Note
	var count int

	notes, count, err = n.Store.Notes().GetLatest(offset, limit)
	if err != nil {
		n.Logger.Printf("error: %s; status: %s; reason: server error", err, http.StatusText(http.StatusInternalServerError))
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Notes []*store.Note `json:"notes"`
		Count int           `json:"count"`
	}{
		Notes: notes,
		Count: count,
	}

	json.NewEncoder(w).Encode(response)
}

//TODO Add more handlers - Hint Map to HTTP VERBS
