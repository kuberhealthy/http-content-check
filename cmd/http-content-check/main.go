package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	nodecheck "github.com/kuberhealthy/kuberhealthy/v3/pkg/nodecheck"
	log "github.com/sirupsen/logrus"
)

// main loads configuration and runs the HTTP content check.
func main() {
	// Enable nodecheck debug output for parity with v2 behavior.
	nodecheck.EnableDebugOutput()

	// Parse configuration from environment variables.
	cfg, err := parseConfig()
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Create a timeout context for readiness checks.
	checkTimeLimit := time.Minute * 1
	ctx, _ := context.WithTimeout(context.Background(), checkTimeLimit)

	// Wait for Kuberhealthy to be reachable before running the check.
	err = nodecheck.WaitForKuberhealthy(ctx)
	if err != nil {
		log.Errorln("Error waiting for kuberhealthy endpoint to be contactable by checker pod with error:", err.Error())
	}

	// Fetch the URL content for inspection.
	log.Infoln("Attempting to fetch content from:", cfg.TargetURL)
	content, err := getURLContent(cfg)
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Search the response body for the target string.
	log.Infoln("Parsing content for string", cfg.TargetString)
	found := findStringInContent(content, cfg.TargetString)
	if !found {
		reportFailureAndExit(fmt.Errorf("could not find string in content"))
		return
	}

	// Report success when the target string is present.
	log.Infoln("Success! Found", cfg.TargetString, "in", cfg.TargetURL)
	err = checkclient.ReportSuccess()
	if err != nil {
		log.Fatalln("error when reporting to kuberhealthy:", err.Error())
	}
	log.Infoln("Successfully reported to Kuberhealthy")
}

// getURLContent fetches the URL response body using the configured timeout.
func getURLContent(cfg *CheckConfig) ([]byte, error) {
	// Build an HTTP client with the configured timeout.
	client := http.Client{Timeout: cfg.TimeoutDuration}

	// Issue the HTTP request.
	resp, err := client.Get(cfg.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url %s: %w", cfg.TargetURL, err)
	}

	// Ensure the response body is closed.
	defer closeResponseBody(resp.Body)

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", cfg.TargetURL, err)
	}

	return body, nil
}

// closeResponseBody closes the response body and logs any error.
func closeResponseBody(body io.ReadCloser) {
	// Close the body when present.
	if body == nil {
		return
	}

	// Close the body and log failures.
	err := body.Close()
	if err != nil {
		log.Errorln("failed to close response body:", err.Error())
	}
}

// findStringInContent returns true when the target string is in the response body.
func findStringInContent(body []byte, target string) bool {
	// Convert the body to a string for substring search.
	content := string(body)

	// Check for the target string.
	return strings.Contains(content, target)
}

// reportFailureAndExit reports a failure to Kuberhealthy and exits the process.
func reportFailureAndExit(err error) {
	// Log the error locally.
	log.Errorln(err)

	// Report the failure to Kuberhealthy.
	reportErr := checkclient.ReportFailure([]string{err.Error()})
	if reportErr != nil {
		log.Fatalln("error when reporting to kuberhealthy:", reportErr.Error())
	}

	// Exit after reporting failure.
	os.Exit(0)
}
