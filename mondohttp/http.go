package mondohttp

import (
	"net/url"
	"strconv"
)

// ProductionAPI is the base URL of Mondo's production API.
const ProductionAPI string = "https://api.getmondo.co.uk/"

// StagingAPI is the usual base URL of Mondo's staging API (commonly available
// during hackathons).
const StagingAPI string = "https://staging-api.gmon.io/"

func auth(auth string) (string, string) {
	return "Authorization", auth
}

func formContentType() (string, string) {
	return "Content-Type", "application/x-www-form-urlencoded"
}

func appendPaginationParams(query *url.Values, since, before string, limit int) {
	if since != "" {
		query.Set("since", since)
	}
	if before != "" {
		query.Set("before", before)
	}
	if limit != 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
}

func appendQueryMap(query *url.Values, prefix, suffix string, params map[string]string) {
	for key, value := range params {
		query.Add(prefix+key+suffix, value)
	}
}
