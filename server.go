package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/eduardoacuna/robowars/engine"
)

func main() {
	hub := engine.NewHub()
	go hub.Run()

	http.HandleFunc("/play", hub.ConnectPlayer())
	http.HandleFunc("/start-game", handleStartGame(hub))
	http.HandleFunc("/stop-game", handleStopGame(hub))

	log.Printf("starting server at localhost:3001\n")
	err := http.ListenAndServe("localhost:3001", nil)
	if err != nil {
		log.Fatalf("[fatal] listen and serve: %v\n", err)
	}
}

func handleStartGame(hub *engine.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Printf("received request to start the game\n")
		secret := r.FormValue("secret")
		if secret != "abretesesamo" {
			log.Printf("request for starting game was declined\n")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		rowsStr := r.FormValue("rows")
		colsStr := r.FormValue("cols")
		wallRootsStr := r.FormValue("wall-roots")
		wallBuildingProbStr := r.FormValue("wall-building-prob")
		povRadiusStr := r.FormValue("pov-radius")

		if rowsStr == "" || colsStr == "" ||
			wallRootsStr == "" || wallBuildingProbStr == "" ||
			povRadiusStr == "" {
			log.Printf("malformed configuration: missing config\n")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		rows, err1 := strconv.ParseInt(rowsStr, 10, 64)
		cols, err2 := strconv.ParseInt(colsStr, 10, 64)
		wallRoots, err3 := strconv.ParseInt(wallRootsStr, 10, 64)
		wallBuildingProb, err4 := strconv.ParseFloat(wallBuildingProbStr, 64)
		povRadius, err5 := strconv.ParseInt(povRadiusStr, 10, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
			log.Printf("[error] malformed configuration: parsing issue\n")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		config := &engine.Config{
			Rows:             int(rows),
			Cols:             int(cols),
			WallRoots:        int(wallRoots),
			WallBuildingProb: wallBuildingProb,
			POVRadius:        int(povRadius),
		}

		log.Printf("request for starting game was accepted\n")
		hub.Playing <- config

		w.Write(nil)
	}
}

func handleStopGame(hub *engine.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Printf("received request to stop the game\n")
		secret := r.FormValue("secret")
		if secret != "cierratesesamo" {
			log.Printf("request for stopping the game was declided\n")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		log.Printf("request for stopping game was accepted\n")
		hub.Playing <- nil

		w.Write(nil)
	}
}
