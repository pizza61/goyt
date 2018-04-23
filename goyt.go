package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func main() {
	apikey := "yourApiKey"
	fmt.Printf("Szukaj po nazwie: ")
	query, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	queryurl := `https://www.googleapis.com/youtube/v3/search`
	req, err := http.NewRequest("GET", queryurl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	qs := req.URL.Query()
	qs.Add("key", apikey)
	qs.Add("part", "snippet,id")
	qs.Add("type", "channel")
	qs.Add("order", "viewCount")
	qs.Add("q", query)
	req.URL.RawQuery = qs.Encode()

	response, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println("Wystąpił błąd podczas wysyłania zapytania do API, czy na pewno masz dostęp do internetu?")
		os.Exit(1)
	}
	defer response.Body.Close()
	res := new(bytes.Buffer)
	res.ReadFrom(response.Body)
	resp := res.String()

	if gjson.Get(resp, "pageInfo.totalResults").Int() == 0 {
		fmt.Println("Nie znaleziono, spróbuj ponownie.")
		os.Exit(1)
	}
	sl := gjson.Get(resp, "items.0.id.channelId")

	reqs, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/channels", nil)
	if err != nil {
		fmt.Println("Wystąpił błąd podczas wysyłania zapytania do API, czy na pewno masz dostęp do internetu?")
		os.Exit(1)
	}
	qsa := reqs.URL.Query()
	qsa.Add("key", apikey)
	qsa.Add("part", "snippet,contentDetails,statistics")
	qsa.Add("id", sl.String())
	reqs.URL.RawQuery = qsa.Encode()
	resps, err := http.Get(reqs.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resps.Body.Close()
	resa := new(bytes.Buffer)
	resa.ReadFrom(resps.Body)
	respa := resa.String()
	if gjson.Get(respa, "items.0.snippet.title").String() == "" {
		fmt.Println("Nie znaleziono, spróbuj ponownie.")
		os.Exit(1)
	}
	fmt.Println("\n** Informacje o kanale " + gjson.Get(respa, "items.0.snippet.title").String() + " **")
	fmt.Println("** Opis: \n" + gjson.Get(respa, "items.0.snippet.description").String() + "\n** Opis")
	fmt.Println("- Subskrybcji: " + gjson.Get(respa, "items.0.statistics.subscriberCount").String())
	fmt.Println("- Wyświetleń: " + gjson.Get(respa, "items.0.statistics.viewCount").String())
	fmt.Println("- Filmów: " + gjson.Get(respa, "items.0.statistics.videoCount").String())
	fmt.Println("\n** Informacje o kanale " + gjson.Get(respa, "items.0.snippet.title").String() + " **")

	bufio.NewReader(os.Stdin).ReadString('\n')
}
