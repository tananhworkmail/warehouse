package controllers

import (
	"log"
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type User_Login_Controller struct {
	*BaseController
}

var User_Login = &User_Login_Controller{}

func (c *User_Login_Controller) Login(ctx *gin.Context) {
	var requestParams types.User_Login
	if err := ctx.ShouldBindJSON(&requestParams); err != nil {
		log.Printf("Invalid form data: %v", err)
		// Trả 200 nhưng báo lỗi
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid form data")
		return
	}

	result, err := services.User_Login.Login(&requestParams)
	if err != nil {
		// Trả 200 nhưng báo lỗi từ err.Error() (đã là tiếng Anh)
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err.Error())
		return
	}

	// Thành công vẫn trả 200
	response.OkWithData(ctx, result)
}
