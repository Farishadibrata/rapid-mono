package main

import (
	"log"
	"os"

	"farishadibrata.com/rapidmono/app"
	"farishadibrata.com/rapidmono/controllers"
	baseView "farishadibrata.com/rapidmono/view/base"
	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

var CASBIN_MODEL = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act) || r.sub == "root"
`

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
				if c.Get("HX-Request") == "true" {
					c.Set("Content-Type", "text/html")
					return baseView.NotFoundPage().Render(c.Context(), c.Response().BodyWriter())
				}
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

	casbinAdapter, err := sqlxadapter.NewAdapter(db, "casbin_rule_test")
	if err != nil {
		logger.Fatal("Failed to initialize casbin adapter", zap.Error(err))
	}

	casbinModel, err := model.NewModelFromString(CASBIN_MODEL)
	if err != nil {
		logger.Fatal("Failed to initialize casbin adapter", zap.Error(err))
	}

	enforcer, err := casbin.NewEnforcer(casbinModel, casbinAdapter)
	if err != nil {
		logger.Fatal("Failed to initialize casbin enforcer.", zap.Error(err))
	}

	if err = enforcer.LoadPolicy(); err != nil {
		logger.Fatal("Failed to initialize casbin LoadPolicy.", zap.Error(err))
	}

	appInstance := &app.AppInstance{
		Fiber:      fiberApp,
		Db:         db,
		Logger:     logger,
		Cache:      rdb,
		Validation: &app.Validation{Validator: validate, Trans: trans},
		Enforcer:   enforcer,
	}
	for _, controller := range controllers {
		controller.New(appInstance)
	}

	log.Fatal(fiberApp.Listen(":4000"))
}
