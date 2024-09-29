package order

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fayleenpc/tj-jeans/internal/types"
)

// signature
// GetOrders() ([]types.Order, error)
// GetOrdersByIDs(ordersIDs []int) ([]types.Order, error)
// GetOrderByID(id int) (*types.Order, error)
// CreateOrder(order types.Order) (int64, error)
// DeleteOrderByID(id int) (int64, error)
// DeleteOrder(order types.Order) (int64, error)
// UpdateOrder(order types.Order) (int64, error)
// GetOrderItems() ([]types.OrderItem, error)
// GetOrderItemsByIDs(ordersItemsIDs []int) ([]types.OrderItem, error)
// CreateOrderItem(orderItem types.OrderItem) error
// DeleteOrderItemByID(id int) (int64, error)
// DeleteOrderItem(orderItem types.OrderItem)
// UpdateOrderItem(orderItem types.OrderItem) (int64, error)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetOrders() ([]types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	orders := make([]types.Order, 0)
	for rows.Next() {
		order, err := scanRowIntoOrders(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (s *Store) GetOrdersByIDs(ordersIDs []int) ([]types.Order, error) {
	placeholders := strings.Repeat(",?", len(ordersIDs)-1)
	query := fmt.Sprintf("SELECT * FROM orders WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(ordersIDs))
	for i, v := range ordersIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	orders := []types.Order{}
	for rows.Next() {
		order, err := scanRowIntoOrders(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func (s *Store) GetOrderByID(id int) (*types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	order := new(types.Order)
	for rows.Next() {
		order, err = scanRowIntoOrders(rows)
		if err != nil {
			return nil, err
		}
	}
	if order.ID != id {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *Store) CreateOrder(order types.Order) (int64, error) {
	res, err := s.db.Exec(
		"INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)",
		order.UserID, order.Total, order.Status, order.Address,
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

func (s *Store) DeleteOrderByID(id int) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from orders WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) DeleteOrder(order types.Order) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from orders WHERE id = ?",
		order.ID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateOrder(order types.Order) (int64, error) {
	res, err := s.db.Exec(
		"UPDATE orders SET status = ?, address = ? WHERE id = ?",
		order.Status, order.Address, order.ID,
	)

	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetOrderItems() ([]types.OrderItem, error) {
	rows, err := s.db.Query("SELECT * FROM order_items")
	if err != nil {
		return nil, err
	}
	orderItems := make([]types.OrderItem, 0)
	for rows.Next() {
		orderItem, err := scanRowIntoOrderItems(rows)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}
	return orderItems, nil
}

func (s *Store) GetOrderItemsByIDs(ordersItemsIDs []int) ([]types.OrderItem, error) {
	placeholders := strings.Repeat(",?", len(ordersItemsIDs)-1)
	query := fmt.Sprintf("SELECT * FROM order_items WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(ordersItemsIDs))
	for i, v := range ordersItemsIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	orderItems := []types.OrderItem{}
	for rows.Next() {
		orderItem, err := scanRowIntoOrderItems(rows)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}

	return orderItems, nil
}

func (s *Store) GetOrderItemsByID(id int) (*types.OrderItem, error) {
	rows, err := s.db.Query("SELECT * FROM order_items WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	orderItem := new(types.OrderItem)
	for rows.Next() {
		orderItem, err = scanRowIntoOrderItems(rows)
		if err != nil {
			return nil, err
		}
	}
	if orderItem.ID != id {
		return nil, fmt.Errorf("order not found")
	}
	return orderItem, nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, qty, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

func (s *Store) DeleteOrderItemByID(id int) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from order_items WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) DeleteOrderItem(orderItem types.OrderItem) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from order_items WHERE id = ?",
		orderItem.ID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateOrderItem(orderItem types.OrderItem) (int64, error) {
	res, err := s.db.Exec(
		"UPDATE order_items SET qty = ?, price = ? WHERE id = ?",
		orderItem.Quantity, orderItem.Price, orderItem.ID,
	)

	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func scanRowIntoOrders(rows *sql.Rows) (*types.Order, error) {
	order := new(types.Order)
	err := rows.Scan(
		&order.ID,
		&order.UserID,
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

func scanRowIntoOrderItems(rows *sql.Rows) (*types.OrderItem, error) {
	orderItems := new(types.OrderItem)
	err := rows.Scan(
		&orderItems.ID,
		&orderItems.OrderID,
		&orderItems.ProductID,
		&orderItems.Quantity,
		&orderItems.Price,
	)
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}
