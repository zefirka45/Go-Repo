package login

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)

func GetLoginHandler() gin.HandlerFunc { 
	return func(c *gin.Context) { 
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":"Войти в систему",
			"error": "Пароль или логин указаны не верно",
		})
	}
}

func PostLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context){
		name := c.PostForm("username")
		password := c.PostForm("password")
		if name == "admin" || password == "admin" {
			fmt.Println("Привет")
		} else {
			fmt.Println("err")
		}
	}
}

	
	
