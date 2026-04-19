package utils

import (
	"reflect"
	"testing"
)

func TestInsertTextAtLine(t *testing.T) {
	got, err := InsertTextAtLine([]string{"a", "d"}, 2, "b\nc")
	if err != nil {
		t.Fatalf("InsertTextAtLine() error = %v", err)
	}
	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("InsertTextAtLine() = %#v, want %#v", got, want)
	}
}

func TestDeleteLineRange(t *testing.T) {
	got, err := DeleteLineRange([]string{"a", "b", "c", "d"}, 2, 3)
	if err != nil {
		t.Fatalf("DeleteLineRange() error = %v", err)
	}
	want := []string{"a", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DeleteLineRange() = %#v, want %#v", got, want)
	}
}

func TestNumberedLines(t *testing.T) {
	got := NumberedLines([]string{"{", "  \"a\": 1", "}"})
	want := "1| {\n2|   \"a\": 1\n3| }"
	if got != want {
		t.Fatalf("NumberedLines() = %q, want %q", got, want)
	}
}
