package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"strconv"
)


type Plate struct {
	PlateId     int     `json:"plateId"`
	State       string  `json:"state"`
	PlateNumber string  `json:"plateNumber"`
	Owner       string  `json:"owner"`
	Address     string  `json:"address"`
	Timestamp   string  `json:"ts"`
	ImageURL    string  `json:"imageURL"`
}

var plate1 = Plate{PlateNumber: "NY 123 AB", Owner: "Mark", State: "NY", Address: "5th Street", ImageURL: "http://google.com"}
var plate2 = Plate{PlateNumber: "NY 124 AB", Owner: "Mark", State: "NY", Address: "6th Street", ImageURL: "http://yahoo.com"}
var plates = []Plate{plate1, plate2}
var plate3 = Plate{PlateNumber: "NY 125 AB", Owner: "Mark", State: "NY", Address: "7th Street", ImageURL: "http://uber.com"}

var ADD_ALL = "add-all"
var LIST_ALL = "list"
var ADD_SINGLE = "add"
var GET_BY_ID = "get/{id}"
var GET_LATEST = "get-latest"
var DELETE_ALL = "delete-all"
var GET_BY_PLATE_NUM = "get-by-plate-number/{plate-number}"

var hostAddress = os.Getenv("WEBLOGIC_HOST_TEST")
var base_url = "http://" + hostAddress + ":7001/licenseplates/rest/plates"

var statuses map[string]string = make(map[string]string)
var STATUS_SUCCESS = "Passed"
var STATUS_FAILED = "Failed"


func main() {

	invokeRest(DELETE_ALL, nil)

	testAddAllPlates()

	testListAllPlates()

	testGetPlateByPlateNum()

	testDeleteAllPlates()

	testAddSinglePlate()

	testGetLatestPlate()

	testGetPlateById()

	fmt.Println(statuses)

	invokeRest(DELETE_ALL, nil)
}

func testGetPlateById() {
	fmt.Println("Start test: " + GET_BY_ID)
	response := invokeRest(GET_BY_ID, []byte(strconv.Itoa(plate3.PlateId)))
	if statuses[GET_BY_ID] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plate := Plate{}
		json.Unmarshal(response, &plate)
		if plate.PlateId == plate3.PlateId {
			statuses[GET_BY_ID] = STATUS_SUCCESS
			fmt.Println("Finish test: " + GET_BY_ID)
		} else {
			fmt.Println("Could not find plate with plate id: " + strconv.Itoa(plate3.PlateId))
			statuses[GET_BY_ID] = STATUS_FAILED
			fmt.Println("Failed test: " + GET_BY_ID)
		}
	} else {
		fmt.Println("Failed test: " + GET_BY_ID)
	}

}

func testGetLatestPlate() {
	fmt.Println("Start test: " + GET_LATEST)
	response := invokeRest(GET_LATEST, nil)
	if statuses[GET_LATEST] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plate := Plate{}
		json.Unmarshal(response, &plate)
		if strings.EqualFold(plate.PlateNumber, plate3.PlateNumber) {
			statuses[GET_LATEST] = STATUS_SUCCESS
			plate3.PlateId = plate.PlateId
			fmt.Println("Finish test: " + GET_LATEST)
		} else {
			fmt.Println("Could not find plate added latest with plate num: " + plate3.PlateNumber)
			statuses[GET_LATEST] = STATUS_FAILED
			fmt.Println("Failed test: " + GET_LATEST)
		}
	} else {
		fmt.Println("Failed test: " + GET_LATEST)
	}

}

func testAddSinglePlate() {
	fmt.Println("Start test: " + ADD_SINGLE)
	jsonData, _ := json.Marshal(plate3)
	fmt.Println("Sending " + string(jsonData))
	response := invokeRest(ADD_SINGLE, jsonData)
	if statuses[ADD_SINGLE] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plates = []Plate{}
		json.Unmarshal(response, &plates)
		if len(plates) != 1 {
			statuses[ADD_SINGLE] = STATUS_FAILED
			fmt.Println("Expected 1 plate to be created, found : " + string(len(plates)))
		} else {
			statuses[ADD_SINGLE] = STATUS_SUCCESS
			fmt.Println("Finish test: " + ADD_SINGLE)
		}
	} else {
		fmt.Println("Failed test: " + ADD_SINGLE)
	}

}

