package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nadhirfr/codefood/helpers"
	"github.com/nadhirfr/codefood/models"

	"github.com/gin-gonic/gin"
)

// RecipeCreate godoc
// @Summary Register a new recipe
// @Description Register a new recipe
// @Tags recipe
// @Accept  json
// @Produce  json
// @Param name body models.RecipeCreate true "name of recipe"
// @Security Bearer
// @Success 201 {object} models.ResponseResult{result=models.RecipeResult201}
// @Failure 400 {object} models.ResponseError{error=models.RecipeError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 500
// @Router /recipe/ [post]
func RecipeCreate(c *gin.Context) {
	var recipeRegister models.RecipeCreate

	if ok, errors := helpers.DefaultValidator(c, &recipeRegister); !ok {
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

	var recipe = models.Recipe{
		Name:             recipeRegister.Name,
		Image:            recipeRegister.Image,
		RecipeCategoryId: recipeRegister.RecipeCategoryId,
		NServing:         recipeRegister.NServing,
		NReactionLike:    0,
		NReactionNeutral: 0,
		NReactionDislike: 0,
	}

	if err := helpers.DB.Save(&recipe).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		for idx, _ := range recipeRegister.IngredientsPerServing {
			recipeRegister.IngredientsPerServing[idx].RecipeID = recipe.ID
		}

		for idx, _ := range recipeRegister.Steps {
			recipeRegister.Steps[idx].RecipeID = recipe.ID
		}

		helpers.DB.Create(recipeRegister.IngredientsPerServing)
		helpers.DB.Create(recipeRegister.Steps)

		c.JSON(http.StatusCreated, models.ResponseResult{Success: true, Message: "Success", Data: models.RecipeResult201{
			ID:                    recipe.ID,
			Name:                  recipe.Name,
			Image:                 recipe.Image,
			RecipeCategoryId:      recipe.RecipeCategoryId,
			NServing:              recipe.NServing,
			NReactionLike:         recipe.NReactionLike,
			NReactionNeutral:      recipe.NReactionNeutral,
			NReactionDislike:      recipe.NReactionDislike,
			IngredientsPerServing: recipeRegister.IngredientsPerServing,
			Steps:                 recipeRegister.Steps,
			CreatedAt:             recipe.CreatedAt,
			UpdatedAt:             recipe.UpdatedAt,
		}})
	}

}

// RecipeGetByRecipeID godoc
// @Summary Get recipe by recipeID from token
// @Description Get recipe by recipeID from token
// @Tags recipe
// @Accept */*
// @Produce  json
// @Param recipe_id path int true "id recipe to get"
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.RecipeResult200}
// @Failure 401 {object} models.ResponseError{Error=string}
// @Failure 404
// @Failure 500
// @Router /recipe/{recipe_id} [get]
func RecipeGetByRecipeID(c *gin.Context) {
	var recipe_id = c.Param("recipe_id")
	var nServing = c.Query("nServing")
	recipe_id_uint64, _ := strconv.ParseUint(recipe_id, 10, 64)
	nServing_float64, _ := strconv.ParseFloat(nServing, 10)

	var recipe models.Recipe

	err := helpers.DB.Model(&recipe).Where(models.Recipe{ID: uint(recipe_id_uint64)}).First(&recipe).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " not found"})
		return
	} else {
		var ingredientsPerServings []models.RecipeIngridient

		err := helpers.DB.Model(&ingredientsPerServings).Where(models.RecipeIngridient{RecipeID: uint(recipe_id_uint64)}).Find(&ingredientsPerServings).Error
		if err != nil {
			c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " has no ingridients"})
			return
		} else {
			if recipe.NServing != nServing_float64 && nServing_float64 > 0 {
				for idx, _ := range ingredientsPerServings {
					ingredientsPerServings[idx].Value = (nServing_float64 / recipe.NServing) * ingredientsPerServings[idx].Value
				}
				recipe.NServing = nServing_float64
			}

			var recipeCategory models.RecipeCategory
			recipeCategory.ID = recipe.RecipeCategoryId
			helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

			c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.RecipeResult200{
				ID:                    recipe.ID,
				Name:                  recipe.Name,
				Image:                 recipe.Image,
				NReactionLike:         recipe.NReactionLike,
				NReactionNeutral:      recipe.NReactionNeutral,
				NReactionDislike:      recipe.NReactionDislike,
				RecipeCategoryId:      recipe.RecipeCategoryId,
				NServing:              recipe.NServing,
				IngredientsPerServing: ingredientsPerServings,
				CreatedAt:             recipe.CreatedAt,
				UpdatedAt:             recipe.UpdatedAt,
				RecipeCategory:        recipeCategory,
			}})
			return
		}

		// var file models.File
		// helpers.DB.Model(&file).Where("id = ?", recipe.ImageID).First(&file)

		// url, _ := helpers.GetFileFromBucket(file.Url)
		// c.JSON(http.StatusOK, models.ResponseResult{Result: models.RecipeResult200{ID: recipe.ID, Judul: recipe.Judul, Image: url, Text: recipe.Text, UndanganID: recipe.UndanganID}})
	}
}

