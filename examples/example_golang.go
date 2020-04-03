package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const url = "https://covid-detc.us.keyayun.net"
const token = "the bearer token"

func newWorkItem(affectedSOPInstanceUID, studyInstanceUID string) error {
	apiUrl := url + "/workitems?" + affectedSOPInstanceUID
	payload := strings.NewReader("{\n    \"00741204\": \"COVID-19\",\n    \"00404021\": {\n        \"0020000D\": \"" + studyInstanceUID + "\"\n    }\n}")

	req, err := http.NewRequest("POST", apiUrl, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New(string(body))
	}

	return nil
}

func getWorkItem(affectedSOPInstanceUID string) error {
	apiUrl := url + "/workitems/" + affectedSOPInstanceUID

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}
	fmt.Println(string(body))
	return nil
}

func main() {
	testSOPUID := "1.3.45.214"
	testStudyUID := "2.16.840.1.113662.2.1.99999.5175439602988854"

	if err := newWorkItem(testSOPUID, testStudyUID); err != nil {
		fmt.Println(err)
		return
	}

	if err := getWorkItem(testSOPUID); err != nil {
		fmt.Println(err)
	}
}
