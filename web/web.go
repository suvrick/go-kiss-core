package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {

	router := gin.Default()
	router.Static("/", "./web/ui")

	router.Run(":8080")
}

func htmlHandel(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "web/index.html")
}
