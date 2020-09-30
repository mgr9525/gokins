package utilService

import "github.com/gin-gonic/gin"

func MidNeedLogin(c *gin.Context) {
	lgusr := CurrUser(c)
	if lgusr == nil {
		c.String(403, "not login")
		c.Abort()
		return
	}
	c.Set("lguser", lgusr)
}
