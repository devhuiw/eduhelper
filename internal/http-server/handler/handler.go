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

	_ "service/internal/docs"

	httpSwagger "github.com/swaggo/http-swagger"
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

	auditLogRepository := repository.NewAuditLogRepository(db)

	userRepository := repository.NewUserRepository(db)
	userHandler := v1.NewUserHandler(userRepository, auditLogRepository)

	authHandler := v1.NewAuthHandler(userRepository, cfg.JwtSecret)

	teacherRepository := repository.NewTeacherRepository(db)
	teacherHandler := v1.NewTeacherHandler(teacherRepository, auditLogRepository)

	permissionRepository := repository.NewPermissionRepository(db)
	permissionHandler := v1.NewPermissionHandler(permissionRepository, auditLogRepository)

	roleRepository := repository.NewRoleRepository(db)
	roleHandler := v1.NewRoleHandler(roleRepository, auditLogRepository)

	userRoleRepository := repository.NewUserRoleRepository(db)
	userRoleHandler := v1.NewUserRoleHandler(userRoleRepository, auditLogRepository)

	rolePermissionRepository := repository.NewRolePermissionRepository(db)
	rolePermissionHandler := v1.NewRolePermissionHandler(rolePermissionRepository)

	studentRepository := repository.NewStudentRepository(db)
	studentHandler := v1.NewStudentHandler(studentRepository, auditLogRepository)

	studentGroupRepository := repository.NewStudentGroupRepository(db)
	studentGroupHandler := v1.NewStudentGroupHandler(studentGroupRepository, auditLogRepository)

	curriculumRepository := repository.NewCurriculumRepository(db)
	curriculumHandler := v1.NewCurriculumHandler(curriculumRepository, auditLogRepository)

	gradeJournalRepository := repository.NewGradeJournalRepository(db)
	gradeJournalHandler := v1.NewGradeJournalHandler(gradeJournalRepository, auditLogRepository)

	attendanceRepository := repository.NewAttendanceRepository(db)
	attendanceHandler := v1.NewAttendanceHandler(attendanceRepository, auditLogRepository)

	semesterRepository := repository.NewSemesterRepository(db)
	semesterHandler := v1.NewSemesterHandler(semesterRepository, auditLogRepository)

	disciplineRepository := repository.NewDisciplineRepository(db)
	disciplineHandler := v1.NewDisciplineHandler(disciplineRepository, auditLogRepository)

	academicYearRepository := repository.NewAcademicYearRepository(db)
	academicYearHandler := v1.NewAcademicYearHandler(academicYearRepository, auditLogRepository)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authHandler.Register(log))
		r.Post("/login", authHandler.Login(log))
	})

	router.Group(func(r chi.Router) {
		r.Use(middle.JWTAuth(cfg.JwtSecret))
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
			rr.With(rbacMiddleware.RequirePermission("teacher:list")).Get("/", teacherHandler.ListTeacher(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:view")).Get("/{id}", teacherHandler.GetTeacherByID(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:update")).Put("/{id}", teacherHandler.UpdateTeacher(log))
			rr.With(rbacMiddleware.RequirePermission("teacher:delete")).Delete("/{id}", teacherHandler.DeleteTeacher(log))
		})

		r.Route("/api/v1/students", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("student:create")).Post("/", studentHandler.CreateStudent(log))
			rr.With(rbacMiddleware.RequirePermission("student:view")).Get("/{id}", studentHandler.GetStudentByID(log))
			rr.With(rbacMiddleware.RequirePermission("student:update")).Put("/{id}", studentHandler.UpdateStudent(log))
			rr.With(rbacMiddleware.RequirePermission("student:delete")).Delete("/{id}", studentHandler.DeleteStudent(log))
			rr.With(rbacMiddleware.RequirePermission("student:list")).Get("/", studentHandler.ListStudent(log))
			rr.With(rbacMiddleware.RequirePermission("student:view_public")).Get("/public/{id}", studentHandler.GetStudentPublicByID(log))
			rr.With(rbacMiddleware.RequirePermission("student:list_public")).Get("/public", studentHandler.ListStudentPublic(log))
		})

		r.Route("/api/v1/student-groups", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("studentgroup:create")).Post("/", studentGroupHandler.CreateStudentGroup(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:view")).Get("/{id}", studentGroupHandler.GetStudentGroupByID(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:update")).Put("/{id}", studentGroupHandler.UpdateStudentGroup(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:delete")).Delete("/{id}", studentGroupHandler.DeleteStudentGroup(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:list")).Get("/", studentGroupHandler.ListStudentGroups(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:view_public")).Get("/public/{id}", studentGroupHandler.GetStudentGroupPublicByID(log))
			rr.With(rbacMiddleware.RequirePermission("studentgroup:list_public")).Get("/public", studentGroupHandler.ListStudentGroupPublic(log))
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
			rr.With(rbacMiddleware.RequirePermission("userrole:view")).Get("/{id}", userRoleHandler.GetRolesByUserID(log))
		})

		r.Route("/api/v1/role-permissions", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("rolepermission:assign")).Post("/assign", rolePermissionHandler.AssignPermission(log))
			rr.With(rbacMiddleware.RequirePermission("rolepermission:remove")).Post("/remove", rolePermissionHandler.RemovePermission(log))
			rr.With(rbacMiddleware.RequirePermission("rolepermission:view")).Get("/{id}", rolePermissionHandler.GetPermissionsByRoleID(log))
		})

		r.Route("/api/v1/curriculums", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("curriculum:create")).Post("/", curriculumHandler.CreateCurriculum(log))
			rr.With(rbacMiddleware.RequirePermission("curriculum:view")).Get("/{id}", curriculumHandler.GetCurriculumByID(log))
			rr.With(rbacMiddleware.RequirePermission("curriculum:update")).Put("/{id}", curriculumHandler.UpdateCurriculum(log))
			rr.With(rbacMiddleware.RequirePermission("curriculum:delete")).Delete("/{id}", curriculumHandler.DeleteCurriculum(log))
			rr.With(rbacMiddleware.RequirePermission("curriculum:list")).Get("/", curriculumHandler.ListCurriculum(log))
		})

		r.Route("/api/v1/gradejournals", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("gradejournal:create")).Post("/", gradeJournalHandler.CreateGradeJournal(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:view")).Get("/{id}", gradeJournalHandler.GetGradeJournalByID(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:update")).Put("/{id}", gradeJournalHandler.UpdateGradeJournal(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:delete")).Delete("/{id}", gradeJournalHandler.DeleteGradeJournal(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:list")).Get("/", gradeJournalHandler.ListGradeJournal(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:list_public")).Get("/public", gradeJournalHandler.ListGradeJournalPublic(log))
			rr.With(rbacMiddleware.RequirePermission("gradejournal:avg")).Get("/average", gradeJournalHandler.GetAverageGrade(log))
		})

		r.Route("/api/v1/attendances", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("attendance:create")).Post("/", attendanceHandler.CreateAttendance(log))
			rr.With(rbacMiddleware.RequirePermission("attendance:view")).Get("/{id}", attendanceHandler.GetAttendanceByID(log))
			rr.With(rbacMiddleware.RequirePermission("attendance:update")).Put("/{id}", attendanceHandler.UpdateAttendance(log))
			rr.With(rbacMiddleware.RequirePermission("attendance:delete")).Delete("/{id}", attendanceHandler.DeleteAttendance(log))
			rr.With(rbacMiddleware.RequirePermission("attendance:list")).Get("/", attendanceHandler.ListAttendance(log))
		})

		r.Route("/api/v1/semesters", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("semester:create")).Post("/", semesterHandler.CreateSemester(log))
			rr.With(rbacMiddleware.RequirePermission("semester:view")).Get("/{id}", semesterHandler.GetSemesterByID(log))
			rr.With(rbacMiddleware.RequirePermission("semester:update")).Put("/{id}", semesterHandler.UpdateSemester(log))
			rr.With(rbacMiddleware.RequirePermission("semester:delete")).Delete("/{id}", semesterHandler.DeleteSemester(log))
			rr.With(rbacMiddleware.RequirePermission("semester:list")).Get("/", semesterHandler.ListSemester(log))
		})

		r.Route("/api/v1/disciplines", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("discipline:create")).Post("/", disciplineHandler.CreateDiscipline(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:view")).Get("/{id}", disciplineHandler.GetDisciplineByID(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:update")).Put("/{id}", disciplineHandler.UpdateDiscipline(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:delete")).Delete("/{id}", disciplineHandler.DeleteDiscipline(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:list")).Get("/", disciplineHandler.ListDiscipline(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:list_public")).Get("/public", disciplineHandler.ListDisciplinePublic(log))
			rr.With(rbacMiddleware.RequirePermission("discipline:view_public")).Get("/public/{id}", disciplineHandler.GetDisciplinePublicByID(log))
		})

		r.Route("/api/v1/academic-years", func(rr chi.Router) {
			rr.With(rbacMiddleware.RequirePermission("academicyear:create")).Post("/", academicYearHandler.CreateAcademicYear(log))
			rr.With(rbacMiddleware.RequirePermission("academicyear:view")).Get("/{id}", academicYearHandler.GetAcademicYearByID(log))
			rr.With(rbacMiddleware.RequirePermission("academicyear:update")).Put("/{id}", academicYearHandler.UpdateAcademicYear(log))
			rr.With(rbacMiddleware.RequirePermission("academicyear:delete")).Delete("/{id}", academicYearHandler.DeleteAcademicYear(log))
			rr.With(rbacMiddleware.RequirePermission("academicyear:list")).Get("/", academicYearHandler.ListAcademicYear(log))
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
