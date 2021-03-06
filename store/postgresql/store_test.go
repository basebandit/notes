package postgresql

import "testing"

const (
	testAddress  = "127.0.0.1:5432"
	testUsername = "note_test_user"
	testPassword = "note_test_password"
	testDatabase = "note_test_database"
)

func getTestStore(t *testing.T) (*Store, func()) {
	conn, err := Connect(testAddress, testUsername, testPassword, testDatabase, "disable", "")
	if err != nil {
		t.Fatalf("failed to connect to the test postgresl database: address=%q, username=%q, password=%q, database=%q: %s", testAddress, testUsername, testPassword, testDatabase, err)
	}

	err = conn.Reset()
	if err != nil {
		t.Fatalf("failed to reset the test postgresql database: %s", err)
	}

	teardown := func() {
		err := conn.Drop()
		if err != nil {
			t.Fatalf("failed to drop tables of the test postgresql database: %s", err)
		}
	}

	return conn, teardown
}

func TestPlaceholders(t *testing.T) {
	testTable := []struct {
		start, count int
		want         string
	}{
		{1, 0, ""},
		{1, 1, "$1"},
		{1, 5, "$1,$2,$3,$4,$5"},
		{3, 0, ""},
		{3, 1, "$3"},
		{3, 5, "$3,$4,$5,$6,$7"},
		{0, 0, ""},
		{-1, -1, ""},
	}
	for _, test := range testTable {
		got := placeholders(test.start, test.count)
		if got != test.want {
			t.Fatalf("got placeholders %q for (%d, %d), want %q", got, test.start, test.count, test.want)
		}
	}
}
