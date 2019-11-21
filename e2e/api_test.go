package main

import (
	"bytes"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// bannersTest
type bannersTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex sync.RWMutex
	stopSignal    chan struct{}

	responseStatusCode int
	responseBody       []byte
	bannerID           string
}

func (test *bannersTest) iSendRequestToWithData(httpMethod, url, contentType string, data *gherkin.DocString) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)
		r, err = http.Post(url, contentType, bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseStatusCode = r.StatusCode

	return
}

func (test *bannersTest) iSendRequestToWithDataAndSelectBanner(
	httpMethod, url,
	contentType string,
	data *gherkin.DocString,
) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)
		r, err = http.Post(url, contentType, bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	test.bannerID = strings.TrimSpace(string(test.responseBody))

	return
}

func (test *bannersTest) iSendRequestToRemoveBanner(httpMethod, url string) error {
	client := &http.Client{}

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	test.responseStatusCode = resp.StatusCode

	return nil
}

func (test *bannersTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *bannersTest) theResponseShouldMatchIdBanner(bannerID string) error {
	if test.bannerID != bannerID {
		return fmt.Errorf("unexpected banner id: %s != %s", test.bannerID, bannerID)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	test := new(bannersTest)

	s.Step(`^1. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^2. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)

	s.Step(`^3. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^4. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^5. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)

	s.Step(`^6. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)

	s.Step(`^7. I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestToRemoveBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
}
