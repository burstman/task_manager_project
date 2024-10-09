package chat

import (
	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router *chi.Mux) {

	router.Group(func(r chi.Router) {
		r.Post("/chatbot/project/", kit.Handler(handlePostProjects))
		r.Get("/chatbot/project/{user_id}_{name}_{deadline}", kit.Handler(handelGetProjectWithDeadline))
		r.Get("/chatbot/project/{user_id}_{name}", kit.Handler(handelGetProject))
		r.Get("/chatbot/project/single/{user_id}_{name}", kit.Handler(handelGetSingleProject))
		r.Put("/chatbot/project/{user_id}", kit.Handler(handelUpdateProjects))
		r.Get("/chatbot/project/all", kit.Handler(handlerGetAllProjects))
		r.Delete("/chatbot/project/{user_id}_{name}", kit.Handler(handelDeleteProjects))
		// r.Post("/chatbot/project/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	w.Write([]byte("test"))
		// }))
	})

}
