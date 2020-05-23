package mysql

import "testing"

func TestNote(t *testing.T) {
	s, teardown := getTestStore(t)
	defer teardown()

	note1, err := s.Notes().New("How the stack grows", "Mr.Parish", "Most people don't know where the top of the stack is. The stack grows towards lower addresses. I will explain why and give concrete examples.")
	if err != nil {
		t.Fatalf("failed to create a new note: %s", err)
	}

	note2, err := s.Notes().New("Pointer vs Value smantics", "Mr.Parish", "How well you manage your memory depends on how good is your pointer/value semantics. How well do you understand this.Do you know when to use either. Let me show you how.")
	if err != nil {
		t.Fatalf("failed to create a new note: %s", err)
	}

	notes, c, err := s.Notes().GetLatest(0, 10)
	if err != nil {
		t.Fatalf("failed to get latest notes: %s", err)
	}

	if len(notes) != 2 {
		t.Fatalf("failed to get all notes: %d", len(notes))
	}

	n1, err := s.Notes().Get(note1)
	if err != nil {
		t.Fatalf("failed to retrieve a note: %s; id: %d", err, note1)
	}

	if c != 2 {
		t.Fatalf("bad note count: %d", c)
	}

	if n1.Content != "How the stack grows" {
		t.Fatalf("bad note1 content: %s", n1.Content)
	}

	newContent := "Value semantics means that you deal directly with values and that you pass copies around. The point here is that when you have a value, you can trust it won't change behind your back. With pointer semantics, you don't have a value, you have an 'address'. Someone else could alter what is there, you can't know."

	err = s.Notes().SetContent(note2, newContent)
	if err != nil {
		t.Fatalf("failed to update note content %s; id: %d", newContent, note2)
	}

	n2, err := s.Notes().Get(note2)
	if err != nil {
		t.Fatalf("failed to retrieve a note: %s; id: %d", err, note2)
	}

	if n2.Content != newContent {
		t.Fatalf("bad note content %s; Expected: %s", n2.Content, newContent)
	}

	notes, count, err := s.Notes().GetByAuthor(n1.Author, 0, 10)
	if err != nil {
		t.Fatalf("failed to retrieve notes by author: %s;", err)
	}

	if len(notes) != 2 {
		t.Fatalf("bad notes len: %d; Expected 2", len(notes))
	}

	if count != 2 {
		t.Fatalf("bad note count: %d; Expected 2", count)
	}
}
