package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	goredis "github.com/redis/go-redis/v9"

	"github.com/mauricio-reportei/taskforge-api-go/internal/middleware"
	"github.com/mauricio-reportei/taskforge-api-go/internal/shared"
)

// Server bundles HTTP dependencies for TaskForge API.
type Server struct {
	addr    string
	handler http.Handler
	pool    *pgxpool.Pool
	rdb     *goredis.Client
}

// New configures routes and middleware.
func New(addr string, pool *pgxpool.Pool, rdb *goredis.Client) *Server {
	s := &Server{
		addr: addr,
		pool: pool,
		rdb:  rdb,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		s.health(w, r)
	})

	stack := middleware.Recoverer(middleware.RequestLogger(middleware.RateLimiter(mux)))
	s.handler = stack
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// Addr returns the listen address (host:port).
func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		slog.Error("health postgres ping failed", "error", err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	if err := s.rdb.Ping(ctx).Err(); err != nil {
		slog.Error("health redis ping failed", "error", err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	if err := shared.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		slog.Error("write health", "error", err)
	}
}
