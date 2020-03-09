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

// IbDeleteARecord creates an a record
func IbDeleteRecord(c *resty.Client, ref string) (int, error) {
	log.Printf("IbDeleteRecord endpoint: /%s", ref)

	r, err := c.R().Delete("/" + ref)
	log.Printf("Response body: \n" + r.String())

	if r.StatusCode() == 200 {
		return r.StatusCode(), nil
	} else if r.StatusCode() == 401 {
		return 401, errors.New("Unauthorised: 401")
	} else if r.StatusCode() == 404 {
		log.Printf("Get request returned 404")
		return 404, nil
	} else if r.StatusCode() == 400 {
		log.Printf("Bad request")
		return 400, errors.New("Bad request: 400" + r.String())
	} else if err != nil {
		log.Printf("Delete request failed")
		err = fmt.Errorf("Error: %s", err)
		return 500, err
	}
	return r.StatusCode(), nil
}
