package order

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (s *Service) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	rows, err := s.db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	orders := make([]*pb.Order, 0)
	for rows.Next() {
		order, err := scanRowIntoOrdersPB(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return &pb.GetOrdersResponse{Orders: orders}, nil
}
func (s *Service) GetOrdersByIDs(ctx context.Context, ordersIDs *pb.GetOrdersByIDsRequest) (*pb.GetOrdersByIDsResponse, error) {
	placeholders := strings.Repeat(",?", len(ordersIDs.GetIds())-1)
	query := fmt.Sprintf("SELECT * FROM orders WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(ordersIDs.GetIds()))
	for i, v := range ordersIDs.GetIds() {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	orders := []*pb.Order{}
	for rows.Next() {
		order, err := scanRowIntoOrdersPB(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return &pb.GetOrdersByIDsResponse{Orders: orders}, nil
}
func (s *Service) GetOrderByID(ctx context.Context, id *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	rows, err := s.db.Query("SELECT * FROM orders WHERE id = ?", id.GetId())
	if err != nil {
		return nil, err
	}
	order := new(pb.Order)
	for rows.Next() {
		order, err = scanRowIntoOrdersPB(rows)
		if err != nil {
			return nil, err
		}
	}
	if order.GetId() != id.GetId() {
		return nil, fmt.Errorf("order not found")
	}
	return &pb.GetOrderByIDResponse{Order: order}, nil
}

func (s *Service) CreateOrder(ctx context.Context, order *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	res, err := s.db.Exec(
		"INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)",
		order.GetOrder().GetId(), order.GetOrder().GetTotal(), order.GetOrder().GetStatus(), order.GetOrder().GetAddress(),
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{Id: id}, nil
}
func (s *Service) DeleteOrderByID(ctx context.Context, id *pb.DeleteOrderByIDRequest) (*pb.DeleteOrderByIDResponse, error) {
	res, err := s.db.Exec(
		"DELETE from orders WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	deletedID, _ := res.LastInsertId()
	return &pb.DeleteOrderByIDResponse{DeletedCount: deletedID}, nil
}
func (s *Service) DeleteOrder(ctx context.Context, order *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	res, err := s.db.Exec(
		"DELETE from orders WHERE id = ?",
		order.GetOrder().GetId(),
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &pb.DeleteOrderResponse{DeletedCount: id}, nil
}
func (s *Service) UpdateOrder(ctx context.Context, order *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	res, err := s.db.Exec(
		"UPDATE orders SET status = ?, address = ? WHERE id = ?",
		order.GetOrder().GetStatus(), order.GetOrder().GetAddress(), order.GetOrder().GetId(),
	)

	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &pb.UpdateOrderResponse{UpdatedCount: id}, nil
}
func scanRowIntoOrdersPB(rows *sql.Rows) (*pb.Order, error) {
	order := new(pb.Order)
	err := rows.Scan(
		&order.Id,
		&order.UserId,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func scanRowIntoOrderItemsPB(rows *sql.Rows) (*pb.OrderItem, error) {
	orderItems := new(pb.OrderItem)
	err := rows.Scan(
		&orderItems.Id,
		&orderItems.OrderId,
		&orderItems.ProductId,
		&orderItems.Quantity,
		&orderItems.Price,
	)
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}
