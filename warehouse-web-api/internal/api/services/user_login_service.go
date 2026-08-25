package services

import (
	"errors"
	"web-api/internal/api/utils"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type User_Login_Service struct {
	*BaseService
}

var User_Login = &User_Login_Service{}

// Login validates the selected warehouse and credentials, then creates an N8d JWT.
func (s *User_Login_Service) Login(req *types.User_Login) (token string, err error) {
	// Panic recovery
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("login failed due to unexpected error")
		}
	}()

	// Connect to DB
	db, errConn := database.LYS_ERP_Connection()
	if errConn != nil {
		return "", errConn
	}
	sqlDB, errDB := db.DB()
	if errDB != nil {
		return "", errDB
	}
	defer sqlDB.Close()

	// Login only validates the account in Busers. GSBH is contextual data,
	// not a login requirement.
	const queryLogin = `
		SELECT USERID, USERNAME, PWD
		FROM Busers
		WHERE USERID = ?
	`
	var user types.User_Login
	if errScan := db.Raw(queryLogin, req.USERID).Scan(&user).Error; errScan != nil {
		return "", errScan
	}

	// Check if account exists
	if user.USERID == "" {
		return "", errors.New("Account does not exist")
	}

	if user.PWD != req.PWD {
		return "", errors.New("Sai mật khẩu")
	}

	// GSBH is retained in the token for module context, but is not validated.
	user.GSBH = req.GSBH

	// Permission lookup is separate from authentication. A missing BLimits row
	// must not prevent a valid Busers account from logging in.
	var permission struct {
		SN8 string `gorm:"column:SN8"`
	}
	db.Raw(`
		SELECT TOP 1 '1' AS SN8
		FROM BLimits
		WHERE USERID = ?
		  AND FORMID = 'SN8'
		  AND (? = '' OR GSBH = ?)
	`, req.USERID, req.GSBH, req.GSBH).Scan(&permission)
	user.SN8 = permission.SN8

	// Generate JWT
	token, errGen := utils.GenerateJWT(user.USERID, user.GSBH, user.USERNAME, user.SN8)
	if errGen != nil {
		return "", errGen
	}
	return token, nil
}
