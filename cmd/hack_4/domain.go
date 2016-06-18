package hack_4

import "time"

type Account struct {
	LogoutToken     string    `json:"-"`
	SpendingToday   int       `json:"spending_today"`
	PendingPayments []Payment `json:"pending_payments"`
}

type Payment struct {
	Origin      string     `json:"origin"`
	Departure   *time.Time `json:"depart,omitempty"`
	Destination string     `json:"dest"`
	Arrival     *time.Time `json:"arrival,omitempty"`
	Cost        int        `json:"cost"`

	// Auto indicates whether an autocomplete icon was shown.
	Auto bool `json:"auto"`
	// Limited indicates whether a price-cap icon was shown.
	Limited bool `json:"limited"`
}

type Sighting struct {
	Place        string
	Time         time.Time
	In           bool
	Out          bool
	Cost         int
	Auto         bool
	LimitCrossed bool
}
