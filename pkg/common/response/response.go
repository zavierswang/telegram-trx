package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var (
	SERVER_CODE         = 40001
	MYSQL_CODE          = 50001
	RPC_SERVER_CODE     = 30001
	TOKEN_INVALID_CODE  = 20001
	QUEREY_INVALID_CODE = 10001
	FORM_INVALID_CODE   = 10002
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"success",
	})
}

func Error(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusOK, Response{
		code,
		nil,
		msg,
	})
	ctx.Abort()
}
