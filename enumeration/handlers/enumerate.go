package handlers

import (
	"bytes"
	"context"
	"io"
	"log"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"strings"
)

func enumerateSubdomains(rootDomain string) ([]string, error) {
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		Sources:            passive.DefaultSources, // Use the default list of passive sources
		AllSources:         passive.DefaultAllSources, // Use the default list of all passive sources
		Recursive:          passive.DefaultRecursiveSources,	// Use the default list of recursive sources
		Providers:          &runner.Providers{}, // Use empty api keys for all providers
	})

	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain(context.Background(), rootDomain, []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}

	subdomains := strings.FieldsFunc(string(data), func(c rune) bool {
		return c == '\n'
	})

	return subdomains, nil
}