package hmac

import (
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	secret := []byte("secret")
	validDateStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	validMapData := map[string]string{
		"date":           validDateStr,
		"(request-line)": "get /hello/world",
	}
	type args struct {
		params    map[string]string
		secretKey []byte
	}
	tests := []struct {
		name        string
		args        args
		wantSigned  string
		wantHeaders string
		wantErr     bool
	}{
		{
			name: "Testcase #1: Signing empty string",
			args: args{
				params:    nil,
				secretKey: secret,
			},
			wantSigned: "",
			wantErr:    true,
		},
		{
			name: "Testcase #2: success",
			args: args{
				params:    validMapData,
				secretKey: secret,
			},
			wantSigned:  "",
			wantHeaders: "date (request-line)",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSigned, gotHeaders, err := Sign(tt.args.params, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSigned == "" && err == nil {
				t.Errorf("Sign() gotSigned = %v, want %v", gotSigned, tt.wantSigned)
			}
			if gotHeaders != tt.wantHeaders {
				t.Errorf("Sign() gotHeaders = [%v], want [%v]", gotHeaders, tt.wantHeaders)
			}
		})
	}
}
