package api

import (
	"database/sql"

	"github.com/KhoirulAziz99/final_project_e_wallet/api/handler"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/app"
	"github.com/KhoirulAziz99/final_project_e_wallet/internal/repository"
	"github.com/KhoirulAziz99/final_project_e_wallet/pkg"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(db *sql.DB) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	userService := app.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userService)
	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	userRouters := apiV1.Group("/users")
	{
		userRouters.POST("/login", pkg.LoginGPTHandler(db))
		userRouters.POST("/", userHandler.InsertUser)
		userRouters.PUT("/:id", userHandler.UpdateUser)
		userRouters.DELETE("/:id", userHandler.DeleteUser)
		userRouters.GET("/:id", userHandler.FindOneUser)
		userRouters.GET("/", userHandler.FindAllUsers)
	}

	return r
}

