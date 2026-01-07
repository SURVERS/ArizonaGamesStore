package handlers

import (
	"arizonagamesstore/backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdCount struct {
	CategoryName string `form:"CategoryName" binding:"required"`
}

// GetAdCount godoc
// @Summary Счетчик объявлений
// @Description Возвращает количество объявлений в категории. Нужно для отображения "Всего объявлений: 420"
// @Tags Объявления
// @Produce json
// @Param CategoryName query string true "Название категории"
// @Success 200 {object} map[string]int "Количество объявлений"
// @Failure 500 {object} map[string]string "Ошибка подсчета"
// @Router /getadcount [get]
func GetAdCount(c *gin.Context) {
	var req AdCount
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := services.GetAdCounts(req.CategoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при получении статистики: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}
