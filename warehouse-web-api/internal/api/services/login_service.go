package services

import (
	"errors"
	"fmt"
	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"
)

type LoginService struct {
	*BaseService
}

var LoginServiceInstance = &LoginService{}

func (s *LoginService) Login(Userid, Password string) (*types.Login, error) {
    db, err := database.LYS_ERP_Connection()
    if err != nil {
        fmt.Println("Database connection error:", err)
        return nil, err
    }
    dbInstance, _ := db.DB()
    defer dbInstance.Close()

    query := `    
    SELECT USERID, PWD, USERNAME, Role FROM dbo.Busers 
    WHERE USERID = ?
    `
    var result types.Login

    err = db.Raw(query, Userid).Scan(&result).Error
    if err != nil {
        fmt.Println("Query error:", err)
        return nil, err
    }

    if result.UserId == "" {
        return nil, errors.New("user not found")
    }

    // ✅ Thêm log để xem DB lưu gì
    fmt.Printf("DB Password: [%s]\n", result.Password)
    fmt.Printf("Input Password: [%s]\n", Password)
    fmt.Printf("Match: %v\n", result.Password == Password)

    if result.Password != Password {
        return nil, errors.New("invalid password")
    }

    result.Password = ""
    return &result, nil
}