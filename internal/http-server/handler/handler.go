package handler

import (
	"database/sql"
	"log/slog"
	"net/http"
	"service/internal/config"
	"service/internal/domain/repository"
	v1 "service/internal/http-server/handler/v1"
	middle "service/internal/http-server/middleware"
	"service/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(log *slog.Logger,
	cfg *config.Config, db *sql.DB) (*http.Server, error) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	userRepository := repository.NewUserRepository(db)
	userHandler := v1.NewUserHandler(userRepository)

	authHandler := v1.NewAuthHandler(userRepository, cfg.JwtSecret)

	router.Route("/api/v1/", func(r chi.Router) {
		r.Post("/register", authHandler.Register(log))
		r.Post("/login", authHandler.Login(log))
	})

	router.Group(func(r chi.Router) {
		r.Use(middle.JWTAuth())      // JWT в контекст
		r.Use(middle.AuthRequired()) // Проверка авторизации
		r.Route("/api/v1/users", func(rr chi.Router) {
			// r.Post("/", userHandler.CreateUser(log))
			r.Get("/", userHandler.ListUsers(log))
			r.Get("/{id}", userHandler.GetUserByID(log))
			r.Put("/{id}", userHandler.UpdateUser(log))
			r.Delete("/{id}", userHandler.DeleteUser(log))
		})
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return srv, nil
}
