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

// ServeCreate godoc
// @Summary Register a new serve
// @Description Register a new serve
// @Tags serve
// @Accept  json
// @Produce  json
// @Param name body models.RecipeCategoryCreate true "name of recipe"
// @Security Bearer
// @Success 201 {object} models.ResponseResult{result=models.ServeResult201}
// @Failure 400 {object} models.ResponseError{error=models.ServeError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 500
// @Router /serve/ [post]
func ServeCreate(c *gin.Context) {
	var serveRegister models.ServeCreate

	if ok, errors := helpers.ValidateServe(c, &serveRegister); !ok {
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

	tokenAuth, err := helpers.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Unauthorized"})
		return
	}

	var recipe = models.Recipe{
		ID: serveRegister.RecipeID,
	}

	if err := helpers.DB.Model(recipe).Where("ID = ?", serveRegister.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serveRegister.RecipeID) + " not found"})
		return
	}

	var serve = models.Serve{
		NServing: serveRegister.NServing,
		UserID:   uint(tokenAuth.UserId),
		RecipeID: serveRegister.RecipeID,
	}

	if err := helpers.DB.Save(&serve).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		var steps []models.RecipeStep
		err := helpers.DB.Model(&steps).Where(models.RecipeStep{RecipeID: serveRegister.RecipeID}).Find(&steps).Error
		if err != nil {
			c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serveRegister.RecipeID) + " has no steps"})
			return
		} else {
			var serveSteps []models.ServeStep
			var serveStepResults []models.ServeStepResult

			for idx, val := range steps {
				done := false
				if idx == 0 {
					done = true
				}
				serveSteps = append(serveSteps, models.ServeStep{
					ServeID:      serve.ID,
					RecipeStepID: val.ID,
					Done:         done,
				})

				serveStepResults = append(serveStepResults, models.ServeStepResult{
					StepOrder:   val.StepOrder,
					Description: val.Description,
					Done:        done,
				})
			}

			helpers.DB.Create(serveSteps)

			var recipeCategory models.RecipeCategory
			recipeCategory.ID = recipe.RecipeCategoryId
			helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

			c.JSON(http.StatusCreated, models.ResponseResult{Success: true, Message: "Success", Data: models.ServeResult201{
				ID:                 serve.ID,
				UserID:             serve.UserID,
				RecipeID:           serve.RecipeID,
				RecipeName:         recipe.Name,
				RecipeCategoryName: recipeCategory.Name,
				RecipeImage:        recipe.Image,
				RecipeCategoryId:   recipe.RecipeCategoryId,
				NServing:           serve.NServing,
				NStep:              float64(len(serveStepResults)),
				NStepDone:          1,
				Reaction:           serve.Reaction,
				Steps:              serveStepResults,
				Status:             "progress",
				CreatedAt:          serve.CreatedAt,
				UpdatedAt:          serve.UpdatedAt,
			}})
		}

	}

}

