package postgresql

import "testing"

func TestNote(t *testing.T) {
	s, teardown := getTestStore(t)
	defer teardown()

	id1, err := s.Notes().New("How the stack grows", "Mr.Parish", "Most people don't know where the top of the stack is. The stack grows towards lower addresses. I will explain why and give concrete examples.")
	if err != nil {
		t.Fatalf("failed to create a new note: %s", err)
	}

	id2, err := s.Notes().New("Pointer vs Value smantics", "Mr.Parish", "How well you manage your memory depends on how good is your pointer/value semantics. How well do you understand this.Do you know when to use either. Let me show you how.")
	if err != nil {
		t.Fatalf("failed to create a new note: %s", err)
	}

	notes, c, err := s.Notes().GetLatest(0, 2)
	if err != nil {
		t.Fatalf("failed to get all notes: %s", err)
	}

	if len(notes) != 2 {
		t.Fatalf("bad note len: %d", len(notes))
	}

	if c != 2 {
		t.Fatalf("bad note count: %d", c)
	}

	n1, err := s.Notes().Get(id1)
	if err != nil {
		t.Fatalf("failed to retrieve a note: %s; id: %d", err, id1)
	}

	content1 := "Most people don't know where the top of the stack is. The stack grows towards lower addresses. I will explain why and give concrete examples."
	if n1.Content != content1 {
		t.Fatalf("bad note1 content: %s", n1.Content)
	}

	content2 := "Value semantics means that you deal directly with values and that you pass copies around. The point here is that when you have a value, you can trust it won't change behind your back. With pointer semantics, you don't have a value, you have an 'address'. Someone else could alter what is there, you can't know."

	err = s.Notes().SetContent(id2, content2)
	if err != nil {
		t.Fatalf("failed to update note content %s; id: %d", content2, id2)
	}

	n2, err := s.Notes().Get(id2)
	if err != nil {
		t.Fatalf("failed to retrieve a note: %s; id: %d", err, id2)
	}

	if n2.Content != content2 {
		t.Fatalf("bad note content %s; Expected: %s", n2.Content, content2)
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
