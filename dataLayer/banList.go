package dataLayer

import (
	"code/Hahachitchat/definition"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func LoadForbiddenConfig() {
	jsonFile, err := os.Open("./definition/forbidden.json")
	if err != nil {
		panic(fmt.Sprintf("jsonFile os.Open err: %v", err))
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &definition.ForbiddenConfig); err != nil {
		panic(fmt.Sprintf("jsonFile Unmarshal err: %v", err))
	}
	jsonFile.Close()
	time.Sleep(5 * time.Second)

	for {
		jsonFile, err := os.Open("./definition/forbidden.json")
		if err != nil {
			Serverlog.Fatalln("jsonFile os.Open err: ", err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var newForbiddenConfig definition.Forbidden
		if err := json.Unmarshal(byteValue, &newForbiddenConfig); err != nil {
			Serverlog.Fatalln("jsonFile Unmarshal err: ", err)
		}
		jsonFile.Close()
		definition.ForbiddenConfig = newForbiddenConfig
		time.Sleep(5 * time.Second)
	}
}