package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prannvs/campus-leave-system/internal/api/handlers"
	"github.com/prannvs/campus-leave-system/internal/api/middleware"
	"github.com/prannvs/campus-leave-system/internal/auth"
	"github.com/prannvs/campus-leave-system/internal/models"
)

type Router struct {
	authHandler       *handlers.AuthHandler
	userHandler       *handlers.UserHandler
	leaveHandler      *handlers.LeaveHandler
	attendanceHandler *handlers.AttendanceHandler
	analyticsHandler  *handlers.AnalyticsHandler
	jwtService        *auth.JWTService
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	leaveHandler *handlers.LeaveHandler,
	attendanceHandler *handlers.AttendanceHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	jwtService *auth.JWTService,
) *Router {
	return &Router{
		authHandler:       authHandler,
		userHandler:       userHandler,
		leaveHandler:      leaveHandler,
		attendanceHandler: attendanceHandler,
		analyticsHandler:  analyticsHandler,
		jwtService:        jwtService,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", middleware.RoleMiddleware(models.RoleAdmin), r.userHandler.GetUsers)
				users.GET("/:id", r.userHandler.GetUser)
				users.DELETE("/:id", middleware.RoleMiddleware(models.RoleAdmin), r.userHandler.DeleteUser)
			}

			// Leave routes
			leaves := protected.Group("/leaves")
			{
				// Student routes
				leaves.POST("/apply", middleware.RoleMiddleware(models.RoleStudent), r.leaveHandler.ApplyLeave)
				leaves.GET("/my", middleware.RoleMiddleware(models.RoleStudent), r.leaveHandler.GetMyLeaves)

				// Faculty/Warden routes
				leaves.GET("/pending",
					middleware.RoleMiddleware(models.RoleFaculty, models.RoleWarden, models.RoleAdmin),
					r.leaveHandler.GetPendingLeaves)
				leaves.PUT("/:id/approve",
					middleware.RoleMiddleware(models.RoleFaculty, models.RoleWarden, models.RoleAdmin),
					r.leaveHandler.ApproveLeave)

				// Admin routes
				leaves.DELETE("/:id", middleware.RoleMiddleware(models.RoleAdmin), r.leaveHandler.DeleteLeave)
			}

			// Attendance routes
			attendance := protected.Group("/attendance")
			{
				attendance.POST("/mark",
					middleware.RoleMiddleware(models.RoleFaculty, models.RoleWarden, models.RoleAdmin),
					r.attendanceHandler.MarkAttendance)
				attendance.GET("/stats", r.attendanceHandler.GetAttendanceStats)
				attendance.GET("/low-attendance",
					middleware.RoleMiddleware(models.RoleFaculty, models.RoleWarden, models.RoleAdmin),
					r.attendanceHandler.GetLowAttendanceStudents)
			}

			// Analytics routes (Admin only)
			analytics := protected.Group("/analytics")
			analytics.Use(middleware.RoleMiddleware(models.RoleAdmin))
			{
				analytics.GET("/summary", r.analyticsHandler.GetAnalyticsSummary)
				analytics.GET("/leave-breakdown", r.analyticsHandler.GetLeaveTypeBreakdown)
			}
		}
	}

	return router
}
