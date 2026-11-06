package password

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type HIBPChecker struct {
	client  *http.Client
	baseURL string
}

func NewHIBPChecker(client *http.Client) (*HIBPChecker, error) {
	if client == nil {
		return nil, fmt.Errorf("HTTP client must not be nil")
	}
	return &HIBPChecker{client: client, baseURL: "https://api.pwnedpasswords.com/range/"}, nil
}

func (c *HIBPChecker) IsCompromised(ctx context.Context, password string) (bool, error) {
	sum := sha1.Sum([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(sum[:]))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+hash[:5], nil)
	if err != nil {
		return false, fmt.Errorf("create password compromise request: %w", err)
	}
	request.Header.Set("Add-Padding", "true")

	response, err := c.client.Do(request)
	if err != nil {
		return false, fmt.Errorf("request password compromise service: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("password compromise service returned HTTP %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return false, fmt.Errorf("read password compromise response: %w", err)
	}
	for _, line := range strings.Split(string(body), "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], hash[5:]) {
			continue
		}
		count, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return false, fmt.Errorf("parse password compromise count: %w", err)
		}
		return count > 0, nil
	}
	return false, nil
}
