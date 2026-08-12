package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/auth"
	"github.com/Mauricio-AFLadeira/taskforge-api-go/internal/shared"
)

func TestRequireAuthSetsContext(t *testing.T) {
	t.Parallel()
	const secret = "middleware-test-secret-here!!!"
	issuer := auth.NewTokenIssuer(secret, 5*time.Minute)
	tok, err := issuer.IssueAccess("22222222-2222-2222-2222-222222222222", "me@example.com")
	if err != nil {
		t.Fatal(err)
	}

	var uid, mail string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		uid, ok = shared.UserIDFromContext(r.Context())
		if !ok {
			t.Error("missing user id in context")
		}
		mail, ok = shared.EmailFromContext(r.Context())
		if !ok {
			t.Error("missing email in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	srv := RequireAuth(secret, inner)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if uid != "22222222-2222-2222-2222-222222222222" || mail != "me@example.com" {
		t.Fatalf("context values %q %q", uid, mail)
	}
}

func TestRequireAuthMissingBearer(t *testing.T) {
	t.Parallel()
	inner := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not run")
	})
	srv := RequireAuth("secret", inner)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("got %d", rr.Code)
	}
}
