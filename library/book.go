package library

import (
	"context"
	"errors"
)

type Book struct {
	ISBN   string `bson:"isbn"`
	Title  string `bson:"title"`
	Author string `bson:"author"`
	Year   string `bson:"year"`
}

var ErrBookNotFound = errors.New("book not found")

type BookRepository interface {
	// Fetch возвращает информацию о книге по ISBN
	// в случае если книги с таким идентификатором нет,
	// то возвращает ErrBookNotFound
	Fetch(ctx context.Context, isbn string) (Book, error)
	// Store сохраняет информацию о книге
	Store(ctx context.Context, b Book) error
	// Update обновляет информацию о книге
	Update(ctx context.Context, isbn string, b Book) error
	// Remove удаляет информацию о книге по ISBN
	Remove(ctx context.Context, isbn string) error
}

type Library struct {
	books BookRepository
}

func NewLibrary(books BookRepository) *Library {
	return &Library{
		books: books,
	}
}

// TODO ks: реальные кейсы из предметки:
//   пользователь: получает информацию о книге по ISBN
//   библиотекарь: добавляет книгу, обновляет, удаляет

// GetBookInfo возвращает информацию о книге по ISBN
// в случае если книги с таким идентификатором нет, то возвращает ErrBookNotFound
func (l *Library) GetBookInfo(ctx context.Context, isbn string) (Book, error) {
	b, err := l.books.Fetch(ctx, isbn)
	return b, err
}

// AddBookInfo сохраняет информацию о книге
func (l *Library) AddBookInfo(ctx context.Context, b Book) error {
	return l.books.Store(ctx, b)
}

// UpdateBookInfo обновляет информацию о книге
func (l *Library) UpdateBookInfo(ctx context.Context, isbn string, b Book) error {
	return l.books.Update(ctx, isbn, b)
}

// RemoveBookInfo удаляет информацию о книге по ISBN
func (l *Library) RemoveBookInfo(ctx context.Context, isbn string) error {
	return l.books.Remove(ctx, isbn)
}
