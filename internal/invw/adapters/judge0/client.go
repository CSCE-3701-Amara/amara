package judge0

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

type Client struct {
	baseURL   string
	token     string
	languages map[string]int
	http      *http.Client
}

// New connects to Judge0 CE using its native token authentication.
func New(baseURL, token string, languages map[string]int, allowLocalHTTP bool) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil || (base.Scheme != "https" && base.Scheme != "http") ||
		base.Hostname() == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" {
		return nil, errors.New("JUDGE0_URL must be an HTTP(S) base URL")
	}
	if base.Scheme == "http" {
		ip := net.ParseIP(base.Hostname())
		if token != "" || !allowLocalHTTP || !ip.IsLoopback() {
			return nil, errors.New("Judge0 requires HTTPS except for explicitly enabled unauthenticated loopback HTTP")
		}
	}
	if languages == nil {
		languages = map[string]int{"javascript": 63, "typescript": 74, "python": 71, "go": 60, "cpp": 54}
	}
	for _, id := range languages {
		if id <= 0 {
			return nil, errors.New("JUDGE0_LANGUAGE_IDS must contain positive language IDs")
		}
	}
	return &Client{
		baseURL: strings.TrimRight(base.String(), "/"), token: token, languages: languages,
		http: &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}},
	}, nil
}

func (c *Client) Run(ctx context.Context, language, source, input string) (string, bool, error) {
	id, ok := c.languages[strings.ToLower(strings.TrimSpace(language))]
	if !ok {
		return "", false, errors.New("saved language has no Judge0 mapping")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	encode := base64.StdEncoding.EncodeToString
	body, err := json.Marshal(map[string]any{
		"language_id": id, "source_code": encode([]byte(source)), "stdin": encode([]byte(input)),
		"enable_network": false, "cpu_time_limit": 2, "wall_time_limit": 5,
		"memory_limit": 128000, "max_file_size": 64,
	})
	if err != nil {
		return "", false, err
	}
	var created struct{ Token string }
	if err := c.request(ctx, http.MethodPost, "/submissions?base64_encoded=true&wait=false", body, &created); err != nil {
		return "", false, err
	}
	if created.Token == "" {
		return "", false, errors.New("Judge0 returned an empty submission token")
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		var result struct {
			Stdout        string
			Stderr        string
			CompileOutput string `json:"compile_output"`
			Status        struct{ ID int }
		}
		path := "/submissions/" + url.PathEscape(created.Token) + "?base64_encoded=true&fields=stdout,stderr,compile_output,status"
		if err := c.request(ctx, http.MethodGet, path, nil, &result); err != nil {
			return "", false, err
		}
		switch result.Status.ID {
		case 1, 2:
			select {
			case <-ctx.Done():
				return "", false, ctx.Err()
			case <-ticker.C:
			}
		case 3:
			out, err := base64.StdEncoding.DecodeString(result.Stdout)
			if err != nil {
				return "", false, errors.New("Judge0 returned invalid output encoding")
			}
			return string(out), true, nil
		case 4, 5, 6, 7, 8, 9, 10, 11, 12:
			diagnostics := ""
			for _, raw := range []string{result.Stdout, result.Stderr, result.CompileOutput} {
				b, err := base64.StdEncoding.DecodeString(raw)
				if err != nil {
					return "", false, errors.New("Judge0 returned invalid output encoding")
				}
				diagnostics += string(b)
			}
			return diagnostics, false, nil
		default:
			return "", false, fmt.Errorf("unexpected Judge0 submission status %d", result.Status.ID)
		}
	}
}

func (c *Client) request(ctx context.Context, method, path string, body []byte, result any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return errors.New("create Judge0 request")
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("Judge0 %s request failed", method)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Judge0 %s returned HTTP %d", method, resp.StatusCode)
	}
	const maxResponse = 256 * 1024
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse+1))
	if err != nil || len(data) > maxResponse || json.Unmarshal(data, result) != nil {
		return errors.New("invalid Judge0 response")
	}
	return nil
}

func (c *Client) SupportedLanguages() []string {
	names := make([]string, 0, len(c.languages))
	for name := range c.languages {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
