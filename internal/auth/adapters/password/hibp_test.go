package password

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHIBPCheckerIsCompromised(t *testing.T) {
	password := "correct horse battery staple"
	sum := sha1.Sum([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(sum[:]))
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/"+hash[:5] {
			t.Fatalf("unexpected path: %s", request.URL.Path)
		}
		_, _ = writer.Write([]byte(hash[5:] + ":7\r\n"))
	}))
	defer server.Close()

	checker, err := NewHIBPChecker(server.Client())
	if err != nil {
		t.Fatal(err)
	}
	checker.baseURL = server.URL + "/"

	compromised, err := checker.IsCompromised(context.Background(), password)
	if err != nil {
		t.Fatal(err)
	}
	if !compromised {
		t.Fatal("expected compromised password")
	}
}

func TestHIBPCheckerReturnsFalseWhenSuffixIsAbsent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = writer.Write([]byte("00000000000000000000000000000000000:1\r\n"))
	}))
	defer server.Close()

	checker, err := NewHIBPChecker(server.Client())
	if err != nil {
		t.Fatal(err)
	}
	checker.baseURL = server.URL + "/"

	compromised, err := checker.IsCompromised(context.Background(), "unique password")
	if err != nil {
		t.Fatal(err)
	}
	if compromised {
		t.Fatal("expected password not to be compromised")
	}
}
