package datadog

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	datadog "github.com/zorkian/go-datadog-api"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATADOG_API_KEY", nil),
			},
			"app_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATADOG_APP_KEY", nil),
			},
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATADOG_HOST", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"datadog_downtime":        resourceDatadogDowntime(),
			"datadog_metric_metadata": resourceDatadogMetricMetadata(),
			"datadog_monitor":         resourceDatadogMonitor(),
			"datadog_timeboard":       resourceDatadogTimeboard(),
			"datadog_screenboard":     resourceDatadogScreenboard(),
			"datadog_user":            resourceDatadogUser(),
			"datadog_integration_gcp": resourceDatadogIntegrationGcp(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type Config struct {
	ApiKey string
	AppKey string
	ApiUrl string
}

func (c *Config) GetClient() (*datadog.Client, error) {
	log.Printf("[INFO] NITRO CLIENT")
	log.Printf(c.ApiKey)
	log.Printf(c.AppKey)
	client := datadog.NewClient(c.ApiKey, c.AppKey)
	if c.ApiUrl != "" {
		client.SetBaseUrl(c.ApiUrl)
	}
	log.Println("[INFO] Datadog client successfully initialized, now validating...")

	ok, err := client.Validate()
	if err != nil {
		log.Printf("[ERROR] Datadog Client validation error: %v", err)
		return client, err
	} else if !ok {
		err := errors.New(`No valid credential sources found for Datadog Provider. Please see https://terraform.io/docs/providers/datadog/index.html for more information on providing credentials for the Datadog Provider`)
		log.Printf("[ERROR] Datadog Client validation error: %v", err)
		return client, err
	}
	log.Printf("[INFO] Datadog Client successfully validated.")

	return client, err
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] NITRO CONFIGURE")
	cfg := &Config{
		ApiKey: d.Get("api_key").(string),
		AppKey: d.Get("app_key").(string),
		ApiUrl: d.Get("api_url").(string),
	}
	cfg.GetClient()

	return cfg, nil
}
