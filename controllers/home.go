package controllers

import (
	"farishadibrata.com/rapidmono/app"
	"github.com/gofiber/fiber/v2"
)

type HomeController struct {
	app.AppInstance

	// if you need additional data for controller, put here
	// SampleData string
}

func (a *HomeController) New(appInstance *app.AppInstance) app.Controller {
	a = &HomeController{
		AppInstance: *appInstance,
	}
	group := a.Fiber.Group("")
	group.Get("/", a.GetHome)

	return a
}

func (a *HomeController) GetHome(ctx *fiber.Ctx) error {
	ctx.Response().SetBodyString("Rapid Mono, Build Web Application Faster, Less build, and less state issue !")
	return nil
}
