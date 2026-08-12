package request

type LoginLoginRequest struct {
    UserId   string `json:"USERID"`  // ← đổi từ "userid"
    Password string `json:"PWD"`     // ← đổi từ "psw"
}