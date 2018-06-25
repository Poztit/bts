package handler

import (
	"apiscpam/model"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// Login
func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if model.CheckLogin(username, password) == false {
		fmt.Println("Erreur de connection")
		return c.String(http.StatusUnauthorized, "Bad login & password")
	}

	t, _ := createJWTToken(username)

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func createJWTToken(username string) (string, error) {

	// Création jeton
	token := jwt.New(jwt.SigningMethodHS512)

	// Paramètres du jeton
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = username
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	// Encodage du jeton et envoi
	t, err := token.SignedString([]byte("poztit"))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Mesures
// GET
func GetMesureCollection(c echo.Context) error {
	return c.JSON(http.StatusOK, model.FindMesureCollection())
}

// POST
func PostMesure(c echo.Context) (err error) {
	m := new(model.Mesure)
	if err = c.Bind(m); err != nil {
		return
	}

	if err = model.AddMesure(m); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, m)
}

// Capteurs
// GET
func GetCapteurCollection(c echo.Context) error {
	return c.JSON(http.StatusOK, model.FindCapteurCollection())
}

func GetCapteur(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, model.FindCapteur(id))
}

func GetCapteurMesures(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, model.FindCapteurMesures(id))
}

// POST
func PostCapteur(c echo.Context) (err error) {
	mesure := new(model.Mesure)
	if err = c.Bind(mesure); err != nil {
		return
	}

	if err = model.AddMesure(mesure); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, mesure)
}
