package products

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fayleenpc/tj-jeans/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) DeleteProduct(product types.Product) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from products WHERE id = ?",
		product.ID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateProduct(product types.Product) (int64, error) {
	res, err := s.db.Exec(
		"UPDATE products SET name = ?, price = ?, image = ?, description = ?, qty = ? WHERE id = ?",
		product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID,
	)

	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) CreateProduct(p types.Product) (int64, error) {
	res, err := s.db.Exec(
		"INSERT INTO products (name, description, merchant , category, currency , image , price, qty) VALUES (?,?,?,?,?,?,?,?)",
		p.Name, p.Description, p.Merchant, p.Category, p.Currency, p.Image, p.Price, p.Quantity,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil

}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Merchant,
		&product.Category,
		&product.Currency,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}
