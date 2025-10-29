package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/services"
	"iris-api/internal/infrastructure/database"
	"iris-api/internal/infrastructure/repositories"
	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/middleware"
	"iris-api/internal/presentation/routes"
	"iris-api/pkg/auth"
	"iris-api/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db.GetDB())
	userService := services.NewUserService(userRepo)
	userUseCase := usecases.NewUserUseCase(userService)
	userHandler := handlers.NewUserHandler(userUseCase)

	recetaRepo := repositories.NewRecetaRepository(db.GetDB())
	recetaService := services.NewRecetaService(recetaRepo)
	recetaUseCase := usecases.NewRecetaUseCase(recetaService)
	recetaHandler := handlers.NewRecetaHandler(recetaUseCase)

	turnoRepo := repositories.NewTurnoRepository(db.GetDB())
	turnoService := services.NewTurnoService(turnoRepo)
	turnoUseCase := usecases.NewTurnoUseCase(turnoService)
	turnoHandler := handlers.NewTurnoHandler(turnoUseCase)

	// Auth setup
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.ExpirationMinutes)*time.Minute)
	clientRepo := repositories.NewClientRepository(db.GetDB())
	authRepo := repositories.NewAuthRepository(db.GetDB())
	authService := services.NewAuthService(authRepo, clientRepo, jwtManager)
	authUseCase := usecases.NewAuthUseCase(authService)
	authHandler := handlers.NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Reporteria setup
	reporteriaRepo := repositories.NewReporteriaRepository(db.GetDB())
	reporteriaService := services.NewReporteriaService(reporteriaRepo)
	reporteriaUseCase := usecases.NewReporteriaUseCase(reporteriaService)
	reporteriaHandler := handlers.NewReporteriaHandler(reporteriaUseCase)

	// Pacientes setup
	pacienteRepo := repositories.NewPacienteRepository(db.GetDB())
	antecedentesMedicosRepo := repositories.NewAntecedentesMedicosRepository(db.GetDB())
	antecedentesVisualesRepo := repositories.NewAntecedentesVisualesRepository(db.GetDB())
	examenVisualRepo := repositories.NewExamenVisualRepository(db.GetDB())
	pacienteService := services.NewPacienteService(pacienteRepo, antecedentesMedicosRepo, antecedentesVisualesRepo, examenVisualRepo)
	pacienteUseCase := usecases.NewPacienteUseCase(pacienteService)
	pacienteHandler := handlers.NewPacienteHandler(pacienteUseCase)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Server is running",
		})
	})

	// Auth routes (public)
	routes.SetupAuthRoutes(router, authHandler, authMiddleware)

	// Protected routes - uncomment when ready to protect
	// For now, routes are public for backward compatibility
	routes.SetupUserRoutes(router, userHandler)
	routes.SetupRecetaRoutes(router, recetaHandler)
	routes.SetupTurnoRoutes(router, turnoHandler)
	routes.SetupReporteriaRoutes(router, reporteriaHandler, authMiddleware)
	routes.SetupPacienteRoutes(router, pacienteHandler, authMiddleware)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
