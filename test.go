package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Distro struct {
	Region   string `json:"region"`
	Snapshot string `json:"snapshot"`
}

type InvokeResponse1 struct {
	Outputs     map[string]interface{}
	Logs        []string
	Distro		interface{}
}

type InvokeRequest1 struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

func simpleHttpTriggerHandlerEventHubOutDistro(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest1
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("The JSON data is:invokeReq metadata......")
	fmt.Println(invokeReq.Metadata)
	fmt.Println("The JSON data is:invokeReq data......")
	fmt.Println(invokeReq.Data)

	var dis Distro
	dis.Region = "uksouth"
	dis.Snapshot = "1.1.1.11111"

	//returnValue := "HelloWorld"
	invokeResponse1 := InvokeResponse1{Logs: []string{"test log1", "test log2"}, Distro: dis}

	js, err := json.Marshal(invokeResponse1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(js)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/HttpTrigger", simpleHttpTriggerHandlerEventHubOutDistro)
}
