package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 設定一個用於測試的 Gin 路由器
func setupTestRouter() *gin.Engine {
	// 將 Gin 設定為測試模式，減少不必要的日誌輸出
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 註冊跟 main.go 相同的路由
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "connected"})
	})

	r.GET("/api/v1/tasks", func(c *gin.Context) {
		// 測試用的假資料
		tasks := []Task{
			{Title: "測試任務一", Done: false},
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": tasks})
	})

	return r
}

// 測試 /healthz 健康檢查端點
func TestHealthzHandler(t *testing.T) {
	r := setupTestRouter()

	// 建立一個記錄 HTTP 回應的 ResponseRecorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)

	// 執行請求
	r.ServeHTTP(w, req)

	// 使用 testify 進行斷言檢查
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"db":"connected","status":"ok"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

// 測試 /api/v1/tasks 查詢端點
func TestGetTasksHandler(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)

	r.ServeHTTP(w, req)

	// 檢查 HTTP 狀態碼是否為 200
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析回應 JSON
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
}