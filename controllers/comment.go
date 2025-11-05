package controllers

import (
	"net/http"
	"strconv"

	"GoBlogService/database"
	"GoBlogService/models"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// CreateComment godoc
// @Summary 创建评论
// @Description 对文章发表评论（需要认证）
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCommentRequest true "评论内容"
// @Success 201 {object} map[string]interface{} "评论成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Router /comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if post exists
	var post models.Post
	if err := database.DB.First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  req.PostID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Preload user data for response
	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetPostComments godoc
// @Summary 获取文章评论
// @Description 获取某篇文章的所有评论
// @Tags 评论
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "成功获取"
// @Router /posts/{id}/comments [get]
func (cc *CommentController) GetPostComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
