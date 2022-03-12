package includes

import (
	"github.com/nadhirfr/codefood/helpers"
	"github.com/nadhirfr/codefood/models"
)

func Migrate() bool {
	err := helpers.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.User{},
		&models.UserLoginFailed{},
		&models.RecipeCategory{},
		&models.Recipe{},
		&models.RecipeStep{},
		&models.RecipeIngridient{},
		&models.Serve{},
		&models.ServeStep{},
	)

	if err == nil {
		// 	var user = []models.User{
		// 		{ID: 1, Username: "user1", Password: "12345678"},
		// 		{ID: 2, Username: "user2", Password: "12345678"},
		// 	}
		// 	helpers.DB.Save(&user)

		// 	var recipeCategory = []models.RecipeCategory{
		// 		{ID: 1, Name: "Kue"},
		// 		{ID: 2, Name: "Sop"},
		// 		{ID: 3, Name: "Minuman"},
		// 	}
		// 	helpers.DB.Save(&recipeCategory)

		// 	var recipe = []models.Recipe{
		// 		{
		// 			ID:               1,
		// 			Name:             "Sop Ayam",
		// 			RecipeCategoryId: 2,
		// 			Image:            "https://i0.wp.com/masakanmama.com/wp-content/uploads/2019/11/resep-sop-ayam-bening.jpg?w=700&ssl=1",
		// 			NReactionLike:    194,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 1},
		// 		{
		// 			ID:               2,
		// 			Name:             "Es Slendang Mayang Nangka",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://resepgulaku.com/wp-content/uploads/2018/03/Resep-Gulaku_2018_800x400px_Es-Selendang-Mayang.jpg",
		// 			NReactionLike:    38,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 2},
		// 		{
		// 			ID:               3,
		// 			Name:             "Es KOPASUS (Kopi Pandan Susu)",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://cdn.idntimes.com/content-images/post/20191220/75489920-1007035202987220-1216100541960430738-n-1fa2ce1e5fa109bef9c092115de69e69.jpg",
		// 			NReactionLike:    13,
		// 			NReactionNeutral: 8,
		// 			NReactionDislike: 3,
		// 		},
		// 		{
		// 			ID:               4,
		// 			Name:             "Kue Pandan Kukus",
		// 			RecipeCategoryId: 1,
		// 			Image:            "https://img-global.cpcdn.com/recipes/642eb0633b28e2c7/751x532cq70/bolu-pandan-kukus-foto-resep-utama.jpg",
		// 			NReactionLike:    5,
		// 			NReactionNeutral: 3,
		// 			NReactionDislike: 2,
		// 		},
		// 		{
		// 			ID:               5,
		// 			Name:             "Sop Ayam 2",
		// 			RecipeCategoryId: 2,
		// 			Image:            "https://i0.wp.com/masakanmama.com/wp-content/uploads/2019/11/resep-sop-ayam-bening.jpg?w=700&ssl=1",
		// 			NReactionLike:    194,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 1},
		// 		{
		// 			ID:               6,
		// 			Name:             "Es Slendang Mayang Nangka 2",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://resepgulaku.com/wp-content/uploads/2018/03/Resep-Gulaku_2018_800x400px_Es-Selendang-Mayang.jpg",
		// 			NReactionLike:    38,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 2},
		// 		{
		// 			ID:               7,
		// 			Name:             "Es KOPASUS (Kopi Pandan Susu) 2",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://cdn.idntimes.com/content-images/post/20191220/75489920-1007035202987220-1216100541960430738-n-1fa2ce1e5fa109bef9c092115de69e69.jpg",
		// 			NReactionLike:    13,
		// 			NReactionNeutral: 8,
		// 			NReactionDislike: 3,
		// 		},
		// 		{
		// 			ID:               8,
		// 			Name:             "Kue Pandan Kukus 2",
		// 			RecipeCategoryId: 1,
		// 			Image:            "https://img-global.cpcdn.com/recipes/642eb0633b28e2c7/751x532cq70/bolu-pandan-kukus-foto-resep-utama.jpg",
		// 			NReactionLike:    5,
		// 			NReactionNeutral: 3,
		// 			NReactionDislike: 2,
		// 		},
		// 		{
		// 			ID:               9,
		// 			Name:             "Sop Ayam 3",
		// 			RecipeCategoryId: 2,
		// 			Image:            "https://i0.wp.com/masakanmama.com/wp-content/uploads/2019/11/resep-sop-ayam-bening.jpg?w=700&ssl=1",
		// 			NReactionLike:    194,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 1},
		// 		{
		// 			ID:               10,
		// 			Name:             "Es Slendang Mayang Nangka 3",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://resepgulaku.com/wp-content/uploads/2018/03/Resep-Gulaku_2018_800x400px_Es-Selendang-Mayang.jpg",
		// 			NReactionLike:    38,
		// 			NReactionNeutral: 4,
		// 			NReactionDislike: 2},
		// 		{
		// 			ID:               11,
		// 			Name:             "Es KOPASUS (Kopi Pandan Susu) 3",
		// 			RecipeCategoryId: 3,
		// 			Image:            "https://cdn.idntimes.com/content-images/post/20191220/75489920-1007035202987220-1216100541960430738-n-1fa2ce1e5fa109bef9c092115de69e69.jpg",
		// 			NReactionLike:    13,
		// 			NReactionNeutral: 8,
		// 			NReactionDislike: 3,
		// 		},
		// 		{
		// 			ID:               12,
		// 			Name:             "Kue Pandan Kukus 3",
		// 			RecipeCategoryId: 1,
		// 			Image:            "https://img-global.cpcdn.com/recipes/642eb0633b28e2c7/751x532cq70/bolu-pandan-kukus-foto-resep-utama.jpg",
		// 			NReactionLike:    5,
		// 			NReactionNeutral: 3,
		// 			NReactionDislike: 2,
		// 		},
		// 	}
		// 	helpers.DB.Save(&recipe)

		// 	//TODO add recipe steps
	}

	return true
}
