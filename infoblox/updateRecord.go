// Package infoblox provides REST actions against an infoblox WAPI
package infoblox

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func init() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// IbUpdateRecord updates a record
func IbUpdateRecord(c *resty.Client, ref string, body []byte) (int, error) {
	log.Printf("IbUpdateRecord endpoint: /%s", ref)
	log.Printf("IbUpdateRecord request body: %s", body)

	r, err := c.R().SetBody(body).Put("/" + ref)
	log.Printf("Response body: \n" + r.String())
	sc := r.StatusCode()

	if sc == 200 {
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
		log.Printf("Put request failed")
		err = fmt.Errorf("Error: %s", err)
		return 500, err
	}
	return sc, nil
}
