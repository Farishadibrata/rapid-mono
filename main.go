package main

import (
	"log"
	"os"

	"farishadibrata.com/rapidmono/app"
	"farishadibrata.com/rapidmono/controllers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func main() {
	godotenv.Load()

	// start Fiber
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if code == fiber.ErrNotFound.Code {
				c.Redirect("/", fiber.StatusTemporaryRedirect)
			}
			return nil
		},
	})
	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	fiberApp.Static("/static", "./dist")
	// To validate HTMX Request
	fiberApp.Use(func(c *fiber.Ctx) error {
		if c.Get("HX-Request") == "true" {
			c.Locals("IsHTMXRequest", true)
		}
		return c.Next()
	})

	// start logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// validation for Request
	validate = validator.New(validator.WithRequiredStructEnabled())
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	controllers := []app.Controller{
		&controllers.AuthController{},
		&controllers.HomeController{},
		&controllers.DashboardController{},
		&controllers.SystemManagementController{},
	}

	appInstance := &app.AppInstance{
		Fiber:      fiberApp,
		Db:         db,
		Logger:     logger,
		Cache:      rdb,
		Validation: &app.Validation{Validator: validate, Trans: trans},
	}
	for _, controller := range controllers {
		controller.New(appInstance)
	}

	log.Fatal(fiberApp.Listen(":4000"))
}
