package main

import (
    "encoding/json"
    "log"
    "net/http"
    "specscope/simulate" 
    "github.com/rs/cors"
)

func main() {
    mux := http.NewServeMux()

    // API route to return simulated data
    mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
        readings := simulate.GenerateSimulatedData(500, 300, 2800) // number of samples
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(readings)
    })

	//Interference route 
	mux.HandleFunc("/interference", func(w http.ResponseWriter, r *http.Request) {
		readings := simulate.GenerateSimulatedData(500, 300, 2800) 
		flagged := simulate.DetectInterference(readings, 30.0, 1.0) // >30 dBm or <1 MHz proximity
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flagged)
	})
	

    // CORS middleware so frontend can access this from another port
    handler := cors.Default().Handler(mux)

    log.Println("🚀 Server running at http://localhost:8080")
    http.ListenAndServe(":8080", handler)
}
