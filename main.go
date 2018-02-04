package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type dCryptResponse struct {
	Success struct {
		Message string   `json:"message"`
		Links   []string `json:"links"`
	} `json:"success"`
}

func main() {
	var dlcFile string
	var linksToFile bool

	flag.StringVar(&dlcFile, "dlc", "", "DLC file to decrypt")
	flag.BoolVar(&linksToFile, "f", false, "Decrypted links should be written to file, not to command line")

	flag.Parse()

	if strings.Compare(dlcFile, "") == 0 {
		fmt.Println("No DLC provided.")
		return
	}

	res, err := upload(dlcFile)
	if err != nil {
		log.Fatal(err)
	}

	processResponse(res, linksToFile)

}

func upload(dlcFile string) (*http.Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(dlcFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fw, err := w.CreateFormFile("dlcfile", dlcFile)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest("POST", "http://dcrypt.it/decrypt/upload", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return nil, err
	}

	return res, nil
}

func processResponse(res *http.Response, linksToFile bool) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	stringBody := string(body)
	jsonResponseString := strings.Split(stringBody, "\n")[1]

	var d dCryptResponse
	if err := json.Unmarshal([]byte(jsonResponseString), &d); err != nil {
		log.Fatal(err)
	}

	links := strings.Join(d.Success.Links, "\n")
	if linksToFile {
		linksFile := "links.txt"
		f, err := os.Create(linksFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		f.WriteString(links)
	} else {
		fmt.Println(links)
	}
}
