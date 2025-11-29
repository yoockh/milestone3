package dto

type UserRequest struct {
	Name string `json:"name" validate:"required,gte=3"` 
	Email string `json:"email" validate:"required,email"` 
	Password string `json:"password" validate:"required,gte=8"` 
}

type UserResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	RoleId int `json:"role_id"`
}

type UserLoginRequest struct {
	Email string `json:"email" validate:"required,email"` 
	Password string `json:"password" validate:"required,gte=8"` 	
}