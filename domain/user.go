package domain

// TODO ks: сейчас не используется
// добавлен для проработки структуры

const (
	Reader Role = "reader"
	Writer Role = "writer"
)

type Role string

type User struct {
	ID   int
	Role Role
}

// func NewUser
// func ChangeUserRole
