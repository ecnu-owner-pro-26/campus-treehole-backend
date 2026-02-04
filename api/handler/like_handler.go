package handler

import (
	"net/http"
	"strconv"

	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/utils"
	"github.com/gin-gonic/gin"
)

// LikeHandler 点赞处理器
type LikeHandler struct {
	likeService *service.LikeService
}

// NewLikeHandler 创建点赞处理器实例
func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// ToggleLike 切换点赞状态（翻转点赞）
// @Summary 切换点赞状态
// @Description 翻转点赞状态，如果已点赞则取消，如果未点赞则点赞
// @Tags 点赞
// @Accept json
// @Produce json
// @Param request body dto.ToggleLikeRequest true "切换点赞请求"
// @Success 200 {object} utils.Response{data=dto.ToggleLikeResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/likes/toggle [post]
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("用户未登录"))
		return
	}

	// 绑定请求参数
	var req dto.ToggleLikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("参数错误: "+err.Error()))
		return
	}

	// 验证目标类型
	if req.TargetType != 1 && req.TargetType != 2 {
		c.JSON(http.StatusBadRequest, utils.BadRequest("目标类型错误，1-记忆 2-留言"))
		return
	}

	// 调用服务层切换点赞状态
	isLiked, err := h.likeService.ToggleLike(userID.(int64), req.TargetID, req.TargetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("操作失败: "+err.Error()))
		return
	}

	// 获取最新的点赞状态和数量
	_, likeCount, err := h.likeService.GetLikeStatus(userID.(int64), req.TargetID, req.TargetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取点赞状态失败: "+err.Error()))
		return
	}

	// 返回结果
	response := dto.ToggleLikeResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}

	message := "取消点赞成功"
	if isLiked {
		message = "点赞成功"
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage(message, response))
}

// GetLikeStatus 获取点赞状态
// @Summary 获取点赞状态
// @Description 获取指定目标的点赞状态和点赞数量
// @Tags 点赞
// @Produce json
// @Param target_id query int true "目标ID"
// @Param target_type query int true "目标类型：1-记忆 2-留言"
// @Success 200 {object} utils.Response{data=dto.LikeStatusResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/likes/status [get]
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("用户未登录"))
		return
	}

	// 获取查询参数
	targetIDStr := c.Query("target_id")
	targetTypeStr := c.Query("target_type")

	if targetIDStr == "" || targetTypeStr == "" {
		c.JSON(http.StatusBadRequest, utils.BadRequest("缺少必要参数"))
		return
	}

	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("目标ID格式错误"))
		return
	}

	targetType, err := strconv.ParseInt(targetTypeStr, 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("目标类型格式错误"))
		return
	}

	// 验证目标类型
	if targetType != 1 && targetType != 2 {
		c.JSON(http.StatusBadRequest, utils.BadRequest("目标类型错误，1-记忆 2-留言"))
		return
	}

	// 获取点赞状态
	isLiked, likeCount, err := h.likeService.GetLikeStatus(userID.(int64), targetID, int8(targetType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取点赞状态失败: "+err.Error()))
		return
	}

	// 返回结果
	response := dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}

	c.JSON(http.StatusOK, utils.Success(response))
}

// ToggleMemoryLike 切换记忆点赞状态（便捷接口）
// @Summary 切换记忆点赞状态
// @Description 翻转记忆点赞状态的便捷接口
// @Tags 点赞
// @Produce json
// @Param id path int true "记忆ID"
// @Success 200 {object} utils.Response{data=dto.ToggleLikeResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/memories/{id}/like [post]
func (h *LikeHandler) ToggleMemoryLike(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("用户未登录"))
		return
	}

	// 获取记忆ID
	memoryIDStr := c.Param("id")
	memoryID, err := strconv.ParseInt(memoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("记忆ID格式错误"))
		return
	}

	// 调用服务层切换点赞状态（目标类型为1表示记忆）
	isLiked, err := h.likeService.ToggleLike(userID.(int64), memoryID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("操作失败: "+err.Error()))
		return
	}

	// 获取最新的点赞数量
	_, likeCount, err := h.likeService.GetLikeStatus(userID.(int64), memoryID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取点赞状态失败: "+err.Error()))
		return
	}

	// 返回结果
	response := dto.ToggleLikeResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}

	message := "取消点赞成功"
	if isLiked {
		message = "点赞成功"
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage(message, response))
}

// ToggleCommentLike 切换留言点赞状态（便捷接口）
// @Summary 切换留言点赞状态
// @Description 翻转留言点赞状态的便捷接口
// @Tags 点赞
// @Produce json
// @Param id path int true "留言ID"
// @Success 200 {object} utils.Response{data=dto.ToggleLikeResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id}/like [post]
func (h *LikeHandler) ToggleCommentLike(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Unauthorized("用户未登录"))
		return
	}

	// 获取留言ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequest("留言ID格式错误"))
		return
	}

	// 调用服务层切换点赞状态（目标类型为2表示留言）
	isLiked, err := h.likeService.ToggleLike(userID.(int64), commentID, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("操作失败: "+err.Error()))
		return
	}

	// 获取最新的点赞数量
	_, likeCount, err := h.likeService.GetLikeStatus(userID.(int64), commentID, 2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerError("获取点赞状态失败: "+err.Error()))
		return
	}

	// 返回结果
	response := dto.ToggleLikeResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}

	message := "取消点赞成功"
	if isLiked {
		message = "点赞成功"
	}

	c.JSON(http.StatusOK, utils.SuccessWithMessage(message, response))
}