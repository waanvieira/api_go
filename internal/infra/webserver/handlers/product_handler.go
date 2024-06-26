package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/waanvieira/api-users/internal/dto"
	"github.com/waanvieira/api-users/internal/entity"
	"github.com/waanvieira/api-users/internal/infra/database"
	entityPkg "github.com/waanvieira/api-users/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

// Aqui é basicamente o nosso construtor, indicando que estamos recebendo a interface, e não a classe concreta
// Isso é inversão de dependencia
func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary      Create product
// @Description  Create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /products [post]
// @Security ApiKeyAuth
// Função que seria o nosso controller, recebe um request e retorna um response
// Seria um método do nosso controler como store(Resquest $request) {//cadastra no banco de dados}
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	// Lembrando aqui, declaramos a variável "erro" com := com declaração e valor porque se der algo de errado
	// retornamos erro, dando tudo certo o Decode(&product) vai hidratar o nosso ponteiro de product assim
	// a variável vai ficar com os valores recebidos no body
	// Aqui está fazendo, pegando os dados vindos da nossa API no r.Body, fazendo um bind com nosso dto
	// Hidratando a nossa variável product
	erro := json.NewDecoder(r.Body).Decode(&product)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro DTO"))
		return
	}

	// Aqui criamos a nossa entidade com os 2 parametros que precisamos
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro para criar"))
		return
	}
	//  Aqui fazemos o cadastro no banco de dados, com a nossa injeção de dependencia do productDB
	// no nosso handler da struct
	// Seria bsicamente fazer igual nos testes
	// "instanciar a classe"
	// productDB := NewProduct(db)
	// Chamar o método de create para salvar no banco passando a entidade
	// err = productDB.Create(product)

	err = h.ProductDB.Create(p)
	// atribuimos o create a erro porque é uma função void, não tem retorno, então validamos
	// se a variável for diferente de nil retornamos um bad request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro para salvar no banco"))

		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Get a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "product ID" Format(uuid)
// @Success      200  {object}  entity.Product
// @Failure      404
// @Failure      500  {object}  Error
// @Router       /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Registro não encontrado"))
		return
	}

	json.NewEncoder(w).Encode(p)
}

// ListAccounts godoc
// @Summary      List products
// @Description  get all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page      query     string  false  "page number"
// @Param        limit     query     string  false  "limit"
// @Success      200       {array}   entity.Product
// @Failure      404       {object}  Error
// @Failure      500       {object}  Error
// @Router       /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	// Automaticamente quando recebemos parametros ele vem como string o pacote strconv é para converter string em number
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	sort := r.URL.Query().Get("sort")
	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        path      string                  true  "product ID" Format(uuid)
// @Success      200
// @Failure      404
// @Failure      500       {object}  Error
// @Router       /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Registro não encontrado"))
		return
	}

	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Registro não encontrado"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Registro deletado com sucesso"))
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id        	path      string                  true  "product ID" Format(uuid)
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      200
// @Failure      404
// @Failure      500       {object}  Error
// @Router       /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Aqui atribuimos a variável como referencia porque o valor já foi setado anteriormente, aqui estamos basicamente atribuindo um novo valor ao err, se mudassemos o valor o nome da variável
	// teriamos que indicar := que seria atribuição do valor na variável err
	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
