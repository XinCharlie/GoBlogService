// @title 个人博客系统 API
// @version 1.0
// @description 基于Gin和GORM开发的个人博客系统后端API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9090
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT认证令牌，格式: "Bearer {token}"
package main

import (
	"log"
	"strconv"

	"GoBlogService/config"
	"GoBlogService/controllers"
	"GoBlogService/database"
	"GoBlogService/middleware"

	"github.com/gin-gonic/gin"

	_ "GoBlogService/docs" // 替换为你的模块路径

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.ConnectDB(
		cfg.Database.Host,
		strconv.Itoa(cfg.Database.Port),
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.LoggerMiddleware())

	// Initialize controllers
	authController := controllers.NewAuthController(cfg)
	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()

	// 文章路由分组
	posts := r.Group("/posts")
	{
		posts.GET("", postController.GetPosts)
		posts.GET("/:id", postController.GetPost)
		posts.GET("/:id/comments", commentController.GetPostComments)

		// 需要认证的文章操作
		authPosts := posts.Group("")
		authPosts.Use(middleware.AuthMiddleware(cfg))
		{
			authPosts.POST("", postController.CreatePost)
			authPosts.PUT("/:id", postController.UpdatePost)
			authPosts.DELETE("/:id", postController.DeletePost)
		}
	}

	// 评论路由分组
	comments := r.Group("/comments")
	comments.Use(middleware.AuthMiddleware(cfg))
	{
		comments.POST("", commentController.CreateComment)
	}

	// 认证路由
	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
