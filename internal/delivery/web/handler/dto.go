package handler

type UserDto struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CarDto struct {
	Name string `json:"name"`
	VIN  string `json:"vin"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CartItemDto struct {
	PartID      string `json:"part_id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Price       int64  `json:"price"`
	Quantity    int64  `json:"quantity"`
	DeliveryDay int    `json:"delivery_day"`
	ImageURL    string `json:"image_url,omitempty"`
}
