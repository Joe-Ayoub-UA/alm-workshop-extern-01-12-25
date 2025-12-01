package main

import (
	"encoding/json"
	"fmt" // Added for printing to console
	"net/http"
	"os" // Added to access environment variables
	"strconv" // Added for converting string to integer
)


func checkScoreRange(score int) bool {
	return score <= 10 && score >= 1
}

var defaultScore int = 5 // Fallback default score

var workshop *Workshop

func init() {
	// Read the environment variable named SWEATER_SCORE
	scoreStr := os.Getenv("SWEATER_SCORE")
	if scoreStr != "" {
		// Attempt to convert the environment variable string to an integer
		if score, err := strconv.Atoi(scoreStr); err == nil {
			// Validate the score range
			if checkScoreRange(score) {
				defaultScore = score
				fmt.Printf("Default SweaterScore set from environment variable to: %d\n", defaultScore)
			} else {
				fmt.Printf("Environment variable SWEATER_SCORE (%s) is out of range (1-10). Using fallback default of %d.\n", scoreStr, defaultScore)
			}
		} else {
			fmt.Printf("Environment variable SWEATER_SCORE ('%s') is not a valid number. Using fallback default of %d.\n", scoreStr, defaultScore)
		}
	} else {
		fmt.Printf("SWEATER_SCORE environment variable not set. Using default of %d.\n", defaultScore)
	}

	// FIX: Initialize the 'workshop' struct AFTER 'defaultScore' has been set
	// based on the environment variable.
	workshop = &Workshop{
		Name:         "ALM Workshop",
		Date:         "1/12/2025",
		Presentator:  "AE Consultants",
		Participants: []string{"John Doe", "Mary Little Lamb", "Chuck Norris", "Joe Ayoub"},
		SweaterScore: defaultScore,
	}
}

type Workshop struct {
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Presentator  string   `json:"presentator"`
	Participants []string `json:"participants"`
	SweaterScore int `json:"sweaterscore"`
}

func getWorkshopHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct to JSON and write it to the response
	json.NewEncoder(w).Encode(workshop)
}

func postWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON data into a new Workshop struct
	var newWorkshop Workshop
	err := json.NewDecoder(r.Body).Decode(&newWorkshop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON data"))
		return
	}

	if !checkScoreRange(newWorkshop.SweaterScore) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid SweaterScore: %d. Score must be between 1 and 10.", newWorkshop.SweaterScore)))
		return
	}

	// Update the workshop details
	*workshop = newWorkshop

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workshop)
}

func WorkshopHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getWorkshopHandler(w, r)
	case "POST":
		postWorkshopHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
