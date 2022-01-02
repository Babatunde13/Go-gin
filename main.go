package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main () {
	// Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    errorLogs, _ := os.Create("gin.error.log")
    gin.DefaultWriter = io.MultiWriter(f)
    gin.DefaultErrorWriter = io.MultiWriter(errorLogs)
	r := gin.Default()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": recovered,
			"status":  false,
			
		})
	}))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/search", func (c *gin.Context) {
		// panic("something wrong")
		searchQuery := c.DefaultQuery("query", "No query provided")
		c.JSON(200, gin.H{
			"message": "search " + searchQuery,
		})
	})
	usersGroup := r.Group("/users")
	// add auth middleware
	usersGroup.Use(func (c *gin.Context) {
		c.Set("type", "user")
		c.Next()
	})
	userRoutes(usersGroup)
	adminGroup := r.Group("/admin")
	adminGroup.Use(func (c *gin.Context) {
		c.Set("type", "admin")
		c.Next()
	})
	adminRoutes(adminGroup)
	r.GET("/logs", func(c *gin.Context) {
		c.File("./gin.log")
	})
	r.GET("/logs/error", func(c *gin.Context) {
		c.File("./gin.error.log")
	})
	r.Use(func (c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "not found",
		})
	})
	
	if os.Getenv("APP_ENV") == "production" && os.Getenv("APP_PORT") != "" {
		r.Run(":" + os.Getenv("APP_PORT"))
	} else {
		r.Run(":8080")
	}
}


func userRoutes (usersGroup *gin.RouterGroup) {
	usersGroup.POST("/login", func (c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"status":  false,
			})
			return
		}
		if login.Username == "user" && login.Password == "password" {
			c.JSON(http.StatusOK, gin.H{
				"message": "login success",
				"status":  true,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "login failed",
				"status":  false,
			})
		}
	})
	usersGroup.GET("/", func (c *gin.Context) {
		userType, _ := c.Get("type")
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Welcome to the  %v home page", userType),
		})
	})
	usersGroup.GET("/:name", func (c *gin.Context) {
		name := c.Param("name")
		c.JSON(200, gin.H{
			"message": "hello User " + name,
		})
	})
}

func adminRoutes (adminGroup *gin.RouterGroup) {
	adminGroup.POST("/login", func (c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"status":  false,
			})
			return
		}
		if login.Username == "user" && login.Password == "password" {
			c.JSON(http.StatusOK, gin.H{
				"message": "login success",
				"status":  true,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "login failed",
				"status":  false,
			})
		}
	})
	adminGroup.GET("/", func (c *gin.Context) {
		userType, _ := c.Get("type")
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Welcome to the  %v home page", userType),
		})
	})
	adminGroup.GET("/:name", func (c *gin.Context) {
		fmt.Println(c.Get("type"))
		name := c.Param("name")
		c.JSON(200, gin.H{
			"message": "hello Admin " + name,
		})
	})
}
