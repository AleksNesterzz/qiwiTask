package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type Valute struct {
	ID       string `xml:",attr"`
	NumCode  int
	CharCode string
	Nominal  int
	Name     string
	Value    float64
}

type ValCurs struct {
	Date string `xml:",attr"`
	//name    string   `xml:",attr"`
	Valutes []Valute `xml:"Valute"`
}

func getXML(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Chrome/23.0.1271.64")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	utf8, err := charset.NewReader(response.Body, response.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func convertXML(body []byte) string {
	s := string(body)
	s = strings.Replace(s, "<?xml version=\"1.0\" encoding=\"windows-1251\"?>", "", -1)
	s = strings.Replace(s, "<?xml version=\"1.0\" encoding=\"windows-1252\"?>", "", -1)
	s = strings.Replace(s, ",", ".", -1)
	return s
}

func main() {
	//r vals []Valute
	var vals *ValCurs
	codeFlag := flag.String("code", "USD", "Краткое именование валюты")
	dateFlag := flag.String("date", time.Now().GoString(), "Дата, с которой смотрим цену")
	flag.Parse()

	code := *codeFlag
	date := *dateFlag

	date_parts := strings.Split(date, "-")
	date = date_parts[2] + "/" + date_parts[1] + "/" + date_parts[0]

	url := "https://www.cbr.ru/scripts/XML_daily.asp?date_req=" + date

	body, err := getXML(url)
	if err != nil {
		fmt.Println("Error getting XML:")
		return
	}

	s := convertXML(body)

	reader := strings.NewReader(s)
	if err = xml.NewDecoder(reader).Decode(&vals); err != nil {
		fmt.Println(err)
		return
	}

	for _, value := range vals.Valutes {
		if value.CharCode == code {
			fmt.Println(code, "(", value.Name, "):", value.Value/float64(value.Nominal))
			return
		}
	}
}