// ServeEditByServeID godoc
// @Summary Edit an serve
// @Description Edit an serve
// @Tags serve
// @Accept  json
// @Produce  json
// @Param stepOrder body models.ServeUpdateStep true "stepOrder"
// @Security Bearer
// @Success 200 {object} models.ResponseResult{result=models.ServeResult201}
// @Failure 400 {object} models.ResponseError{error=models.ServeError400}
// @Failure 401 {object} models.ResponseError{error=string}
// @Failure 406 {object} models.ResponseError{error=string}
// @Failure 404
// @Failure 500
// @Router /serve-histories/{serve_id}/done-step [post]
func ServeEditStepByServeID(c *gin.Context) {
	var serve_id = c.Param("serve_id")
	serve_id_uint64, _ := strconv.ParseUint(serve_id, 10, 64)

	var serveUpdatestep models.ServeUpdateStep

	if ok, errors := helpers.ValidateServe(c, &serveUpdatestep); !ok {
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

	tokenAuth, err := helpers.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Unauthorized"})
		return
	}

	var serve = models.Serve{
		ID: uint(serve_id_uint64),
	}

	if err := helpers.DB.Model(serve).Where("ID = ?", serve_id).First(&serve).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serve_id) + " not found"})
		return
	}

	if tokenAuth.UserId != uint64(serve.UserID) {
		c.JSON(http.StatusForbidden, models.ResponseError{Success: false, Message: "Forbidden"})
		return
	}

	var steps []models.ServeRecipeStep
	var serveSteps []models.ServeStep
	err = helpers.DB.Model(&serveSteps).
		Select("serve_steps.*", "recipe_steps.step_order", "recipe_steps.description").
		Joins("INNER JOIN recipe_steps ON serve_steps.recipe_step_id = recipe_steps.id").
		Where(models.ServeStep{ServeID: serve.ID}).
		Find(&steps).
		Order("step_order asc").Error
	fmt.Println(steps)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serve.ID) + " has no steps"})
		return
	} else {
		var updateId uint
		for _, val := range steps {

			if val.StepOrder < serveUpdatestep.StepOrder && !val.Done {
				c.JSON(http.StatusConflict, models.ResponseError{Success: false, Message: "Some steps before " + fmt.Sprint(serveUpdatestep.StepOrder) + " is not done yet"})
				return
			}

			if val.StepOrder == serveUpdatestep.StepOrder && !val.Done {
				updateId = val.ID
			}

		}

		if updateId > 0 {
			var serveStep = models.ServeStep{ID: updateId}

			err = helpers.DB.Model(&serveStep).Update("done", true).Error
			if err != nil {
				c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Failed to update"})
				return
			} else {
				var recipe = models.Recipe{
					ID: serve.RecipeID,
				}

				if err := helpers.DB.Model(recipe).Where("ID = ?", serve.RecipeID).First(&recipe).Error; err != nil {
					c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serve.RecipeID) + " not found"})
					return
				}

				var stepsUpdated []models.ServeRecipeStep
				var serveStepsUpdated []models.ServeStep
				err = helpers.DB.Model(&serveStepsUpdated).
					Select("serve_steps.*", "recipe_steps.step_order", "recipe_steps.description").
					Joins("INNER JOIN recipe_steps ON serve_steps.recipe_step_id = recipe_steps.id").
					Where(models.ServeStep{ServeID: serve.ID}).
					Find(&stepsUpdated).
					Order("step_order asc").Error
				fmt.Println(stepsUpdated)

				if err != nil {
					c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe.ID) + " has no steps"})
					return
				} else {
					var serveStepResults []models.ServeStepResult

					var undoneCount = 0
					for _, val := range stepsUpdated {
						serveStepResults = append(serveStepResults, models.ServeStepResult{
							StepOrder:   val.StepOrder,
							Description: val.Description,
							Done:        val.Done,
						})

						if !val.Done {
							undoneCount++
						}
					}

					var recipeCategory models.RecipeCategory
					recipeCategory.ID = recipe.RecipeCategoryId
					helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

					nStep := float64(len(stepsUpdated))
					nStepDone := float64(len(stepsUpdated) - undoneCount)
					var status string
					if nStepDone >= nStep {
						status = "done"
					} else {
						status = "progress"
					}

					c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.ServeResult201{
						ID:                 serve.ID,
						UserID:             serve.UserID,
						RecipeID:           serve.RecipeID,
						RecipeName:         recipe.Name,
						RecipeCategoryName: recipeCategory.Name,
						RecipeImage:        recipe.Image,
						RecipeCategoryId:   recipe.RecipeCategoryId,
						NServing:           serve.NServing,
						NStep:              nStep,
						NStepDone:          nStepDone,
						Reaction:           serve.Reaction,
						Steps:              serveStepResults,
						Status:             status,
						CreatedAt:          serve.CreatedAt,
						UpdatedAt:          serve.UpdatedAt,
					}})
				}

			}

		}

	}

}

