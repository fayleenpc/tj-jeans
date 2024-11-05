package types

import (
	"context"
	"time"

	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
)

type FinanceReport struct {
	OrderCount     int `json:"order_count"`
	TotalItemsSold int `json:"total_items_sold"`
	TotalRevenue   int `json:"total_revenue"`
}

type UserStore interface {
	GetUsers() ([]User, error)
	GetUsersByIDs([]int) ([]User, error)
	UpdateVerifiedUserByEmail(string) error
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*User, error)
	DeleteUserByID(int) (int64, error)
	DeleteUser(User) (int64, error)
	UpdateUser(User) (int64, error)
	CreateUser(User) error
}

type UserService interface {
	GetUsers(context.Context, *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	GetUsersByIDs(context.Context, *pb.GetUsersByIDsRequest) (*pb.GetUsersByIDsResponse, error)
	UpdateVerifiedUserByEmail(context.Context, *pb.UpdateVerifiedUserByEmailRequest) error
	GetUserByEmail(context.Context, *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error)
	GetUserByID(context.Context, *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error)
	DeleteUserByID(context.Context, *pb.DeleteUserByIDRequest) (*pb.DeleteUserByIDResponse, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductsByIDs([]int) ([]Product, error)
	GetProductByID(int) (*Product, error)
	CreateProduct(Product) (int64, error)
	DeleteProductByID(int) (int64, error)
	DeleteProduct(Product) (int64, error)
	UpdateProduct(Product) (int64, error)
}

type ProductService interface {
	GetProducts(context.Context, *pb.GetProductsRequest) (*pb.GetProductsResponse, error)
	GetProductsByIDs(context.Context, *pb.GetProductsByIDsRequest) (*pb.GetProductsByIDsResponse, error)
	GetProductByID(context.Context, *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error)
	CreateProduct(context.Context, *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	DeleteProductByID(context.Context, *pb.DeleteProductByIDRequest) (*pb.DeleteProductByIDResponse, error)
	DeleteProduct(context.Context, *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)
	UpdateProduct(context.Context, *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error)
}

type OrderStore interface {
	GetOrders() ([]Order, error)
	GetOrdersByIDs([]int) ([]Order, error)
	GetOrderByID(int) (*Order, error)
	CreateOrder(Order) (int64, error)
	DeleteOrderByID(int) (int64, error)
	DeleteOrder(Order) (int64, error)
	UpdateOrder(Order) (int64, error)
	GetOrderItems() ([]OrderItem, error)
	GetOrderItemsByIDs([]int) ([]OrderItem, error)
	GetOrderItemsByID(int) (*OrderItem, error)
	CreateOrderItem(OrderItem) error
	DeleteOrderItemByID(int) (int64, error)
	DeleteOrderItem(OrderItem) (int64, error)
	UpdateOrderItem(OrderItem) (int64, error)
}

type OrderService interface {
	GetOrders(context.Context, *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error)
	GetOrdersByIDs(context.Context, *pb.GetOrdersByIDsRequest) (*pb.GetOrdersByIDsResponse, error)
	GetOrderByID(context.Context, *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error)
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	DeleteOrderByID(context.Context, *pb.DeleteOrderByIDRequest) (*pb.DeleteOrderByIDResponse, error)
	DeleteOrder(context.Context, *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
	UpdateOrder(context.Context, *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error)
	// GetOrderItems() ([]OrderItem, error)
	// GetOrderItemsByIDs([]int) ([]OrderItem, error)
	// GetOrderItemsByID(int) (*OrderItem, error)
	// CreateOrderItem(OrderItem) error
	// DeleteOrderItemByID(int) (int64, error)
	// DeleteOrderItem(OrderItem) (int64, error)
	// UpdateOrderItem(OrderItem) (int64, error)
}

type TokenStore interface {
	GetBlacklistedTokens() ([]Token, error)
	CreateBlacklistTokens(Token) (*Token, error)
	GetBlacklistTokenByString(string) (*Token, error)
}

type TokenService interface {
	GetBlacklistedTokens(context.Context, *pb.GetBlacklistedTokensRequest) (*pb.GetBlacklistedTokensResponse, error)
	CreateBlacklistTokens(context.Context, *pb.CreateBlacklistTokenRequest) (*pb.CreateBlacklistTokenResponse, error)
	GetBlacklistTokenByString(context.Context, *pb.GetBlacklistTokenByStringRequest) (*pb.GetBlacklistTokenByStringResponse, error)
}

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"qty"`
}

type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Verified    bool      `json:"-"`
	Role        string    `json:"-"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"qty"`
	Price     float64 `json:"price"`
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
	CreatedAt   time.Time `json:"created_at"`
}

type ResponseProduct struct {
	CreatedProduct Product `json:"created_product"`
	UpdatedProduct Product `json:"updated_product"`
	DeletedProduct Product `json:"deleted_product"`
}

type RegisterUserPayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=3,max=130"`
	PhoneNumber string `json:"phone_number" validate:"required,min=12,max=12"`
	Address     string `json:"address" validate:"required"`
}

type ResponseRegister struct {
	URL   string `json:"verify_url"`
	Error string `json:"error"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResponseLogin struct {
	AccessToken string `json:"access_token"`
	SecretToken string `json:"secret_token"`
	Error       string `json:"error"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}
type ResponseCart struct {
	Total   float64   `json:"total_price"`
	OrderID int       `json:"order_id"`
	Items   []Product `json:"items"`
}

type InvoicePayload struct {
	Duration time.Duration `json:"duration"`
	Payment  struct {
		Type   string  `json:"payment_type"`
		Amount float64 `json:"amount"`
	} `json:"payment"`
	Customer struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	} `json:"customer"`
	Items []Product `json:"items" validate:"required"`
}

type InvoiceResponse struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Number      string    `json:"number"`
	InvoiceDate time.Time `json:"invoice_date"`
	DueDate     time.Time `json:"due_date"`
	PaidAt      time.Time `json:"paid_at"`
	Items       []struct {
		ID          int       `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DeletedAt   time.Time `json:"deleted_at"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Category    string    `json:"category"`
		Merchant    string    `json:"merchant_name"`
		Currency    string    `json:"currency"`
		UnitPrice   float64   `json:"unit_price"`
		Quantity    float64   `json:"qty"`
	} `json:"items"`
	Payment struct {
		ID            int       `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		DeletedAt     time.Time `json:"deleted_at"`
		Gateway       string    `json:"gateway"`
		Type          string    `json:"payment_type"`
		Token         string    `json:"token"`
		RedirectURL   string    `json:"redirect_url"`
		TransactionID string    `json:"transaction_id"`
	} `json:"payment"`
	BillingAddress struct {
		ID          int       `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DeletedAt   time.Time `json:"deleted_at"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone_number"`
	} `json:"billing_address"`
	SuccessRedirectURL string `json:"success_redirect_url"`
	FailureRedirectURL string `json:"failure_redirect_url"`
	Title              string `json:"title"`
	State              string `json:"state"`
	TransactionValues  struct {
		Currency       string  `json:"currency"`
		Total          float64 `json:"total_amount"`
		SubTotal       float64 `json:"sub_total_amount"`
		Discount       float64 `json:"discount_amount"`
		Tax            float64 `json:"tax_amount"`
		AdminFee       float64 `json:"admin_fee_amount"`
		InstallmentFee float64 `json:"installment_fee_amount"`
	} `json:"transaction_values"`
}

type Token struct {
	ID        int
	Token     string `json:"token"`
	CreatedAt time.Time
}

type RefreshTokenPayload struct {
	AccessToken string `json:"access_token"`
	SecretToken string `json:"secret_token"`
}

type ResponseRefreshToken struct {
	AccessToken string `json:"access_token"`
}
