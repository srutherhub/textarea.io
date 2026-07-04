package main

import (
	c "github.com/srutherhub/web-app/controller"
	s "github.com/srutherhub/web-app/server"
)

func main() {
	cfg := s.InitServerCfg("5555")
	app := s.New()

	dbservice := InitDB()
	appservice := InitAppService(dbservice)

	app.RegisterController(*BaseController(appservice))
	app.RegisterController(*ArticleController(appservice))

	app.Start(cfg)
}

func BaseController(app *AppService) *c.Controller {
	base := c.New()
	base.SetBase("/")
	base.RegisterRoute(c.Route{Method: "GET", Path: "", Handler: BaseViewHandler()})
	return base
}

func ArticleController(app *AppService) *c.Controller {
	article := c.New()
	article.SetBase("/article")
	article.RegisterRoute(c.Route{Method: "GET", Path: "/create", Handler: CreateArticleViewHandler()})
	article.RegisterRoute(c.Route{Method: "POST", Path: "/create", Handler: CreateArticleHandler(app)})
	article.RegisterRoute(c.Route{Method: "GET", Path: "/{id}/edit", Handler: EditArticleViewHandler(app)})
	article.RegisterRoute(c.Route{Method: "POST", Path: "/{id}/save", Handler: SaveArticleHandler(app)})
	return article
}
