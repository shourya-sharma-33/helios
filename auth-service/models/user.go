package models

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