func ServeGetByServeID(c *gin.Context) {
	var serve_id = c.Param("serve_id")
	serve_id_uint64, _ := strconv.ParseUint(serve_id, 10, 64)

	tokenAuth, err := helpers.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Unauthorized"})
		return
	}

	var serve = models.Serve{
		ID: uint(serve_id_uint64),
	}

	if err := helpers.DB.Model(serve).Where("ID = ?", serve_id).First(&serve).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serve_id) + " not found"})
		return
	}

	if tokenAuth.UserId != uint64(serve.UserID) {
		c.JSON(http.StatusForbidden, models.ResponseError{Success: false, Message: "Forbidden"})
		return
	}

	var recipe = models.Recipe{
		ID: serve.RecipeID,
	}

	if err := helpers.DB.Model(recipe).Where("ID = ?", serve.RecipeID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serve.RecipeID) + " not found"})
		return
	}

	var stepsUpdated []models.ServeRecipeStep
	var serveStepsUpdated []models.ServeStep
	err = helpers.DB.Model(&serveStepsUpdated).
		Select("serve_steps.*", "recipe_steps.step_order", "recipe_steps.description").
		Joins("INNER JOIN recipe_steps ON serve_steps.recipe_step_id = recipe_steps.id").
		Where(models.ServeStep{ServeID: serve.ID}).
		Find(&stepsUpdated).
		Order("step_order asc").Error
	fmt.Println(stepsUpdated)

	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(recipe.ID) + " has no steps"})
		return
	} else {
		var serveStepResults []models.ServeStepResult

		var undoneCount = 0
		for _, val := range stepsUpdated {
			serveStepResults = append(serveStepResults, models.ServeStepResult{
				StepOrder:   val.StepOrder,
				Description: val.Description,
				Done:        val.Done,
			})

			if !val.Done {
				undoneCount++
			}
		}

		var recipeCategory models.RecipeCategory
		recipeCategory.ID = recipe.RecipeCategoryId
		helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

		nStep := float64(len(stepsUpdated))
		nStepDone := float64(len(stepsUpdated) - undoneCount)
		var status string
		if nStepDone >= nStep {
			if serve.Reaction == models.ReactionUnknown {
				status = "need-rating"
			} else {
				status = "done"
			}
		} else {
			status = "progress"
		}

		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.ServeResult201{
			ID:                 serve.ID,
			UserID:             serve.UserID,
			RecipeID:           serve.RecipeID,
			RecipeName:         recipe.Name,
			RecipeCategoryName: recipeCategory.Name,
			RecipeImage:        recipe.Image,
			RecipeCategoryId:   recipe.RecipeCategoryId,
			NServing:           serve.NServing,
			NStep:              nStep,
			NStepDone:          nStepDone,
			Reaction:           serve.Reaction,
			Steps:              serveStepResults,
			Status:             status,
			CreatedAt:          serve.CreatedAt,
			UpdatedAt:          serve.UpdatedAt,
		}})
	}
}

func ServeGetAll(c *gin.Context) {
	var skip = c.Query("skip")
	var limit = c.Query("limit")
	var sort = c.Query("sort")
	var q = c.Query("q")
	var statusFilter = c.Query("status")
	var categoryId = c.Query("categoryId")

	var serves []models.Serve
	var servesResult []models.ServeResultGetAll
	var servesResultFiltered []models.ServeResultGetAll

	query := helpers.DB.Model(&serves).
		Select(
			"serves.id as id", "serves.n_serving as n_serving", "serves.reaction as reaction", "serves.created_at as created_at", "serves.updated_at as updated_at",
			"recipes.id as recipe_id", "recipes.name as recipe_name", "recipes.image as recipe_image", "recipes.recipe_category_id as recipe_category_id",
			"recipe_categories.name as recipe_category_name",
		).
		Joins("INNER JOIN recipes ON serves.recipe_id = recipes.id").
		Joins("INNER JOIN recipe_categories ON recipes.recipe_category_id = recipe_categories.id")

	categoryId_uint64, _ := strconv.ParseUint(categoryId, 10, 64)
	if categoryId_uint64 > 0 {
		query.Where("recipes.recipe_category_id", categoryId_uint64)
	}

	if q != "" {
		query.Where("recipes.name LIKE ?", "%"+q+"%")
	}

	if sort != "" {
		//newest|oldest|nserve_asc|nserve_desc
		s := strings.Split(sort, "_")

		if s[0] != "" {
			if s[0] == "nserve" {
				s[0] = "serves.n_serving"
				query.Order(s[0] + " " + s[1])
			} else {
				s[0] = "created_at"
				if s[0] == "oldest" {
					query.Order(s[0] + " desc")
				} else {
					query.Order(s[0] + " asc")
				}
			}
		}

	}

	if limit != "" {
		limit_uint64, _ := strconv.ParseInt(limit, 10, 64)
		query.Limit(int(limit_uint64))
	}

	if skip != "" {
		skip_uint64, _ := strconv.ParseInt(skip, 10, 64)
		query.Offset(int(skip_uint64))
	}

	query.Find(&servesResult)

	fmt.Println(servesResult)

	if query.Error != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe not found"})
		return
	} else {
		for idx, serveResult := range servesResult {
			var steps []models.ServeRecipeStep
			var serveSteps []models.ServeStep

			err := helpers.DB.Model(&serveSteps).
				Select("serve_steps.*", "recipe_steps.step_order", "recipe_steps.description").
				Joins("INNER JOIN recipe_steps ON serve_steps.recipe_step_id = recipe_steps.id").
				Where(models.ServeStep{ServeID: serveResult.ID}).
				Find(&steps).
				Order("step_order asc").Error
			if err != nil {
				c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serveResult.ID) + " has no steps"})
				return
			} else {
				var recipe = models.Recipe{
					ID: serveResult.RecipeID,
				}

				if err := helpers.DB.Model(recipe).Where("ID = ?", serveResult.RecipeID).First(&recipe).Error; err != nil {
					c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serveResult.RecipeID) + " not found"})
					return
				}

				var nStep = len(steps)
				var nStepDone = 0

				for _, val := range steps {
					if val.Done {
						nStepDone++
					}
				}

				status := ""
				if nStepDone >= nStep {
					if serveResult.Reaction == models.ReactionUnknown {
						status = "need-rating"
					} else {
						status = "done"
					}
				} else {
					status = "progress"
				}

				servesResult[idx].NStep = float64(nStep)
				servesResult[idx].NStepDone = float64(nStepDone)
				servesResult[idx].Status = status

				if statusFilter != "" {
					if statusFilter == status {
						servesResultFiltered = append(servesResultFiltered, servesResult[idx])
					}
				} else {
					servesResultFiltered = append(servesResultFiltered, servesResult[idx])
				}
			}
		}

		data := struct {
			Total   int                        `json:"total"`
			History []models.ServeResultGetAll `json:"history"`
		}{
			Total:   len(servesResultFiltered),
			History: servesResultFiltered,
		}

		c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: data})
		return
	}
}

