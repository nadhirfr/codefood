package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nadhirfr/codefood/helpers"
	"github.com/nadhirfr/codefood/models"
)

// UserGetByUserID godoc
// @Summary Get user by userID from token
// @Description Get user by userID from token
// @Tags user
// @Accept */*
// @Produce  json
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.UserResult201}
// @Failure 401 {object} models.ResponseError{Message=string}
// @Failure 404
// @Failure 500
// @Router /user/detail [get]
func UserGetByUserID(c *gin.Context) {
	var user models.User

	tokenAuth, err := helpers.ExtractTokenMetadata(c.Request)
	if err != nil || tokenAuth.UserRole != helpers.ROLE_PERSONAL {
		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Unauthorized"})
		return
	}

	err = helpers.DB.Where("id = ?", tokenAuth.UserId).First(&user).Error
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		c.JSON(http.StatusOK, models.ResponseResult{Data: models.UserResult201{ID: user.ID, Username: user.Username}})
	}
}

// UserLogin godoc
// @Summary Login user to get auth token
// @Description Login user to get auth token
// @Tags user
// @Accept  json
// @Produce  json
// @Param username body models.UserLogin true "username"
// @Param password body models.UserLogin true "password"
// @Success 200 {object} models.ResponseResult{result=helpers.TokenResult200}
// @Failure 400 {object} models.ResponseError{error=models.UserError400}
// @Failure 406,401,422 {object} models.ResponseError{error=string}
// @Failure 500
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var user models.UserLogin

	if ok, errors := helpers.ValidateUser(c, &user); !ok {
		if _, ok := errors.(gin.H); !ok {
			c.JSON(http.StatusNotAcceptable, models.ResponseError{Success: false, Message: fmt.Sprintf("%v", errors)})
			return
		}

		_error := ""
		for _, val := range errors.(gin.H) {
			_error = _error + val.(string)
		}

		c.JSON(http.StatusBadRequest, models.ResponseError{Success: false, Message: _error})
		return
	}

	var _user models.User
	err := helpers.DB.Where("username = ?", user.Username).First(&_user).Error
	if err != nil {
		log.Print(err)
	}

	var userLoginFaileds []models.UserLoginFailed
	err = helpers.DB.Model(userLoginFaileds).Where(models.UserLoginFailed{UserID: _user.ID}).Find(&userLoginFaileds).Limit(3).Order("created_at desc").Error
	if err == nil && len(userLoginFaileds) >= 3 {
		now := time.Now()
		oneMinuteAgo := now.Add(time.Duration(-1) * time.Minute)
		count := len(userLoginFaileds)
		if oneMinuteAgo.Before(userLoginFaileds[count-1].CreatedAt) {
			c.JSON(http.StatusForbidden, models.ResponseError{Success: false, Message: "Too many invalid login, please wait for 1 minute"})
			return
		}
	}

	//compare the user from the request, with the one we defined:
	if user.Username != _user.Username || helpers.CheckPassword(_user.Password, user.Password) != nil {
		var userLoginFailed models.UserLoginFailed
		userLoginFailed.UserID = _user.ID
		helpers.DB.Save(&userLoginFailed)

		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Invalid username or Password"})
		return
	}

	ts, err := helpers.CreateToken(_user.ID, helpers.ROLE_PERSONAL)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ResponseError{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ResponseResult{
		Success: true,
		Message: "Success",
		Data: struct {
			Token string `json:"token"`
		}{Token: ts.AccessToken}})
}

func UserRegister(c *gin.Context) {
	var user models.User

	if ok, errors := helpers.ValidateUser(c, &user); !ok {
		if _, ok := errors.(gin.H); !ok {
			c.JSON(http.StatusNotAcceptable, models.ResponseError{Success: false, Message: fmt.Sprintf("%v", errors)})
			return
		}

		_error := ""
		for _, val := range errors.(gin.H) {
			_error = _error + val.(string)
		}

		c.JSON(http.StatusBadRequest, models.ResponseError{Success: false, Message: _error})
		return
	}

	var _user models.User
	err := helpers.DB.Where("username = ?", user.Username).First(&_user).Error
	if err != nil {
		log.Print(err)
	}

	if _user.Username != user.Username {
		if err := helpers.DB.Save(&user).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusCreated, models.ResponseResult{Success: true, Message: "Success", Data: models.UserResult201{ID: user.ID, Username: user.Username}})
		}
	} else {
		c.JSON(http.StatusBadRequest, models.ResponseError{Success: false, Message: "username " + fmt.Sprint(user.Username) + " already registered"})
	}

}