func RecipeGetAll(c *gin.Context) {
	var skip = c.Query("skip")
	var limit = c.Query("limit")
	var sort = c.Query("sort")
	var q = c.Query("q")
	var categoryId = c.Query("categoryId")

	var recipes []models.Recipe
	var recipesResult []models.RecipeResultGetAll

	query := helpers.DB.Model(&recipes)

	categoryId_uint64, _ := strconv.ParseUint(categoryId, 10, 64)
	if categoryId_uint64 > 0 {
		query.Where("recipe_category_id", categoryId_uint64)
	}

	if q != "" {
		query.Where("name LIKE ?", "%"+q+"%")
	}

	if sort != "" {
		//name_asc|name_desc|like_desc
		s := strings.Split(sort, "_")

		if s[0] != "" {
			if s[0] == "like" {
				s[0] = "n_reaction_like"
			}
		}

		query.Order(s[0] + " " + s[1])

	}

	if limit != "" {
		limit_uint64, _ := strconv.ParseInt(limit, 10, 64)
		query.Limit(int(limit_uint64))
	}

	if skip != "" {
		skip_uint64, _ := strconv.ParseInt(skip, 10, 64)
		query.Offset(int(skip_uint64))
	}

	query.Find(&recipes)

	if query.Error != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe not found"})
		return
	} else {
		for _, recipe := range recipes {
			var recipeCategory models.RecipeCategory
			recipeCategory.ID = recipe.RecipeCategoryId
			helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

			recipesResult = append(recipesResult, models.RecipeResultGetAll{
				ID:               recipe.ID,
				Name:             recipe.Name,
				Image:            recipe.Image,
				NReactionLike:    recipe.NReactionLike,
				NReactionNeutral: recipe.NReactionNeutral,
				NReactionDislike: recipe.NReactionDislike,
				RecipeCategoryId: recipe.RecipeCategoryId,
				CreatedAt:        recipe.CreatedAt,
				UpdatedAt:        recipe.UpdatedAt,
				RecipeCategory:   recipeCategory,
			})
		}

		data := struct {
			Total   int                         `json:"total"`
			Recipes []models.RecipeResultGetAll `json:"recipes"`
		}{
			Total:   len(recipesResult),
			Recipes: recipesResult,
		}

		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: data})
		return
	}

}

func RecipeSearch(c *gin.Context) {
	var limit = c.Query("limit")
	var q = c.Query("q")

	var recipes []models.Recipe
	var recipesResult []models.RecipeResultSearch

	query := helpers.DB.Model(&recipes).Select("id", "name")

	if q != "" && len(q) >= 2 {
		query.Where("name LIKE ?", "%"+q+"%")
	}

	if limit != "" {
		limit_uint64, _ := strconv.ParseInt(limit, 10, 64)
		query.Limit(int(limit_uint64))
	} else {
		query.Limit(5)
	}

	query.Find(&recipesResult)

	if query.Error != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe not found"})
		return
	} else {
		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: recipesResult})
		return
	}

}

// RecipeStepsGetByRecipeID godoc
// @Summary Get recipe by recipeID from token
// @Description Get recipe by recipeID from token
// @Tags recipe
// @Accept */*
// @Produce  json
// @Param recipe_id path int true "id recipe to get"
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.RecipeResult200}
// @Failure 401 {object} models.ResponseError{Error=string}
// @Failure 404
// @Failure 500
// @Router /recipes/{recipe_id}/steps [get]
func RecipeStepsGetByRecipeID(c *gin.Context) {
	var recipe_id = c.Param("recipe_id")
	recipe_id_uint64, _ := strconv.ParseUint(recipe_id, 10, 64)

	if recipe_id_uint64 <= 0 {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " not found"})
		return
	}

	var recipe models.Recipe

	err := helpers.DB.Model(&recipe).Where(models.Recipe{ID: uint(recipe_id_uint64)}).First(&recipe).Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " not found"})
		return
	} else {
		var steps []models.RecipeStep
		err := helpers.DB.Model(&steps).Where(models.RecipeStep{RecipeID: uint(recipe_id_uint64)}).Find(&steps).Error
		if err != nil {
			c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " has no steps"})
			return
		} else {
			c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: steps})
			return
		}
	}
}

