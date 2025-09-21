package dto

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID int64
	Name   string
	Email  string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type GetMeOutput struct {
	UserID int64
	Name   string
	Email  string
}
