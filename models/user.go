package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Hash     string `json:"-"`
	Email    string `json:"email"`
	UID      string `json:"uid"`
}

type NewUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required"`
}
