package routes

import (
	"gotry/controllers"
	"gotry/middleware"
	"gotry/ws"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		ws.HandleConnections(c.Writer, c.Request)
	})

	r.POST("/login", controllers.Login)
	r.POST("/sign-up", controllers.SignUp)
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserById)
	r.GET("/posts", controllers.GetPost)
	r.GET("/posts/:id", controllers.GetPostId)
	r.GET("/followers/:id", controllers.GetFollowers)
	r.GET("/following/:id", controllers.GetFollowing)

	userRoutes := r.Group("/users", middleware.AuthMiddleware())
	{
		userRoutes.POST("/logout", controllers.Logout)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
	}

	postsRoutes := r.Group("/posts", middleware.AuthMiddleware())
	{
		postsRoutes.POST("", controllers.CreatePost)
		postsRoutes.PUT("/:id", controllers.UpdatePost)
		postsRoutes.DELETE("/:id", controllers.DeletePost)
	}

	followRoutes := r.Group("/follow", middleware.AuthMiddleware())
	{
		followRoutes.POST("", controllers.FollowUser)
		followRoutes.DELETE("/unfollow", controllers.UnfollowUser)
	}

	msgRoutes := r.Group("/messages", middleware.AuthMiddleware())
	{
		msgRoutes.GET("", controllers.GetMessages)
		msgRoutes.GET("/:id", controllers.GetMessageByID)
		msgRoutes.POST("", controllers.CreateMessage)
		msgRoutes.PUT("/:id", controllers.UpdateMessage)
		msgRoutes.DELETE("/:id", controllers.DeleteMessage)
	}

	return r
}
