package main

import (
	"fmt"
	"main/database"
	"main/functions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})

	database.ConnectDatabase()

	route.POST("/login", functions.LoginUser)
	route.POST("/send", functions.AddNewTransaction)
	route.GET("/users", functions.RetrieveAllAccount)
	route.GET("/transactions/:userid", functions.GetTransacctionPerAccount)
	route.PATCH("/withdraw/:accountno", functions.WithdrawBalance)

	err := route.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