func testDeleteAllPlates() {
	fmt.Println("Start test: " + DELETE_ALL)
	response := invokeRest(DELETE_ALL, nil)
	if statuses[DELETE_ALL] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plates = []Plate{}
		json.Unmarshal(response, &plates)
		if len(plates) != 0 {
			statuses[DELETE_ALL] = STATUS_FAILED
			fmt.Println("Expected no plates remaining, found : " + string(len(plates)))
		} else {
			statuses[DELETE_ALL] = STATUS_SUCCESS
			fmt.Println("Finish test: " + DELETE_ALL)
		}

	} else {
		fmt.Println("Failed test: " + DELETE_ALL)
	}

}

func testGetPlateByPlateNum() {
	fmt.Println("Start test: " + GET_BY_PLATE_NUM)
	fmt.Println("Query for plate with plate num: " + plate1.PlateNumber)
	response := invokeRest(GET_BY_PLATE_NUM, []byte(plate1.PlateNumber))
	if statuses[GET_BY_PLATE_NUM] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plates = []Plate{}
		json.Unmarshal(response, &plates)
		plate := plates[0]
		if strings.EqualFold(plate.PlateNumber, plate1.PlateNumber) {
			statuses[GET_BY_PLATE_NUM] = STATUS_SUCCESS
			fmt.Println("Finish test: " + GET_BY_PLATE_NUM)
		} else {
			fmt.Println("Could not find plate with plate num: " + plate1.PlateNumber + " , returned plate with plate num: " + plate.PlateNumber)
			statuses[GET_BY_PLATE_NUM] = STATUS_FAILED
			fmt.Println("Failed test: " + GET_BY_PLATE_NUM)
		}
	} else {
		fmt.Println("Failed test: " + GET_BY_PLATE_NUM)
	}

}

func testAddAllPlates() {
	fmt.Println("Start test: " + ADD_ALL)
	jsonData, _ := json.Marshal(plates)
	fmt.Println("Sending " + string(jsonData))
	response := invokeRest(ADD_ALL, jsonData)
	if statuses[ADD_ALL] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plates = []Plate{}
		json.Unmarshal(response, &plates)
		if len(plates) == 2 {
			statuses[ADD_ALL] = STATUS_SUCCESS
			fmt.Println("Finish test: " + ADD_ALL)
		} else {
			fmt.Println("Expected 2 plates to be created, found " + string(len(plates)))
			statuses[ADD_ALL] = STATUS_FAILED
			fmt.Println("Failed test: " + ADD_ALL)
		}
	} else {
		fmt.Println("Failed test: " + ADD_ALL)
	}

}

func testListAllPlates() {
	var plate1Found bool
	var plate2Found bool

	fmt.Println("Start test: " + LIST_ALL)
	response := invokeRest(LIST_ALL, nil)
	if statuses[LIST_ALL] != STATUS_FAILED {
		fmt.Println("Response: " + string(response))
		plates = []Plate{}
		json.Unmarshal(response, &plates)
		if len(plates) != 2 {
			statuses[LIST_ALL] = STATUS_FAILED
			fmt.Println("Expected two plates, found : " + string(len(plates)))
		}
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
		fmt.Println("Finish test: " + LIST_ALL)
	} else {
		fmt.Println("Failed test: " + LIST_ALL)
	}

}

func invokeRest(api string, data []byte) []byte {

	var output []byte
	var url string

	switch api {
	case ADD_ALL:
		url = base_url + "/" + ADD_ALL
		fmt.Println("Invoke " + url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		output = processResponse(api, response, err)
	case LIST_ALL:
		url = base_url + "/" + LIST_ALL
		fmt.Println("Invoke " + url)
		response, err := http.Get(url)
		output = processResponse(api, response, err)
	case ADD_SINGLE:
		url = base_url + "/" + ADD_SINGLE
		fmt.Println("Invoke " + url)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		output = processResponse(api, response, err)
	case GET_BY_ID:
		url = base_url + "/" + strings.Replace(GET_BY_ID, "{id}", string(data), -1)
		fmt.Println("Invoke " + url)
		response, err := http.Get(url)
		output = processResponse(api, response, err)
	case GET_LATEST:
		url = base_url + "/" + GET_LATEST
		fmt.Println("Invoke " + url)
		response, err := http.Get(url)
		output = processResponse(api, response, err)
	case DELETE_ALL:
		url = base_url + "/" + DELETE_ALL
		fmt.Println("Invoke " + url)
		response, err := http.Get(url)
		output = processResponse(api, response, err)
	case GET_BY_PLATE_NUM:
		url = base_url + "/" + strings.Replace(GET_BY_PLATE_NUM, "{plate-number}", string(data), -1)
		fmt.Println("Invoke " + url)
		response, err := http.Get(url)
		output = processResponse(api, response, err)
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
	defer response.Body.Close()
	return data
}
