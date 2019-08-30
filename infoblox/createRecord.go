// Package infoblox provides REST actions against an infoblox WAPI
package infoblox

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"log"
	"strings"
)

func init() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// IbCreateRecord creates a record
func IbCreateRecord(c *resty.Client, rcdType string, body []byte) (int, error) {
	var url strings.Builder
	switch rcdType {
	case "a":
		url.WriteString("/record:a")
	case "txt":
		url.WriteString("/record:txt")
	case "cname":
		url.WriteString("/record:cname")
	default:
		return 500, errors.New("Unsupported record type")
	}
	log.Printf("IbCreateRecord endpoint: %s", url.String())
	log.Printf("IbCreateRecord request body: %s", body)

	r, err := c.R().SetBody(body).Post(url.String())
	log.Printf("Response body: \n" + r.String())
	sc := r.StatusCode()

	if sc == 201 {
		return sc, nil
	} else if sc == 401 {
		return 401, errors.New("Unauthorised: 401")
	} else if sc == 404 {
		log.Printf("Get request returned 404")
		return 404, nil
	} else if sc == 400 {
		log.Printf("Bad request")
		return 400, errors.New("Bad request: 400" + r.String())
	} else if err != nil {
		log.Printf("Post request failed")
		err = fmt.Errorf("Error: %s", err)
		return 500, err
	}
	return sc, nil
}