func ServeCreateReactionByServeID(c *gin.Context) {
	var serve_id = c.Param("serve_id")
	serve_id_uint64, _ := strconv.ParseUint(serve_id, 10, 64)

	var serveUpdateReaction models.ServeUpdateReaction
	if ok, errors := helpers.ValidateServe(c, &serveUpdateReaction); !ok {
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

	if models.GetReactionId(serveUpdateReaction.Reaction) == models.ReactionUnknown {
		c.JSON(http.StatusBadRequest, models.ResponseError{Success: false, Message: "reaction is invalid"})
		return
	}

	tokenAuth, err := helpers.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ResponseError{Success: false, Message: "Unauthorized"})
		return
	}

	var serve = models.Serve{
		ID: uint(serve_id_uint64),
	}

	if err := helpers.DB.Model(serve).Where("ID = ?", serve_id).First(&serve).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serve_id) + " not found"})
		return
	}

	if tokenAuth.UserId != uint64(serve.UserID) {
		c.JSON(http.StatusForbidden, models.ResponseError{Success: false, Message: "Forbidden"})
		return
	}

	var steps []models.ServeRecipeStep
	var serveSteps []models.ServeStep
	err = helpers.DB.Model(&serveSteps).
		Select("serve_steps.*", "recipe_steps.step_order", "recipe_steps.description").
		Joins("INNER JOIN recipe_steps ON serve_steps.recipe_step_id = recipe_steps.id").
		Where(models.ServeStep{ServeID: serve.ID}).
		Find(&steps).
		Order("step_order asc").Error
	if err != nil {
		c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Serve history with id " + fmt.Sprint(serve.ID) + " has no steps"})
		return
	} else {
		var recipe = models.Recipe{
			ID: serve.RecipeID,
		}

		if err := helpers.DB.Model(recipe).Where("ID = ?", serve.RecipeID).First(&recipe).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ResponseError{Success: false, Message: "Recipe with id " + fmt.Sprint(serve.RecipeID) + " not found"})
			return
		}

		var serveStepResults []models.ServeStepResult
		for _, val := range steps {
			if !val.Done {
				c.JSON(http.StatusBadRequest, models.ResponseError{Success: false, Message: "Invalid status, status need to be need-reaction"})
				return
			}

			serveStepResults = append(serveStepResults, models.ServeStepResult{
				StepOrder:   val.StepOrder,
				Description: val.Description,
				Done:        val.Done,
			})

		}

		serve.Reaction = models.GetReactionId(serveUpdateReaction.Reaction)
		if err := helpers.DB.Save(&serve).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			var recipeCategory models.RecipeCategory
			recipeCategory.ID = recipe.RecipeCategoryId
			helpers.DB.Model(&recipeCategory).Where(models.RecipeCategory{ID: recipe.RecipeCategoryId}).First(&recipeCategory)

			c.JSON(http.StatusOK, models.ResponseResult{Success: true, Message: "Success", Data: models.ServeResult201{
				ID:                 serve.ID,
				UserID:             serve.UserID,
				RecipeID:           serve.RecipeID,
				RecipeName:         recipe.Name,
				RecipeCategoryName: recipeCategory.Name,
				RecipeImage:        recipe.Image,
				RecipeCategoryId:   recipe.RecipeCategoryId,
				NServing:           serve.NServing,
				NStep:              float64(len(serveSteps)),
				NStepDone:          float64(len(serveSteps)),
				Steps:              serveStepResults,
				Status:             "done",
				CreatedAt:          serve.CreatedAt,
				UpdatedAt:          serve.UpdatedAt,
			}})

		}

	}

}
