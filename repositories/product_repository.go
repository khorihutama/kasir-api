package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db           *sql.DB
	categoryRepo *CategoryRepository
}

func NewProductRepository(db *sql.DB, categoryRepo *CategoryRepository) *ProductRepository {
	return &ProductRepository{
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"

	var args []interface{}

	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	query += " ORDER BY id asc"

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	exists, err := repo.categoryRepo.Exists(product.CategoryID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Category not found")
	}

	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err = repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetById(id int) (*models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, c.name as category_name 
			  FROM products p 
			  LEFT JOIN categories c ON p.category_id = c.id 
			  WHERE p.id = $1`

	var p models.Product
	var categoryName sql.NullString
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryName)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}
	if err != nil {
		return nil, err
	}
	if categoryName.Valid {
		p.CategoryName = categoryName.String
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	exists, err := repo.categoryRepo.Exists(product.CategoryID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Category not found")
	}

	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"

	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
	}

	return err
}
