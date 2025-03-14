package auth

import "github.com/gin-gonic/gin"

func (ah *Handler) getWebsocketToken(c *gin.Context) string {
	token := c.Request.Header.Get("Sec-WebSocket-Protocol")
	return token
}
