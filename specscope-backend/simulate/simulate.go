package simulate

import (
    "math"
    "math/rand"
    "time"
)

// SpectrumReading represents one signal sample
type SpectrumReading struct {
    Timestamp time.Time `json:"timestamp"`
    Frequency float64   `json:"frequency"` // in MHz
    Power     float64   `json:"power"`     // in dBm
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

// GenerateSimulatedData returns a slice of fake readings
func GenerateSimulatedData(samples int, startFreq, endFreq float64) []SpectrumReading {
    readings := make([]SpectrumReading, samples)
    now := time.Now()
	freqStep := (endFreq - startFreq) / float64(samples)

    for i := 0; i < samples; i++ {
        freq := startFreq + float64(i)*freqStep // Covers 300 MHz to ~2800 MHz
        t := float64(i) / 10.0
		lat := 37.0 + rand.Float64()*0.1 //around sf
		lon := -122.0 + rand.Float64()*0.1

        var power float64

        switch {
        // FM band: 88–108 MHz 
        case freq >= 88 && freq <= 108:
            power = 35 + rand.Float64()*5 // 35 to 40 dBm

        // Wi-Fi band: 2400–2485 MHz 
        case freq >= 2400 && freq <= 2485:
            power = 25 + 10*math.Sin(0.2*t) + rand.Float64()*5
		
		// AM Radio (530–1700 kHz converted to MHz)
		case freq >= 0.53 && freq <= 1.7:
			power = 10 + rand.Float64()*5

		// Bluetooth 
		case freq >= 2402 && freq <= 2480 && rand.Float64() < 0.3:
			power = 20 + rand.Float64()*10

		// GSM cellular band
		case (freq >= 850 && freq <= 900) || (freq >= 1800 && freq <= 1900):
			power = 30 + 5*math.Sin(0.5*t) + rand.Float64()*5

		// GPS
		case freq >= 1574 && freq <= 1576:
			power = -60 + rand.Float64()*3

		// NOAA Weather 
		case freq >= 162.4 && freq <= 162.55:
			power = 45 + rand.Float64()*2


        // Quiet zone
        case freq >= 300 && freq <= 350:
            power = -80 + rand.Float64()*5

        // Random interference spikes (anywhere)
        case rand.Float64() < 0.01: // 1% chance per sample
            power = 50 + rand.Float64()*10

        // General low-noise background
        default:
            power = -50 + 10*math.Sin(0.1*t) + rand.Float64()*10
        }

        readings[i] = SpectrumReading{
            Timestamp: now.Add(time.Duration(i) * time.Millisecond),
            Frequency: freq,
            Power:     power,
			Latitude:  lat,
			Longitude: lon,
        }
    }

    return readings
}

//Interference Detector
func DetectInterference(data []SpectrumReading, powerThreshold float64, proximityMHz float64) []SpectrumReading{
	var result []SpectrumReading

	for i := 0; i < len(data)-1; i++{
		a := data[i]
		b := data[i+1]

		//power spike 
		if a.Power > powerThreshold {
			result = append(result, a)
			continue
		}

		//overlapping signals
		if math.Abs(a.Frequency-b.Frequency) < proximityMHz && math.Abs(a.Power-b.Power) > 10 {
			result = append(result, a)
		}
	}

	return result
}
