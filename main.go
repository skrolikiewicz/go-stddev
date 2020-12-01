package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type CalculateMeanQuery struct {
	Requests int `schema:"requests"`
	Length   int `schema:"length"`
}

type CalculationResult struct {
	StdDev float64 `json:"stddev"`
	Data   []int   `json:"data"`
}

func main() {
	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/random/mean", meanHandler)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func meanHandler(w http.ResponseWriter, r *http.Request) {
	var calculateMeanParams CalculateMeanQuery
	decoder := schema.NewDecoder()
	err := decoder.Decode(&calculateMeanParams, r.URL.Query())
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if calculateMeanParams.Length <= 0 {
		http.Error(w, "Length must be greater than 0", http.StatusBadRequest)
		return
	}

	if calculateMeanParams.Length > 10000 {
		http.Error(w, "Length must be less than or equal 10000", http.StatusBadRequest)
		return
	}

	if calculateMeanParams.Requests <= 0 {
		http.Error(w, "Requests must be greater than 0", http.StatusBadRequest)
		return
	}

	if calculateMeanParams.Requests > 100 {
		http.Error(w, "Requests must be less than or equal 100", http.StatusBadRequest)
		return
	}

	randomNumbersSets, err := retrieveRandomNumbersSets(r.Context(), calculateMeanParams)
	if err != nil {
		http.Error(w, "Error while retrieving data from random.org API. "+err.Error(), http.StatusInternalServerError)
		return
	}

	var appendedSets []int
	result := make([]CalculationResult, calculateMeanParams.Requests+1)
	for i, set := range randomNumbersSets {
		result[i] = CalculationResult{StdDev: calculateStdDev(set), Data: set}
		appendedSets = append(appendedSets, set...)
	}
	result[calculateMeanParams.Requests] = CalculationResult{StdDev: calculateStdDev(appendedSets), Data: appendedSets}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func retrieveRandomNumbersSets(ctx context.Context, calculateMeanParams CalculateMeanQuery) (randomNumbersSets [][]int, err error) {
	errors := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(calculateMeanParams.Requests)

	go func() {
		wg.Wait()
		close(errors)
	}()

	ctx, cancel := context.WithCancel(ctx)
	randomNumbersSets = make([][]int, calculateMeanParams.Requests)
	for i := 0; i < calculateMeanParams.Requests; i++ {
		go func(ctx context.Context, i int) {
			defer wg.Done()

			result, err := getRandomNumbers(ctx, calculateMeanParams.Length)
			if err != nil {
				defer cancel()
				errors <- err
				return
			}
			randomNumbersSets[i] = result
		}(ctx, i)
	}

	err = <-errors
	if err != nil {
		return nil, err
	}

	return randomNumbersSets, nil
}
