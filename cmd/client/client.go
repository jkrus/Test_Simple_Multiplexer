package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jkrus/Test_Simple_Multuplexor/pkg/models"
)

func main() {
	data := models.InfoDTO{Urls: []string{"https://ya.ru", "https://ya.ru"}}
	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}
	mData := bytes.NewReader(marshal)

	resp, err := http.Post("http://localhost:8080/api/v1/info", "application/json", mData)
	if err != nil {
		log.Println("ERR = ", err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERR = ", err)
	}

	var urls []models.Info
	if err = json.Unmarshal(bytes, &urls); err != nil {
		log.Println("ERR = ", err)
	}

	log.Println(urls)
}
