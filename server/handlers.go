package server

import "github.com/gin-gonic/gin"

func welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Yeah-yeah, I'm up!",
	})
}
