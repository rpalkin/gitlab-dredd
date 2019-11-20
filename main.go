package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

var (
	insecureClient = flag.Bool("k", false, "Disable SSL verification")
	configFile     = flag.String("config", "/etc/gitlab-dredd/gitlab-dredd.yaml", "Path to configuration file")
	dryRun         = flag.Bool("dry-run", false, "Runs without making changes")
	workMode       = flag.String("mode", "plugin", "Work mode (plugin, standalone)")
	logLevel       = flag.String("log-level", "INFO", "Level of logging (trace, debug, info, warning, error, fatal, panic)")

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

	config, err := LoadFromFile(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	mode, err := ParseWorkMode(*workMode)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.WithFields(logrus.Fields{
		"mode":     mode,
		"endpoint": config.GitLabEndpoint,
		"insecure": *insecureClient,
	}).Info("Starting gitlab-dredd...")

	if *insecureClient {
		netTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		httpClient.Transport = netTransport
	}

	client := gitlab.NewClient(httpClient, config.GitLabToken)
	err = client.SetBaseURL(config.GitLabEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}

	dredd := &Dredd{
		GitLab: client,
		Config: config,
		DryRun: *dryRun,
	}

	err = dredd.Run(mode)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Done")
}
