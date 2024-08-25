package chat

import (
	"taskManager/plugins/auth"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router *chi.Mux) {
	auth.InitializeRoutes(router)
	authConfig := kit.AuthenticationConfig{
		AuthFunc:    auth.AuthenticateUser,
		RedirectURL: "/login",
	}
	router.Group(func(r chi.Router) {
		r.Use(kit.WithAuthentication(authConfig, false))
		r.Post("/chatbot/project/", kit.Handler(handlePostProjects))
		r.Get("/chatbot/project/{user_id}-{name}-{searchterm}", kit.Handler(handelGetProject))
		r.Put("/chatbot/project/", kit.Handler(handelUpdateProjects))
		r.Get("/chatbot/project/all", kit.Handler(handlerGetAllProjects))
	})

}
