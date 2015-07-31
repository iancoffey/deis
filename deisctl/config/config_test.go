package config

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/deis/deis/deisctl/etcdclient"
	"github.com/deis/deis/deisctl/test/mock"
)

func TestGetConfig(t *testing.T) {
	t.Parallel()

	testMock := mock.Client{Expected: []*etcdclient.ServiceKey{{Key: "/deis/controller/testing", Value: "foo"}, {Key: "/deis/controller/port", Value: "8000"}}}
	testWriter := bytes.Buffer{}

	err := doConfig("controller", "get", []string{"testing", "port"}, testMock, &testWriter)

	if err != nil {
		t.Fatal(err)
	}

	expected := "foo\n8000\n"
	output := testWriter.String()
	if output != expected {
		t.Error(fmt.Errorf("Expected: '%s', Got:'%s'", expected, output))
	}
}

func TestGetConfigError(t *testing.T) {
	t.Parallel()

	testMock := mock.Client{Expected: []*etcdclient.ServiceKey{{Key: "/deis/controller/testing", Value: "foo"}}}
	testWriter := bytes.Buffer{}

	err := doConfig("controller", "get", []string{"port"}, testMock, &testWriter)

	if err == nil {
		t.Fatal("Error Expected")
	}
}

func TestSetConfig(t *testing.T) {
	t.Parallel()

	testMock := mock.Client{Expected: []*etcdclient.ServiceKey{{Key: "/deis/controller/testing", Value: "foo"}, {Key: "/deis/controller/port", Value: "8000"}}}
	testWriter := bytes.Buffer{}

	err := doConfig("controller", "set", []string{"testing=bar", "port=1000"}, testMock, &testWriter)

	if err != nil {
		t.Fatal(err)
	}

	expected := "bar\n1000\n"
	output := testWriter.String()
	if output != expected {
		t.Error(fmt.Errorf("Expected: '%s', Got:'%s'", expected, output))
	}
}

func TestDeleteConfig(t *testing.T) {
	t.Parallel()

	testMock := mock.Client{Expected: []*etcdclient.ServiceKey{{Key: "/deis/controller/testing", Value: "foo"}, {Key: "/deis/controller/port", Value: "8000"}}}
	testWriter := bytes.Buffer{}

	err := doConfig("controller", "rm", []string{"testing", "port"}, testMock, &testWriter)

	if err != nil {
		t.Fatal(err)
	}

	expected := "testing\nport\n"
	output := testWriter.String()
	if output != expected {
		t.Error(fmt.Errorf("Expected: '%s', Got:'%s'", expected, output))
	}
}

// TestConfigSSHPrivateKey ensures private keys are base64 encoded from file path
func TestConfigSSHPrivateKey(t *testing.T) {
	t.Parallel()

	f, err := writeTempFile("private-key")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/platform/sshPrivateKey", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte("private-key"))

	if val != encoded {
		t.Fatalf("expected: %v, got: %v", encoded, val)
	}
}

func TestConfigRouterKey(t *testing.T) {
	t.Parallel()

	f, err := writeTempFile("router-key")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/router/sslKey", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if val != "router-key" {
		t.Fatalf("expected: router-key, got: %v", val)
	}

}

func TestConfigRouterCert(t *testing.T) {
	t.Parallel()

	f, err := writeTempFile("router-cert")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/router/sslCert", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if val != "router-cert" {
		t.Fatalf("expected: router-cert, got: %v", val)
	}

}

func writeTempFile(data string) (*os.File, error) {
	f, err := ioutil.TempFile("", "deisctl")
	if err != nil {
		return nil, err
	}

	f.Write([]byte(data))
	defer f.Close()

	return f, nil
}
