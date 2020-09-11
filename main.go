package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

var (
	insecureClient = flag.Bool("k", false, "Disable SSL verification")
	configFile     = flag.String("config", "/etc/gitlab-dredd/gitlab-dredd.yaml", "Path to configuration file")
	dryRun         = flag.Bool("dry-run", false, "Runs without making changes")
	workMode       = flag.String("mode", "plugin", "Work mode (plugin, standalone, webhook)")
	logLevel       = flag.String("log-level", "INFO", "Level of logging (trace, debug, info, warning, error, fatal, panic)")
	logFormat      = flag.String("log-format", "text", "Output logs format (text or json)")

	netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		DisableKeepAlives:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       20 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport,
	}
)

func main() {
	flag.Parse()

	level, err := logrus.ParseLevel(*logLevel)
	if err == nil {
		logrus.SetLevel(level)
	}

	switch strings.ToLower(*logFormat) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
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

	client, err := gitlab.NewClient(config.GitLabToken, gitlab.WithBaseURL(config.GitLabEndpoint), gitlab.WithHTTPClient(httpClient))
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
