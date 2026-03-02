package input

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegister struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type TokenRefresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenVerify struct {
	AccessToken string `json:"access_token" validate:"required"`
}
