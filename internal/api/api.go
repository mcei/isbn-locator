package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"go.mongodb.org/mongo-driver/mongo"

	"isbn-locator/internal"
	"isbn-locator/library"
	"isbn-locator/storage"
)

type Handler struct {
	lib *library.Library
}

func NewHandler(lib *library.Library) *Handler {
	return &Handler{
		lib: lib,
	}
}

func Run(client *mongo.Client) {
	db := client.Database("books")
	repo := storage.NewBookStorage(db)
	lib := library.NewLibrary(repo)
	router := Router(lib)
	port := os.Getenv("SERVER_PORT")
	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Router(lib *library.Library) *chi.Mux {
	handler := NewHandler(lib)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(1, time.Second))

	r.Get("/books/{id}", handler.GetBook)
	r.Post("/books", handler.AddBook)
	r.Put("/books/{id}", handler.UpdateBook)
	r.Delete("/books/{id}", handler.RemoveBook)

	return r
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b library.Book
	id := chi.URLParam(r, "id")
	b, err := h.lib.GetBookInfo(ctx, id)
	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(b)
}

func (h *Handler) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b library.Book
	_ = json.NewDecoder(r.Body).Decode(&b)
	err := internal.CheckISBN(b.ISBN)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	err = h.lib.AddBookInfo(ctx, b)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	w.WriteHeader(201)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var b library.Book
	_ = json.NewDecoder(r.Body).Decode(&b)
	err := h.lib.UpdateBookInfo(ctx, id, b)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	w.WriteHeader(204)
}

func (h *Handler) RemoveBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	err := h.lib.RemoveBookInfo(ctx, id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	w.WriteHeader(200)
}
