package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func myIPWithTimeout() string {
        timeout := time.Duration(3 * time.Second)
        client := http.Client{
                Timeout: timeout,
        }
        resp, err := client.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
        if err != nil {
                fmt.Println("Something went wrong in myIPWithTimeout")
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        return string(body)
}

func simulatedExternalUserService(w http.ResponseWriter, r *http.Request) {
	t := simulateComputation()
	ip := myIPWithTimeout()
	writeReturnedJSON(w, t, ip)
}

func simulateComputation() int {
	//rand.Seed(time.Now().Unix())
	const maxSimulatedComputation = 16
	t := rand.Intn(maxSimulatedComputation)
	time.Sleep(time.Second * time.Duration(t))
	return t
}

type returnedJSON struct {
	Results int
	Status  string
}

func writeReturnedJSON(w http.ResponseWriter, t int, ip string) {
	p := returnedJSON{t, ip}
	b, err := json.Marshal(p)
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
        mux := http.NewServeMux()
//        mux.HandleFunc("/", home)
//        mux.HandleFunc("/work", work)
//        mux.HandleFunc("/status", status)
        mux.HandleFunc("/external", simulatedExternalUserService)
        http.ListenAndServe(":8077", mux)
}

