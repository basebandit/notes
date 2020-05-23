package store

import "time"

//Note is a single note entity
type Note struct {
	ID        int64     `json:"id"`
	Author    string    `json:author"`
	Name      string    `json:name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
