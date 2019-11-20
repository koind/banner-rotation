package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	queueName    = "banners"
	exchangeName = "statistics_of_rotation"
)

var amqpDSN = os.Getenv("TESTS_AMQP_DSN")

// Statistic model
type Statistic struct {
	ID       int       `json:"id" db:"id"`
	Type     int       `json:"type" db:"type"`
	BannerID int       `json:"bannerId" db:"banner_id"`
	SlotID   int       `json:"slotId" db:"slot_id"`
	GroupID  int       `json:"groupId" db:"group_id"`
	CreateAt time.Time `json:"createAt" db:"create_at"`
}

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

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (test *bannersTest) startConsuming(interface{}) {
	test.messages = make([][]byte, 0)
	test.messagesMutex = sync.RWMutex{}
	test.stopSignal = make(chan struct{})

	var err error

	test.conn, err = amqp.Dial(amqpDSN)
	panicOnErr(err)

	test.ch, err = test.conn.Channel()
	panicOnErr(err)

	// Consume
	_, err = test.ch.QueueDeclare(queueName, true, false, true, false, nil)
	panicOnErr(err)

	err = test.ch.QueueBind(queueName, "", exchangeName, false, nil)
	panicOnErr(err)

	events, err := test.ch.Consume(queueName, "", true, false, false, false, nil)
	panicOnErr(err)

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
}

func (test *bannersTest) stopConsuming(interface{}, error) {
	test.stopSignal <- struct{}{}

	panicOnErr(test.ch.Close())
	panicOnErr(test.conn.Close())
	test.messages = nil
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
	test.bannerID = string(test.responseBody)

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

func (test *bannersTest) iReceiveEventWithIdBanner(bannerID string) error {
	time.Sleep(1 * time.Second)

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	for _, msg := range test.messages {
		var statistic Statistic

		err := json.Unmarshal(msg, &statistic)
		if err != nil {
			return err
		}

		if strconv.Itoa(statistic.BannerID) == bannerID {
			return nil
		}
	}

	return fmt.Errorf("banner with id '%s' was not found in %s", bannerID, test.messages)
}

func FeatureContext(s *godog.Suite) {
	test := new(bannersTest)

	s.BeforeScenario(test.startConsuming)

	s.Step(`^1. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^2. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)
	s.Step(`^I receive event with id banner "([^"]*)"$`, test.iReceiveEventWithIdBanner)

	s.Step(`^3. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^I receive event with id banner "([^"]*)"$`, test.iReceiveEventWithIdBanner)

	s.Step(`^4. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^5. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)
	s.Step(`^I receive event with id banner "([^"]*)"$`, test.iReceiveEventWithIdBanner)

	s.Step(`^6. I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithDataAndSelectBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match id banner "([^"]*)"$`, test.theResponseShouldMatchIdBanner)
	s.Step(`^I receive event with id banner "([^"]*)"$`, test.iReceiveEventWithIdBanner)

	s.Step(`^7. I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestToRemoveBanner)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.AfterScenario(test.stopConsuming)
}
