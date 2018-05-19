package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/lifei", lifei)

	http.ListenAndServe(":8080", router)
}

func lifei(context *gin.Context) {
	
	context.String(http.StatusOK, "hahahah")
}
