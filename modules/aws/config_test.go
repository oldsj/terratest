package aws

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAwsEndpointsOverridesShouldNotAllowSetCustomoEndpointsMoreThanOnce(t *testing.T) {
	t.Parallel()
	var localEndpoints = map[string]string{
		"apigateway":     "http://localhost:4566",
		"apigatewayv2":   "http://localhost:4566",
		"cloudformation": "http://localhost:4566",
		"cloudwatch":     "http://localhost:4566",
		"dynamodb":       "http://localhost:4566",
		"es":             "http://localhost:4566",
		"firehose":       "http://localhost:4566",
		"iam":            "http://localhost:4566",
		"kinesis":        "http://localhost:4566",
		"lambda":         "http://localhost:4566",
		"route53":        "http://localhost:4566",
		"redshift":       "http://localhost:4566",
		"s3":             "http://localhost:4566",
		"secretsmanager": "http://localhost:4566",
		"ses":            "http://localhost:4566",
		"sns":            "http://localhost:4566",
		"sqs":            "http://localhost:4566",
		"ssm":            "http://localhost:4566",
		"stepfunctions":  "http://localhost:4566",
		"sts":            "http://localhost:4566",
		"ec2":            "http://localhost:4566",
	}

	// Run the crashing code when DUPLICATE_ENDPOINTS is set
	if os.Getenv("DUPLICATE_ENDPOINTS") == "1" {
		SetAwsEndpointsOverrides(localEndpoints)

		// Subject under testing is try to force scenario to double set custom endpoints
		SetAwsEndpointsOverrides(localEndpoints)
		return
	}

	// Start the actual test in a different subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestSetAwsEndpointsOverrides")
	cmd.Env = append(os.Environ(), "DUPLICATE_ENDPOINTS=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Check that the log  message is what we expected
	gotBytes, _ := ioutil.ReadAll(stdout)
	got := string(gotBytes)
	messageExpected := "The aws custom endpoints have been already setup, check your configuration. The global endpoint override is just allow one time"
	messageReceived := got[:len(got)-1]

	assert.True(t, strings.HasSuffix(messageReceived, messageExpected), messageReceived)

}
