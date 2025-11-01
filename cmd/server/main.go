package main

import (
	"log"

	"github.com/prannvs/campus-leave-system/internal/api/handlers"
	"github.com/prannvs/campus-leave-system/internal/api/routes"
	"github.com/prannvs/campus-leave-system/internal/auth"
	"github.com/prannvs/campus-leave-system/internal/core"
	"github.com/prannvs/campus-leave-system/internal/repositories"
	"github.com/prannvs/campus-leave-system/internal/services"
	"github.com/prannvs/campus-leave-system/pkg/db"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbConfig := db.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	if err := db.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	database := db.GetDB()

	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry)

	userRepo := repositories.NewUserRepository(database)
	leaveRepo := repositories.NewLeaveRepository(database)
	attendanceRepo := repositories.NewAttendanceRepository(database)

	notificationService := services.NewNotificationService(cfg.SMTP)
	userService := services.NewUserService(userRepo)
	leaveService := services.NewLeaveService(leaveRepo, attendanceRepo, notificationService)
	attendanceService := services.NewAttendanceService(attendanceRepo)

	authHandler := handlers.NewAuthHandler(userService, jwtService)
	userHandler := handlers.NewUserHandler(userService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	analyticsHandler := handlers.NewAnalyticsHandler(leaveService, attendanceService)

	router := routes.NewRouter(
		authHandler,
		userHandler,
		leaveHandler,
		attendanceHandler,
		analyticsHandler,
		jwtService,
	)

	engine := router.Setup()

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
