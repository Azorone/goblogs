package handler

import (
	"github.com/gin-gonic/gin"
	"mastcat/internal/database"
	"mastcat/internal/model"
	"mastcat/pkg/util"
	"net/http"
	"strconv"
)

func AddCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to create category", nil))
		return
	}
	var id int
	database.DB.Raw("select LAST_INSERT_ID() as id").Pluck("ID", &id)
	ID := make(map[string]string)
	ID["ID"] = strconv.Itoa(id)
	ID["Name"] = category.Name
	c.JSON(http.StatusOK, util.NewResponse(200, "Category added successfully", ID))
}

func UpdateCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	updates := util.StructToMap(category)
	if err := database.DB.Model(&model.Category{}).Where("id = ?", category.ID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to update category", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Category updated successfully", nil))
}

func GetCategories(c *gin.Context) {
	var categories []model.Category
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	database.DB.Offset(offset).Limit(pageSize).Find(&categories)
	c.JSON(http.StatusOK, util.NewResponse(200, "Categories retrieved successfully", categories))
}
