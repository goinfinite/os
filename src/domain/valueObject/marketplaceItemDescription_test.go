package valueObject

import (
	"testing"
)

func TestMarketplaceItemDescription(t *testing.T) {
	t.Run("ValidMarketplaceItemDescription", func(t *testing.T) {
		validMarketplaceItemDescriptions := []string{
			"Build and grow your website with the best way to WordPress. Lightning-fast hosting, intuitive, flexible editing, and everything you need to grow your site and audience, baked right in.",
			"It's comprised of Elasticsearch, Kibana, Beats, and Logstash (also known as the ELK Stack) and more. Reliably and securely take data from any source, in any format, then search, analyze, and visualize.",
			"RabbitMQ is a reliable and mature messaging and streaming broker, which is easy to deploy on cloud environments, on-premises, and on your local machine. It is currently used by millions worldwide.",
		}
		for _, mid := range validMarketplaceItemDescriptions {
			_, err := NewMarketplaceItemDescription(mid)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mid, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemDescription", func(t *testing.T) {
		invalidMarketplaceItemDescriptions := []string{
			"",
			"a",
		}
		for _, mid := range invalidMarketplaceItemDescriptions {
			_, err := NewMarketplaceItemDescription(mid)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mid)
			}
		}
	})
}
