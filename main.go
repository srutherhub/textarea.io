package main

import (
	"app/handlers"
	"app/services"

	c "github.com/srutherhub/web-app/controller"
	s "github.com/srutherhub/web-app/server"
)

func main() {
	cfg := s.InitServerCfg("5555")

	db := services.InitDB()
	db.InitTables()

	app := services.NewAppService(db)
	auth := services.NewAuthService(db)

	SpaceController := InitSpaceController(app, auth)
	BaseController := InitBaseController()
	UIController := InitUiController()

	server := s.New()

	server.RegisterController(*SpaceController)
	server.RegisterController(*BaseController)
	server.RegisterController(*UIController)

	server.Start(cfg)
}

func InitSpaceController(app *services.AppService, auth *services.AuthService) *c.Controller {
	m := handlers.Middleware{}
	controller := c.New()
	controller.SetBase("/area")
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/create", Handler: handlers.GetCreateSpace()})
	controller.RegisterRoute(c.Route{Method: "POST", Path: "/create", Handler: handlers.CreateSpace(app, auth)})
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/{id}", Handler: m.NoCache(m.HasAccess(auth, handlers.GetSpace(app, auth)))})
	controller.RegisterRoute(c.Route{Method: "POST", Path: "/{id}/save", Handler: m.HasAccess(auth, handlers.SaveSpace(app, auth))})
	controller.RegisterRoute(c.Route{Method: "DELETE", Path: "/{id}/delete", Handler: m.HasAccess(auth, handlers.DeleteSpace(app, auth))})
	controller.RegisterRoute(c.Route{Method: "GET", Path: "", Handler: handlers.Area(auth)})
	controller.RegisterRoute(c.Route{Method: "POST", Path: "/auth", Handler: handlers.AreaAuth(auth)})
	return controller
}

func InitBaseController() *c.Controller {
	controller := c.New()
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/", Handler: handlers.Base()})
	return controller
}

func InitUiController() *c.Controller {
	controller := c.New()
	controller.SetBase("/ui")
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/action/forget", Handler: handlers.UIForgetArea()})
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/action/delete", Handler: handlers.UIDeleteArea()})
	controller.RegisterRoute(c.Route{Method: "GET", Path: "/clear", Handler: handlers.UIClear()})
	return controller
}
