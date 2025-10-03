package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"iris-api/internal/application/usecases"
	"iris-api/internal/domain/services"
	"iris-api/internal/infrastructure/database"
	"iris-api/internal/infrastructure/repositories"
	"iris-api/internal/presentation/handlers"
	"iris-api/internal/presentation/routes"
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

	routes.SetupUserRoutes(router, userHandler)
	routes.SetupRecetaRoutes(router, recetaHandler)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
