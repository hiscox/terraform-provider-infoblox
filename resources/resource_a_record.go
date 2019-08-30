package resources

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hiscox/terraform-provider-infoblox/infoblox"
	"log"
)

func resourceARecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceARecordCreate,
		Read:   resourceARecordRead,
		Update: resourceARecordUpdate,
		Delete: resourceARecordDelete,

		Schema: map[string]*schema.Schema{
			"ipv4addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"view": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFOBLOX_VIEW", nil),
				Description: "Infoblox view, case sensitive",
			},
		},
	}
}

func resourceARecordCreate(d *schema.ResourceData, m interface{}) error {
	ipv4addr := d.Get("ipv4addr").(string)
	name := d.Get("name").(string)
	comment := d.Get("comment").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)
	body := []byte(fmt.Sprintf(`{"ipv4addr":%q, "name":%q, "comment":%q, "view":%q}`, ipv4addr, name, comment, view))
	// view cannot be updated so require special body for syncing remote state
	bodyUp := []byte(fmt.Sprintf(`{"ipv4addr":%q, "name":%q, "comment":%q}`, ipv4addr, name, comment))

	// this handles a record pre-existing to terraform being used
	log.Printf("Does remote record:a exist for %s ?", name)
	r, i, err := infoblox.IbReadRecord(client, name, "a")
	if r == 404 {
		log.Printf("Remote record:a %s does not exist", name)
		d.SetId("")
		log.Printf("Creating record:a %s", name)
		r, err = infoblox.IbCreateRecord(client, "a", body)
		if err != nil {
			return err
		}
		if r == 201 {
			log.Printf("Setting state references...")
			d.Set("ipv4addr", ipv4addr)
			d.Set("name", name)
			d.Set("comment", comment)
			d.Set("view", view)
			d.SetId(ipv4addr + name + comment + view)
		}
		return resourceARecordRead(d, m)
	} else if r == 200 { // already exists, update local state and call update func
		log.Printf("Record:a %s already exists", name)
		log.Printf("Updating remote...")
		r, err = infoblox.IbUpdateRecord(client, i.Ref, bodyUp)
		if err != nil {
			return err
		}
		log.Printf("Updating local state...")
		d.Set("ipv4addr", i.Ipv4addr)
		d.Set("name", i.Name)
		d.Set("comment", i.Comment)
		d.Set("view", i.View)
		d.SetId(ipv4addr + name + comment + view)
		return nil
	}
	if err != nil {
		return err
	}
	return resourceARecordRead(d, m)
}

func resourceARecordRead(d *schema.ResourceData, m interface{}) error {
	ipv4addr := d.Get("ipv4addr").(string)
	name := d.Get("name").(string)
	comment := d.Get("comment").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)

	log.Printf("Retrieving remote record:a for %s", name)
	r, i, err := infoblox.IbReadRecord(client, name, "a")
	// 404 indicates resource doesn't exist
	if r == 404 {
		log.Printf("Resource not found")
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	// remote state doesn't match local (manual updates to ib)
	if ipv4addr != i.Ipv4addr || name != i.Name || comment != i.Comment || view != i.View {
		log.Printf("Remote state doesn't match local")
		log.Printf("Local:\n" + ipv4addr + " " + name + " " + comment + " " + view)
		log.Printf("Remote:\n" + i.Ipv4addr + " " + i.Name + " " + i.Comment + " " + i.View)
		return nil
	}
	return nil
}

func resourceARecordUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("ipv4addr") || d.HasChange("name") || d.HasChange("comment") {
		ipv4addr := d.Get("ipv4addr").(string)
		name := d.Get("name").(string)
		comment := d.Get("comment").(string)
		view := d.Get("view").(string)
		client := m.(*resty.Client)
		body := []byte(fmt.Sprintf(`{"ipv4addr":%q, "name":%q, "comment":%q}`, ipv4addr, name, comment))

		// we need the _ref of the record to update it
		r, i, err := infoblox.IbReadRecord(client, name, "a")
		if err != nil {
			return err
		}
		// note that view cannot be updated
		r, err = infoblox.IbUpdateRecord(client, i.Ref, body)
		if err != nil {
			return err
		}
		if r == 200 {
			log.Printf("Setting state references...")
			d.Set("ipv4addr", ipv4addr)
			d.Set("name", name)
			d.Set("comment", comment)
			d.Set("view", view)
			d.SetId(ipv4addr + name + comment + view)
			return nil
		}
	}
	return resourceARecordRead(d, m)
}

func resourceARecordDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := m.(*resty.Client)

	// we need the _ref of the record to delete it
	r, i, err := infoblox.IbReadRecord(client, name, "a")
	if err != nil {
		return err
	}

	r, err = infoblox.IbDeleteRecord(client, i.Ref)
	if err != nil {
		return err
	}
	if r == 200 {
		return nil
	}
	return nil
}
