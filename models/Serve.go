package models

import (
	"strings"
	"time"
)

type Serve struct {
	ID         uint        `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	NServing   float64     `form:"nServing" json:"nServing" binding:"required"`
	RecipeID   uint        `form:"recipeId" json:"recipeId"`
	UserID     uint        `form:"userId" json:"userId"`
	Reaction   Reaction    `form:"reaction" json:"reaction"`
	ServeSteps []ServeStep `gorm:"foreignKey:ServeID"`
	CreatedAt  time.Time   `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt  time.Time   `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt  *time.Time  `form:"deletedAt" json:"deletedAt" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

// Reaction represent given rating
type Reaction int

const (
	ReactionUnknown Reaction = iota
	ReactionLike
	ReactionNeutral
	ReactionDislike
)

func (g Reaction) String() string {
	return []string{"", "like", "neutral", "dislike"}[g]
}

func GetReactionId(val string) Reaction {
	var r Reaction
	switch strings.ToLower(val) {
	case "like":
		r = ReactionLike
	case "neutral":
		r = ReactionNeutral
	case "dislike":
		r = ReactionDislike
	default:
		r = ReactionUnknown
	}
	return r
}

type ServeCreate struct {
	RecipeID uint     `form:"recipeId" json:"recipeId" binding:"required"`
	NServing *float64 `form:"nServing" json:"nServing,omitempty" binding:"required,min=1"`
}

type ServeUpdateStep struct {
	StepOrder int `form:"stepOrder" json:"stepOrder" binding:"required"`
}

type ServeUpdateReaction struct {
	Reaction string `form:"reaction" json:"reaction" binding:"required"`
}

type ServeStep struct {
	ID           uint       `gorm:"primaryKey" json:"-" form:"-" swaggertype:"integer"`
	ServeID      uint       `form:"serveId" json:"-"`
	RecipeStepID uint       `form:"recipeStepId" json:"recipeStepId"`
	Done         bool       `json:"done"`
	CreatedAt    time.Time  `form:"createdAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt    time.Time  `form:"updatedAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt    *time.Time `form:"deletedAt" json:"-" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type ServeRecipeStep struct {
	ID           uint       `gorm:"primaryKey" json:"-" form:"-" swaggertype:"integer"`
	ServeID      uint       `form:"serveId" json:"-"`
	RecipeStepID uint       `form:"recipeStepId" json:"recipeStepId"`
	Done         bool       `json:"done"`
	Description  string     `json:"description" form:"description"`
	StepOrder    int        `form:"stepOrder" json:"stepOrder" binding:"required"`
	CreatedAt    time.Time  `form:"createdAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt    time.Time  `form:"updatedAt" json:"-" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt    *time.Time `form:"deletedAt" json:"-" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

type ServeStepResult struct {
	StepOrder   int    `json:"stepOrder" form:"stepOrder"`
	Description string `json:"description" form:"description"`
	Done        bool   `json:"done"`
}

type ServeResultGetAll struct {
	ID                 uint      `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	RecipeID           uint      `form:"recipeId" json:"recipeId"`
	RecipeName         string    `form:"recipeName" json:"recipeName" `
	RecipeCategoryName string    `form:"recipeCategoryName" json:"recipeCategoryName" `
	RecipeImage        string    `form:"recipeImage" json:"recipeImage" `
	RecipeCategoryId   uint      `form:"recipeCategoryId" json:"recipeCategoryId"`
	NServing           float64   `form:"nServing" json:"nServing" binding:"required"`
	NStep              float64   `form:"nStep" json:"nStep" binding:"required"`
	NStepDone          float64   `form:"nStepDone" json:"nStepDone" binding:"required"`
	Reaction           Reaction  `form:"reaction" json:"reaction" `
	Status             string    `form:"status" json:"status" `
	CreatedAt          time.Time `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt          time.Time `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}
type ServeResult201 struct {
	ID                 uint              `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	UserID             uint              `form:"userId" json:"userId"`
	RecipeID           uint              `form:"recipeId" json:"recipeId"`
	RecipeName         string            `form:"recipeName" json:"recipeName" `
	RecipeCategoryName string            `form:"recipeCategoryName" json:"recipeCategoryName" `
	RecipeImage        string            `form:"recipeImage" json:"recipeImage" `
	RecipeCategoryId   uint              `form:"recipeCategoryId" json:"recipeCategoryId"`
	NServing           float64           `form:"nServing" json:"nServing" binding:"required"`
	NStep              float64           `form:"nStep" json:"nStep" binding:"required"`
	NStepDone          float64           `form:"nStepDone" json:"nStepDone" binding:"required"`
	Steps              []ServeStepResult `form:"steps" json:"steps" binding:"required"`
	Reaction           Reaction          `form:"reaction" json:"reaction" `
	Status             string            `form:"status" json:"status" `
	CreatedAt          time.Time         `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt          time.Time         `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}
