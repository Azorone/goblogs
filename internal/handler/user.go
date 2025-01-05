package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"mastcat/internal/database"
	"mastcat/internal/model"
	"mastcat/pkg/util"
	"net/http"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to hash password", nil))
		return
	}
	user.Password = string(hashedPassword)
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to create user", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "User registered successfully", nil))
}

func Login(c *gin.Context) {
	var user model.User
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponse(400, err.Error(), nil))
		return
	}
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, util.NewResponse(404, "Invalid username or password", nil))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, util.NewResponse(404, "Invalid username or password", nil))
		return
	}
	token, err := util.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponse(500, "Failed to generate token", nil))
		return
	}
	c.JSON(http.StatusOK, util.NewResponse(200, "Login successful", gin.H{"token": token}))
}
