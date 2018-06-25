#!/bin/sh

curl -i -X POST -H "Content-Type:application/json" http://localhost:8080/api/mesure -d '{"capteur": 1, "temperature": 26.134023666381836, "humidite": 58.13821792602539, "date": "2018-03-20" }'

