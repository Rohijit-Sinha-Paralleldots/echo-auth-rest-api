package main

import (
	"github.com/Rohijit-Sinha-Paralleldots/echo-auth-rest-api/storage"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	storage.InitDB()
	e.Logger.Fatal(e.Start(":8080"))
}
