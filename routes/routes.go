package routes

import (
	"gotry/controllers"
	"gotry/middleware"
	"gotry/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ganti "*" dengan domain FE kalau sudah production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.GET("/ws", func(c *gin.Context) {
		ws.HandleConnections(c.Writer, c.Request)
	})

	r.POST("/login", controllers.Login)
	r.POST("/sign-up", controllers.SignUp)
	r.GET("/users", controllers.GetUsers)
	r.GET("/posts", controllers.GetPost)
	r.GET("/posts/:id", controllers.GetPostId)

	userRoutes := r.Group("/users", middleware.AuthMiddleware())
	{
		userRoutes.GET("/:id", controllers.GetUserById)
		userRoutes.POST("/logout", controllers.Logout)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.POST("/find", controllers.GetUserByUsername)
		userRoutes.GET("/me", controllers.Me)
	}

	postsRoutes := r.Group("/posts", middleware.AuthMiddleware())
	{
		postsRoutes.POST("", controllers.CreatePost)
		postsRoutes.PUT("/:id", controllers.UpdatePost)
		postsRoutes.DELETE("/:id", controllers.DeletePost)
	}

	followRoutes := r.Group("/follow", middleware.AuthMiddleware())
	{
		followRoutes.GET("/followers", controllers.GetFollowers)
		followRoutes.GET("/following", controllers.GetFollowing)
		followRoutes.POST("", controllers.FollowUser)
		followRoutes.DELETE("/unfollow", controllers.UnfollowUser)
	}

	msgRoutes := r.Group("/messages", middleware.AuthMiddleware())
	{
		msgRoutes.GET("", controllers.GetMessages)
		msgRoutes.GET("/:id", controllers.GetMessageByID)
		msgRoutes.GET("user/:id", controllers.GetMessageUser)
		msgRoutes.POST("", controllers.CreateMessage)
		msgRoutes.PUT("/:id", controllers.UpdateMessage)
		msgRoutes.DELETE("/:id", controllers.DeleteMessage)
	}

	return r
}
