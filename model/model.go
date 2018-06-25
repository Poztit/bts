package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math"
)

var database *sql.DB

type MesureCollection struct {
	Mesures []Mesure `json:"mesures"`
}

type Mesure struct {
	Capteur     int     `json:"capteur"`
	Temperature float64 `json:"temperature"`
	Humidite    float64 `json:"humidite"`
	Date        string  `json:"date"`
}

type CapteurCollection struct {
	Capteurs []Capteur `json:"capteurs"`
}

type Capteur struct {
	Nom       string  `json:"nom"`
	Lieu      string  `json:"lieu"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Etat      string  `json:"etat"`
}

func InitDB(dbpath string) {
	conn, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	} else {
		database = conn
	}
}

func CloseDB() {
	database.Close()
}

func CheckLogin(username string, password string) bool {
	var name string

	//fmt.Println("hello", "hello")
	//fmt.Println(username, password)

	err := database.QueryRow(`
		SELECT nom FROM utilisateur
		WHERE utilisateur.nom='admin' AND utilisateur.motdepasse='admin'
		`).Scan(&name)

	if err != nil {
		fmt.Println("Erreur de requete")
		fmt.Println(err)
		return false
	}

	return true
}

// Mesure
// GET
func FindMesureCollection() *MesureCollection {

	rows, err := database.Query(`
		SELECT fk_capteur, temperature, humidite, date FROM mesures
		ORDER BY date DESC LIMIT 10
		`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := new(MesureCollection)
	for rows.Next() {
		mesure := Mesure{}
		err := rows.Scan(&mesure.Capteur, &mesure.Temperature, &mesure.Humidite, &mesure.Date)
		if err != nil {
			panic(err)
		}
		m.Mesures = append(m.Mesures, mesure)
	}
	return m
}

// POST
func AddMesure(m *Mesure) error {
	m.Humidite = round(m.Humidite, 0.001)
	m.Temperature = round(m.Temperature, 0.001)

	_, err := database.Exec(`
		INSERT INTO mesures (fk_capteur, temperature, humidite, date)
		VALUES (? , ?, ?, datetime('now'))
		`, m.Capteur, m.Temperature, m.Humidite)
	return err
}

// Capteur
// GET
func FindCapteurCollection() *CapteurCollection {

	rows, err := database.Query(`
		SELECT nom, lieu, latitude, longitude, etat FROM capteurs
		`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	c := new(CapteurCollection)
	for rows.Next() {
		capteur := Capteur{}
		err := rows.Scan(&capteur.Nom, &capteur.Lieu, &capteur.Latitude, &capteur.Longitude, &capteur.Etat)
		if err != nil {
			panic(err)
		}
		c.Capteurs = append(c.Capteurs, capteur)
	}
	return c
}

func FindCapteur(id string) *CapteurCollection {

	rows, err := database.Query(`
		SELECT nom, lieu, latitude, longitude, etat FROM capteurs
		WHERE capteurs.id=?
		`, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	c := new(CapteurCollection)
	for rows.Next() {
		capteur := Capteur{}
		err := rows.Scan(&capteur.Nom, &capteur.Lieu, &capteur.Latitude, &capteur.Longitude, &capteur.Etat)
		if err != nil {
			panic(err)
		}
		c.Capteurs = append(c.Capteurs, capteur)
	}
	return c
}

func FindCapteurMesures(id string) *MesureCollection {

	rows, err := database.Query(`
		SELECT fk_capteur, temperature, humidite, date
		FROM mesures
		WHERE mesures.fk_capteur=?
		ORDER BY date ASC LIMIT 8
		`, id)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := new(MesureCollection)
	for rows.Next() {
		mesure := Mesure{}
		err := rows.Scan(&mesure.Capteur, &mesure.Temperature, &mesure.Humidite, &mesure.Date)
		if err != nil {
			panic(err)
		}
		m.Mesures = append(m.Mesures, mesure)
	}
	return m
}

// POST
func AddCapteur(s *Capteur) error {
	_, err := database.Exec(`
		INSERT INTO mesures (nom, lieu, latitude, longitude)
		VALUES (?, ?, ?, ?)
		`, s.Nom, s.Lieu, s.Latitude, s.Longitude)
	return err
}

// Utils

func round(n, nearest float64) float64 {
	return math.Round(n/nearest) * nearest
}
