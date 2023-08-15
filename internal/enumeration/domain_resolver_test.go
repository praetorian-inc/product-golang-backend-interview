package enumeration

import (
	"errors"
	"net"
	"testing"
)

// todo: remove this file

func Test_domainResolves(t *testing.T) {
	t.Skip("This tests functionality that isn't present in the application")
	type args struct {
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Successful Resolve 1",
			args:    args{domain: "google.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Successful Resolve 2",
			args:    args{domain: "scan.praetorianlabs.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Unsuccessful Resolve 1",
			args:    args{domain: "case.praetorianlabs.com"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Unsuccessful Resolve 2",
			args:    args{domain: "vpn.praetorianlabs.com"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domainResolves(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("isLive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isLive() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
