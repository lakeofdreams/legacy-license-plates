package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func floattostr(input_num float64) string {
    return strconv.FormatFloat(input_num, 'g', 1, 64)
 }

type Plate struct {
	PlateId             float64 `json:"plateId"`
	State               string `json:"state"`
	PlateNumber         string `json:"plateNumber"`
	Owner               string `json:"owner"`
	Address             string `json:"address"`
	Timestamp           string `json:"ts"`
	ImageURL            string `json:"imageURL"`
}

var plate1 = Plate{PlateNumber:"NY 123 AB", Owner:"Mark", State:"NY", Address:"5th Street", ImageURL:"http://google.com"}
var plate2 = Plate{PlateNumber:"NY 124 AB", Owner:"Mark", State:"NY", Address:"6th Street", ImageURL:"http://yahoo.com"}
var plates = []Plate{plate1,plate2}
var plate3 = Plate{PlateNumber:"NY 125 AB", Owner:"Mark", State:"NY", Address:"7th Street", ImageURL:"http://uber.com"}

var ADD_ALL = "add-all";
var LIST_ALL = "list";
var ADD_SINGLE = "add";
var GET_BY_ID = "get/{id}";
var GET_LATEST = "get-latest";
var DELETE_ALL = "delete-all";
var GET_BY_PLATE_NUM = "get-by-plate-number/{plate-number}";

var base_url = "http://weblogic:7001/licenseplates/rest/plates";

var statuses map[string]string = make(map[string]string)
var STATUS_SUCCESS = "Passed"
var STATUS_FAILED = "Failed"


func main() {

	testAddAllPlates()

	testListAllPlates()

	fmt.Println(statuses)

	/* ==WIP==
	var platesData map[string]interface{};

	fmt.Println("Testing listPlates()  http://weblogic:7001/licenseplates/rest/plates/list")
	response, err := http.Get("http://weblogic:7001/licenseplates/rest/plates/list")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		fmt.Println("http://weblogic:7001/licenseplates/rest/plates/list succeeded")
	}

	fmt.Println("Testing addPlate  http://weblogic:7001/licenseplates/rest/plates/add")
	jsonData := map[string]string{"state": "NY", "plateNumber": "NY 123 AB", "owner": "Mark", "address": "NY", "imageURL": "http://google.com"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err = http.Post("http://weblogic:7001/licenseplates/rest/plates/add", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	fmt.Println("Testing getLatestPlate()  http://weblogic:7001/licenseplates/rest/plates/get-latest")
	response, err = http.Get("http://weblogic:7001/licenseplates/rest/plates/get-latest")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		json.Unmarshal(data, &platesData)
		fmt.Println("http://weblogic:7001/licenseplates/rest/plates/get-latest succeeded and returned plateId ", platesData["plateId"])
	}
	
	fmt.Println("Testing getPlate() by Id  http://weblogic:7001/licenseplates/rest/plates/get/",platesData["plateId"])
	response, err = http.Get("http://weblogic:7001/licenseplates/rest/plates/get/"+floattostr(platesData["plateId"].(float64)))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		json.Unmarshal(data, &platesData)
		fmt.Println("http://weblogic:7001/licenseplates/rest/plates/get/"+ floattostr(platesData["plateId"].(float64)) + " succeeded")
	}
	
*/
}
func testAddAllPlates() {
	fmt.Println("Start test: "+ADD_ALL)
	jsonData, _ := json.Marshal(plates)
	fmt.Println("Sending "+string(jsonData))
	response := invokeRest(ADD_ALL,jsonData)
	fmt.Println("Response: " + string(response))
	statuses[ADD_ALL] = STATUS_SUCCESS
	fmt.Println("Finish test: "+ADD_ALL)
}

func testListAllPlates() {
	var plate1Found bool
	var plate2Found bool

	fmt.Println("Start test: "+LIST_ALL)
	response := invokeRest(LIST_ALL,nil)
	fmt.Println(response)
	fmt.Println("Response: " + string(response))
	plates = []Plate{}
	json.Unmarshal(response, &plates)
	for _, v := range plates {
		if v.PlateNumber == plate1.PlateNumber {
			plate1Found = true
		}
		if v.PlateNumber == plate2.PlateNumber {
			plate2Found = true
		}
	}
	if !plate1Found || !plate2Found {
		fmt.Println("One of the plates is missing ")
		statuses[LIST_ALL] = STATUS_FAILED
	} else {
		statuses[LIST_ALL] = STATUS_SUCCESS
	}
	fmt.Println("Finish test: "+LIST_ALL)
}

func invokeRest(api string, data []byte) ([]byte){

	var output [] byte
	var url string

	switch api {
	case ADD_ALL:
		url = base_url+"/"+ADD_ALL
		fmt.Println("Invoke "+url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		output =  processResponse(api, response, err)
	case LIST_ALL:
		url = base_url+"/"+LIST_ALL
		fmt.Println("Invoke "+url)
		response, err := http.Get(url)
		output =  processResponse(api, response, err)
	case ADD_SINGLE:
		url = base_url+"/"+ADD_SINGLE
		fmt.Println("Invoke "+url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		output =  processResponse(api, response, err)
	case GET_BY_ID:
		url = base_url+"/"+strings.Replace(GET_BY_ID,"{id}",string(data),-1)
		fmt.Println("Invoke "+url)
		response, err := http.Get(url)
		output =  processResponse(api, response, err)
	case GET_LATEST:
		url = base_url+"/"+GET_LATEST
		fmt.Println("Invoke "+url)
		response, err := http.Get(url)
		output =  processResponse(api, response, err)
	case DELETE_ALL:
		url = base_url+"/"+DELETE_ALL
		fmt.Println("Invoke "+url)
		response, err := http.Get(url)
		output =  processResponse(api, response, err)
	case GET_BY_PLATE_NUM:
		url = base_url+"/"+strings.Replace(GET_BY_PLATE_NUM,"{plate-number}",string(data),-1)
		fmt.Println("Invoke "+url)
		response, err := http.Get(url)
		output =  processResponse(api, response, err)
	}

	return output
}
func processResponse(api string, response *http.Response, err error) []byte {
	var data []byte
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		statuses[api] = STATUS_FAILED
	} else {
		data, _ = ioutil.ReadAll(response.Body)
	}
	defer  response.Body.Close()
	return data
}