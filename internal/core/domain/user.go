package domain

type User struct {
	UserID    int64
	Name      string
	Email     string
	Password  string // hashed
	CreatedAt int64
	UpdatedAt int64
}
