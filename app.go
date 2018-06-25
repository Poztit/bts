package main

import (
	"apiscpam/handler"
	"apiscpam/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	dbpath = "./storage.db"
)

/*
----/--> index.html @
	|
	static --> static @
	|
	api/mesure [GET, POST]
	|
	api/capteur [GET, POST]
		|
		api/capteur/:id [GET]
		|
		api/capteur/:id/mesure [GET]
*/

func main() {

	// Ouverture de la BDD
	model.InitDB(dbpath)

	// Fermeture de la BDD
	defer model.CloseDB()

	e := echo.New()

	// Application d'accès à l'API
	e.File("/", "index.html")
	e.Static("/static", "static")
	e.GET("/login", handler.Login)

	// API
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("poztit"),
	}))

	// Mesures
	apiGroup.POST("/mesure", handler.PostMesure)
	apiGroup.GET("/mesure", handler.GetMesureCollection)

	// Capteurs
	apiGroup.POST("/capteur", handler.PostCapteur)
	apiGroup.GET("/capteur", handler.GetCapteurCollection)
	apiGroup.GET("/capteur/:id", handler.GetCapteur)
	apiGroup.GET("/capteur/:id/mesure", handler.GetCapteurMesures)

	// Démarre du serveur de service web
	e.Start(":8080")
}
