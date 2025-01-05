package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"mastcat/internal/database"
	"mastcat/pkg/util"
	"net/http"
	"os"
	"path/filepath"
)

type Image struct {
	ID   string `gorm:"type:varchar(255);primary_key"`
	Path string `gorm:"type:varchar(255)"`
}

func (image *Image) BeforeCreate(tx *gorm.DB) (err error) {
	image.ID = uuid.New().String()
	return
}

var imageFolderPath = "/absolute/path/to/images"

func init_() {
	if _, err := os.Stat(imageFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(imageFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal("Failed to create image folder: ", err)
		}
	}
}

func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, "Failed to get form data", nil))
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, "No images uploaded", nil))
		return
	}

	for _, file := range files {
		image := Image{}
		filePath := filepath.Join(imageFolderPath, image.ID+filepath.Ext(file.Filename))
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to save image", nil))
			return
		}
		image.Path = filePath
		if err := database.DB.Create(&image).Error; err != nil {
			c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to save image record", nil))
			return
		}
	}

	c.JSON(http.StatusOK, util.NewResponse(200, "Images uploaded successfully", nil))
}

func GetImage(c *gin.Context) {
	imageID := c.Param("id")
	var image Image
	if err := database.DB.First(&image, "id = ?", imageID).Error; err != nil {
		c.JSON(http.StatusNotFound, util.NewResponse(404, "Image not found", nil))
		return
	}

	c.File(image.Path)
}
