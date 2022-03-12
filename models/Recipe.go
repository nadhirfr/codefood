package models

import (
	"sort"
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	ID                uint               `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name              string             `form:"name" json:"name" `
	Image             string             `form:"image" json:"image" `
	NReactionLike     int                `form:"nReactionLike" json:"nReactionLike" `
	NReactionNeutral  int                `form:"nReactionNeutral" json:"nReactionNeutral" `
	NReactionDislike  int                `form:"nReactionDislike" json:"nReactionDislike" `
	NServing          float64            `form:"nServing" json:"nServing" binding:"required"`
	RecipeCategoryId  uint               `form:"recipeCategoryId" json:"recipeCategoryId" binding:"required"`
	RecipeSteps       []RecipeStep       `gorm:"foreignKey:RecipeID"`
	Serves            []Serve            `gorm:"foreignKey:RecipeID"`
	RecipeIngridients []RecipeIngridient `gorm:"foreignKey:RecipeID"`
	CreatedAt         time.Time          `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt         time.Time          `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt         *time.Time         `form:"deletedAt" json:"deletedAt" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

func (u *Recipe) AfterFind(tx *gorm.DB) (err error) {
	if len(u.RecipeSteps) > 0 {
		sort.Slice(u.RecipeSteps, func(i, j int) bool { return u.RecipeSteps[i].StepOrder < u.RecipeSteps[j].StepOrder })
	}
	return
}

// AfterDelete hook defined for cascade delete
func (recipe *Recipe) AfterDelete(tx *gorm.DB) error {
	err := tx.Model(&RecipeStep{}).Where(RecipeStep{RecipeID: recipe.ID}).Unscoped().Delete(&RecipeStep{}).Error
	if err != nil {
		return err
	} else {
		return tx.Model(&RecipeIngridient{}).Where(RecipeIngridient{RecipeID: recipe.ID}).Unscoped().Delete(&RecipeIngridient{}).Error

	}
}

type RecipeCreate struct {
	Name                  string             `form:"name" json:"name" binding:"required"`
	RecipeCategoryId      uint               `form:"recipeCategoryId" json:"recipeCategoryId" binding:"required"`
	Image                 string             `form:"image" json:"image" binding:"required"`
	NServing              float64            `form:"nServing" json:"nServing" binding:"required"`
	IngredientsPerServing []RecipeIngridient `form:"ingredientsPerServing" json:"ingredientsPerServing" binding:"required"`
	Steps                 []RecipeStep       `form:"steps" json:"steps" binding:"required"`
}

type RecipeCategory struct {
	ID        uint       `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name      string     `form:"name" json:"name"`
	Recipes   []Recipe   `gorm:"foreignKey:RecipeCategoryId" json:"-"`
	CreatedAt time.Time  `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt time.Time  `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt *time.Time `form:"deletedAt" json:"-" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type RecipeCategoryCreate struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type RecipeCategoryResult201 struct {
	ID        uint      `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name      string    `form:"name" json:"name"`
	CreatedAt time.Time `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt time.Time `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type RecipeStep struct {
	ID          uint        `gorm:"primaryKey" json:"-" form:"-" swaggertype:"integer"`
	RecipeID    uint        `gorm:"recipeId" json:"-"`
	StepOrder   int         `json:"stepOrder" form:"stepOrder"`
	ServeSteps  []ServeStep `gorm:"foreignKey:RecipeStepID"`
	Description string      `json:"description" form:"description"`
	CreatedAt   time.Time   `form:"createdAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt   time.Time   `form:"updatedAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt   *time.Time  `form:"deletedAt" json:"-" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type RecipeIngridient struct {
	ID        uint       `gorm:"primaryKey" json:"-" form:"-" swaggertype:"integer"`
	RecipeID  uint       `gorm:"recipeId" json:"-"`
	Value     float64    `json:"value" form:"value"`
	Unit      string     `json:"unit" form:"unit"`
	Item      string     `json:"item" form:"item"`
	CreatedAt time.Time  `form:"createdAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt time.Time  `form:"updatedAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt *time.Time `form:"deletedAt" json:"-" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type RecipeResult201 struct {
	ID                    uint               `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name                  string             `form:"name" json:"name" `
	Image                 string             `form:"image" json:"image" `
	NReactionLike         int                `form:"nReactionLike" json:"nReactionLike" `
	NReactionNeutral      int                `form:"nReactionNeutral" json:"nReactionNeutral" `
	NReactionDislike      int                `form:"nReactionDislike" json:"nReactionDislike" `
	RecipeCategoryId      uint               `form:"recipeCategoryId" json:"recipeCategoryId"`
	NServing              float64            `form:"nServing" json:"nServing" binding:"required"`
	IngredientsPerServing []RecipeIngridient `form:"ingredientsPerServing" json:"ingredientsPerServing" binding:"required"`
	Steps                 []RecipeStep       `form:"steps" json:"steps" binding:"required"`
	CreatedAt             time.Time          `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt             time.Time          `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type RecipeResult200 struct {
	ID                    uint               `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name                  string             `form:"name" json:"name" `
	Image                 string             `form:"image" json:"image" `
	NReactionLike         int                `form:"nReactionLike" json:"nReactionLike" `
	NReactionNeutral      int                `form:"nReactionNeutral" json:"nReactionNeutral" `
	NReactionDislike      int                `form:"nReactionDislike" json:"nReactionDislike" `
	RecipeCategoryId      uint               `form:"recipeCategoryId" json:"recipeCategoryId"`
	NServing              float64            `form:"nServing" json:"nServing" binding:"required"`
	IngredientsPerServing []RecipeIngridient `form:"ingredientsPerServing" json:"ingredientsPerServing" binding:"required"`
	CreatedAt             time.Time          `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt             time.Time          `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	RecipeCategory        RecipeCategory     `form:"recipeCategory" json:"recipeCategory"`
}

type RecipeResultGetAll struct {
	ID               uint           `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name             string         `form:"name" json:"name" `
	Image            string         `form:"image" json:"image" `
	NReactionLike    int            `form:"nReactionLike" json:"nReactionLike" `
	NReactionNeutral int            `form:"nReactionNeutral" json:"nReactionNeutral" `
	NReactionDislike int            `form:"nReactionDislike" json:"nReactionDislike" `
	RecipeCategoryId uint           `form:"recipeCategoryId" json:"recipeCategoryId"`
	CreatedAt        time.Time      `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt        time.Time      `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	RecipeCategory   RecipeCategory `form:"recipeCategory" json:"recipeCategory"`
}

type RecipeResultSearch struct {
	ID   uint   `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Name string `form:"name" json:"name" `
}
