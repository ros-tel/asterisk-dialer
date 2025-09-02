package asterisk

import (
	"github.com/gin-gonic/gin"
	"github.com/heltonmarx/goami/ami"
)

func Inject(ami *ami.Socket) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ami", ami)
	}
}
