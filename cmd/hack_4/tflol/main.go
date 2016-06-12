package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/icio/mondo/cmd/hack_4"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stderr)

	jar, _ := cookiejar.New(nil)
	client := &TFLClient{
		Client: &http.Client{
			Jar: jar,
		},
	}

	updates := make(chan *hack_4.Account)
	user := os.Getenv("TFL_USERNAME")
	pass := os.Getenv("TFL_PASSWORD")
	go client.Monitor(user, pass, time.Second, updates)

	for update := range updates {
		jsonUpdate, _ := json.Marshal(update)
		log.Printf("%s", jsonUpdate)
		http.Post(Getenv("NERVE_URL", "http://localhost:8080/journeys"), "application/json", bytes.NewReader(jsonUpdate))
	}
}

func Getenv(name, _default string) string {
	env := os.Getenv(name)
	if env == "" {
		return _default
	}
	return env
}

type TFLClient struct {
	Client *http.Client
}

func (c *TFLClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://contactless.tfl.gov.uk")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.110 Safari/537.36")

	return c.Client.Do(req)
}

func (client *TFLClient) Monitor(user, pass string, wait time.Duration, accountStates chan<- *hack_4.Account) {
	// Bake our first batch of cookies.
	client.Do(NewHomepageRequest())

	for {
		time.Sleep(wait)

		// Login and fetch dashboard.
		accountResp, err := client.Do(NewLoginRequest(user, pass))
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		// Parse account info.
		account, err := parseAccountPage(accountResp)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		accountStates <- account

		// Log out to kick their cache.
		resp, err := client.Do(NewLogoutRequest(account.LogoutToken))
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		resp.Body.Close()
	}
}

func NewHomepageRequest() *http.Request {
	home, _ := http.NewRequest("GET", "https://contactless.tfl.gov.uk/", nil)
	return home
}

func NewLoginRequest(username, password string) *http.Request {
	loginBody := url.Values{
		"AppId":     {"a3ac81d4-80e8-4427-b348-a3d028dfdbe7"},
		"App":       {"a3ac81d4-80e8-4427-b348-a3d028dfdbe7"},
		"ReturnUrl": {"https://contactless.tfl.gov.uk/DashBoard"},
		"ErrorUrl":  {"https://contactless.tfl.gov.uk/"},
		"UserName":  {username},
		"Password":  {password},
	}

	login, _ := http.NewRequest("POST", "https://account.tfl.gov.uk/Login", strings.NewReader(loginBody.Encode()))
	login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return login
}

func NewLogoutRequest(verificationToken string) *http.Request {
	logoutBody := url.Values{
		"__RequestVerificationToken": {verificationToken},
	}
	logout, _ := http.NewRequest("POST", "https://contactless.tfl.gov.uk/HomePage/SignOut", strings.NewReader(logoutBody.Encode()))
	logout.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return logout
}

func parseAccountPage(resp *http.Response) (*hack_4.Account, error) {
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	spendingToday, _ := parseCost(doc.Find("#intradayTotalCharge").Text())

	var payments []hack_4.Payment
	doc.Find(".csc-payment-row").Each(func(i int, s *goquery.Selection) {
		cost, _ := parseCost(s.Find("span[data-pageobject=journey-price]").Text())
		departure, arrival := parseJourneyTime(s.Find("span[data-pageobject=journey-time]").Text())
		payments = append(payments, hack_4.Payment{
			Origin:      s.Find("span[data-pageobject=journey-from]").Text(),
			Departure:   departure,
			Destination: s.Find("span[data-pageobject=journey-to]").Text(),
			Arrival:     arrival,
			Cost:        cost,
			Auto:        s.Find(".autocompleted-icon").Length() > 0,
			Limited:     s.Find(".capped-icon-day").Length() > 0,
		})
	})

	return &hack_4.Account{
		SpendingToday:   spendingToday,
		PendingPayments: payments,
		LogoutToken:     doc.Find("[name=__RequestVerificationToken]").AttrOr("value", ""),
	}, nil
}

func parseJourneyTime(journeyTime string) (*time.Time, *time.Time) {
	// "09:16 - 09:24"
	// "09:16 - --:--"
	if len(journeyTime) != 13 {
		return nil, nil
	}
	return parseTime(journeyTime[:5]), parseTime(journeyTime[8:])
}

var london, _ = time.LoadLocation("Europe/London")

func parseTime(timeText string) *time.Time {
	hour, hErr := strconv.Atoi(timeText[:2])
	min, mErr := strconv.Atoi(timeText[3:])
	if hErr != nil || mErr != nil {
		return nil
	}

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, london)
	return &t
}

func parseCost(price string) (int, error) {
	price = strings.Trim(price, " \n")
	if len(price) > 2 {
		cost, err := strconv.ParseFloat(price[2:], 64)
		if err != nil {
			return 0, err
		}

		return int(cost * 100), nil
	}

	return 0, nil
}
