package http

import (
	url2 "net/url"
	"testing"
)

func TestHTTP_IsURLValid(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestHTTP_IsURLValid 1",
			fields: fields{
				url: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "TestHTTP_IsURLValid 2",
			fields: fields{
				url: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "TestHTTP_IsURLValid 3",
			fields: fields{
				url: "https://example.com/file/path/",
			},
			wantErr: false,
		},
		{
			name: "TestHTTP_IsURLValid 4",
			fields: fields{
				url: "https://example.com/file/path",
			},
			wantErr: false,
		},
		{
			name: "TestHTTP_IsURLValid 4",
			fields: fields{
				url: "https://example.com/file/path?a=b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := url2.Parse(tt.fields.url); (err != nil) != tt.wantErr {
				t.Errorf("IsURLValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapper_Reachable(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestWrapper_Reachable 1",
			fields: fields{
				url: "https://curl.se/dlwiz?type=bin",
			},
			wantErr: true, // NOTE: it will return 404
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wrapper{
				url: tt.fields.url,
			}
			if err := w.Reachable(); (err != nil) != tt.wantErr {
				t.Errorf("Reachable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapper_ReachableByCurl(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestWrapper_ReachableByCurl 1",
			fields: fields{
				url: "https://curl.se/dlwiz?type=bin",
			},
			wantErr: true, // NOTE: it will return 404
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wrapper{
				url: tt.fields.url,
			}
			if err := w.ReachableByCurl(); (err != nil) != tt.wantErr {
				t.Errorf("ReachableByCurl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
