package controllers

import (
	"net/http"
	"strconv"

	"GoBlogService/database"
	"GoBlogService/models"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreatePost godoc
// @Summary 创建文章
// @Description 创建新的博客文章（需要认证）
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "文章内容"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未认证"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /posts [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Preload user data for response
	database.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// GetPosts godoc
// @Summary 获取文章列表
// @Description 获取所有博客文章列表
// @Tags 文章
// @Produce json
// @Success 200 {object} map[string]interface{} "成功获取"
// @Router /posts [get]
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	if err := database.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetPost godoc
// @Summary 获取单篇文章
// @Description 根据ID获取单篇博客文章详情
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "成功获取"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Router /posts/{id} [get]
func (pc *PostController) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (pc *PostController) UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check ownership
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	database.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

func (pc *PostController) DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check ownership
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
