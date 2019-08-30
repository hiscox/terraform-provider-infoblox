// Package infoblox provides REST actions against an infoblox WAPI
package infoblox

import (
	"crypto/tls"
	"errors"
	"github.com/go-resty/resty"
	"log"
	"time"
)

// Cfg config to construct client
type Cfg struct {
	Host        string
	WAPIVersion string
	Username    string
	Password    string
	TLSVerify   bool
	Timeout     int
}

func init() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// ClientInit establishes default settings on the rest client
func ClientInit(c *Cfg) (*resty.Client, error) {
	client := resty.New()

	if c.Host == "" {
		return nil, errors.New("Invalid Host setting")
	}
	if c.Username == "" {
		return nil, errors.New("Invalid Username setting")
	}
	if c.Password == "" {
		return nil, errors.New("Invalid Password setting")
	}

	if c.TLSVerify == false {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	}

	client.SetBasicAuth(c.Username, c.Password)
	client.SetHeader("Content-Type", "application/json")
	client.SetTimeout(time.Duration(c.Timeout) * time.Second)
	client.SetHostURL("https://" + c.Host + "/wapi/v" + c.WAPIVersion)

	return client, nil
}

// IbGetTest sends a GET request to infoblox
func IbGetTest(c *resty.Client) error {
	r, err := c.R().Get("")
	if r.StatusCode() == 401 {
		return errors.New("Unauthorised: 401")
	}
	if err != nil {
		log.Printf("Get request failed")
		return err
	}
	log.Printf("[DEBUG] " + r.String())
	return nil
}
