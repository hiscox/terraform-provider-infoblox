package resources

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hiscox/terraform-provider-infoblox/infoblox"
)

func resourceTxtRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceTxtRecordCreate,
		Read:   resourceTxtRecordRead,
		Update: resourceTxtRecordUpdate,
		Delete: resourceTxtRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"text": &schema.Schema{
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

func resourceTxtRecordCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	text := d.Get("text").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)
	body := []byte(fmt.Sprintf(`{"name":%q, "text":%q, "view":%q}`, name, text, view))

	// this handles a record pre-existing to terraform being used
	log.Printf("Does remote record:txt exist for %s ?", name)
	r, i, err := infoblox.IbReadRecord(client, name, "txt")
	if r == 404 {
		log.Printf("Remote record:txt %s does not exist", name)
		d.SetId("")
		log.Printf("Creating record:txt %s", name)
		r, err = infoblox.IbCreateRecord(client, "txt", body)
		if err != nil {
			return err
		}
		if r == 201 {
			log.Printf("Setting state references...")
			d.Set("name", name)
			d.Set("text", text)
			d.Set("view", view)
			d.SetId(name + text + view)
		}
		return resourceTxtRecordRead(d, m)
	} else if r == 200 { // already exists, update local state and call update func
		log.Printf("Record:txt %s already exists", name)
		log.Printf("Updating remote...")
		r, err = infoblox.IbUpdateRecord(client, i.Ref, body)
		if err != nil {
			return err
		}
		log.Printf("Updating local state...")
		d.Set("name", i.Name)
		d.Set("text", i.Text)
		d.Set("view", i.View)
		d.SetId(name + text + view)
		return nil
	}
	if err != nil {
		return err
	}
	return resourceTxtRecordRead(d, m)
}

func resourceTxtRecordRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	text := d.Get("text").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)

	log.Printf("Retrieving remote record:txt for %s", name)
	r, i, err := infoblox.IbReadRecord(client, name, "txt")
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
	if name != i.Name || text != i.Text || view != i.View {
		log.Printf("Remote state doesn't match local")
		log.Printf("Local:\n" + " " + name + " " + text + " " + view)
		log.Printf("Remote:\n" + " " + i.Name + " " + i.Text + " " + i.View)
		return nil
	}
	return nil
}

func resourceTxtRecordUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("text") {
		name := d.Get("name").(string)
		text := d.Get("text").(string)
		view := d.Get("view").(string)
		client := m.(*resty.Client)
		body := []byte(fmt.Sprintf(`{"name":%q, "text":%q, "view":%q}`, name, text, view))

		// we need the _ref of the record to update it
		r, i, err := infoblox.IbReadRecord(client, name, "txt")
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
			d.Set("name", name)
			d.Set("text", text)
			d.Set("view", view)
			d.SetId(name + text + view)
			return nil
		}
	}
	return resourceTxtRecordRead(d, m)
}

func resourceTxtRecordDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := m.(*resty.Client)

	// we need the _ref of the record to delete it
	r, i, err := infoblox.IbReadRecord(client, name, "txt")
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
