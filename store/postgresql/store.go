package postgresql

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lib/pq"
	"github.com/parish/notes/store"
)

//Store is a mysql implementation of store (see store.go in package store)
type Store struct {
	db        *sql.DB
	noteStore *noteStore
	//Other stores go here as well
}

//Notes returns a note store.
func (s *Store) Notes() store.NoteStore {
	return s.noteStore
}

//Compile time check that we satisfy the Store interface
var _ store.Store = (*Store)(nil)

type scanner interface {
	Scan(v ...interface{}) error
}

//Connect connects to the store's mysql database
func Connect(address, username, password, database, sslmode, sslrootcert string) (*Store, error) {
	connstr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s&connect_timeout=10",
		username, password, address, database, sslmode,
	)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{
		db:        db,
		noteStore: &noteStore{db: db},
		//other store initializations go here
	}

	err = s.Migrate()
	if err != nil {
		return nil, err
	}
	return s, nil
}

//Migrate migrates the store's mysql database.
func (s *Store) Migrate() error {
	for _, q := range migrate {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

//Drop drops the store's mysql database.
func (s *Store) Drop() error {
	for _, q := range drop {
		_, err := s.db.Exec(q)
		if err != nil {
			return fmt.Errorf("sql exec error: %s; query: %q", err, q)
		}
	}
	return nil
}

//Reset reverts the store's mysql database to the initial state.
func (s *Store) Reset() error {
	err := s.Drop()
	if err != nil {
		return err
	}
	return s.Migrate()
}

func placeholders(start, count int) string {
	buf := new(bytes.Buffer)
	for i := start; i < start+count; i++ {
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(i))
		if i < start+count-1 {
			buf.WriteByte(',')
		}
	}
	return buf.String()
}

func isUniqueConstraintError(err error) bool {
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}
