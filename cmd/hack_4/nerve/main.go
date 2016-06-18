package main

import (
	"encoding/json"
	"fmt"
	"github.com/icio/mondo"
	"github.com/icio/mondo/cmd/hack_4"
	"github.com/icio/mondo/mondodomain"
	"github.com/icio/mondo/mondohttp"
	"io/ioutil"
	"log"
	// "math"
	"net/http"
	"os"
	"sync"
)

func main() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stderr)

	m := &mondo.Client{
		HTTPClient: &mondo.HTTPClient{
			Host:      Getenv("MONDO_API", "api.getmondo.co.uk"),
			Client:    http.DefaultClient,
			UserAgent: "hackathon-iv-tfl/0.1 (+https://github.com/icio/mondo)",
		},
		Auth: mondo.NewAccessTokenAuth(os.Getenv("MONDO_ACCESS_TOKEN")),
	}

	http.HandleFunc("/accounts", mondoAuth(m, accounts))
	http.HandleFunc("/journeys", mondoAuth(m, journeys))
	log.Fatal(http.ListenAndServe(":8080", httpLogger(http.DefaultServeMux)))
}

func Getenv(name, _default string) string {
	env := os.Getenv(name)
	if env == "" {
		return _default
	}
	return env
}

func mondoAuth(m *mondo.Client, handler func(m *mondo.Client, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(m, w, r)
	}
}

var allSightings = make([]hack_4.Sighting, 0)
var allCost int = 0
var journeysLock = sync.Mutex{} // awww yeeaaah.

// journeys parses the account information to determine what connections a user
// has made.
func journeys(m *mondo.Client, w http.ResponseWriter, r *http.Request) {
	journeysLock.Lock()
	defer journeysLock.Unlock()

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR reading body: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	account := new(hack_4.Account)
	err = json.Unmarshal(body, account)
	if err != nil {
		log.Printf("ERROR during unmarshalling: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Account: %#v", account)

	sightings := extrapolateSightings(account)

	for _, sighting := range sightings {
		log.Printf("    Sighting: %#v", sighting)
	}

	// Merge the time-ordered list of journeys
	added := 0
	spent := float64(0)

	allSightings, added = mergeSightings(allSightings, sightings)
	if added > 0 {
		lastSeen := allSightings[len(allSightings)-1]
		spent = float64(lastSeen.Cost - allCost)
		allCost = lastSeen.Cost

		if lastSeen.Out {
			suffix := ""
			if lastSeen.LimitCrossed {
				suffix = " You've just hit your spending limit for the day. Future travel won't cost a thing! \U0001f389\U0001f4b8"
				// Unless you travel into a different zone?
			}
			resp, err := m.Do(mondohttp.NewCreateURLFeedItemRequest(
				"",
				os.Getenv("MONDO_ACCOUNT_ID"),
				"http://www.nyan.cat/",
				fmt.Sprintf("Welcome to %s. This journey cost you £%.2f.%s", lastSeen.Place, spent/100, suffix),
				"https://tfl.gov.uk/cdn/static/assets/icons/favicon-160x160.png",
			))
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Printf("Feed notification sent. [%d]", resp.StatusCode)
			}
		}
	}

	log.Printf("  + %d sightings added", added)

	for _, sighting := range allSightings {
		log.Printf("--- Sighting: %#v", sighting)
	}

	w.WriteHeader(http.StatusNoContent)
}

func mergeSightings(target []hack_4.Sighting, extra []hack_4.Sighting) ([]hack_4.Sighting, int) {
	added := 0

	i := 0
	for j := 0; j < len(extra); j++ {
		a := extra[j]
		// log.Printf(" a: [%d] %#v", j, a)

		for {
			// Append to end.
			if i >= len(target) {
				// log.Printf("  - Insert at front")
				added++
				target = append(target, extra[j])
				break
			}

			// Insert in time-order elsewhere.
			b := target[i]
			// log.Printf(" b: [%d] %#v", i, b)
			if a.Time.Equal(b.Time) && a.Place == b.Place && a.In == b.In && a.Out == b.Out {
				// log.Printf("  - Skip")
				break // Dedupe.
			} else if a.Time.Before(b.Time) {
				// log.Printf("  - Insert ")
				added++
				target = append(target[:i], append([]hack_4.Sighting{a}, target[i:]...)...)
				break
			} else {
				// log.Printf("  .")
			}
			i++
		}
	}

	return target, added
}

func extrapolateSightings(acc *hack_4.Account) []hack_4.Sighting {
	sightings := make([]hack_4.Sighting, 0)
	cost := 0
	limited := false

	for _, payment := range acc.PendingPayments {
		// Log the origin.
		sightings = append(sightings, hack_4.Sighting{
			Place: payment.Origin,
			Time:  *payment.Departure,
			In:    true,
			Cost:  cost,
		})

		// Skip guesses TfL make right after a user taps in. Also,
		// unfortunately, skips cases where the user has legitimately forgotten
		// to tap out.
		if payment.Arrival == nil || (*payment.Departure).Equal(*payment.Arrival) {
			continue
		}
		if payment.Destination == "Unknown" || payment.Destination == "" {
			continue
		}

		cost += payment.Cost

		sightings = append(sightings, hack_4.Sighting{
			Place:        payment.Destination,
			Time:         *payment.Arrival,
			Out:          true,
			Cost:         cost,
			Auto:         payment.Auto,
			LimitCrossed: !limited && payment.Limited, // FIXME: This needs to be determined _after_ merging all of the sightings together.
		})
		limited = limited || payment.Limited
	}

	return sightings
}

// accounts writes the list of a user's mondo accounts in the response.
func accounts(m *mondo.Client, w http.ResponseWriter, r *http.Request) {
	// Get the accounts listing.
	accounts := new(mondodomain.AccountsResponse)
	err := m.DoInto(mondohttp.NewAccountsRequest(""), accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response to the client.
	body, err := json.Marshal(accounts.Accounts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// httpLogger wraps http handler functions and logs the kind of request made.
func httpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP %s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
