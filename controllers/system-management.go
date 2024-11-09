package controllers

import (
	"farishadibrata.com/rapidmono/app"
	"farishadibrata.com/rapidmono/model"
	siteSettingView "farishadibrata.com/rapidmono/view/system-management/site-settings"
	userManagementView "farishadibrata.com/rapidmono/view/system-management/user-management"
	"github.com/gofiber/fiber/v2"
)

type SystemManagementController struct {
	app.AppInstance
}

func (a *SystemManagementController) New(appInstance *app.AppInstance) app.Controller {
	a = &SystemManagementController{
		AppInstance: *appInstance,
	}
	group := a.Fiber.Group("/system-management")
	group.Get("/site-settings", a.GetSiteSettings)
	group.Get("/user-management", a.GetUserManagement)

	return a
}

func (a *SystemManagementController) GetSiteSettings(ctx *fiber.Ctx) error {
	return a.Render(ctx, siteSettingView.Index())
}

func (a *SystemManagementController) GetUserManagement(ctx *fiber.Ctx) error {
	var users []model.User
	a.Db.Get(&users, "SELECT * FROM users LIMIT 10")
	return a.Render(ctx, userManagementView.Index())
}
