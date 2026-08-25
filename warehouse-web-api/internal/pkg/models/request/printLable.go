package request

type PrintLableRequest struct {
	IP         string `json:"ip" binding:"required"`         // 添加 IP 字段
	ZPLCommand string `json:"zplCommand" binding:"required"` // 添加 ZPL 指令字段
}
type CheckPrinterRequest struct {
	IP string `json:"ip" binding:"required"` // 添加 IP 字段
}
