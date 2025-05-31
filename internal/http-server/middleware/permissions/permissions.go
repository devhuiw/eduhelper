package permissions

import (
	"log/slog"
	"net/http"
	"service/internal/domain/repository"
	"service/internal/http-server/middleware"
	"service/internal/lib/api/response"
	"strings"

	"github.com/go-chi/render"
)

type RBACMiddleware struct {
	userRoleRepo   *repository.UserRoleRepository
	rolePermRepo   *repository.RolePermissionRepository
	permissionRepo *repository.PermissionRepository
	logger         *slog.Logger
}

func NewRBACMiddleware(
	userRoleRepo *repository.UserRoleRepository,
	rolePermRepo *repository.RolePermissionRepository,
	permissionRepo *repository.PermissionRepository,
	logger *slog.Logger,
) *RBACMiddleware {
	return &RBACMiddleware{
		userRoleRepo:   userRoleRepo,
		rolePermRepo:   rolePermRepo,
		permissionRepo: permissionRepo,
		logger:         logger,
	}
}

func (m *RBACMiddleware) RequirePermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := middleware.GetUserClaims(r)
			idClaim, ok := claims["id"]
			if !ok {
				m.logger.Info("user id not found in claims")
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, response.Error("unauthorized"))
				return
			}

			var userID int64
			switch v := idClaim.(type) {
			case int64:
				userID = v
			case float64:
				userID = int64(v)
			}

			roles, err := m.userRoleRepo.GetRolesByUserID(r.Context(), userID)
			if err != nil {
				m.logger.Error("failed to get user roles", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal error"))
				return
			}
			permsSet := make(map[string]struct{})
			for _, role := range roles {
				perms, err := m.rolePermRepo.GetPermissionsByRoleID(r.Context(), role.RoleID)
				if err != nil {
					m.logger.Error("failed to get role permissions", slog.String("err", err.Error()))
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, response.Error("internal error"))
					return
				}
				for _, perm := range perms {
					permsSet[strings.ToLower(perm.PermissionName)] = struct{}{}
				}
			}
			if _, ok := permsSet[strings.ToLower(permissionName)]; !ok {
				m.logger.Info("permission denied", slog.String("permission", permissionName))
				w.WriteHeader(http.StatusForbidden)
				render.JSON(w, r, response.Error("permission denied"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
