package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestData struct {
	Flights [][]string `json:"flights"`
	Source  string      `json:"source"`
}

func main() {

	e := echo.New()
	e.POST("/track", Track)

	fmt.Printf("Running app on port 8000")
	log.Fatal(http.ListenAndServe(":8000", e))
}

// Controller to Track user
func Track(c echo.Context) (err error) {
	ctx := c.Request().Context()
	timeoutContext := time.Duration(60) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeoutContext)
	defer cancel()

	jsonBody := make(map[string]interface{})
	if err := json.NewDecoder(c.Request().Body).Decode(&jsonBody); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	allFlights, ok := jsonBody["flights"].([]interface{})
	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, errors.New("error"))
	}
	flights := make([][]string, len(allFlights))
	for i, v := range allFlights {
		flights[i] = make([]string, 2)
		for y, x := range v.([]interface{}) {
			flights[i][y] = x.(string)
		}
	}

	source, ok := jsonBody["source"].(string)
	if !ok {
		return c.JSON(http.StatusUnprocessableEntity, errors.New("error"))
	}
	request := &RequestData{
		Flights: flights,
		Source:  source,
	}

	flightPaths := findRoutes(*request)
	fmt.Println(flightPaths)

	return c.JSON(http.StatusOK, flightPaths)
}

// Find Route function
func findRoutes(data RequestData) []string {
	m := createRouteMap(data.Flights)

	route := []string{}
	DFSearch(m, &route, data.Source)

	return reverseStrings(route)
}

// Depth Frist Search
func DFSearch(m map[string][]string, route *[]string, curr string) {
	for len(m[curr]) > 0 {
		next := m[curr][0]
		m[curr] = m[curr][1:]

		DFSearch(m, route, next)
	}

	*route = append(*route, curr)
}

// Create Route Map
func createRouteMap(tickets [][]string) map[string][]string {
	m := make(map[string][]string)

	for _, ticket := range tickets {
		to, from := ticket[0], ticket[1]
		m[to] = append(m[to], from)
	}

	for key, conns := range m {
		sort.Slice(m[key], func(i, j int) bool {
			return strings.Compare(conns[i], conns[j]) == -1
		})
	}

	return m
}

// Reverse Routes
func reverseStrings(s []string) []string {
	i, j := 0, len(s)-1
	for i < j {
		s[i], s[j] = s[j], s[i]

		i++
		j--
	}

	return s
}
