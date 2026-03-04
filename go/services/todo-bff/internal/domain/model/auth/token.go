package auth

type Tokens struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type AccessToken struct {
	Token   string
	Expires int64
}

type RefreshToken struct {
	Token   string
	Expires int64
}
