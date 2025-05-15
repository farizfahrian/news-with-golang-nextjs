package entity

type UserEntity struct {
	ID       int64  `gorm:"id"`
	Name     string `gorm:"name"`
	Email    string `gorm:"email"`
	Password string `gorm:"password"`
}
