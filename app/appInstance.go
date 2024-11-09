package app

import (
	"github.com/a-h/templ"
	"github.com/casbin/casbin/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Validation struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

type AppInstance struct {
	Fiber      *fiber.App
	Db         *sqlx.DB
	Logger     *zap.Logger
	Cache      *redis.Client
	Validation *Validation
	Enforcer   *casbin.Enforcer
}

type Controller interface {
	New(AppInstance *AppInstance) Controller
}

type AuthorizationMethod int

const (
	Read AuthorizationMethod = iota + 1
	Create
	Update
	Delete
)

func (app *AppInstance) ValidateStruct(ctx *fiber.Ctx, input interface{}, component func(string, string) templ.Component) error {
	err := app.Validation.Validator.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg := err.Translate(app.Validation.Trans)
			app.Render(ctx, component("error", errMsg))
			return err
		}
	}
	return nil
}

func (app *AppInstance) Render(ctx *fiber.Ctx, component templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	ctx.Locals("Redirect", ctx.Redirect)
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func (app *AppInstance) Enforce() error {
	return nil
}
