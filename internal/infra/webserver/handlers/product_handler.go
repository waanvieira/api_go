package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/waanvieira/api-users/internal/dto"
	"github.com/waanvieira/api-users/internal/entity"
	"github.com/waanvieira/api-users/internal/infra/database"
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
