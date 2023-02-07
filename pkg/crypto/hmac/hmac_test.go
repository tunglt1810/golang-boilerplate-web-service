package hmac

import (
	"testing"
	"time"
)

func TestHMAC_validateClockSkew(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		dateStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Testcase #1: valid clock skew",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
				},
			},
			args: args{
				dateStr: time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Testcase #2: time expired",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
				},
			},
			args: args{
				dateStr: "Mon, 02 Jan 2022 15:04:05 GMT",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Testcase #3: time parse error",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
				},
			},
			args: args{
				dateStr: "Mon, 02an 202215:04:05 GMT",
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HMAC{
				cfg: tt.fields.cfg,
			}
			got, err := h.validateClockSkew(tt.args.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateClockSkew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateClockSkew() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHMAC_ValidateSignature(t *testing.T) {
	validDateStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	validMapData := map[string]string{
		"date":           validDateStr,
		"(request-line)": "get /hello/world",
	}

	invalidDateStr := "Mon, 02 Jan 2022 15:04:05 GMT"
	invalidMapData := map[string]string{
		"date":           validDateStr,
		"(request-line)": "get /hello/world",
	}
	secret := "secret"
	strStr, _, _ := Sign(validMapData, []byte(secret))
	type fields struct {
		cfg *Config
	}
	type args struct {
		signatureStr string
		dateStr      string
		params       map[string]string
		secret       []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Testcase #1: valid signature",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
					Secret:           secret,
				},
			},
			args: args{
				signatureStr: strStr,
				dateStr:      validDateStr,
				params:       validMapData,
				secret:       []byte(secret),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Testcase #2: invalid signature: wrong secret",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
					Secret:           secret,
				},
			},
			args: args{
				signatureStr: strStr,
				dateStr:      validDateStr,
				params:       invalidMapData,
				secret:       []byte("othersecret"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Testcase #3: invalid signature: expired time",
			fields: fields{
				cfg: &Config{
					AllowedClockSkew: 10,
					Secret:           secret,
				},
			},
			args: args{
				signatureStr: strStr,
				dateStr:      invalidDateStr,
				params:       invalidMapData,
				secret:       []byte("secret"),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HMAC{
				cfg: tt.fields.cfg,
			}
			got, err := h.ValidateSignature(tt.args.signatureStr, tt.args.dateStr, tt.args.params, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}
