package model

import "testing"

func TestBuildOwnerFilter_ProjectScope(t *testing.T) {
	filter, err := buildOwnerFilter("user-1", "project-1")
	if err != nil {
		t.Fatalf("buildOwnerFilter returned error: %v", err)
	}

	expected := "owner_type == 'PROJECT' and owner_id == 'project-1'"
	if filter != expected {
		t.Fatalf("unexpected filter, got=%q want=%q", filter, expected)
	}
}

func TestBuildOwnerFilter_UserScope(t *testing.T) {
	filter, err := buildOwnerFilter("user-1", "")
	if err != nil {
		t.Fatalf("buildOwnerFilter returned error: %v", err)
	}

	expected := "owner_type == 'PERSONAL' and owner_id == 'user-1'"
	if filter != expected {
		t.Fatalf("unexpected filter, got=%q want=%q", filter, expected)
	}
}

func TestBuildOwnerFilter_EmptyScope(t *testing.T) {
	_, err := buildOwnerFilter(" ", "")
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if err != ErrSearchScopeRequired {
		t.Fatalf("unexpected error, got=%v want=%v", err, ErrSearchScopeRequired)
	}
}

func TestBuildOwnerFilter_EscapeSpecialChars(t *testing.T) {
	filter, err := buildOwnerFilter("", "proj\\a'b")
	if err != nil {
		t.Fatalf("buildOwnerFilter returned error: %v", err)
	}

	expected := "owner_type == 'PROJECT' and owner_id == 'proj\\\\a\\'b'"
	if filter != expected {
		t.Fatalf("unexpected filter, got=%q want=%q", filter, expected)
	}
}

