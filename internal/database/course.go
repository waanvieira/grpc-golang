package database

import (
	"database/sql"

	"github.com/google/uuid"
)

// Aqui seria uma struc de BANCO DE DADOS, estamos trabalhando direto com banco de dados, então seria uma mistura de model com DB, aqui recebemos o sql.DB para manipular o banco de dados
type Course struct {
	// Aqui temos acesso a todas as interfaces da nossa lib de sql.db que seria type DB struct
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

// Criamos aqui como se fosse um construtor, para quando chamarmos a struct course no nosso main, será "instanciado" com o DB
func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, categoryID, description string) (Course, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO courses(id, category_id, name, description) VALUES($1, $2, $3, $4)", id, categoryID, name, description)
	if err != nil {
		return Course{}, err
	}

	return Course{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c *Course) GetAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, category_id, name, description FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []Course{}
	for rows.Next() {
		var id, name, description, categoryID string
		// Ele vai hidratar minhas variáveis na mesma ordem dos campos retornados do banco, nesse caso
		// id, category_id, name, description
		if err := rows.Scan(&id, &categoryID, &name, &description); err != nil {
			return nil, err
		}

		categories = append(categories, Course{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}

	return categories, nil
}

func (c *Course) GetAllByCategoryId(categoryID string) ([]Course, error) {
	rows, err := c.db.Query("SELECT id, category_id, name, description FROM courses where category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := []Course{}
	for rows.Next() {
		var id, name, description, categoryID string
		// Ele vai hidratar minhas variáveis na mesma ordem dos campos retornados do banco, nesse caso
		// id, category_id, name, description
		if err := rows.Scan(&id, &categoryID, &name, &description); err != nil {
			return nil, err
		}

		categories = append(categories, Course{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}

	return categories, nil
}
