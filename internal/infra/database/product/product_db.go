package database

import (
	"github.com/waanvieira/api-users/internal/entity"
	"gorm.io/gorm"
)

// Indica que nossa struct Product que seria do DB recebe a variável DB do gorm
type Product struct {
	DB *gorm.DB
}

// Cria uma struct para compor a nossa "classe" e o restante é do próprio ORM
func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

// Não retorna a nossa entity, retorna apenas um erro, então em algum lugar podemos chamar essa função e verifica apenas se tem um erro
func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindAll(page int, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	// Iniciamos a variavel de erro, se por acaso der algum erro retorna um erro
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		// Aqui informamos que na paginação o page -1 para sempre subtrair 1 e passando o sort, se encontra algum registro hidrata a variavel "products" se não retorna um erro
		// Nesse caso se existe registros e deu tudo certo a nossa variável "products" que vai ser hidratada, se der algum erro vai hidratar a variável error
		// em outra linguagem seria basicamente isso
		// product = faz_a_consulta_sql
		//  if (!product) { retorna um erro }
		// return product
		// No caso do GO e o GORM se vem o erro hidrata a variável err para retornar, porque podemos retornar 2 parametros na mesma função
		// Nesse caso a variável error vai retornar como nil, que seria em branco
		err = p.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		// Aqui usa da mesma base porém aqui faz um find e apenas ordena, retorna todos os dados apenas ordenado
		err = p.DB.Order("created_at " + sort).Find(&products).Error
	}

	return products, err
}

// (u *Product) - indica que a função é dessa nossa struct
// (id string) Nossao paramaetro que é uma string
// (*entity.Product, error) - Significa que retorna um ponteiro de Product da nossa entity ou retorna um erro
func (p *Product) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	// Os dados são preenchidos no Firs(&product), significa que não deu nenhum erro e vai hidratar o nosso ponteiro
	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) Update(product *entity.Product) error {
	// Verifica se o registro existe, se deixar apenas com o "Save" se o registro não existir ele vai gravar de qualquer forma
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(&product).Error
}

func (p *Product) Delete(id string) error {
	product, err := p.FindByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}
