package dto

type SignUpDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	Token string `json:"token"`
}

type AssignRoleDTO struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}
