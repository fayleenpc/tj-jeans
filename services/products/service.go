package products

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]*pb.Product, 0)
	for rows.Next() {
		productPB, err := scanRowIntoProductPB(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, productPB)
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
func (s *Service) GetProductsByIDs(ctx context.Context, productIDs *pb.GetProductsByIDsRequest) (*pb.GetProductsByIDsResponse, error) {
	placeholders := strings.Repeat(",?", len(productIDs.GetIds())-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIDs.GetIds()))
	for i, v := range productIDs.GetIds() {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []*pb.Product{}
	for rows.Next() {
		productPB, err := scanRowIntoProductPB(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, productPB)
	}

	return &pb.GetProductsByIDsResponse{Products: products}, nil
}
func (s *Service) GetProductByID(ctx context.Context, id *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERe id = ?", id.GetId())
	if err != nil {
		return nil, err
	}
	productPB := new(pb.Product)
	for rows.Next() {
		productPB, err = scanRowIntoProductPB(rows)
		if err != nil {
			return nil, err
		}

	}
	return &pb.GetProductByIDResponse{Product: productPB}, nil
}

func (s *Service) CreateProduct(ctx context.Context, p *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	res, err := s.db.Exec(
		"INSERT INTO products (name, description, merchant , category, currency , image , price, qty) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		p.GetProduct().GetName(), p.GetProduct().GetDescription(), p.Product.GetMerchant(), p.GetProduct().GetCategory(), p.GetProduct().GetCurrency(), p.GetProduct().GetImage(), p.GetProduct().GetPrice(), p.GetProduct().GetQuantity(),
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Id: id}, nil
}
func (s *Service) DeleteProductByID(ctx context.Context, id *pb.DeleteProductByIDRequest) (*pb.DeleteProductByIDResponse, error) {
	res, err := s.db.Exec(
		"DELETE from products WHERE id = ?",
		id.GetId(),
	)
	if err != nil {
		return nil, err
	}
	deletedProductID, _ := res.LastInsertId()
	return &pb.DeleteProductByIDResponse{DeletedCount: deletedProductID}, nil
}
func (s *Service) DeleteProduct(ctx context.Context, product *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	res, err := s.db.Exec(
		"DELETE from products WHERE id = ?",
		product.GetProduct().GetId(),
	)
	if err != nil {
		return nil, err
	}
	deletedProductID, _ := res.LastInsertId()
	return &pb.DeleteProductResponse{DeletedCount: deletedProductID}, nil
}
func (s *Service) UpdateProduct(ctx context.Context, product *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	res, err := s.db.Exec(
		"UPDATE products SET name = ?, description = ?, merchant = ?, category = ?, currency = ?, image = ?, price = ?, qty = ? WHERE id = ?",
		product.GetProduct().GetName(), product.GetProduct().GetDescription(), product.GetProduct().GetMerchant(), product.GetProduct().GetCategory(), product.GetProduct().GetCurrency(), product.GetProduct().GetImage(), product.GetProduct().GetPrice(), product.GetProduct().GetQuantity(), product.GetProduct().GetId(),
	)

	if err != nil {
		return nil, err
	}
	updatedProductID, _ := res.LastInsertId()
	return &pb.UpdateProductResponse{UpdatedCount: updatedProductID}, nil
}

//	{
//		"error": "sql: Scan error on column index 9, name \"createdAt\":
//		converting driver.Value type time.Time (\"2024-09-06 15:00:45 +0000 UTC\")
//		to a int64: invalid syntax"
//	}
func scanRowIntoProductPB(rows *sql.Rows) (*pb.Product, error) {
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
	productPB := new(pb.Product)
	productPB.Category = product.Category
	productPB.CreatedAt = product.CreatedAt.Unix()
	productPB.Currency = product.Currency
	productPB.Description = product.Description
	productPB.Id = int32(product.ID)
	productPB.Image = product.Image
	productPB.Merchant = product.Merchant
	productPB.Name = product.Name
	productPB.Price = product.Price
	productPB.Quantity = int32(product.Quantity)
	return productPB, nil
}
