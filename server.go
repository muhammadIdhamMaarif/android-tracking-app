package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const logFile = "locations.csv"

type Location struct {
	Timestamp int64    `json:"timestamp"`           // ms dari epoch
	DeviceID  string   `json:"device_id,omitempty"` // default nya ke "unknown" kalau kosong
	Lat       *float64 `json:"lat"`                 // required (jan diapus)
	Lon       *float64 `json:"lon"`                 // required (jan diapus)
	Accuracy  *float64 `json:"accuracy,omitempty"`
	Speed     *float64 `json:"speed,omitempty"`
}

var (
	csvMu sync.Mutex 
)

func ensureCSV() error {
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		return w.Write([]string{"timestamp_ms", "received_iso", "device_id", "lat", "lon", "accuracy", "speed"})
	}
	return nil
}

func appendCSV(rec []string) error {
	csvMu.Lock()
	defer csvMu.Unlock()

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	return w.Write(rec)
}

func locHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}

	var loc Location
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	
	if loc.Lat == nil || loc.Lon == nil {
		http.Error(w, "lat and lon are required", http.StatusBadRequest)
		return
	}
	device := loc.DeviceID
	if device == "" {
		device = "unknown"
	}

	receivedISO := time.Now().UTC().Format(time.RFC3339)
	ts := ""
	if loc.Timestamp != 0 {
		ts = int64ToString(loc.Timestamp)
	}
	acc := floatPtrToString(loc.Accuracy)
	spd := floatPtrToString(loc.Speed)

	rec := []string{
		ts,
		receivedISO,
		device,
		float64ToString(*loc.Lat),
		float64ToString(*loc.Lon),
		acc,
		spd,
	}

	if err := appendCSV(rec); err != nil {
		log.Printf("write error: %v", err)
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func int64ToString(v int64) string { return json.Number(json.MarshalNumber(v)).String() }

func float64ToString(f float64) string {	
	return strconvFormatFloat(f)
}
func floatPtrToString(p *float64) string {
	if p == nil {
		return ""
	}
	return strconvFormatFloat(*p)
}
func strconvFormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'g', 15, 64)
}

func main() {
	if err := ensureCSV(); err != nil {
		log.Fatalf("failed to init csv: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/loc", locHandler)

	addr := ":5000"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		log.Fatal(err)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
