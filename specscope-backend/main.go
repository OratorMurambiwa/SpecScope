package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"

    
    "specscope/simulate"

    
    "github.com/rs/cors"
)

// Struct to match the expected input for the ML model (FastAPI)
type Reading struct {
    Frequency float64 `json:"frequency"`
    Power     float64 `json:"power"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Hour      int     `json:"hour"`
}

// Struct to decode the ML model's prediction response
type PredictionResponse struct {
    Interference bool    `json:"interference"`
    Confidence   float64 `json:"confidence"`
}

// This function sends one reading to the FastAPI ML model and gets the prediction
func checkInterferenceML(reading Reading) (PredictionResponse, error) {
    url := "http://localhost:8000/predict" 

    // convert the reading to JSON
    payload, err := json.Marshal(reading)
    if err != nil {
        return PredictionResponse{}, err
    }

    // make a POST request to FastAPI
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return PredictionResponse{}, err
    }
    defer resp.Body.Close()

    // decode the JSON response
    var result PredictionResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return PredictionResponse{}, err
    }

    return result, nil
}

func main() {
    mux := http.NewServeMux() // new HTTP router

    // Just a homepage to confirm the backend is running
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte("SpecScope backend is running!\nTry visiting /data, /interference or /ml-interference"))
    })

    // This gives me the raw spectrum data (generated using math/rand etc.)
    mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
        readings := simulate.GenerateSimulatedData(500, 300, 2800)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(readings)
    })

    // This uses my rule-based detector (hardcoded) to find interference
    mux.HandleFunc("/interference", func(w http.ResponseWriter, r *http.Request) {
        readings := simulate.GenerateSimulatedData(500, 300, 2800)
        flagged := simulate.DetectInterference(readings, 30.0, 1.0) // power > 30dBm or freq overlap < 1MHz
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(flagged)
    })

    //I use the AI model I trained (RandomForest) to detect interference
    mux.HandleFunc("/ml-interference", func(w http.ResponseWriter, r *http.Request) {
        readings := simulate.GenerateSimulatedData(100, 300, 2800) // smaller batch so it's faster
        var results []map[string]interface{}
        now := time.Now()

        // for every signal, send it to the ML model and add the prediction to the response
        for _, r := range readings {
            reading := Reading{
                Frequency: r.Frequency,
                Power:     r.Power,
                Latitude:  r.Latitude,
                Longitude: r.Longitude,
                Hour:      now.Hour(), 
            }

            pred, err := checkInterferenceML(reading)
            if err != nil {
                log.Printf("Prediction error: %v", err)
                continue // if prediction fails, just skip it
            }

            // add the original reading + prediction to the output
            results = append(results, map[string]interface{}{
                "timestamp":    r.Timestamp,
                "frequency":    r.Frequency,
                "power":        r.Power,
                "latitude":     r.Latitude,
                "longitude":    r.Longitude,
                "interference": pred.Interference,
                "confidence":   pred.Confidence,
            })
        }

        // return all the results as JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(results)
    })

    // allow frontend to access this backend from a different port
    handler := cors.Default().Handler(mux)

    log.Println("Server running at http://localhost:8081")
    if err := http.ListenAndServe(":8081", handler); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
