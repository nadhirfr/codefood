package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nadhirfr/codefood/controllers"
	"github.com/nadhirfr/codefood/helpers"
)

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS for * and origin request it self (origin == origin) origins, allowing:
	// - POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "HEAD", "PATCH", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	user := r.Group("/auth")
	{
		user.POST("/login", controllers.UserLogin)
		user.POST("/register", controllers.UserRegister)
		user.GET("/detail", helpers.TokenAuthMiddleware(), controllers.UserGetByUserID)
	}

	recipe := r.Group("/recipes")
	{
		recipe.POST("", controllers.RecipeCreate)
		recipe.GET("", controllers.RecipeGetAll)
		recipe.DELETE("/:recipe_id", controllers.RecipeDeleteByRecipeID)
		recipe.PUT("/:recipe_id", controllers.RecipeEditByRecipeID)
		recipe.GET("/:recipe_id", controllers.RecipeGetByRecipeID)
		recipe.GET("/:recipe_id/steps", controllers.RecipeStepsGetByRecipeID)
	}

	search := r.Group("/search")
	{
		search.GET("/recipes", controllers.RecipeSearch)
	}

	recipeCategories := r.Group("/recipe-categories")
	{
		recipeCategories.POST("", controllers.RecipeCategoryCreate)
		recipeCategories.GET("", controllers.RecipeCategoryGetAll)
		recipeCategories.DELETE("/:recipeCategory_id", controllers.RecipeCategoryDeleteByRecipeCategoryID)
		recipeCategories.PUT("/:recipeCategory_id", controllers.RecipeCategoryEditByRecipeCategoryID)
	}

	serveHistories := r.Group("/serve-histories")
	{
		serveHistories.POST("", helpers.TokenAuthMiddleware(), controllers.ServeCreate)
		serveHistories.GET("", controllers.ServeGetAll)
		serveHistories.PUT("/:serve_id/done-step", helpers.TokenAuthMiddleware(), controllers.ServeEditStepByServeID)
		serveHistories.POST("/:serve_id/reaction", helpers.TokenAuthMiddleware(), controllers.ServeCreateReactionByServeID)
		serveHistories.GET("/:serve_id", controllers.ServeGetByServeID)

	}

	return r
}
