package client

// UserResponse mirrors the server DTO.
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// LoginResponse mirrors dto.LoginResponse.
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// Login authenticates and returns the JWT token.
func (c *Client) Login(email, password string) (*LoginResponse, error) {
	body := map[string]string{"email": email, "password": password}
	return post[*LoginResponse](c, "/api/v1/auth/login", body)
}

// Me returns the currently authenticated user.
func (c *Client) Me() (*UserResponse, error) {
	return get[*UserResponse](c, "/api/v1/auth/me")
}
