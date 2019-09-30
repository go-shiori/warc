package archiver

import (
	nurl "net/url"
	"testing"
)

func Test_toAbsoluteURI(t *testing.T) {
	baseURL, _ := nurl.ParseRequestURI("http://localhost:8080/absolute/")

	tests := []struct {
		name string
		uri  string
		want string
	}{{
		name: "hash url",
		uri:  "#here",
		want: "#here",
	}, {
		name: "relative url based on root",
		uri:  "/test/123",
		want: "http://localhost:8080/test/123",
	}, {
		name: "relative url based on current page",
		uri:  "test/123",
		want: "http://localhost:8080/absolute/test/123",
	}, {
		name: "relative url based on current scheme",
		uri:  "//www.google.com",
		want: "http://www.google.com",
	}, {
		name: "absolute url",
		uri:  "https://www.google.com",
		want: "https://www.google.com",
	}, {
		name: "absolute url with non-http protocol",
		uri:  "ftp://ftp.server.com",
		want: "ftp://ftp.server.com",
	}, {
		name: "domain name only without protocol",
		uri:  "www.google.com",
		want: "http://localhost:8080/absolute/www.google.com",
	}, {
		name: "absolute url but missing colon",
		uri:  "http//www.google.com",
		want: "http://localhost:8080/absolute/http//www.google.com",
	}, {
		name: "relative url based on current page, go up first",
		uri:  "../hello/relative",
		want: "http://localhost:8080/hello/relative",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toAbsoluteURI(tt.uri, baseURL); got != tt.want {
				t.Errorf("toAbsoluteURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
