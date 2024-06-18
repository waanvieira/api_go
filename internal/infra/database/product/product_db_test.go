package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waanvieira/api-users/internal/entity"

	// Banco em memória sqlite
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	product, _ := entity.NewProduct("product test", 10)
	// Basicamente iniciamos a struct
	productDB := NewProduct(db)

	// Nesse caso seria nosso repositorio recebendo a nossa entity para salvar no banco
	err = productDB.Create(product)
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error

	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
	assert.NotEmpty(t, productFound.CreatedAt)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	product, _ := entity.NewProduct("product test", 10)
	// Basicamente iniciamos a struct
	productDB := NewProduct(db)

	// Nesse caso seria nosso repositorio recebendo a nossa entity para salvar no banco
	err = productDB.Create(product)
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)

	product.Name = "name updated"
	product.Price = 20

	err = productDB.Update(product)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error

	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, "name updated", productFound.Name)
	assert.Equal(t, 20, productFound.Price)
}

func TestUpdateProductProductNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	product, _ := entity.NewProduct("product test", 10)
	fmt.Println(product)
	// Basicamente iniciamos a struct
	productDB := NewProduct(db)
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)

	product.Name = "name updated"
	product.Price = 20

	err = productDB.Update(product)
	assert.Error(t, err)

	product, err = productDB.FindByID(product.ID.String())
	assert.Nil(t, product)
}

func TestFindByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	product, _ := entity.NewProduct("product test", 20)
	// Basicamente iniciamos a struct
	productDB := NewProduct(db)

	// Nesse caso seria nosso repositorio recebendo a nossa entity para salvar no banco
	err = productDB.Create(product)
	assert.Nil(t, err)

	// Na nossa função retornamos a entidade do nosso "repository" ou um erro
	productByID, err := productDB.FindByID(product.ID.String())
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productByID.ID)
	assert.Equal(t, product.Name, productByID.Name)
	assert.Equal(t, product.Price, productByID.Price)
}

func TestDelete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	product, _ := entity.NewProduct("product test", 20)
	// Basicamente iniciamos a struct
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.Nil(t, err)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	product, err = productDB.FindByID(product.ID.String())
	assert.Nil(t, product)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	// Cria a nossa entity de product
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		// Nesse caso estamos criando o produto direto no banco, sem passar pelo nossa struct e usando o nativo do GORM
		db.Create(product)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	// Verificando a paginação, se está retornando 10 registros na 1° pagina
	assert.Len(t, products, 10)
	// Como o nosso FindAll retorna um slice de Products, podemos verificar igual um array para saber a 1° posição
	// Como criamos dinamicamente no for os nosso produtos sabemos que vai de 0 a 24 produtos criados
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	// Verificando a paginação, se está retornando 10 registros na 1° pagina
	assert.Len(t, products, 10)
	// Como o nosso FindAll retorna um slice de Products, podemos verificar igual um array para saber a 1° posição
	// Como criamos dinamicamente no for os nosso produtos sabemos que vai de 0 a 23 produtos criados
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(1, 10, "desc")
	assert.NoError(t, err)

	assert.Len(t, products, 10)
	assert.Equal(t, "Product 23", products[0].Name)
	assert.Equal(t, "Product 14", products[9].Name)

	products, err = productDB.FindAll(1, 15, "desc")
	assert.NoError(t, err)
	assert.Len(t, products, 15)
}
