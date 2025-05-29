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
	"service/internal/http-server/middleware/permissions"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(
	log *slog.Logger,
	cfg *config.Config,
	db *sql.DB,
) (*http.Server, error) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	rbacMiddleware := permissions.NewRBACMiddleware(
		repository.NewUserRoleRepository(db),
		repository.NewRolePermissionRepository(db),
		repository.NewPermissionRepository(db),
		log,
	)

	userRepository := repository.NewUserRepository(db)
	userHandler := v1.NewUserHandler(userRepository)

	authHandler := v1.NewAuthHandler(userRepository, cfg.JwtSecret)

	teacherRepository := repository.NewTeacherRepository(db)
	teacherHandler := v1.NewTeacherHandler(teacherRepository)

	permissionRepository := repository.NewPermissionRepository(db)
	permissionHandler := v1.NewPermissionHandler(permissionRepository)

	roleRepository := repository.NewRoleRepository(db)
	roleHandler := v1.NewRoleHandler(roleRepository)

	userRoleRepository := repository.NewUserRoleRepository(db)
	userRoleHandler := v1.NewUserRoleHandler(userRoleRepository)

	rolePermissionRepository := repository.NewRolePermissionRepository(db)
	rolePermissionHandler := v1.NewRolePermissionHandler(rolePermissionRepository)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register(log))
		r.Post("/login", authHandler.Login(log))
	})

	router.Group(func(r chi.Router) {
		r.Use(middle.JWTAuth())
		r.Use(middle.AuthRequired())

		r.Route("/api/v1/users", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("user:list")).Get("/", userHandler.ListUsers(log))
			rr.With(rbacMiddleware.RequirePermission("user:view")).Get("/{id}", userHandler.GetUserByID(log))
			rr.With(rbacMiddleware.RequirePermission("user:update")).Put("/{id}", userHandler.UpdateUser(log))
			rr.With(rbacMiddleware.RequirePermission("user:delete")).Delete("/{id}", userHandler.DeleteUser(log))
		})

		r.Route("/api/v1/teacher", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("teacher:view_self")).Get("/me", teacherHandler.GetMyTeacherProfile(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:view_self")).Get("/{id}", teacherHandler.GetTeacherPublicByID(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:update_self")).Put("/me", teacherHandler.UpdateMyTeacherProfile(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:create")).Post("/", teacherHandler.CreateTeacher(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:view")).Get("/{id}", teacherHandler.GetTeacherByID(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:update")).Put("/{id}", teacherHandler.UpdateTeacher(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:delete")).Delete("/{id}", teacherHandler.DeleteTeacher(log))
		})

		r.Route("/api/v1/permissions", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("permission:list")).Get("/", permissionHandler.ListPermissions(log))
			rr.With(rbacMiddleware.RequirePermission("permission:create")).Post("/", permissionHandler.CreatePermission(log))
			rr.With(rbacMiddleware.RequirePermission("permission:view")).Get("/{id}", permissionHandler.GetPermissionByID(log))
			rr.With(rbacMiddleware.RequirePermission("permission:update")).Put("/{id}", permissionHandler.UpdatePermission(log))
			rr.With(rbacMiddleware.RequirePermission("permission:delete")).Delete("/{id}", permissionHandler.DeletePermission(log))
		})

		r.Route("/api/v1/roles", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("role:list")).Get("/", roleHandler.ListRoles(log))
			rr.With(rbacMiddleware.RequirePermission("role:create")).Post("/", roleHandler.CreateRole(log))
			rr.With(rbacMiddleware.RequirePermission("role:view")).Get("/{id}", roleHandler.GetRoleByID(log))
			rr.With(rbacMiddleware.RequirePermission("role:update")).Put("/{id}", roleHandler.UpdateRole(log))
			rr.With(rbacMiddleware.RequirePermission("role:delete")).Delete("/{id}", roleHandler.DeleteRole(log))
		})

		r.Route("/api/v1/user-roles", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("userrole:assign")).Post("/assign", userRoleHandler.AssignRole(log))
			rr.With(rbacMiddleware.RequirePermission("userrole:remove")).Post("/remove", userRoleHandler.RemoveRole(log))
			// TODO: GET для просмотра ролей пользователя
		})

		r.Route("/api/v1/role-permissions", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("rolepermission:assign")).Post("/assign", rolePermissionHandler.AssignPermission(log))
			rr.With(rbacMiddleware.RequirePermission("rolepermission:remove")).Post("/remove", rolePermissionHandler.RemovePermission(log))
			// TODO: GET для просмотра прав роли
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
