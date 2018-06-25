package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"
)

var database *sql.DB
var seed rand.Source
var gen *rand.Rand

type Mesure struct {
	Sensor      int
	Temperature float32
	Humidite    float32
	Date        string
}

func main() {
	seed = rand.NewSource(time.Now().UnixNano())
	gen = rand.New(seed)
	now := time.Now()

	fmt.Println("Now: " + now.String())

	InitDB("./storage.db")

	_, err := database.Exec(`DELETE FROM mesures`)
	if err != nil {
		fmt.Println(err)
	}

	var m Mesure
	var maxCount = 30

	for i := 0; i < maxCount; i++ {
		fmt.Printf("\r%d / %d", i, maxCount)

		now = now.Add(3 * time.Hour)
		m.Date = now.Format(time.RFC3339)

		for j := 1; j <= 8; j++ {
			m.Sensor = j
			m.Temperature = generateTemperature()
			m.Humidite = generateHumidite()
			postMesureModel(&m)
		}
	}

	fmt.Printf("\r%d / %d\n", maxCount, maxCount)

	CloseDB()
}

func generateHumidite() float32 {
	n := gen.Float32()
	n = mapRange(n, 0.0, 1.0, 30.0, 60.0)
	return n
}

func generateTemperature() float32 {
	n := gen.Float32()
	n = mapRange(n, 0.0, 1.0, 15, 30)
	return n
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

func postMesureModel(m *Mesure) error {
	_, err := database.Exec(`
		INSERT INTO mesures (fk_capteur, temperature, humidite, date)
		VALUES (? , ?, ?, ?)
		`, m.Sensor, m.Temperature, m.Humidite, m.Date)
	return err
}

func mapRange(x float32, inMin float32, inMax float32, outMin float32, outMax float32) float32 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
