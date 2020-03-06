package model

type User struct {
	UUID      string `json:"uid" db:"uid" valid:"-"`
	Email     string `json:"email" db:"email" valid:"required,email"`
	Password  string `json:"password" db:"password" valid:"-"`
	FirstName string `json:"first_name" db:"first_name" valid:"-"`
	LastName  string `json:"last_name" db:"last_name" valid:"-"`
	CreatedAt string `json:"-" db:"created_at" valid:"-"`
}

type Login struct {
	Email    string `json:"email" db:"email" valid:"required,email"`
	Password string `json:"password" db:"password" valid:"-"`
}
