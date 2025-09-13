package simulate

import (
	"math"
	"math/rand"
	"time"
)

// SpectrumReading represents one signal sample
type SpectrumReading struct {
	Timestamp      time.Time `json:"timestamp"`
	Frequency      float64   `json:"frequency"`       // in MHz
	Power          float64   `json:"power"`           // in dBm
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	Hour           int       `json:"hour"`            // extracted from timestamp
	TimePeriod     string    `json:"time_period"`     // morning, afternoon, evening, night
	WifiProximity  int       `json:"wifi_proximity"`  // 1 if near Wi-Fi band
	Interference   bool      `json:"interference"`    // (optional) default false
	Confidence     float64   `json:"confidence"`      // (optional) ML confidence
}

// getTimePeriod returns a label for part of day
func getTimePeriod(hour int) string {
	switch {
	case hour >= 6 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 18:
		return "afternoon"
	case hour >= 18 && hour < 24:
		return "evening"
	default:
		return "night"
	}
}

// GenerateSimulatedData creates spectrum readings with derived fields
func GenerateSimulatedData(samples int, startFreq, endFreq float64) []SpectrumReading {
	readings := make([]SpectrumReading, samples)
	now := time.Now()
	freqStep := (endFreq - startFreq) / float64(samples)

	for i := 0; i < samples; i++ {
		freq := startFreq + float64(i)*freqStep
		t := float64(i) / 10.0

		lat := 37.02 + rand.Float64()*0.05 // tighter range (SF-like area)
		lon := -121.93 + rand.Float64()*0.05

		// Simulated power value
		var power float64

		switch {
		case freq >= 88 && freq <= 108:
			power = 35 + rand.Float64()*5
		case freq >= 2400 && freq <= 2485:
			power = 25 + 10*math.Sin(0.2*t) + rand.Float64()*5
		case freq >= 0.53 && freq <= 1.7:
			power = 10 + rand.Float64()*5
		case freq >= 2402 && freq <= 2480 && rand.Float64() < 0.3:
			power = 20 + rand.Float64()*10
		case (freq >= 850 && freq <= 900) || (freq >= 1800 && freq <= 1900):
			power = 30 + 5*math.Sin(0.5*t) + rand.Float64()*5
		case freq >= 1574 && freq <= 1576:
			power = -60 + rand.Float64()*3
		case freq >= 162.4 && freq <= 162.55:
			power = 45 + rand.Float64()*2
		case freq >= 300 && freq <= 350:
			power = -80 + rand.Float64()*5
		case rand.Float64() < 0.01:
			power = 50 + rand.Float64()*10
		default:
			power = -50 + 10*math.Sin(0.1*t) + rand.Float64()*10
		}

		timestamp := now.Add(time.Duration(i) * time.Millisecond)
		hour := timestamp.Hour()
		timePeriod := getTimePeriod(hour)
		wifiProximity := 0
		if freq >= 2400 && freq <= 2485 {
			wifiProximity = 1
		}

		readings[i] = SpectrumReading{
			Timestamp:     timestamp,
			Frequency:     freq,
			Power:         power,
			Latitude:      lat,
			Longitude:     lon,
			Hour:          hour,
			TimePeriod:    timePeriod,
			WifiProximity: wifiProximity,
			Interference:  false,  // updated by ML or detection logic
			Confidence:    0.0,    // updated by ML
		}
	}

	return readings
}

// DetectInterference still works as your rule-based fallback
func DetectInterference(data []SpectrumReading, powerThreshold float64, proximityMHz float64) []SpectrumReading {
	var result []SpectrumReading

	for i := 0; i < len(data)-1; i++ {
		a := data[i]
		b := data[i+1]

		if a.Power > powerThreshold {
			result = append(result, a)
			continue
		}

		if math.Abs(a.Frequency-b.Frequency) < proximityMHz && math.Abs(a.Power-b.Power) > 10 {
			result = append(result, a)
		}
	}

	return result
}
