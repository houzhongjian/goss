package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"goss.io/goss/lib/ini"
)

type Admin struct {
	WebPort string
}

func NewAdmin() *Admin {
	return &Admin{
		WebPort: fmt.Sprintf(":%d", ini.GetInt("node_web_port")),
	}
}

//Start .
func (this *Admin) Start() {
	r := gin.Default()
	r.Static("/img", "./admin/static/img/")
	r.Static("/css", "./admin/static/css/")
	r.LoadHTMLGlob("./admin/views/*")

	r.GET("/login", this.handleLogin)
	r.POST("/login", this.handleLogin)
	if err := r.Run(this.WebPort); err != nil {
		log.Panicln(err)
	}
}

//handleLogin .
func (this *Admin) handleLogin(c *gin.Context) {
	if c.Request.Method == "POST" {
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}
