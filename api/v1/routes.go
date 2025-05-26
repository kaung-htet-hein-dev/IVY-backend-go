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
	userRoutes.POST("/register", utils.BindAndValidateDecorator(userHandler.RegisterUser))
	userRoutes.POST("/login", utils.BindAndValidateDecorator(userHandler.LoginUser))
	userRoutes.GET("/me", userHandler.GetMe)
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
