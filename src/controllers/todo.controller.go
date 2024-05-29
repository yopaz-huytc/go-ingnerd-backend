package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
	"github.com/yopaz-huytc/go-ingnerd-backend/src/config"
	"github.com/yopaz-huytc/go-ingnerd-backend/src/models"
	"gorm.io/gorm"
)

var db *gorm.DB = config.ConnectDB()
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Todo struct for request body
type todoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsDone      int    `json:"is_done" validate:"gte=0,lte=1"`
}

// Defining the struct for the response body
type todoResponse struct {
	todoRequest
	ID uint `json:"id"`
}

// CreateTodo Create todo data to database
func CreateTodo(context *gin.Context) {
	var data todoRequest
	// Binding request body json to request body struct
	if err := context.ShouldBindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validate request body
	validationErr := validate.Struct(data)
	if validationErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return

	}
	// Matching todo models struct with todo request struct
	todo := models.Todo{}
	todo.Name = data.Name
	todo.Description = data.Description
	todo.IsDone = 0

	// Querying to database
	result := db.Create(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Matching result to create response
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description
	response.IsDone = todo.IsDone

	//create http response
	context.JSON(http.StatusCreated, response)
}

func GetAllTodos(context *gin.Context) {
	var todos []models.Todo
	// Querying to find todo data
	err := db.Find(&todos)
	if err.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error.Error()})
		return
	}
	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    todos,
	})
}

func UpdateTodo(context *gin.Context) {
	var data todoRequest
	// Defining request parameter to get todo id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)
	// Binding request body json to request body struct
	if err := context.ShouldBindJSON(&data); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validate request body
	validationErr := validate.Struct(data)
	if validationErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// Initiate models todo
	todo := models.Todo{}
	// Querying to find todo data by id from request parameter
	result := db.Where("id = ?", idTodo).First(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Updating todo with data from request
	todo.Name = data.Name
	todo.Description = data.Description
	todo.IsDone = data.IsDone
	// Saving updated todo to the database
	result = db.Save(&todo)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	// Matching todo request with todo models
	var response todoResponse
	response.ID = todo.ID
	response.Name = todo.Name
	response.Description = todo.Description
	response.IsDone = todo.IsDone
	//  Creating http response
	context.JSON(http.StatusOK, response)
}

// Delete todo data by id
func DeleteTodo(context *gin.Context) {
	// Initiate todo models
	todo := models.Todo{}
	// getting request parameter id
	reqParamId := context.Param("idTodo")
	idTodo := cast.ToUint(reqParamId)
	// Querying to delete todo data by id
	result := db.Where("id = ?", idTodo).Unscoped().Delete(&todo)
	fmt.Println(result) // print the result of the Delete operation
	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    idTodo,
	})
}
