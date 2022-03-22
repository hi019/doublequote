package http

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	dq "doublequote"
	"doublequote/pkg/domain"
	"doublequote/pkg/utils"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/fx"
	"golang.org/x/crypto/acme/autocert"
)

const ShutdownTimeout = 5

type Server struct {
	ln     net.Listener
	server *http.Server
	router *chi.Mux
	now    func() time.Time

	userService       domain.UserService
	cryptoService     domain.CryptoService
	sessionService    domain.SessionService
	collectionService domain.CollectionService
	feedService       domain.FeedService
	entryService      domain.EntryService
	storageService    domain.StorageService

	config domain.Config
}

func NewServer(
	lc fx.Lifecycle,
	userService domain.UserService,
	cryptoService domain.CryptoService,
	sessionService domain.SessionService,
	collectionService domain.CollectionService,
	feedService domain.FeedService,
	entryService domain.EntryService,
	storageService domain.StorageService,
	config domain.Config,
) *Server {
	s := &Server{
		userService:       userService,
		cryptoService:     cryptoService,
		sessionService:    sessionService,
		collectionService: collectionService,
		feedService:       feedService,
		entryService:      entryService,
		storageService:    storageService,
		config:            config,

		server: &http.Server{},
		router: chi.NewRouter(),
		now:    time.Now,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.open()
		},
		OnStop: func(ctx context.Context) error {
			return s.close()
		},
	})

	s.setupRouter()

	return s
}

func (s *Server) open() (err error) {
	if s.config.HTTP.Domain != "" {
		s.ln = autocert.NewListener(s.config.HTTP.Domain)
	} else {
		if s.ln, err = net.Listen("tcp", ":"+s.config.HTTP.Port); err != nil {
			return err
		}
	}

	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) close() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) setupRouter() {
	s.router.Use(chimw.Logger)

	// TODO change this?
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	s.router.Route("/api/pub", func(r chi.Router) {
		r.Use(setContentType)

		s.registerPublicUserRoutes(r)
	})

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.requireAuth)
		r.Use(setContentType)

		s.registerUserRoutes(r)
		s.registerCollectionRoutes(r)
		s.registerFeedRoutes(r)
	})

	s.router.Route("/", func(r chi.Router) {
		staticFs, err := fs.Sub(dq.Assets, "assets/frontend")
		if err != nil {
			log.Fatalf("http: %s", err)
		}

		fserver := http.FileServer(http.FS(staticFs))

		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			// TODO currently we open index.html for all frontend requests
			indexFile, err := staticFs.Open("index.html")
			if err != nil {
				Error(w, req, err)
				return
			}

			_, err = staticFs.Open(strings.TrimPrefix(req.URL.Path, "/"))
			if errors.Is(err, fs.ErrNotExist) {
				io.Copy(w, indexFile)
				return
			}

			fserver.ServeHTTP(w, req)
		})
	})

	s.server.Handler = s.router
}

// requireAuth is a middleware for requiring authentication
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session, make sure it's not empty.
		sess, err := s.sessionService.Get(r)
		if err != nil {
			Error(w, r, err)
			return
		}
		if sess == nil {
			Error(w, r, domain.Errorf(domain.EUNAUTHORIZED, domain.ErrUnauthorized))
			return
		}

		// Find user associated with session, then add it to the request context.
		u, err := s.userService.FindUser(r.Context(), domain.UserFilter{ID: utils.IntPtr(sess.UserID())}, domain.UserInclude{})
		if err != nil {
			Error(w, r, err)
			return
		}

		r = r.WithContext(domain.NewContextWithUser(r.Context(), u))
		next.ServeHTTP(w, r)
	})
}

// setContentType is a middleware that sets the response Content-type to application/json.
func setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}
