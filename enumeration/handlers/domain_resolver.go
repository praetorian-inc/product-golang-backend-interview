package handlers

import (
	"errors"
	"net"
)

func domainResolves(domain string) (bool, error) {
	iprecords, err := net.LookupIP(domain)
	var dnsError *net.DNSError

	if err != nil {
		if errors.As(err, &dnsError) {
			if dnsError.IsNotFound {
				return false, nil
			}
		}
		return false, err
	}

	if len(iprecords) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
