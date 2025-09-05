package api

import (
	"KaungHtetHein116/IVY-backend/api/middleware"
	v1 "KaungHtetHein116/IVY-backend/api/v1"
	"KaungHtetHein116/IVY-backend/config"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/clerk/clerk-sdk-go/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func StartServer() {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	db := config.ConnectDB()

	db.AutoMigrate(
		&entity.Booking{},
		&entity.Branch{},
		&entity.Category{},
		&entity.Service{},
		&entity.User{},
	)

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	middleware.RegisterBasicMiddlewares(e)
	middleware.RegisterAuthMiddleware(e)

	v1.RegisterUserRoutes(e, db)
	v1.RegisterBranchRoutes(e, db)
	v1.RegisterCategoryRoutes(e, db)
	v1.RegisterServiceRoutes(e, db)
	v1.RegisterBookingRoutes(e, db)

	port := ":" + os.Getenv("APP_PORT")

	log.Fatal(e.Start(port))
}
