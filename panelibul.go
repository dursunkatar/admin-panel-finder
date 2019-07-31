package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	url string

	isaretler = []string{
		" name=\"pwd\" ",
		" name=\"pass\" ",
		" name=\"password\" ",
		" value=\"Giriş Yap\" ",
		" value=\"Login\" ",
	}
)

func main() {

	fmt.Println(`
 ##################################### 
 #        Author: Dursun Katar       #
 #       github.com/dursunkatar      #
 #####################################`)

	url = os.Args[1]
	panelFile := os.Args[2]

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	paneller, err := panelleriYukle(panelFile)

	if err != nil {
		fmt.Println("Hata: ", err.Error())
		return
	}

	fmt.Println("")
	fmt.Println(" Bakıyor...")
	for _, panel := range *paneller {
		if baglan(panel) {
			fmt.Println(" Panel : ", panel)
			break
		}
	}
	fmt.Println(" Bitti")
}

func baglan(panel string) bool {
	_url := url + panel
	req, err := http.NewRequest("GET", _url, nil)

	if err != nil {
		return false
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Content-Type", "text/html; charset=utf-8")
	req.Header.Add("Accept-Language", "tr-TR,tr;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Connection", "keep-alive")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	result := string(body)

	for _, isaret := range isaretler {
		if strings.Contains(result, isaret) {
			fmt.Println(" İşaret: ", strings.Trim(isaret, " "))
			return true
		}
	}
	return false
}

func panelleriYukle(dosyaYolu string) (*[]string, error) {

	var paneller []string
	file, err := os.Open(dosyaYolu)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		paneller = append(paneller, strings.Trim(scanner.Text(), " /"))
	}
	file.Close()
	return &paneller, nil
}
