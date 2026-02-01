package auth

type Claims struct {
	UserId string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

type User struct {
}
