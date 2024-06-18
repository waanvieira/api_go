package database

import (
	"github.com/waanvieira/api-users/internal/entity"
	"gorm.io/gorm"
)

// Indica que nossa struct User que seria do DB recebe a variável DB do gorm
type User struct {
	DB *gorm.DB
}

// Cria uma struct para compor a nossa "classe" e o restante é do próprio ORM
func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

// Não retorna a nossa entity, retorna apenas um erro, então em algum lugar podemos chamar essa função e verifica apenas se tem um erro
func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

// (u *User) - indica que a função é dessa nossa struct
// (email string) Nossao paramaetro que é uma string
// (*entity.User, error) - Significa que retorna um ponteiro de User da nossa entity ou retorna um erro
func (u *User) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	// Os dados são preenchidos no Firs(&user), significa que não deu nenhum erro e vai hidratar o nosso ponteiro
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
