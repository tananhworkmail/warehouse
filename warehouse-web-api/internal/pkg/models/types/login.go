package types

type Login struct {
	UserId   string `gorm:"column:USERID"`
	UserName string `gorm:"column:USERNAME"`
	Role     string `gorm:"column:Role"`
	Password string `gorm:"column:PWD"` // Không nên trả về mật khẩu
}
