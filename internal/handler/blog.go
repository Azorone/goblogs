package handler

import (
	"github.com/gin-gonic/gin"
	"mastcat/internal/database"
	"mastcat/internal/model"
	"mastcat/pkg/util"
	"net/http"
	"strconv"
)

func AddBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	if err := database.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to create blog", nil))
		return
	}
	var id int
	database.DB.Raw("select LAST_INSERT_ID() as id").Pluck("ID", &id)
	ID := make(map[string]int)
	ID["ID"] = id
	c.JSON(http.StatusOK, util.NewResponse(200, "Blog added successfully", ID))
}

func UpdateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	updates := util.StructToMap(blog)
	if err := database.DB.Model(&model.Blog{}).Where("id = ?", blog.ID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to update blog", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Blog updated successfully", nil))
}

func DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	var blog model.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, util.NewResponse(404, "Blog not found", nil))
		return
	}
	blog.Deleted = true
	if err := database.DB.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to delete blog", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Blog deleted successfully", nil))
}

func GetBlogs(c *gin.Context) {
	var blogs []model.Blog
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	database.DB.Select([]string{"Title", "des", "created_at", "updated_at", "id"}).Where("is_Push = ? and deleted = ?", true, false).Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&blogs)
	c.JSON(http.StatusOK, util.NewResponse(200, "Blogs retrieved successfully", blogs))
}

func GetBlogManger(c *gin.Context) {
	var blogs []model.Blog
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize
	database.DB.Select("Title,is_push,id,Created_at,Updated_at,des").Offset(offset).Limit(pageSize).Find(&blogs)
	c.JSON(http.StatusOK, util.NewResponse(200, "Blogs retrieved successfully", blogs))
}

func GetBlogsByCategory(c *gin.Context) {
	var blogs []model.Blog
	categoryID := c.Param("categoryID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	if err := database.DB.Select([]string{"Title", "des", "created_at", "updated_at", "id"}).Where("is_Push = ? and deleted = ? and category_id = ?", true, false, categoryID).Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&blogs).Error; err != nil {
		c.JSON(http.StatusNotFound, util.NewResponse(404, "No blogs found for this category", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Blogs retrieved successfully", blogs))
}

func GetBlogByID(c *gin.Context) {
	var blog model.Blog
	blogID := c.Param("id")
	if err := database.DB.First(&blog, blogID).Error; err != nil {
		c.JSON(http.StatusNotFound, util.NewResponse(404, "Blog not found", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Blog retrieved successfully", blog))
}
