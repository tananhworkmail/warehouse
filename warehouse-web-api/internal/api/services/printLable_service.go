package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CheckPrinterExists 檢查印表機是否存在於給定的 IP 地址
func CheckPrinterExists(ip string) (bool, error) {
	url := fmt.Sprintf("http://%s/server", ip)

	// 創建一個 HTTP 客戶端並設置超時
	client := &http.Client{
		Timeout: 3 * time.Second, // 設置超時為 3 秒
	}

	// 發送 GET 請求
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to reach printer at %s: %w", ip, err)
	}
	defer resp.Body.Close()

	// 檢查狀態碼是否為 200 OK
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// PrintZPL 向指定的印表機 IP 發送 ZPL 指令的函數
func PrintZPL(ip string, zplCommand string) error {
	// 檢查印表機是否存在
	exists, err := CheckPrinterExists(ip)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("printer not found at IP: %s", ip)
	}

	// 準備 POST 請求的 URL
	url := fmt.Sprintf("http://%s/pstprnt", ip)
	payload := []byte(zplCommand)

	// 建立 POST 請求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 設置請求的 Content-Type 標頭
	req.Header.Set("Content-Type", "text/plain")

	// 執行請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() // 確保請求結束後關閉響應主體

	// 檢查響應狀態碼是否為 200 OK
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body) // 使用 io.ReadAll
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("error response from printer: %s, body: %s", resp.Status, body)
	}

	return nil
}
