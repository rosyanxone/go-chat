package dto

type UserDataResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	NIK         *string `json:"nik"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Role        string  `json:"role"`
}

type UserContactRow struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
