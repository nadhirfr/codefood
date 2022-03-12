package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nadhirfr/codefood/helpers"
	"github.com/nadhirfr/codefood/models"

	"github.com/gin-gonic/gin"
)

// RecipeCategoryCreate godoc
// @Summary Register a new recipeCategory
// @Description Register a new recipeCategory
// @Tags recipeCategory
// @Accept  json
// @Produce  json
// @Param name body models.RecipeCategoryCreate true "name of category"
// @Success 201 {object} models.ResponseResult{result=models.RecipeCategoryResult201}
// @Failure 400 {object} models.ResponseError{error=models.RecipeCategoryError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 500
// @Router /recipeCategory/ [post]
func RecipeCategoryCreate(c *gin.Context) {
	var recipeCategoryRegister models.RecipeCategoryCreate

	if ok, errors := helpers.DefaultValidator(c, &recipeCategoryRegister); !ok {
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

	var recipeCategory = models.RecipeCategory{
		Name: recipeCategoryRegister.Name,
	}

	if err := helpers.DB.Save(&recipeCategory).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusCreated, models.ResponseResult{Success: true, Message: "Success", Data: models.RecipeCategoryResult201{
			ID:        recipeCategory.ID,
			Name:      recipeCategory.Name,
			CreatedAt: recipeCategory.CreatedAt,
			UpdatedAt: recipeCategory.UpdatedAt,
		}})
	}

}

func RecipeCategoryGetAll(c *gin.Context) {
	var recipeCategories []models.RecipeCategory

	err := helpers.DB.Model(&recipeCategories).Find(&recipeCategories).Error
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: recipeCategories})
	}

}

// RecipeCategoryEditByRecipeCategoryID godoc
// @Summary Edit an recipeCategory
// @Description Edit an recipeCategory
// @Tags recipeCategory
// @Accept  json
// @Produce  json
// @Param name body models.RecipeCategoryCreate true "name of category"
// @Param recipeCategory_id path integer true "recipeCategory_id"
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.RecipeCategoryResult201}
// @Failure 400 {object} models.ResponseError{error=models.RecipeCategoryError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 404
// @Failure 500
// @Router /recipeCategory/{recipeCategory_id} [post]
func RecipeCategoryEditByRecipeCategoryID(c *gin.Context) {
	var recipeCategory_id = c.Param("recipeCategory_id")
	recipeCategory_id_uint64, _ := strconv.ParseUint(recipeCategory_id, 10, 64)

	var recipeCategoryRegister models.RecipeCategoryCreate

	if ok, errors := helpers.DefaultValidator(c, &recipeCategoryRegister); !ok {
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

	var _recipeCategory = models.RecipeCategory{
		ID: uint(recipeCategory_id_uint64),
	}

	if err := helpers.DB.Model(_recipeCategory).Where("ID = ?", recipeCategory_id_uint64).First(&_recipeCategory).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe Category with id " + fmt.Sprint(recipeCategory_id_uint64) + " not found"})
		return
	}

	var recipeCategory = models.RecipeCategory{
		ID:        uint(recipeCategory_id_uint64),
		Name:      recipeCategoryRegister.Name,
		CreatedAt: _recipeCategory.CreatedAt,
	}

	if err := helpers.DB.Save(&recipeCategory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
	} else {
		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.RecipeCategoryResult201{
			ID:        recipeCategory.ID,
			Name:      recipeCategory.Name,
			CreatedAt: recipeCategory.CreatedAt,
			UpdatedAt: recipeCategory.UpdatedAt,
		}})

	}

}

// RecipeCategoryDeleteByRecipeCategoryID godoc
// @Summary Delete recipeCategory by recipeCategory id
// @Description Delete recipeCategory by recipeCategory id
// @Tags recipeCategory
// @Accept  */*
// @Produce  json
// @Param recipeCategory_id path int true "id recipeCategory to delete"
// @Security Bearer
// @Success 200
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 404
// @Router /recipeCategory/{recipeCategory_id} [delete]
func RecipeCategoryDeleteByRecipeCategoryID(c *gin.Context) {
	var recipeCategory_id = c.Param("recipeCategory_id")
	recipeCategory_id_uint64, _ := strconv.ParseUint(recipeCategory_id, 10, 64)

	var recipeCategory models.RecipeCategory
	if err := helpers.DB.Model(recipeCategory).Where("ID = ?", recipeCategory_id_uint64).First(&recipeCategory).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe Category with id " + fmt.Sprint(recipeCategory_id_uint64) + " not found"})
		return
	}

	if err := helpers.DB.Model(recipeCategory).Delete("ID = ?", recipeCategory_id_uint64).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
	} else {
		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: struct{}{}})
	}

}
