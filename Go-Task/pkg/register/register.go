package register

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetRegisterHandler() gin.HandlerFunc { 
	return func(c *gin.Context) { 
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title":"Зарегистрироваться",
		})
	}
}

	
	
