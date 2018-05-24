package main

import (
"net/http"
"blog/models"
"fmt"
)


func main() {

	err := models.InitBaseConfig("./conf/blog.ini")
	if err != nil{
		fmt.Println(err)
		return
	}

	router := initRouter()

	http.ListenAndServe(":8080", router)

}
