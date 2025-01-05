package domain

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt string
	UpdatedAt string
}

type UserDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
