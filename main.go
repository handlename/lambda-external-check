package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var version = "0.0.1"

const EnvPrefix = "EXTERNAL_CHECK_"

type Config struct {
	Target  string
	Timeout time.Duration
}

func main() {
	var flagVersion bool
	flag.BoolVar(&flagVersion, "v", false, "show version")
	flag.Parse()

	if flagVersion {
		log.Printf("Version: v%s", version)
		os.Exit(0)
	}

	config, err := initConfig()
	if err != nil {
		log.Printf("failed to init config: %s", err)
		os.Exit(1)
	}

	handler := &LambdaHandler{
		Config: config,
	}

	if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambda.Start(handler.HandleRequest)
	} else {
		res, err := handler.HandleRequest(context.Background())
		if err != nil {
			log.Printf("Error: %s", err)
			os.Exit(1)
		}

		log.Printf("Result: %s", res)
	}
}

func initConfig() (*Config, error) {
	target := os.Getenv(EnvPrefix + "TARGET")
	timeout := os.Getenv(EnvPrefix + "TIMEOUT")

	if target == "" || timeout == "" {
		return nil, fmt.Errorf("Missing required environment variables")
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}

	return &Config{
		Target:  target,
		Timeout: timeoutDuration,
	}, nil
}

type LambdaHandler struct {
	Config *Config
}

func (h *LambdaHandler) HandleRequest(ctx context.Context) (string, error) {
	client := http.Client{
		Timeout: h.Config.Timeout,
	}

	res, err := client.Get(h.Config.Target)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 200 {
		return "OK", nil
	}

	return "", fmt.Errorf("Error: StatusCode=%d", res.StatusCode)
}
