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

func TestNewHIBPCheckerRequiresClient(t *testing.T) {
	if _, err := NewHIBPChecker(nil); err == nil {
		t.Fatal("expected nil client error")
	}
}

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

func TestHIBPCheckerErrors(t *testing.T) {
	t.Run("non-OK response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			http.Error(writer, "unavailable", http.StatusServiceUnavailable)
		}))
		defer server.Close()
		checker, err := NewHIBPChecker(server.Client())
		if err != nil {
			t.Fatal(err)
		}
		checker.baseURL = server.URL + "/"
		if _, err := checker.IsCompromised(context.Background(), "password"); err == nil {
			t.Fatal("expected service status error")
		}
	})

	t.Run("malformed count", func(t *testing.T) {
		password := "malformed count password"
		sum := sha1.Sum([]byte(password))
		hash := strings.ToUpper(hex.EncodeToString(sum[:]))
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			_, _ = writer.Write([]byte(hash[5:] + ":not-a-number\n"))
		}))
		defer server.Close()
		checker, err := NewHIBPChecker(server.Client())
		if err != nil {
			t.Fatal(err)
		}
		checker.baseURL = server.URL + "/"
		if _, err := checker.IsCompromised(context.Background(), password); err == nil {
			t.Fatal("expected malformed count error")
		}
	})

	t.Run("request creation failure", func(t *testing.T) {
		checker, err := NewHIBPChecker(&http.Client{})
		if err != nil {
			t.Fatal(err)
		}
		checker.baseURL = "://invalid"
		if _, err := checker.IsCompromised(context.Background(), "password"); err == nil {
			t.Fatal("expected request creation error")
		}
	})
}
