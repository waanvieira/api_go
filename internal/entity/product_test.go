package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("test", 15)
	// Verifica que o erro está em branco
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.NotEmpty(t, product.Name)
	assert.NotEmpty(t, product.CreatedAt)
	assert.Equal(t, "test", product.Name)
	assert.Equal(t, 15.0, product.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	// Criando um produto com nome zero para disparar o erro
	p, err := NewProduct("", 0)
	// Testando se a variavel está em branco, o inverso do test anterior que verificamos se o erro está em branco
	assert.Nil(t, p)
	// Verifica qual é o erro que apresentou, nesse caso tem que ser o erro de nome obrigatório
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequiredAndInvalid(t *testing.T) {
	// Criando um produto com nome zero para disparar o erro
	p, err := NewProduct("test", 0)
	// Testando se a variavel está em branco, o inverso do test anterior que verificamos se o erro está em branco
	assert.Nil(t, p)
	// Verifica qual é o erro que apresentou, nesse caso tem que ser o erro de nome obrigatório
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("test", -10)
	// Testando se a variavel está em branco, o inverso do test anterior que verificamos se o erro está em branco
	assert.Nil(t, p)
	// Verifica qual é o erro que apresentou, nesse caso tem que ser o erro de nome obrigatório
	assert.Equal(t, ErrInvalidPrice, err)

}
