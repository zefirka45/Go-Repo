package main

import (
    "github.com/gin-gonic/gin"
    "Go-Task/pkg/login"
    "Go-Task/pkg/register"
)

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.GET("/login", login.GetLoginHandler())
    r.GET("/register", register.GetRegisterHandler()) 
    r.Run(":8080")
}