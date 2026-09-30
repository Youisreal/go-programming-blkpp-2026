package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadWriteStudentsJSON(t *testing.T) {
	dir := t.TempDir()
	oldFile := jsonFile
	jsonFile = filepath.Join(dir, "students.json")
	t.Cleanup(func() { jsonFile = oldFile })

	want := []Student{
		{
			ID:    1,
			Name:  "Dewi",
			Email: "dewi@example.com",
			Profile: Profile{
				City:  "Balikpapan",
				Phone: "0812-7777-8888",
			},
			Courses: []Course{
				{Code: "GO201", Title: "File I/O", Grade: "A"},
			},
		},
	}

	writeStudents(want)
	got := readStudents()

	if len(got) != 1 || got[0].Profile.City != "Balikpapan" || got[0].Courses[0].Code != "GO201" {
		t.Fatalf("got %#v", got)
	}

	if _, err := os.Stat(jsonFile); err != nil {
		t.Fatal(err)
	}
}
