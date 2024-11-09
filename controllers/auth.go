package controllers

import (
	"context"
	"fmt"
	"time"

	"farishadibrata.com/rapidmono/app"
	"farishadibrata.com/rapidmono/app/hash"
	"farishadibrata.com/rapidmono/model"
	authView "farishadibrata.com/rapidmono/view/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthController struct {
	app.AppInstance
}

func (a *AuthController) New(appInstance *app.AppInstance) app.Controller {
	a = &AuthController{
		AppInstance: *appInstance,
	}

	group := a.Fiber.Group("/auth")

	group.Use(a.AuthMiddleware)
	group.Get("/login", a.GetLogin)
	group.Get("/register", a.GetRegister)
	group.Get("/forgot-password", a.GetForgotPassword)
	group.Post("/login", a.PostLogin)
	group.Post("/register", a.PostRegister)
	return a
}

func (app *AuthController) AuthMiddleware(ctx *fiber.Ctx) error {
	// if ctx.Cookies("session") != "" {
	// 	return ctx.Redirect("/dashboard/")
	// }
	return ctx.Next()
}

func (app *AuthController) GetLogin(ctx *fiber.Ctx) error {
	return app.Render(ctx, authView.Home())
}

func (app *AuthController) GetRegister(ctx *fiber.Ctx) error {
	return app.Render(ctx, authView.Register())
}

func (app *AuthController) GetForgotPassword(ctx *fiber.Ctx) error {
	return app.Render(ctx, authView.ForgotPassword())
}

func (app *AuthController) PostLogin(ctx *fiber.Ctx) error {
	userDataInput := model.User{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	if err := app.ValidateStruct(ctx, userDataInput, authView.FormLogin); err != nil {
		return err
	}

	var parsedUser model.User

	if err := app.Db.Get(&parsedUser, `SELECT id, "uuid", email, "password", user_email_verified FROM users WHERE email = $1`, userDataInput.Email); err != nil {
		app.Logger.Error("Failed to get user", zap.Error(err))
		return app.Render(ctx, authView.FormLogin("error", "Invalid email or password"))
	}

	if match, err := hash.NewArgonParams().CheckPassword(userDataInput.Password, parsedUser.Password); err != nil || !match {
		app.Logger.Error("Failed to get user", zap.Error(err))
		return app.Render(ctx, authView.FormLogin("error", "Invalid email or password"))
	}

	sessionId := fmt.Sprintf("session-%s", uuid.New().String())

	if err := app.Cache.Set(context.Background(), sessionId, parsedUser.UUID, 24*time.Hour).Err(); err != nil {
		app.Logger.Error("Failed to set session", zap.Error(err))
		return app.Render(ctx, authView.FormLogin("error", "Invalid email or password"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "session",
		Value: sessionId,
	})

	ctx.Response().Header.Set("HX-Redirect", "/dashboard/")

	return app.Render(ctx, authView.FormLogin("success", "Login success, please refresh the page if you're not redirected"))
}

func (app *AuthController) PostRegister(ctx *fiber.Ctx) error {
	userDataInput := model.User{
		Email:                ctx.FormValue("email"),
		Password:             ctx.FormValue("password"),
		PasswordConfirmation: ctx.FormValue("repeat_password"),
	}

	if err := app.ValidateStruct(ctx, userDataInput, authView.FormRegister); err != nil {
		return err
	}

	isEmailExist := false
	if err := app.Db.Get(&isEmailExist, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", ctx.FormValue("email")); err != nil {
		app.Logger.Error("Failed to check email", zap.Error(err))
		return app.Render(ctx, authView.FormRegister("error", "Internal server error"))
	}

	if isEmailExist {
		return app.Render(ctx, authView.FormRegister("error", "Email already exist"))
	}

	passwordHash, err := hash.NewArgonParams().GeneratePasswordHash(userDataInput.Password)
	if err != nil {
		app.Logger.Error("Failed to hash password", zap.Error(err))
		return app.Render(ctx, authView.FormRegister("error", "Internal server error"))
	}

	userDataInput.Password = passwordHash

	if _, err := app.Db.NamedExec("INSERT INTO users (email, password) VALUES (:email, :password)", userDataInput); err != nil {
		app.Logger.Error("Failed to insert user", zap.Error(err))
		return app.Render(ctx, authView.FormRegister("error", "Internal server error"))
	}

	return app.Render(ctx, authView.FormRegister("success", "Register success. Please check your email for verification"))
}
