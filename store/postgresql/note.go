package postgresql

import (
	"database/sql"
	"time"

	"github.com/parish/notes/store"
)

type noteStore struct {
	db *sql.DB
}

//New creates a new note.
func (n *noteStore) New(name, author, content string) (int64, error) {
	var id int64
	now := time.Now()

	tx, err := n.db.Begin()
	if err != nil {
		return 0, err
	}

	err = n.db.QueryRow(`insert into notes(name,author,content,created_at) values ($1,$2,$3,$4) returning id`, name, author, content, now).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}

//Get retrieves a note that matches the given id
func (n *noteStore) Get(id int64) (*store.Note, error) {
	row := n.db.QueryRow(`select id ,name, author,content, created_at from notes where deleted=false and id=$1`, id)
	return n.scanNote(row)
}

//GetByAuthor queries searches for a given note entity using its name and limits the number of rows returned using the given limit and offset for pagination.
func (n *noteStore) GetByAuthor(author string, offset, limit int) ([]*store.Note, int, error) {
	var count int
	err := n.db.QueryRow(`select count(*) from notes where deleted=false and  author=$1`, author).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	if limit <= 0 || offset > count {
		return []*store.Note{}, count, nil
	}

	rows, err := n.db.Query(`select id,name,author,content,created_at from notes where deleted=false and name=$1 order by created_at, id limit $2 offset $3`, author, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := []*store.Note{}
	for rows.Next() {
		note, err := n.scanNote(rows)
		if err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return notes, count, nil
}

//GetLatest retrieves the latest notes based off their created date and id.
func (n *noteStore) GetLatest(offset, limit int) ([]*store.Note, int, error) {
	var count int
	err := n.db.QueryRow(`select count(*) from notes where deleted=false`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	if limit <= 0 || offset > count {
		return []*store.Note{}, count, nil
	}

	rows, err := n.db.Query(`select id,name,author,content,created_at from notes where deleted=false order by created_at desc, id desc limit $1 offset $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := []*store.Note{}
	for rows.Next() {
		note, err := n.scanNote(rows)
		if err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return notes, count, nil
}

//SetContent updates content of the note whose id matches the given id
func (n *noteStore) SetContent(id int64, content string) error {
	_, err := n.db.Exec(`update notes set content=$1 where id=$2`, content, id)
	return err
}

//SetName updates name of the note whose id matches the given id
func (n *noteStore) SetName(id int64, name string) error {
	_, err := n.db.Exec(`update notes set name=$1 where id=$2`, name, id)
	return err
}

//SetAuthor updates the author of the note whose id matches the given id
func (n *noteStore) SetAuthor(id int64, author string) error {
	_, err := n.db.Exec(`update notes set author=$1 where id=$2`, author, id)
	return err
}

//Delete soft deletes a note that matches the given id
func (n *noteStore) Delete(id int64) error {
	tx, err := n.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`update notes set deleted=true where id=$1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (n *noteStore) scanNote(scanner scanner) (*store.Note, error) {
	note := new(store.Note)
	err := scanner.Scan(&note.ID, &note.Name, &note.Author, &note.Content, &note.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return note, nil
}
