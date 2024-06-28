package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/waanvieira/api-users/configs"
	"github.com/waanvieira/api-users/internal/entity"
	databaseUser "github.com/waanvieira/api-users/internal/infra/database"
	databaseProduct "github.com/waanvieira/api-users/internal/infra/database/product"
	"github.com/waanvieira/api-users/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// iniciando as nossas config
	configs, err := configs.LoadConfig(".")

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
	// Cria logs em cada requisição
	r.Use(middleware.Logger)
	// Esse middleware serve para se der algum problema, algum panic no sistema o sistema continuar a execução e não parar
	// No caso quando acontece algum erro inesperado, algum panico a nossa aplicação ficaria fora do ar
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExperesIn", configs.JwtExpiresIn))	

	// // Estamos iniciando a struct de "classe" indicando qual banco de dados vamos usar
	productDB := databaseProduct.NewProduct(db)
	// Passamos a nossa "classe" concreta da nossa classe de manipulação de dados para o nosso handler (controller)
	// fazer as tratativas criando a entidade e salvando no banco
	produductHandler := handlers.NewProductHandler(productDB)

	userDB := databaseUser.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	// Injetamos o nosso método "CreateProduct" quando bater na rota de products
	r.Route("/products", func(r chi.Router) {
		// Middleware pega o nosso token e injeta no nosso contexto para poder pegar em qualquer rota, assim como fizemos no jwt
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		// Esse middleware que vai verificar se o nosso token é valido, dentro do tempo certo entre outras validações
		r.Use(jwtauth.Authenticator)
		r.Post("/", produductHandler.CreateProduct)
		r.Get("/", produductHandler.GetAllProducts)
		r.Get("/{id}", produductHandler.FindByID)
		r.Put("/{id}", produductHandler.UpdateProduct)
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

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		// r.Get("/{email}", userHandler.FindByEmail)
		r.Post("/generate_token", userHandler.GetJWT)

	})

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Teste"))
	})

	http.ListenAndServe(":8001", r)
}
