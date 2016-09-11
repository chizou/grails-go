package main

import (
	"encoding/json"
	"strings"
)

func booleanOfDotCompletion() bool {
	return stringOfDotStatusEachRepresentsAPiazzaJob() == strings.Repeat("4", numColorfulDisplayDots)
}

func stringOfDotStatusEachRepresentsAPiazzaJob() string {
	s := ""
	for i := 0; i < numColorfulDisplayDots; i++ {
		s += testVectorStatus(q[i])
	}
	return s
}

func testVectorStatus(p *testVector) string {
	if p.id4 != "" {
		return "4"
	}
	if p.id3 != "" {
		return "3"
	}
	if p.id2 != "" {
		return "2"
	}
	if p.id1 != "" {
		return "1"
	}
	return "0"
}
func stringOfDotDurationEachRepresentsAPiazzaJob() string {
	s := ""
	for i := 0; i < numColorfulDisplayDots; i++ {
		s += parseResultsFieldFromID4(q[i].id4)
	}
	return s
}
func parseResultsFieldFromID4(id4 string) string {
	temp := "5"
	if id4 != "" {
		var dat map[string]interface{}
		if err3 := json.Unmarshal([]byte(id4), &dat); err3 != nil {
			panic(err3)
		}
		if dat["Results"] != nil {
			duration := dat["Results"].(float64)
			if duration <= 16 {
				temp = "3"
			}
			if duration <= 12 {
				temp = "2"
			}
			if duration <= 8 {
				temp = "1"
			}
			if duration <= 4 {
				temp = "0"
			}
		}
	}
	return temp
}
func stringOfDotColor() string {
	s := ""
	for i := 0; i < numColorfulDisplayDots; i++ {
		ip := extractIPFromID4ResultsString(q[i].id4)
		s += translateIPToColor(ip)
	}
	return s
}
func extractIPFromID4ResultsString(s string) string {
	//Temporary function to read a field (Status) from a JSON string, e.g.
	//  {"Results":1,"Status":"52.88.226.0"}
	//Long term would definitely prefer to receive a JSON object instead.
	var r string

	ix := strings.Index(s, `"Status":"`)
	if ix != -1 {
		r = s[ix+len(`"Status":"`) : len(s)-len(`"}`)]
	}
	return r
}

func translateIPToColor(s string) string {
	if len(eusips) == 2 {
		if s == eusips[0] {
			return "G"
		}
		if s == eusips[1] {
			return "B"
		}
	}
	if len(eusips) == 1 {
		if s == eusips[0] {
			return "G"
		}
	}
	return "X"
}

func stringOfDotCompletion() string {
	if booleanOfDotCompletion() {
		return "active... completed"
	}
	return "active"
}
