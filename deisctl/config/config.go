package config

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deis/deis/deisctl/etcdclient"
	"github.com/deis/deis/deisctl/utils"
)

// fileKeys define config keys to be read from local files
var fileKeys = []string{
	"/deis/platform/sshPrivateKey",
	"/deis/router/sslCert",
	"/deis/router/sslKey",
	"/deis/router/sslDhparam"}

// b64Keys define config keys to be base64 encoded before stored
var b64Keys = []string{"/deis/platform/sshPrivateKey"}

// Config runs the config subcommand
func Config(target string, action string, key []string) error {
	client, err := etcdclient.GetEtcdClient()
	if err != nil {
		return err
	}

	return doConfig(target, action, key, client, os.Stdout)
}

// CheckConfig looks for a value at a keyspace path
// and returns an error if a value is not found
func CheckConfig(root string, k string) error {

	client, err := etcdclient.GetEtcdClient()
	if err != nil {
		return err
	}

	_, err = doConfigGet(client, root, []string{k})
	if err != nil {
		return err
	}

	return nil
}

func doConfig(target string, action string, key []string, client etcdclient.Client, w io.Writer) error {
	rootPath := "/deis/" + target + "/"

	var vals []string
	var err error

	switch action {
	case "rm":
		vals, err = doConfigRm(client, rootPath, key)
	case "set":
		vals, err = doConfigSet(client, rootPath, key)
	default:
		vals, err = doConfigGet(client, rootPath, key)
	}
	if err != nil {
		return err
	}

	// print results
	for _, v := range vals {
		fmt.Fprintf(w, "%v\n", v)
	}
	return nil
}

func doConfigSet(client etcdclient.Client, root string, kvs []string) ([]string, error) {
	var result []string

	for _, kv := range kvs {

		// split k/v from args
		split := strings.SplitN(kv, "=", 2)
		k, v := split[0], split[1]

		// prepare path and value
		path := root + k
		val, err := valueForPath(path, v)
		if err != nil {
			return result, err
		}

		// set key/value in etcd
		ret, err := client.Set(path, val)
		if err != nil {
			return result, err
		}
		result = append(result, ret)

	}
	return result, nil
}

func doConfigGet(client etcdclient.Client, root string, keys []string) ([]string, error) {
	var result []string
	for _, k := range keys {
		val, err := client.Get(root + k)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

func doConfigRm(client etcdclient.Client, root string, keys []string) ([]string, error) {
	var result []string
	for _, k := range keys {
		err := client.Delete(root + k)
		if err != nil {
			return result, err
		}
		result = append(result, k)
	}
	return result, nil
}

// valueForPath returns the canonical value for a user-defined path and value
func valueForPath(path string, v string) (string, error) {

	// check if path is part of fileKeys
	for _, p := range fileKeys {

		if path == p {

			// read value from filesystem
			bytes, err := ioutil.ReadFile(utils.ResolvePath(v))
			if err != nil {
				return "", err
			}

			// see if we should return base64 encoded value
			for _, pp := range b64Keys {
				if path == pp {
					return base64.StdEncoding.EncodeToString(bytes), nil
				}
			}

			return string(bytes), nil
		}
	}

	return v, nil

}
