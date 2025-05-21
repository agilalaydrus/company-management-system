package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"metro-backend/internal/letter"
	"metro-backend/internal/payslip"

	"metro-backend/internal/attendance"
	"metro-backend/internal/auth"
	"metro-backend/internal/employee"
	"metro-backend/internal/leave"
	"metro-backend/internal/middleware"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// Public routes
	router.POST("/api/register", auth.RegisterHandler(db))
	router.POST("/api/login", auth.LoginHandler(db))

	// Protected routes
	protected := router.Group("/api", middleware.JWTmiddleware())
	{
		protected.GET("/dashboard", auth.DashboardHandler)

		// Attendance
		protected.POST("/attendance", attendance.CreateAttendanceHandlers(db))

		// Employee
		protected.POST("/employees", employee.CreateEmployee(db))
		protected.GET("/employees", employee.GetAllEmployees(db))
		protected.GET("/employees/:id", employee.GetEmployeeByID(db))
		protected.PUT("/employees/:id", employee.UpdateEmployee(db))
		protected.DELETE("/employees/:id", employee.DeleteEmployee(db))

		// Leave
		protected.POST("/leave", leave.ApplyLeave(db))
		protected.GET("/leave", leave.ListLeave(db))
		protected.PUT("/leave/:id/approve", leave.ApproveLeave(db))
		protected.PUT("/leave/:id/reject", leave.RejectLeave(db))

		// Letter
		protected.POST("/letters", letter.CreateLetter(db))
		protected.GET("/letters", letter.GetLettersFiltered(db))
		protected.GET("/letters/:id/html", letter.GenerateLetterHTML(db))
		protected.GET("/letters/:id/pdf", letter.ExportLetterPDF(db, "http://localhost:8080"))

		// Payslip
		protected.POST("/payslips", payslip.CreatePayslip(db))
		protected.GET("/payslips", payslip.GetPayslips(db))
	}

}
