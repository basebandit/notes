package store

import "errors"

var (
	// ErrNotFound means the requested entity is not found.
	ErrNotFound = errors.New("store: item not found")
	// ErrConflict means the operation failed because of a conflict between entities.
	ErrConflict = errors.New("store: item conflict")
)

//Store is a our application's data store interface. It abstracts the different logic for different databases that the user might
//decide to choose/implement.
type Store interface {
	Notes() NoteStore
}

//NoteStore is our application's notes data store interface.We are abstracting the note's table.So it doesn't matter whether you are
//using mongodb collections or postgesql tables or mysql tables as long as you satisfy this interface you can mainipulate the data in the
//underlying table.
type NoteStore interface {
	New(name, author, content string) (int64, error)
	Get(id int64) (*Note, error)
	GetByAuthor(author string, offset, limit int) ([]*Note, int, error)
	GetLatest(offset, limit int) ([]*Note, int, error)
	SetContent(id int64, content string) error
	SetName(id int64, name string) error
	SetAuthor(id int64, author string) error
	Delete(id int64) error
}
