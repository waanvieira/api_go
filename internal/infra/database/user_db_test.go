package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waanvieira/api-users/internal/entity"

	// Banco em memória sqlite
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.User{})
	// Cria a nossa entity de user
	user, _ := entity.NewUser("user test", "user@teste.com", "123456")
	// Basicamente iniciamos a struct
	userDB := NewUser(db)

	// Nesse caso seria nosso repositorio recebendo a nossa entity para salvar no banco
	err = userDB.Create(user)
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error

	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)

}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	// Fazendo um migrate da tabela de usuário
	db.AutoMigrate(&entity.User{})
	// Cria a nossa entity de user
	user, _ := entity.NewUser("user test", "user@teste.com", "123456")
	// Basicamente iniciamos a struct
	userDB := NewUser(db)

	// Nesse caso seria nosso repositorio recebendo a nossa entity para salvar no banco
	err = userDB.Create(user)
	assert.Nil(t, err)

	// Na nossa função retornamos a entidade do nosso "repository" ou um erro
	userByEmail, err := userDB.FindByEmail(user.Email)
	// verificamos se não deu nenhum erro
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userByEmail.ID)
	assert.Equal(t, user.Name, userByEmail.Name)
	assert.Equal(t, user.Email, userByEmail.Email)
}
