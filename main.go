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
	configFile     = flag.String("config", "/etc/gitlab-dredd/gitlab-dredd.yaml", "Path to configuration file")
	dryRun         = flag.Bool("dry-run", false, "Runs without making changes")
	pluginMode     = flag.Bool("plugin", true, "Runs as a GitLab plugin.")
	logLevel       = flag.String("log-level", "INFO", "Level of logging (trace, debug, info, warning, error, fatal, panic).")

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

	level, err := logrus.ParseLevel(*logLevel)
	if err == nil {
		logrus.SetLevel(level)
	}

	token := os.Getenv("GITLAB_TOKEN")
	endpoint := os.Getenv("GITLAB_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = defaultBaseURL
	}

	logrus.WithFields(logrus.Fields{
		"plugin":   *pluginMode,
		"endpoint": endpoint,
		"insecure": *insecureClient,
	}).Info("Starting gitlab-dredd...")

	if *insecureClient {
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

	if *pluginMode {
		err = dredd.Hook()
		if err != nil {
			logrus.Fatal(err)
		}
		return
	}

	err = dredd.Run()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Done")
}
