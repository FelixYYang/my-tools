// gin框架的一些工具函数

package ginx

import (
	"github.com/FelixYYang/my-tools/pkg/errorx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSON(c *gin.Context, data any) {
	res := gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	}
	c.JSON(http.StatusOK, res)
}

func JSONFail(c *gin.Context, msg string) {
	JSONFailCode(c, msg, errorx.ReqErr)
}

func JSONFailCode(c *gin.Context, msg string, code int) {
	res := gin.H{
		"code": code,
		"msg":  msg,
	}
	c.JSON(http.StatusOK, res)
}

func JSONErr(c *gin.Context, err error) {
	var res gin.H
	if appErr, ok := err.(errorx.AppError); ok {
		res = gin.H{
			"code": appErr.Code(),
			"msg":  appErr.Msg(),
		}
		if appErr.Code() < 0 {
			_ = c.Error(err)
		}
	} else {
		res = gin.H{
			"code": errorx.UNKNOWN,
			"msg":  "未知错误",
		}
		_ = c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}
