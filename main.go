package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

const (
	defaultBaseURL = "https://gitlab.com/"
)

var (
	listenAddress  = flag.String("listen-address", ":8080", "Address to serve HTTP requests")
	insecureClient = flag.Bool("k", false, "Disable SSL verification")
	configFile     = flag.String("config", "dredd.yaml", "Path to configuration file")
	dryRun         = flag.Bool("dry-run", false, "Runs without making changes")

	netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
)

func main() {
	flag.Parse()

	token := os.Getenv("GITLAB_TOKEN")
	endpoint := os.Getenv("GITLAB_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = defaultBaseURL
	}

	logrus.Info("Starting dredd...")
	logrus.Infof("GitLab endpoint: %s", endpoint)

	if *insecureClient {
		logrus.Warn("SSL verification disabled")
		netTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		httpClient.Transport = netTransport
	}

	config, err := LoadFromFile(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	client := gitlab.NewClient(httpClient, token)
	err = client.SetBaseURL(endpoint)
	if err != nil {
		logrus.Fatal(err)
	}

	dredd := &Dredd{
		GitLab: client,
		Config: config,
		DryRun: *dryRun,
	}
	err = dredd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Done")
}
