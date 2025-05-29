package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"service/internal/domain/models"
	"service/internal/lib/api/response"
	"service/internal/lib/jwt"
	"time"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  UserRepository
	jwtSecret string
}

func NewAuthHandler(userRepo UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(log *slog.Logger) http.HandlerFunc {
	const op = "auth.Login"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op))
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("invalid login request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}
		user, err := h.userRepo.GetClientByEmail(r.Context(), req.Email)
		if err != nil || user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("invalid credentials"))
			return
		}
		// bcrypt сравнение
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("invalid credentials"))
			return
		}

		//создание токена
		token, err := jwt.NewToken(*user, time.Duration(time.Duration.Hours(24)), h.jwtSecret)
		if err != nil {
			log.Error("failed to sign jwt", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}
		render.JSON(w, r, map[string]string{"token": token})
	}
}

func (h *AuthHandler) Register(log *slog.Logger) http.HandlerFunc {
	const op = "auth.Register"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op))
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info("invalid register request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		if req.Email == "" || req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("email and password required"))
			return
		}

		existingUser, err := h.userRepo.GetClientByEmail(r.Context(), req.Email)
		fmt.Printf("DEBUG GetByEmail: user=%+v, err=%v\n", existingUser, err)
		if existingUser != nil {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Error("email already exists"))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		user := &models.User{
			Email:      req.Email,
			Password:   hashedPassword,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
		}
		if err := h.userRepo.CreateClient(r.Context(), user); err != nil {
			log.Error("failed to create user", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		token, err := jwt.NewToken(*user, time.Duration(time.Duration.Hours(24)), h.jwtSecret)
		if err != nil {
			log.Error("failed to sign jwt", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}
		render.JSON(w, r, map[string]string{"token": token})
	}
}
