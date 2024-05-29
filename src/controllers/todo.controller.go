package controllers

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yopaz-huytc/go-ingnerd-backend/src/config"
    "github.com/yopaz-huytc/go-ingnerd-backend/src/models"
    "github.com/spf13/cast"
    "gorm.io/gorm"
)

var db *gorm.DB = config.ConnectDB()

// Todo struct for request body
type todoRequest struct {
    Name string `json:"name"`
    Description string `json:"description"`
}

// Defining the struct for the response body
type todoResponse struct {
    todoRequest
    ID uint `json:"id"`
}

// Create todo data to database
func CreateTodo (context *gin.Context) {
    var data todoRequest
     // Binding request body json to request body struct
    if err := context.ShouldBindJSON(&data); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Matching todo models struct with todo request struct
    todo := models.Todo{}
    todo.Name = data.Name
    todo.Description = data.Description

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

    //create http response
    context.JSON(http.StatusCreated, response)
}

func GetAllTodos(context *gin.Context) {
    var todos []models.Todo
    // Querying to find todo data
    err:= db.Find(&todos)
    if err.Error != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error.Error()})
        return
    }
    // Creating http response
    context.JSON(http.StatusOK, gin.H{
        "status" : "200",
        "message": "Success",
        "data": todos})
}

func UpdateTodo (context *gin.Context) {
    var data todoRequest
    // Defining request parameter to get todo id
    reqParamId := context.Param("idTodo")
    idTodo := cast.ToUint(reqParamId)
    // Binding request body json to request body struct
    if err := context.ShouldBindJSON(&data); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Initiate models todo
    todo := models.Todo{}
    // Querying to find todo data by id from request parameter
    todoById := db.Where("id = ?", idTodo).First(&todo)
    if todoById.Error != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": todoById.Error.Error()})
        return
    }
    // Matching todo request with todo models
    var response todoResponse
    response.ID = todo.ID
    response.Name = todo.Name
    response.Description = todo.Description
    //  Creating http response
    context.JSON(http.StatusCreated, response)
}

// Delete todo data by id
func DeleteTodo (context *gin.Context) {
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
        "status": "200",
        "message": "Success",
        "data": idTodo,
    })
}
