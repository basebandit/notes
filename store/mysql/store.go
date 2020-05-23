package mysql

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
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
func Connect(address, username, password, database string) (*Store, error) {
	connstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, address, database)

	db, err := sql.Open("mysql", connstr)
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

func placeholders(count int) string {
	buf := new(bytes.Buffer)
	for i := 0; i < count; i++ {
		buf.WriteByte('?')
		if i < count-1 {
			buf.WriteByte(',')
		}
	}
	return buf.String()
}

func isUniqueConstraintError(err error) bool {
	if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1062 {
		return true
	}
	return false
}
