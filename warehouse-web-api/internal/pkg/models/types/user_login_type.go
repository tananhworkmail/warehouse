package types

type User_Login struct {
	GSBH     string `json:"GSBH"`
	USERID   string `json:"USERID" binding:"required"`
	USERNAME string `json:"USERNAME"`
	PWD      string `json:"PWD" binding:"required"`
	SN8      string `json:"SN8"`
}
