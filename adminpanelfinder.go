package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	print "github.com/fatih/color"
)

var (
	index      int = -1
	url        string
	panelPaths []string

	marks = []string{
	        " type=\"password\" ",
		" name=\"pwd\" ",
		" name=\"pass\" ",
		" name=\"password\" ",
		" name=\"username\" ",
		" value=\"Giriş Yap\" ",
		" value=\"Login\" ",
		" action=\"/login.php\" ",
	}
)

func main() {

	print.HiRed(`
	##################################### 
	#        Admin Panel Finder         #
	#-----------------------------------#
	#       Author: Dursun Katar        #
	#-----------------------------------#
	#       github.com/dursunkatar      #
	#####################################`)

	url = os.Args[1]
	panelPath := os.Args[2]

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	err := loadPanels(panelPath)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	goCount := 5
	chFound := make(chan string, 2)
	chThisIsNot := make(chan bool, goCount)
	chFinished := make(chan bool, 1)

	fmt.Println("")
	fmt.Println("Started...")

	for i := 0; i < goCount; i++ {
		go doControl(chFound, chThisIsNot, chFinished)
	}

exitLOOP:
	for {
		select {
		case mark := <-chFound:
			fmt.Print("MARK  : ")
			print.Cyan(strings.Trim(mark, " "))
			fmt.Print("FOUND : ")
			print.Green(<-chFound)
			break exitLOOP
		case <-chThisIsNot:
			go doControl(chFound, chThisIsNot, chFinished)
		case <-chFinished:
			break exitLOOP
		}
	}

	close(chFound)
	close(chThisIsNot)
	close(chFinished)
	fmt.Println("Finish")
}

func doControl(found chan string, thisIsNot, finish chan bool) {
	index++
	if index >= len(panelPaths) {
		finish <- true
		return
	}

	if ok, path, mark := connectUrl(); ok {
		found <- mark
		found <- path
	} else {
		thisIsNot <- true
	}
}

func connectUrl() (bool, string, string) {

	path := panelPaths[index]
	_url := url + path

	req, err := http.NewRequest("GET", _url, nil)

	if err != nil {
		return false, "", ""
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Content-Type", "text/html; charset=utf-8")
	req.Header.Add("Accept-Language", "tr-TR,tr;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Connection", "keep-alive")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, "", ""
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	result := string(body)

	if ok, mark := isThis(&result); ok {
		return true, path, mark
	}
	return false, "", ""
}

func isThis(source *string) (bool, string) {
	for _, mark := range marks {
		if strings.Contains(*source, mark) {
			return true, mark
		}
	}
	return false, ""
}

func loadPanels(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		panelPaths = append(panelPaths, strings.Trim(scanner.Text(), " /"))
	}
	file.Close()
	return nil
}
