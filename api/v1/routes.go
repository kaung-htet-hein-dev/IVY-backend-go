package v1

import (
	"KaungHtetHein116/IVY-backend/api/v1/handler"
	"KaungHtetHein116/IVY-backend/internal/repository"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUserRoutes(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	userRoutes := e.Group("/api/v1/user")
	userRoutes.GET("", userHandler.GetAllUsers)
	// userRoutes.POST("/clerk-webhook", utils.BindAndValidateDecorator(userHandler.ClerkWebhook))
	userRoutes.POST("/register", utils.BindAndValidateDecorator(userHandler.RegisterUser))
	userRoutes.POST("/login", utils.BindAndValidateDecorator(userHandler.LoginUser))
	userRoutes.GET("/me", userHandler.GetMe)
	userRoutes.POST("/logout", userHandler.Logout)
	userRoutes.PUT("/:id", utils.BindAndValidateDecorator(userHandler.UpdateUser))
}

func RegisterBranchRoutes(e *echo.Echo, db *gorm.DB) {
	branchRepo := repository.NewBranchRepository(db)
	branchUsecase := usecase.NewBranchUsecase(branchRepo)
	branchHandler := handler.NewBranchHandler(branchUsecase)

	branchRoutes := e.Group("/api/v1/branch")
	branchRoutes.POST("", utils.BindAndValidateDecorator(branchHandler.CreateBranch))
	branchRoutes.GET("", branchHandler.GetAllBranches)
	branchRoutes.GET("/:id", branchHandler.GetBranchByID)
	branchRoutes.PUT("/:id", utils.BindAndValidateDecorator(branchHandler.UpdateBranch))
	branchRoutes.DELETE("/:id", branchHandler.DeleteBranch)
}

func RegisterCategoryRoutes(e *echo.Echo, db *gorm.DB) {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)

	categoryRoutes := e.Group("/api/v1/category")
	categoryRoutes.POST("", utils.BindAndValidateDecorator(categoryHandler.CreateCategory))
	categoryRoutes.GET("", categoryHandler.GetAllCategories)
	categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
	categoryRoutes.PUT("/:id", utils.BindAndValidateDecorator(categoryHandler.UpdateCategory))
	categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
}

func RegisterServiceRoutes(e *echo.Echo, db *gorm.DB) {
	serviceRepo := repository.NewServiceRepository(db)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo)
	serviceHandler := handler.NewServiceHandler(serviceUsecase)

	serviceRoutes := e.Group("/api/v1/service")
	serviceRoutes.POST("", utils.BindAndValidateDecorator(serviceHandler.CreateService))
	serviceRoutes.GET("", serviceHandler.GetAllServices)
	serviceRoutes.GET("/:id", serviceHandler.GetServiceByID)
	serviceRoutes.PUT("/:id", utils.BindAndValidateDecorator(serviceHandler.UpdateService))
	serviceRoutes.DELETE("/:id", serviceHandler.DeleteService)
}

func RegisterBookingRoutes(e *echo.Echo, db *gorm.DB) {
	bookingRepo := repository.NewBookingRepository(db)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)

	bookingRoutes := e.Group("/api/v1/booking")
	bookingRoutes.POST("", utils.BindAndValidateDecorator(bookingHandler.CreateBooking))
	bookingRoutes.GET("", bookingHandler.GetAllBookings)
	bookingRoutes.GET("/slots", bookingHandler.GetAvailableSlots)
	bookingRoutes.GET("/me", bookingHandler.GetUserBookings)
	bookingRoutes.GET("/:id", bookingHandler.GetBookingByID)
	bookingRoutes.PUT("/:id", utils.BindAndValidateDecorator(bookingHandler.UpdateBooking))
	bookingRoutes.DELETE("/:id", bookingHandler.DeleteBooking)
}
