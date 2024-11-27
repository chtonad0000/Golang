package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Record struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type Medals struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type MedalsByYear struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type AthleteResponse struct {
	Athlete      string                  `json:"athlete"`
	Country      string                  `json:"country"`
	Medals       Medals                  `json:"medals"`
	MedalsByYear map[string]MedalsByYear `json:"medals_by_year"`
}

type CountryResponse struct {
	Country string `json:"country"`
	Medals
}

var dataPath string
var records []Record

func loadData() error {
	file, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			return
		}
	}(file)

	if err := json.NewDecoder(file).Decode(&records); err != nil {
		return err
	}
	return nil
}

func handleAthleteInfo(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "no parameter name", http.StatusBadRequest)
		return
	}

	athleteRecords := make([]Record, 0)
	for _, record := range records {
		if record.Athlete == name {
			athleteRecords = append(athleteRecords, record)
		}
	}

	if len(athleteRecords) == 0 {
		http.Error(w, "athlete not found", http.StatusNotFound)
		return
	}

	country := athleteRecords[0].Country
	medals := Medals{}
	medalsByYear := make(map[string]MedalsByYear)

	for _, record := range athleteRecords {
		medals.Gold += record.Gold
		medals.Silver += record.Silver
		medals.Bronze += record.Bronze
		medals.Total += record.Total

		year := strconv.Itoa(record.Year)
		if _, exists := medalsByYear[year]; !exists {
			medalsByYear[year] = MedalsByYear{}
		}
		medalsByYear[year] = MedalsByYear{
			Gold:   medalsByYear[year].Gold + record.Gold,
			Silver: medalsByYear[year].Silver + record.Silver,
			Bronze: medalsByYear[year].Bronze + record.Bronze,
			Total:  medalsByYear[year].Total + record.Total,
		}
	}

	response := AthleteResponse{
		Athlete:      name,
		Country:      country,
		Medals:       medals,
		MedalsByYear: medalsByYear,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func handleTopAthletesInSport(w http.ResponseWriter, r *http.Request) {
	sport := r.URL.Query().Get("sport")
	if sport == "" {
		http.Error(w, "sport parameter is required", http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 3
	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	athletes := make(map[string][]Record)
	for _, record := range records {
		if record.Sport == sport {
			athletes[record.Athlete] = append(athletes[record.Athlete], record)
		}
	}

	if len(athletes) == 0 {
		http.Error(w, fmt.Sprintf("sport '%s' not found", sport), http.StatusNotFound)
		return
	}

	type Athlete struct {
		Name         string
		Country      string
		Medals       Medals
		MedalsByYear map[string]MedalsByYear
	}

	allAthletes := make([]Athlete, 0)
	for name, records := range athletes {
		country := records[0].Country
		medals := Medals{}
		medalsByYear := make(map[string]MedalsByYear)

		for _, record := range records {
			medals.Gold += record.Gold
			medals.Silver += record.Silver
			medals.Bronze += record.Bronze
			medals.Total += record.Total

			year := strconv.Itoa(record.Year)
			if _, exists := medalsByYear[year]; !exists {
				medalsByYear[year] = MedalsByYear{}
			}
			medalsByYear[year] = MedalsByYear{
				Gold:   medalsByYear[year].Gold + record.Gold,
				Silver: medalsByYear[year].Silver + record.Silver,
				Bronze: medalsByYear[year].Bronze + record.Bronze,
				Total:  medalsByYear[year].Total + record.Total,
			}
		}

		allAthletes = append(allAthletes, Athlete{
			Name:         name,
			Country:      country,
			Medals:       medals,
			MedalsByYear: medalsByYear,
		})
	}

	sort.Slice(allAthletes, func(i, j int) bool {
		if allAthletes[i].Medals.Gold != allAthletes[j].Medals.Gold {
			return allAthletes[i].Medals.Gold > allAthletes[j].Medals.Gold
		}
		if allAthletes[i].Medals.Silver != allAthletes[j].Medals.Silver {
			return allAthletes[i].Medals.Silver > allAthletes[j].Medals.Silver
		}
		if allAthletes[i].Medals.Bronze != allAthletes[j].Medals.Bronze {
			return allAthletes[i].Medals.Bronze > allAthletes[j].Medals.Bronze
		}
		return allAthletes[i].Name < allAthletes[j].Name
	})

	response := make([]AthleteResponse, 0, limit)
	for i := 0; i < limit && i < len(allAthletes); i++ {
		a := allAthletes[i]
		response = append(response, AthleteResponse{
			Athlete:      a.Name,
			Country:      a.Country,
			Medals:       a.Medals,
			MedalsByYear: a.MedalsByYear,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func handleTopCountriesInYear(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	if yearParam == "" {
		http.Error(w, "year parameter is required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil || year <= 0 {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 3
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	countryMedals := make(map[string]Medals)
	for _, record := range records {
		if record.Year == year {
			medals := countryMedals[record.Country]
			medals.Gold += record.Gold
			medals.Silver += record.Silver
			medals.Bronze += record.Bronze
			medals.Total += record.Total
			countryMedals[record.Country] = medals
		}
	}

	if len(countryMedals) == 0 {
		http.Error(w, fmt.Sprintf("no records found for year %d", year), http.StatusNotFound)
		return
	}

	response := make([]CountryResponse, 0, len(countryMedals))
	for country, medals := range countryMedals {
		response = append(response, CountryResponse{
			Country: country,
			Medals:  medals,
		})
	}

	sort.Slice(response, func(i, j int) bool {
		if response[i].Gold != response[j].Gold {
			return response[i].Gold > response[j].Gold
		}
		if response[i].Silver != response[j].Silver {
			return response[i].Silver > response[j].Silver
		}
		if response[i].Bronze != response[j].Bronze {
			return response[i].Bronze > response[j].Bronze
		}
		return response[i].Country < response[j].Country
	})

	if len(response) > limit {
		response = response[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	errEnc := json.NewEncoder(w).Encode(response)
	if errEnc != nil {
		return
	}
}

func main() {
	port := flag.Int("port", 8080, "port")
	data := flag.String("data", "./testdata/olympicWinners.json", "file")
	flag.Parse()

	dataPath = *data
	if err := loadData(); err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "failed to load data: %v\n", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

	http.HandleFunc("/athlete-info", handleAthleteInfo)
	http.HandleFunc("/top-athletes-in-sport", handleTopAthletesInSport)
	http.HandleFunc("/top-countries-in-year", handleTopCountriesInYear)

	fmt.Printf("Server listening on port %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		return
	}
}
