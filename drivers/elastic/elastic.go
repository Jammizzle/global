package elastic

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

var (
	User     string = os.Getenv("ELASTIC_USERNAME")
	Password string = os.Getenv("ELASTIC_PASSWORD")
	Host     string = os.Getenv("ELASTIC_HOST")
)

type ElasticConnector struct {
	Client *elasticsearch.Client
}

// Connect initalised conenction to
func (e *ElasticConnector) Connect() (err error) {
	config := elasticsearch.Config{
		Addresses: []string{
			Host,
		},
		Username: User,
		Password: Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}

	e.Client, err = elasticsearch.NewClient(config)

	if err != nil {
		logrus.Error(err)
		return err
	}

	res, err := e.Client.Info()
	if res.IsError() {
		logrus.Fatalf("Error: %s", res.String())
	}

	return nil
}
