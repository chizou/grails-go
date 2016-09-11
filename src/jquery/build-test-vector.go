package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type testVector struct {
	externalUserService string
	piazzaPrimeBox      string
	id1                 string
	id2                 string
	id3                 string
	id4                 string
}

func newTestVector(piazzaBox string, externalUserService string) *testVector {
	p := new(testVector)
	p.externalUserService = "http://" + externalUserService + ":8076/external"
	p.piazzaPrimeBox = piazzaBox
	p.id1 = ""
	p.id2 = ""
	p.id3 = ""
	p.id4 = ""
	return p
}

func nextStep(p testVector) *testVector {
	if p.id1 == "" {
		p.id1 = pz1(p)
		if p.id1 == "" {
			return &p
		}
	}

	if p.id2 == "" {
		p.id2 = pz2(p)
		if p.id2 == "" {
			return &p
		}
	}

	if p.id3 == "" {
		p.id3 = pz3(p)
		if p.id3 == "" {
			return &p
		}
	}

	if p.id4 == "" {
		p.id4 = pz4(p)
		if p.id4 == "" {
			return &p
		}
	}

	return &p
}

func pz1(p testVector) string {
	body2 := `{"url":"REPLACEME","method":"GET","contractUrl":"REPLACEME/","resourceMetadata":{"name":"Hello World Test","description":"Hello world test","classType":{"classification":"UNCLASSIFIED"}}}`
	body2 = strings.Replace(body2, "REPLACEME", p.externalUserService, -1)
	pz1curlresult := pz1curl(p, body2)
	return pz1curlresult
}

func pz1curl(p testVector, body2 string) string {

	body := strings.NewReader(body2)
	req, err := http.NewRequest("POST", "http://"+p.piazzaPrimeBox+":8081/service", body)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body3, err8 := ioutil.ReadAll(resp.Body)
	if err8 != nil {
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal(body3, &dat); err3 != nil {
		panic(err3)
	}
	data := dat["data"].(map[string]interface{})
	foo := data["serviceId"]
	return foo.(string)
}

func pz2(p testVector) string {
	body2 := `{"type":"execute-service","data":{"serviceId":"REPLACEME","dataInputs":{},"dataOutput":[{"mimeType":"application/json","type":"text"}]}}`
	body2 = strings.Replace(body2, "REPLACEME", p.id1, -1)
	pz2curlresult := pz2curl(p, body2)
	return pz2curlresult
}

func pz2curl(p testVector, body2 string) string {
	body := strings.NewReader(body2)
	req, err := http.NewRequest("POST", "http://"+p.piazzaPrimeBox+":8081/job", body)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body3, err8 := ioutil.ReadAll(resp.Body)
	if err8 != nil {
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal(body3, &dat); err3 != nil {
		panic(err3)
	}
	data := dat["data"].(map[string]interface{})
	foo := data["jobId"]
	return foo.(string)
}

func pz3(p testVector) string {
	body2 := `{"type":"execute-service","data":{"serviceId":"REPLACEME","dataInputs":{},"dataOutput":[{"mimeType":"application/json","type":"text"}]}}`
	body2 = strings.Replace(body2, "REPLACEME", p.id2, -1)
	pz3curlresult := pz3curl(p, body2)
	return pz3curlresult
}

func pz3curl(p testVector, body2 string) string {
	req, err := http.NewRequest("GET", "http://"+p.piazzaPrimeBox+":8081/job/"+p.id2, nil)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body3, err8 := ioutil.ReadAll(resp.Body)
	if err8 != nil {
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal(body3, &dat); err3 != nil {
		panic(err3)
	}
	data := dat["data"].(map[string]interface{})
	if data["result"] == nil {
		return ""
	}

	result := data["result"].(map[string]interface{})
	foo := result["dataId"]
	return foo.(string)
}

func pz4(p testVector) string {
	body2 := `{"type":"execute-service","data":{"serviceId":"REPLACEME","dataInputs":{},"dataOutput":[{"mimeType":"application/json","type":"text"}]}}`
	body2 = strings.Replace(body2, "REPLACEME", p.id3, -1)
	pz4curlresult := pz4curl(p, body2)
	return pz4curlresult
}

func pz4curl(p testVector, body2 string) string {
	temp := "http://" + p.piazzaPrimeBox + ":8081/data/" + p.id3
	req, err := http.NewRequest("GET", temp, nil)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body3, err8 := ioutil.ReadAll(resp.Body)
	if err8 != nil {
	}
	var dat map[string]interface{}
	if err3 := json.Unmarshal(body3, &dat); err3 != nil {
		panic(err3)
	}
	if dat["data"] == nil {
		return ""
	}
	data := dat["data"].(map[string]interface{})
	if data["dataType"] == nil {
		fmt.Println("pz4curl The dataType is nil")
		return ""
	}

	result := data["dataType"].(map[string]interface{})
	foo := result["content"]
	return foo.(string)
}
