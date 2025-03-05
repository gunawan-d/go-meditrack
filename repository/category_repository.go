package repository

import (
	"database/sql"
	"meditrack/structs"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) CreateCategory(category structs.Category) (int, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	var id int
	err := r.DB.QueryRow(query, category.Name).Scan(&id)
	return id, err
}

func (r *CategoryRepository) GetCategories() ([]structs.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []structs.Category
	for rows.Next() {
		var c structs.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
