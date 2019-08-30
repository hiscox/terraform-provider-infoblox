package resources

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hiscox/terraform-provider-infoblox/infoblox"
)

func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_HOST", nil),
				Description: "Infoblox host",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_USERNAME", nil),
				Description: "Infoblox username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_PASSWORD", nil),
				Description: "Infoblox password",
			},
			"wapi_version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_WAPI_VERSION", nil),
				Description: "Infoblox Web API version",
			},
			"tls_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_TLS_VERIFY", true),
				Description: "Infoblox SSL validation",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_TIMEOUT", 30),
				Description: "Request timeout",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"infoblox_a_record":     resourceARecord(),
			"infoblox_txt_record":   resourceTxtRecord(),
			"infoblox_cname_record": resourceCnameRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	params := infoblox.Cfg{
		Host:        d.Get("host").(string),
		WAPIVersion: d.Get("wapi_version").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		TLSVerify:   d.Get("tls_verify").(bool),
		Timeout:     d.Get("timeout").(int),
	}

	client, err := infoblox.ClientInit(&params)
	if err != nil {
		return nil, err
	}

	err = infoblox.IbGetTest(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
