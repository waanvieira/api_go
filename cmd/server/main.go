package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/waanvieira/api-users/configs"
	"github.com/waanvieira/api-users/internal/entity"
	database "github.com/waanvieira/api-users/internal/infra/database/product"
	"github.com/waanvieira/api-users/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// iniciando as nossas config
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic("erro ao criar user")
	}
	// // Indicando qual banco vamor usar, nesse caso o sqlite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Criando as nossas migracoes
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// // Estamos iniciando a struct de "classe" indicando qual banco de dados vamos usar
	productDB := database.NewProduct(db)
	// Passamos a nossa "classe" concreta da nossa classe de manipulação de dados para o nosso handler (controller)
	// fazer as tratativas criando a entidade e salvando no banco
	produductHandler := handlers.NewProductHandler(productDB)
	// // // Injetamos o nosso método "CreateProduct" quando bater na rota de products
	r.Route("/products", func(r chi.Router) {
		r.Post("/", produductHandler.CreateProduct)
		r.Get("/", produductHandler.GetAllProducts)
		r.Get("/{id}", produductHandler.FindByID)
		// r.Put("/{id}", produductHandler.UpdateProduct)
		// userID := chi.URLParam(r, "userID")
		r.Delete("/{id}", produductHandler.DeleteProduct)
		// Subrouters:
		// r.Route("/{id}", func(r chi.Router) {
		// 	r.Use(ArticleCtx)
		// 	r.Get("/", getArticle)                                          // GET /articles/123
		// 	r.Put("/", updateArticle)                                       // PUT /articles/123
		// 	r.Delete("/", deleteArticle)                                    // DELETE /articles/123
		// 	})

	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Teste"))
	})
	http.ListenAndServe(":8000", r)
}