// RecipeEditByRecipeID godoc
// @Summary Edit an recipe
// @Description Edit an recipe
// @Tags recipe
// @Accept  json
// @Produce  json
// @Param name body models.RecipeCreate true "name of recipe"
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.RecipeResult201}
// @Failure 400 {object} models.ResponseError{error=models.RecipeError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 404
// @Failure 500
// @Router /recipe/{recipe_id} [post]
func RecipeEditByRecipeID(c *gin.Context) {
	var recipe_id = c.Param("recipe_id")
	recipe_id_uint64, _ := strconv.ParseUint(recipe_id, 10, 64)

	var recipeRegister models.RecipeCreate

	if ok, errors := helpers.DefaultValidator(c, &recipeRegister); !ok {
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

	var _recipe = models.Recipe{
		ID: uint(recipe_id_uint64),
	}

	if err := helpers.DB.Model(_recipe).Where("ID = ?", recipe_id_uint64).First(&_recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " not found"})
		return
	}

	var recipe = models.Recipe{
		ID:               uint(recipe_id_uint64),
		Name:             recipeRegister.Name,
		Image:            recipeRegister.Image,
		RecipeCategoryId: recipeRegister.RecipeCategoryId,
		NServing:         recipeRegister.NServing,
		NReactionLike:    _recipe.NReactionLike,
		NReactionNeutral: _recipe.NReactionNeutral,
		NReactionDislike: _recipe.NReactionDislike,
		CreatedAt:        _recipe.CreatedAt,
	}

	if err := helpers.DB.Save(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
	} else {
		err := helpers.DB.Model(&models.RecipeStep{}).Where(models.RecipeStep{RecipeID: recipe.ID}).Unscoped().Delete(&models.RecipeStep{}).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
		} else {
			err = helpers.DB.Model(&models.RecipeIngridient{}).Where(models.RecipeIngridient{RecipeID: recipe.ID}).Unscoped().Delete(&models.RecipeIngridient{}).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
			} else {
				for idx, _ := range recipeRegister.IngredientsPerServing {
					recipeRegister.IngredientsPerServing[idx].RecipeID = recipe.ID
				}

				for idx, _ := range recipeRegister.Steps {
					recipeRegister.Steps[idx].RecipeID = recipe.ID
				}

				helpers.DB.Create(recipeRegister.IngredientsPerServing)
				helpers.DB.Create(recipeRegister.Steps)

				c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.RecipeResult201{
					ID:                    recipe.ID,
					Name:                  recipe.Name,
					Image:                 recipe.Image,
					RecipeCategoryId:      recipe.RecipeCategoryId,
					NServing:              recipe.NServing,
					NReactionLike:         recipe.NReactionLike,
					NReactionNeutral:      recipe.NReactionNeutral,
					NReactionDislike:      recipe.NReactionDislike,
					IngredientsPerServing: recipeRegister.IngredientsPerServing,
					Steps:                 recipeRegister.Steps,
					CreatedAt:             recipe.CreatedAt,
					UpdatedAt:             recipe.UpdatedAt,
				}})
			}

		}

	}

}

// RecipeDeleteByRecipeID godoc
// @Summary Delete recipe by recipe id
// @Description Delete recipe by recipe id
// @Tags recipe
// @Accept  */*
// @Produce  json
// @Param recipe_id path int true "id recipe to delete"
// @Security Bearer
// @Success 200
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 404
// @Router /recipe/{recipe_id} [delete]
func RecipeDeleteByRecipeID(c *gin.Context) {
	var recipe_id = c.Param("recipe_id")
	recipe_id_uint64, _ := strconv.ParseUint(recipe_id, 10, 64)

	var recipe models.Recipe
	if err := helpers.DB.Model(recipe).Where("ID = ?", recipe_id_uint64).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe_id_uint64) + " not found"})
		return
	}

	//TODO how to edit constraint ON DELETE
	err := helpers.DB.Model(&models.RecipeStep{}).Where(models.RecipeStep{RecipeID: recipe.ID}).Unscoped().Delete(&models.RecipeStep{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
	} else {
		err = helpers.DB.Model(&models.RecipeIngridient{}).Where(models.RecipeIngridient{RecipeID: recipe.ID}).Unscoped().Delete(&models.RecipeIngridient{}).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
		} else {
			if err := helpers.DB.Model(recipe).Delete("ID = ?", recipe_id_uint64).Error; err != nil {
				c.JSON(http.StatusInternalServerError, models.ResponseError{Success: false, Message: "Delete failed " + err.Error()})
			} else {
				c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: struct{}{}})
			}
		}

	}

}
