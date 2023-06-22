package SystemController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func Audit(c *gin.Context) {
	fmt.Println(c.Request.URL.String())
	fmt.Println(c.Request.Header)
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Println(string(body))
}
