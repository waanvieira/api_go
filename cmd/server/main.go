package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/waanvieira/api-users/configs"
	_ "github.com/waanvieira/api-users/docs"
	"github.com/waanvieira/api-users/internal/entity"
	databaseUser "github.com/waanvieira/api-users/internal/infra/database"
	databaseProduct "github.com/waanvieira/api-users/internal/infra/database/product"
	"github.com/waanvieira/api-users/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wanderson Vieira

// @license.name   W

// @host      localhost:8001
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	// r.Use(LogOwner) // Aqui estamos utilizando o nosso próprio middleware dentro do ichi
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
		// Nesse caso todas as rotas dentro de /products estão com o middleware do JWT, se quisermos rotas sem jwt colocar fora
		// Como é com as rotas de usuários
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

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8001/docs/doc.json")))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Teste"))
	})

	http.ListenAndServe(":8001", r)
}

// exemplo de um middleware próprio, seria um exmeplo de middleware para validar por exemplo ACL, verificar permissionamento de um usuário dependendo
// Do seu perfil
// Next -> seria apenas uma convenção usada para referenciar a variável recebida na request, seria um meio campo para acessar a nossa aplicação
// nesse caso next recebe um handler http e retorna esse handler, na função podemos fazer alguma validação, alteração e ou retornar um erro
// Direcionar para outra função ou middleware ou passar para frente essa requisição
func LogOwner(next http.Handler) http.Handler {
	// Nesse caso retornamos um handler como recebemos pegando uma resposta response e uma request (requisição)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Nesse caso executamos uma ação na nossa aplicação que seria imprimir na tela o método e a URL que recebemos a requisição
		// Nesse caos nó estamos pegando o contexto do usuário, podemos recuperar alguma informação no contexto com o middleware
		log.Println(r.Method, r.URL.Path)
		// Aqui damos continuidade na aplicação com next.ServeHTTP que é da própria linguagem
		next.ServeHTTP(w, r)
	})
}
