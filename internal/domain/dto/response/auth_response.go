package response

type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type RegisterResponse struct {
	User *UserResponse `json:"user"`
}
