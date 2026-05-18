package auth

import (
	"testing"
	"time"
)

func TestIssueAndParseAccess(t *testing.T) {
	t.Parallel()
	issuer := NewTokenIssuer("unit-test-secret-key-ok", time.Minute)
	token, err := issuer.IssueAccess("11111111-1111-1111-1111-111111111111", "user@example.com")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	id, email, err := ParseAccess([]byte("unit-test-secret-key-ok"), token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if id != "11111111-1111-1111-1111-111111111111" || email != "user@example.com" {
		t.Fatalf("claims: got %q %q", id, email)
	}
}

func TestParseAccessWrongSecret(t *testing.T) {
	t.Parallel()
	issuer := NewTokenIssuer("secret-a", time.Minute)
	token, err := issuer.IssueAccess("id", "a@b.com")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = ParseAccess([]byte("secret-b"), token)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHashRefreshTokenStable(t *testing.T) {
	t.Parallel()
	const raw = "opaque-test"
	a := HashRefreshToken(raw)
	b := HashRefreshToken(raw)
	if a != b || len(a) != 64 {
		t.Fatalf("hash=%q len=%d", a, len(a))
	}
}
