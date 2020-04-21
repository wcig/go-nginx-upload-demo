package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		path := c.PostForm(".path")
		name := c.PostForm(".name")

		md5 := c.PostForm(".md5")
		size := c.PostForm(".size")

		fmt.Printf("file path:%s, name:%s, md5:%s, size:%s\n", path, name, md5, size)
		c.JSON(200, result{0, nil, "ok"})
	})
	router.Run(":8080")
}

type result struct {
	Code int
	Data interface{}
	Msg string
}
