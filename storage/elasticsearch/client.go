package elasticsearch

import "github.com/elastic/go-elasticsearch/v7"

// NewDefaultClient creates a new client with default options.
//
// It will use http://localhost:9200 as the default address.
//
// It will use the ELASTICSEARCH_URL environment variable, if set,
// to configure the addresses; use a comma to separate multiple URLs.
func NewDefaultClient() (*elasticsearch.Client, error) {
	return elasticsearch.NewDefaultClient()
}
