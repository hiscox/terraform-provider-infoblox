// Package infoblox provides REST actions against an infoblox WAPI
package infoblox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"log"
	"strings"
)

// Result contains the json fields from a Get request
type Result struct {
	Ref       string `json:"_ref"`
	Comment   string `json:"comment"`
	Text      string `json:"text"`
	Ipv4addr  string `json:"ipv4addr"`
	Name      string `json:"name"`
	Canonical string `json:"canonical"`
	View      string `json:"view"`
}

func init() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// IbReadRecord returns data about a record
func IbReadRecord(c *resty.Client, name string, rcdType string) (int, Result, error) {
	var url strings.Builder
	switch rcdType {
	case "a":
		url.WriteString("/record:a?name=")
		url.WriteString(name)
		url.WriteString("&_return_fields=ipv4addr,name,view,comment")
	case "txt":
		url.WriteString("/record:txt?name=")
		url.WriteString(name)
		url.WriteString("&_return_fields=name,view,text")
	case "cname":
		url.WriteString("/record:cname?name=")
		url.WriteString(name)
		url.WriteString("&_return_fields=name,view,comment,canonical")
	default:
		return 500, Result{}, errors.New("Unsupported record type")
	}
	log.Printf("IbReadRecord endpoint: %s", url.String())

	r, err := c.R().Get(url.String())
	log.Printf("Response body: \n" + r.String())

	// requires struct array as response returns a list of json[]
	log.Printf("Initialise Result struct")
	var result []Result

	log.Printf("Populating Result struct...")
	err = json.Unmarshal(r.Body(), &result)
	if err != nil {
		log.Printf("Error unmarshalling response into struct")
		return 500, Result{}, err
	}

	if r.String() == "[]" {
		log.Printf("Empty response body")
		return 404, Result{}, nil
	}

	if r.StatusCode() == 401 {
		return 401, Result{}, errors.New("Unauthorised: 401")
	} else if r.StatusCode() == 404 {
		log.Printf("Get request returned 404")
		return 404, Result{}, nil
	} else if err != nil {
		log.Printf("Get request failed")
		err = fmt.Errorf("Error: %s", err)
		return 500, Result{}, err
	}
	log.Println("Struct:")
	log.Println(result)

	log.Println("Status code: ", r.StatusCode())

	return r.StatusCode(), result[0], nil
}
