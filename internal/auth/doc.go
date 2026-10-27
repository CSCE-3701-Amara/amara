// Package auth owns accounts, sessions, roles, and company-scoped access.
package auth

import (
	"bufio"
	"strings"
	"strconv"
	"net/mail"
	"math/rand"
	"crypto/sha1"
	"golang.org/x/crypto/argon2"
	"github.com/CSCE-3701-Amara/amara/db"
)

func hashPassword(password string, []byte salt) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 1, 32)
}

func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

func isUniqueEmail(email string) bool {
	_, err := db.GetEmail(email)
	return errors.Is(err, pgx.ErrNoRows)
}

func isCompromised(password string) (bool, err) {
	sum := sha1.Sum([]byte(password))
	hash := strings.ToUpper(hex.EncodeToString(sum[:]))

	prefix = hash[:5]
	suffix = hash[5:]

	url = "https://api.pwnedpasswords.com/range/" + prefix
	req := http.NewRequest("Get", url, nil)
	req.Header.Set("Add-Padding", "true")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HIBP returned code %d", res.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 || !strings.EqualFold(parts[0], suffix) {
			continue
		}
		count, err := strconv.Atoi(parts[1])
		return count > 0, err
	}
	return false, nil
}

func isValidPassword(password string) bool {
	length := len(password)
	return length >= 8 && length <= 128 && !isCompromised(password)
}

func CreateCandidateAccount(name, email, password) (db.Candidate, error) {
	if !isValidEmail(name) {
		return nil, fmt.Error("Invalid email format")
	}

	if !isUniqueEmail(email) {
		return nil, fmt.Error("Email already exists")
	}

	if !isValidPassword(password) {
		return nil, fmt.Error("Invalid password")
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	password_hash := hashPassword(password, salt)

	return db.CreateCandidate(name, email, password_hash, salt)
}

func CreateInterviewerAccount(user InterviewerInfo) (Interviewer, error) {
	if !isValidEmail(user.Email) {
		return nil, fmt.Error("Invalid email format")
	}

	if !isUniqueEmail(user.Email) {
		return nil, fmt.Error("Email already exists")
	}

	if !isValidPassword(user.Password) {
		return nil, fmt.Error("Invalid password")
	}

	// TODO:
	interviewer, err := db.CreateInterviewer()

	return interviewer, err
}