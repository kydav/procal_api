package entity

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	IsActive       bool   `json:"is_active"`
	ProfilePicture string `json:"profile_picture"`
}
