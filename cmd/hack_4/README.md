# Mondo Hackathon 4: Realtime TfL Notifications

Assuming you have a defined `GOPATH` envvar:

```bash
# Check out the code:
git clone https://github.com/icio/mondo/ $GOPATH/src/github.com/icio/mondo --branch=hack4-realtime-tfl

# TfL Authentication: (sign up at https://contactless.tfl.gov.uk/ and register your Mondo card)
export TFL_USERNAME=
export TFL_PASSWORD=

# Run the TfL Scraper:
go run $GOPATH/src/github.com/icio/mondo/cmd/hack_4/tflol/main.go &

# Mondo Authentication: (copy from https://developers.getmondo.co.uk)
export MONDO_ACCOUNT_ID=
export MONDO_ACCESS_TOKEN=

# Run the journey-conflation service:
go run $GOPATH/src/github.com/icio/mondo/cmd/hack_4/nerve/main.go
```

Please see **[Mondo Hackathon 4: Realtime TfL Notifications][article]** for information about how the data is consumed from O&C and the behaviours of O&C that necessitate this.

##### Notes

* The current authentication mechanism uses an access token, each of whose lifetime is limited to a couple of days.

* The notification logic is particularly naive at the moment -- see [the writeup][article] for a preferrable implementation. Nerve only stores its state in memory. This means it is unable to distinguish new sightings received because you travelled on the underground from sightings already seen and re-sent from the TfL scraper after starting/restarting nerve. If you already have journeys in O&C when you start nerve, you'll get a notification of the latest journey with cost equal to the total spent today.

* O&C is only showing time-slices of events, so care should be given to investigate how our applications behave when we tick over into the next slice. Our own sighting and journey reverse-engineering should be resilient to this, assuming we use accurate timestamps, but this requires that the scraper dates all times as the current day, in the current timezone. No investigation has been made into what happens when journeys span the billing point of TfL.

* A large (and deliberate) oversight of the project is that we assume only tube journeys. What do bus journeys look like? A quick test suggests they don’t actually show up on O&C during the day, and so our service won't faulter on them, but this behaviour -- if correct -- is a confusing limitation that users would stumble on in the way they currently do with all of TfL.

* O&C allows you to register multiple cards with O&C, but we assume only one. The cards registered can be given recognisable nicknames, and the “Today’s journeys” features allows you to switch between cards. These features could be used to handle multiple cards.

[article]: https://medium.com/@paulscott/c9ff8084724c
