package controllers

import (
	"farishadibrata.com/rapidmono/app"
	dashboardView "farishadibrata.com/rapidmono/view/dashboard"
	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	app.AppInstance
}

func (a *DashboardController) New(appInstance *app.AppInstance) app.Controller {
	a = &DashboardController{
		AppInstance: *appInstance,
	}
	group := a.Fiber.Group("/dashboard")
	group.Get("/", a.GetDashboard)
	group.Get("/notification", a.GetNotification)
	return a
}

func (app *DashboardController) GetDashboard(ctx *fiber.Ctx) error {
	return app.Render(ctx, dashboardView.Index())
}

func (app *DashboardController) GetNotification(ctx *fiber.Ctx) error {
	return app.Render(ctx, dashboardView.NotificationPage())
}
