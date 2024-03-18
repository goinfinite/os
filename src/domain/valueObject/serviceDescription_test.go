package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestServiceDescription(t *testing.T) {
	t.Run("ValidServiceDescription", func(t *testing.T) {
		validServiceDescriptions := []string{
			"Build and grow your website with the best way to WordPress. Lightning-fast hosting, intuitive, flexible editing, and everything you need to grow your site and audience, baked right in.",
			"It's comprised of Elasticsearch, Kibana, Beats, and Logstash (also known as the ELK Stack) and more. Reliably and securely take data from any source, in any format, then search, analyze, and visualize.",
			"RabbitMQ is a reliable and mature messaging and streaming broker, which is easy to deploy on cloud environments, on-premises, and on your local machine. It is currently used by millions worldwide.",
		}
		for _, svcDesc := range validServiceDescriptions {
			_, err := NewServiceDescription(svcDesc)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", svcDesc, err.Error())
			}
		}
	})

	t.Run("InvalidServiceDescription", func(t *testing.T) {
		invalidLength := 600
		invalidServiceDescriptions := []string{
			"",
			"a",
			testHelpers.GenerateString(invalidLength),
		}
		for _, svcDesc := range invalidServiceDescriptions {
			_, err := NewServiceDescription(svcDesc)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", svcDesc)
			}
		}
	})
}
