package resources

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hiscox/terraform-provider-infoblox/infoblox"
	"log"
)

func resourceCnameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceCnameRecordCreate,
		Read:   resourceCnameRecordRead,
		Update: resourceCnameRecordUpdate,
		Delete: resourceCnameRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"canonical": &schema.Schema{
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

func resourceCnameRecordCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	canonical := d.Get("canonical").(string)
	comment := d.Get("comment").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)
	body := []byte(fmt.Sprintf(`{"name":%q, "canonical":%q, "comment":%q, "view":%q}`, name, canonical, comment, view))

	// this handles a record pre-existing to terraform being used
	log.Printf("Does remote record:cname exist for %s ?", name)
	r, i, err := infoblox.IbReadRecord(client, name, "cname")
	if r == 404 {
		log.Printf("Remote record:cname %s does not exist", name)
		d.SetId("")
		log.Printf("Creating record:cname %s", name)
		r, err = infoblox.IbCreateRecord(client, "cname", body)
		if err != nil {
			return err
		}
		if r == 201 {
			log.Printf("Setting state references...")
			d.Set("name", name)
			d.Set("canonical", canonical)
			d.Set("comment", comment)
			d.Set("view", view)
			d.SetId(name + canonical + comment + view)
		}
		return resourceCnameRecordRead(d, m)
	} else if r == 200 { // already exists, update local state and call update func
		log.Printf("Record:cname %s already exists", name)
		log.Printf("Updating remote...")
		r, err = infoblox.IbUpdateRecord(client, i.Ref, body)
		if err != nil {
			return err
		}
		log.Printf("Updating local state...")
		d.Set("name", i.Name)
		d.Set("canonical", i.Canonical)
		d.Set("comment", i.Comment)
		d.Set("view", i.View)
		d.SetId(name + canonical + comment + view)
		return nil
	}
	if err != nil {
		return err
	}
	return resourceCnameRecordRead(d, m)
}

func resourceCnameRecordRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	canonical := d.Get("canonical").(string)
	comment := d.Get("comment").(string)
	view := d.Get("view").(string)
	client := m.(*resty.Client)

	log.Printf("Retrieving remote record:cname for %s", name)
	r, i, err := infoblox.IbReadRecord(client, name, "cname")
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
	if name != i.Name || canonical != i.Canonical || comment != i.Comment || view != i.View {
		log.Printf("Remote state doesn't match local")
		log.Printf("Local:\n" + " " + name + " " + canonical + " " + comment + " " + view)
		log.Printf("Remote:\n" + " " + i.Name + " " + i.Canonical + " " + i.Comment + " " + i.View)
		return nil
	}
	return nil
}

func resourceCnameRecordUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("name") || d.HasChange("canonical") || d.HasChange("comment") {
		name := d.Get("name").(string)
		canonical := d.Get("canonical").(string)
		comment := d.Get("comment").(string)
		view := d.Get("view").(string)
		client := m.(*resty.Client)
		body := []byte(fmt.Sprintf(`{"name":%q, "canonical":%q, "comment":%q, "view":%q}`, name, canonical, comment, view))

		// we need the _ref of the record to update it
		r, i, err := infoblox.IbReadRecord(client, name, "cname")
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
			d.Set("canonical", canonical)
			d.Set("comment", comment)
			d.Set("view", view)
			d.SetId(name + canonical + comment + view)
			return nil
		}
	}
	return resourceCnameRecordRead(d, m)
}

func resourceCnameRecordDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := m.(*resty.Client)

	// we need the _ref of the record to delete it
	r, i, err := infoblox.IbReadRecord(client, name, "cname")
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
