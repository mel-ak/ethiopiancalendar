package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	ethiopiancalendar "github.com/mel-ak/ethiopiancalendar/pkg"
)

// APIRequest structs for JSON payloads
type ConvertRequest struct {
	Type  string `json:"type"` // "etToGreg" or "gregToEt"
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
}

type FormatRequest struct {
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Day    int    `json:"day"`
	Layout string `json:"layout"`
}

type ArithmeticRequest struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Operation string `json:"operation"` // "days", "months", "years"
	Value     int    `json:"value"`
}

type LeapRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

// APIResponse for all endpoints
type APIResponse struct {
	Year        int    `json:"year,omitempty"`
	Month       int    `json:"month,omitempty"`
	Day         int    `json:"day,omitempty"`
	Result      string `json:"result,omitempty"`
	IsLeap      bool   `json:"isLeap,omitempty"`
	DaysInMonth *int   `json:"daysInMonth,omitempty"`
	Error       string `json:"error,omitempty"`
}

func main() {
	// Serve static index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "index.html"))
	})

	// Conversion endpoint
	http.HandleFunc("/api/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ConvertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid JSON")
			return
		}
		var resp APIResponse
		if req.Type == "etToGreg" {
			date := ethiopiancalendar.EtDate{Year: req.Year, Month: req.Month, Day: req.Day}
			if err := date.Validate(); err != nil {
				sendError(w, err.Error())
				return
			}
			gy, gm, gd, err := date.ToGregorian()
			if err != nil {
				sendError(w, err.Error())
				return
			}
			resp = APIResponse{Year: gy, Month: gm, Day: gd}
		} else if req.Type == "gregToEt" {
			date, err := ethiopiancalendar.FromGregorian(req.Year, req.Month, req.Day)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			resp = APIResponse{Year: date.Year, Month: date.Month, Day: date.Day}
		} else {
			sendError(w, "Invalid conversion type")
			return
		}
		sendJSON(w, resp)
	})

	// Format endpoint
	http.HandleFunc("/api/format", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req FormatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid JSON")
			return
		}
		date := ethiopiancalendar.EtDate{Year: req.Year, Month: req.Month, Day: req.Day}
		if err := date.Validate(); err != nil {
			sendError(w, err.Error())
			return
		}
		result := date.Format(req.Layout)
		sendJSON(w, APIResponse{Result: result})
	})

	// Arithmetic endpoint
	http.HandleFunc("/api/arithmetic", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ArithmeticRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid JSON")
			return
		}
		date := ethiopiancalendar.EtDate{Year: req.Year, Month: req.Month, Day: req.Day}
		if err := date.Validate(); err != nil {
			sendError(w, err.Error())
			return
		}
		var newDate ethiopiancalendar.EtDate
		var err error
		switch req.Operation {
		case "days":
			newDate, err = date.AddDays(req.Value)
		case "months":
			newDate = date.AddMonths(req.Value)
		case "years":
			newDate = date.AddYears(req.Value)
		default:
			sendError(w, "Invalid operation")
			return
		}
		if err != nil {
			sendError(w, err.Error())
			return
		}
		sendJSON(w, APIResponse{Year: newDate.Year, Month: newDate.Month, Day: newDate.Day})
	})

	// Leap year and days in month endpoint
	http.HandleFunc("/api/leap", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req LeapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "Invalid JSON")
			return
		}
		if req.Year <= 0 {
			sendError(w, "Year must be positive")
			return
		}
		resp := APIResponse{IsLeap: ethiopiancalendar.IsLeap(req.Year)}
		if req.Month != 0 {
			days := ethiopiancalendar.DaysInMonth(req.Year, req.Month)
			if days == 0 {
				sendError(w, "Invalid month")
				return
			}
			resp.DaysInMonth = &days
		}
		sendJSON(w, resp)
	})

	// Current Ethiopian Date (optional, for future expansion)
	http.HandleFunc("/api/current", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		now := time.Now()
		et, err := ethiopiancalendar.FromGregorian(now.Year(), int(now.Month()), now.Day())
		if err != nil {
			sendError(w, err.Error())
			return
		}
		sendJSON(w, APIResponse{Year: et.Year, Month: et.Month, Day: et.Day})
	})

	fmt.Println("Server starting at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Helper functions
func sendError(w http.ResponseWriter, msg string) {
	sendJSON(w, APIResponse{Error: msg})
}

func sendJSON(w http.ResponseWriter, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
