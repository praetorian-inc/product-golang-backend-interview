package enumeration

import (
	"testing"
)

func Test_enumerateSubdomains(t *testing.T) {
	type args struct {
		rootDomain string
	}
	tests := []struct {
		name       string
		args       args
		minResults int
		wantErr    bool
	}{
		{
			name:       "Successful Enumeration",
			args:       args{rootDomain: "praetorianlabs.com"},
			minResults: 5,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subdomains(tt.args.rootDomain)
			if len(got) < tt.minResults {
				t.Errorf("enumerateSubdomains() got = %v, minResults %d", got, tt.minResults)
			}
		})
	}
}
