package main

import (
    "github.com/yopaz-huytc/go-ingnerd-backend/src/config"
    "github.com/yopaz-huytc/go-ingnerd-backend/src/routes"
    "gorm.io/gorm"
)

var (
    db *gorm.DB = config.ConnectDB()
)

func main() {
    defer config.DisconnectDB(db)
    routes.Routes()
}
