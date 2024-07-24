package database

import (
	"database/sql"

	"github.com/google/uuid"
)

// Aqui seria uma struc de BANCO DE DADOS, estamos trabalhando direto com banco de dados, então seria uma mistura de model com DB, aqui recebemos o sql.DB para manipular o banco de dados
type Category struct {
	// Aqui temos acesso a todas as interfaces da nossa lib de sql.db que seria type DB struct
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

// Criamos aqui como se fosse um construtor, para quando chamarmos a struct category no nosso main, será "instanciado" com o DB
func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) Create(name, description string) (Category, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories(id, name, description) VALUES($1, $2, $3)", id, name, description)
	if err != nil {
		return Category{}, err
	}

	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) GetAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []Category{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return nil, err
		}

		categories = append(categories, Category{ID: id, Name: name, Description: description})
	}

	return categories, nil
}

// Aqui estamos fazendo a nível didatico, no mundo real poderiamos fazer um join direto na tabela de courses que ganhariamos performance gastando menos recursos
func (c *Category) GetCategoryByCourseID(courseID string) (Category, error) {
	// No query row nós podemos retornar apenas 1 registro
	var id, name, description string
	err := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = $1", courseID).Scan(&id, &name, &description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) GetCategoryByID(ID string) (Category, error) {
	// No query row nós podemos retornar apenas 1 registro
	var id, name, description string
	err := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", ID).Scan(&id, &name, &description)
	if err != nil {
		return Category{}, err
	}
	return Category{ID: id, Name: name, Description: description}, nil
}
