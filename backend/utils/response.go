package utils

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码：0表示成功，其他表示失败
	Message string      `json:"message"` // 提示消息
	Data    interface{} `json:"data,omitempty"`
}

// VersionConflictResponse 版本冲突响应结构
type VersionConflictResponse struct {
	Code            int         `json:"code"`
	Message         string      `json:"message"`
	CurrentVersion  int         `json:"current_version"`
	ProvidedVersion int         `json:"provided_version"`
	LatestData      interface{} `json:"latest_data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应（通用）
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 错误请求
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

// NotFound 404 未找到
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// Conflict 409 冲突（用于乐观锁冲突）
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
	})
}

// VersionConflict 版本冲突响应（包含最新数据）
func VersionConflict(c *gin.Context, conflictErr *services.VersionConflictError) {
	c.JSON(http.StatusConflict, VersionConflictResponse{
		Code:            http.StatusConflict,
		Message:         conflictErr.Message,
		CurrentVersion:  conflictErr.CurrentVersion,
		ProvidedVersion: conflictErr.ProvidedVersion,
		LatestData:      conflictErr.LatestData,
	})
}

// InternalServerError 500 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

// HandleServiceError 统一处理 Service 层错误
// 根据错误类型自动返回合适的 HTTP 状态码
func HandleServiceError(c *gin.Context, err error) {
	// 处理版本冲突错误
	if conflictErr, ok := err.(*services.VersionConflictError); ok {
		VersionConflict(c, conflictErr)
		return
	}

	// 处理其他错误
	errMsg := err.Error()

	// 根据错误消息判断错误类型
	switch {
	case contains(errMsg, "not found"):
		NotFound(c, errMsg)
	case contains(errMsg, "invalid"), contains(errMsg, "required"), contains(errMsg, "cannot exceed"):
		BadRequest(c, errMsg)
	case contains(errMsg, "conflict"):
		Conflict(c, errMsg)
	default:
		InternalServerError(c, errMsg)
	}
}

// contains 判断字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
