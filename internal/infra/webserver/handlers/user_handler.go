package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/waanvieira/api-users/internal/dto"
	user_dto "github.com/waanvieira/api-users/internal/dto/users"
	"github.com/waanvieira/api-users/internal/entity"
	"github.com/waanvieira/api-users/internal/infra/database"
)

// Struct que representa como vamos retornar um erro caso tenha algum problema
type Error struct {
	Message string `json:"message"`
}
type UserHandler struct {
	UserDB database.UserInterface
}

// Aqui é basicamente o nosso construtor, indicando que estamos recebendo a interface, e não a classe concreta
// Isso é inversão de dependencia
func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWT godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     user_dto.GetJWTInput  true  "user credentials"
// @Success      200  {object}  user_dto.GetJWTOutput
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExperesIn").(int)
	var user user_dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	accessToken := user_dto.GetJWTOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /users [post]
// Função que seria o nosso controller, recebe um request e retorna um response
// Seria um método do nosso controler como store(Resquest $request) {//cadastra no banco de dados}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	erro := json.NewDecoder(r.Body).Decode(&user)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro DTO"))
		return
	}

	userDB, _ := h.UserDB.FindByEmail(user.Email)
	if userDB != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Usuário já registrado"))
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) FindByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := h.UserDB.FindByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Registro não encontrado"))
		return
	}

	json.NewEncoder(w).Encode(p)
}

// func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
// 	page := r.URL.Query().Get("page")
// 	limit := r.URL.Query().Get("limit")
// 	// Automaticamente quando recebemos parametros ele vem como string o pacote strconv é para converter string em number
// 	pageInt, err := strconv.Atoi(page)
// 	if err != nil {
// 		pageInt = 0
// 	}
// 	limitInt, err := strconv.Atoi(limit)
// 	if err != nil {
// 		limitInt = 0
// 	}

// 	sort := r.URL.Query().Get("sort")
// 	users, err := h.UserDB.FindAll(pageInt, limitInt, sort)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(users)
// }

// func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	if id == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	_, err := h.UserDB.FindByID(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("Registro não encontrado"))
// 		return
// 	}

// 	err = h.UserDB.Delete(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("Registro não encontrado"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// 	w.Write([]byte("Registro deletado com sucesso"))
// }

// func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	if id == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	var user entity.User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	user.ID, err = entityPkg.ParseID(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	_, err = h.UserDB.FindByID(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	// Aqui atribuimos a variável como referencia porque o valor já foi setado anteriormente, aqui estamos basicamente atribuindo um novo valor ao err, se mudassemos o valor o nome da variável
// 	// teriamos que indicar := que seria atribuição do valor na variável err
// 	err = h.UserDB.Update(&user)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
