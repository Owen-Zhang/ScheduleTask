package api

import "github.com/gin-gonic/gin"

//服务端和客户端的心跳服务
func (this *ApiServer) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}
