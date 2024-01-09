package webhooks

type WebhookVendor string

const (
	Gitea  WebhookVendor = "gitea"
	Github WebhookVendor = "github"
	Gitlab WebhookVendor = "gitlab"
)

type WebhookConfig struct {
	Vendor WebhookVendor `yaml:"vendor"`
	Secret string        `yaml:"secret"`
	Events []string      `yaml:"events"`
}
