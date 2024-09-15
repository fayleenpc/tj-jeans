package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	UpdateVerifiedUserByEmail(string) error
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductsByIDs([]int) ([]Product, error)
	CreateProduct(Product) (int64, error)
	UpdateProduct(Product) (int64, error)
	DeleteProduct(Product) (int64, error)
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"qty"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Merchant    string    `json:"merchant"`
	Category    string    `json:"category"`
	Currency    string    `json:"currency"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"qty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"qty"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Verified  bool      `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
