package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"metro-backend/internal/atk"
	"metro-backend/internal/attendance"
	"metro-backend/internal/auth"
	"metro-backend/internal/company"
	"metro-backend/internal/employee"
	"metro-backend/internal/inventory"
	"metro-backend/internal/leave"
	"metro-backend/internal/letter"
	"metro-backend/internal/middleware"
	"metro-backend/internal/payslip"
	"metro-backend/internal/warehouse"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	// ============================
	// 📂 Public routes (no auth)
	// ============================
	router.POST("/api/register", auth.RegisterHandler(db))
	router.POST("/api/login", auth.LoginHandler(db))

	// Serve static files (uploads/photos)
	router.Static("/uploads", "./uploads")

	// ============================
	// 🔐 Protected routes (JWT)
	// ============================
	protected := router.Group("/api", middleware.JWTmiddleware())
	{
		// Dashboard
		protected.GET("/dashboard", auth.DashboardHandler)

		// Attendance
		protected.POST("/attendance", attendance.CreateAttendanceHandler(db))
		protected.POST("/attendance/clockout", attendance.ClockOutHandler(db))
		protected.GET("/attendance/:employee_id", attendance.GetAttendanceHistoryByEmployee(db))

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

		// Companies
		protected.POST("/companies", company.Create(db))
		protected.GET("/companies", company.List(db))
		protected.GET("/companies/:id", company.Detail(db))
		protected.PUT("/companies/:id", company.Update(db))
		protected.DELETE("/companies/:id", company.Delete(db))

		// Warehouse
		protected.POST("/warehouses", warehouse.Create(db))
		protected.GET("/warehouses", warehouse.List(db))
		protected.GET("/warehouses/:id", warehouse.Detail(db))
		protected.PUT("/warehouses/:id", warehouse.Update(db))
		protected.DELETE("/warehouses/:id", warehouse.Delete(db))

		// Inventory Items
		protected.POST("/inventory-items", inventory.Create(db))
		protected.GET("/inventory-items", inventory.List(db))
		protected.GET("/inventory-items/:id", inventory.Detail(db))
		protected.PUT("/inventory-items/:id", inventory.Update(db))
		protected.DELETE("/inventory-items/:id", inventory.Delete(db))

		// ATK Items
		protected.POST("/atk-items", atk.Create(db))
		protected.GET("/atk-items", atk.List(db))
		protected.GET("/atk-items/:id", atk.Detail(db))
		protected.PUT("/atk-items/:id", atk.Update(db))
		protected.DELETE("/atk-items/:id", atk.Delete(db))
	}
}
