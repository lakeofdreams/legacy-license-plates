package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func floattostr(input_num float64) string {
    return strconv.FormatFloat(input_num, 'g', 1, 64)
 }


func main() {

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
	

}
