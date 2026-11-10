// sandbox is a loopback-only Docker execution service for local development.
// It implements the subset of Judge0 used by Amara; deploy Judge0 on a dedicated
// sandbox host for production and additional compiled languages.
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type result struct {
	Status struct {
		ID int `json:"id"`
	} `json:"status"`
	Stdout  string    `json:"stdout"`
	Stderr  string    `json:"stderr"`
	Created time.Time `json:"-"`
}

type service struct {
	mu    sync.Mutex
	jobs  map[string]result
	slots chan struct{}
}

type submission struct {
	Language int    `json:"language_id"`
	Source   string `json:"source_code"`
	Input    string `json:"stdin"`
	Network  bool   `json:"enable_network"`
}

type output struct{ data []byte }

const (
	maxRequestBytes         = 128 << 10
	maxSourceBytes          = 64 << 10
	maxInputBytes           = 16 << 10
	maxOutputBytes          = 128 << 10
	jobRetention            = 10 * time.Minute
	maxJobs                 = 1000
	executionTimeout        = 10 * time.Second
	containerCleanupTimeout = 5 * time.Second
)

const (
	languageJavaScript = 63
	languagePython     = 71
	languageTypeScript = 74
	statusInQueue      = 1
	statusProcessing   = 13
	statusAccepted     = 3
	statusTimeLimit    = 5
	statusRuntimeError = 11
)

type languageConfig struct {
	image    string
	program  []string
	filename string
}

var languageConfigs = map[int]languageConfig{
	languagePython: {
		image:    "python:3.12-slim",
		program:  []string{"python", "-I", "/code.py"},
		filename: "/code.py",
	},
	languageJavaScript: {
		image:    "node:24-alpine",
		program:  []string{"node", "/code.js"},
		filename: "/code.js",
	},
	languageTypeScript: {
		image:    "node:24-alpine",
		program:  []string{"node", "--experimental-strip-types", "/code.ts"},
		filename: "/code.ts",
	},
}

func (w *output) Write(p []byte) (int, error) {
	if len(w.data)+len(p) > maxOutputBytes {
		return 0, errors.New("output limit exceeded")
	}
	w.data = append(w.data, p...)
	return len(p), nil
}

func (s *service) submit(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "" {
		http.Error(w, "browser requests must use Amara", http.StatusForbidden)
		return
	}

	input, source, stdin, status := decodeSubmission(w, r)
	if status != 0 {
		http.Error(w, statusMessage(status), status)
		return
	}
	if _, ok := languageConfigs[input.Language]; !ok {
		http.Error(w, "language unavailable in development runner", http.StatusUnprocessableEntity)
		return
	}
	if !s.tryAcquireSlot() {
		http.Error(w, "execution busy", http.StatusServiceUnavailable)
		return
	}

	token := rand.Text()
	job := newJob(statusInQueue)
	if !s.storeJob(token, job) {
		<-s.slots
		http.Error(w, "queue full", http.StatusServiceUnavailable)
		return
	}

	go s.execute(token, input.Language, source, stdin)
	writeJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func decodeSubmission(w http.ResponseWriter, r *http.Request) (submission, []byte, []byte, int) {
	var input submission
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBytes)).Decode(&input) != nil || input.Network {
		return submission{}, nil, nil, http.StatusBadRequest
	}
	source, err := decodeInput(input.Source, maxSourceBytes)
	if err != nil {
		return submission{}, nil, nil, http.StatusBadRequest
	}
	stdin, err := decodeInput(input.Input, maxInputBytes)
	if err != nil {
		return submission{}, nil, nil, http.StatusBadRequest
	}
	return input, source, stdin, 0
}

func decodeInput(encoded string, maxBytes int) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(decoded) > maxBytes {
		return nil, errors.New("input exceeds limit")
	}
	return decoded, nil
}

func statusMessage(status int) string {
	if status == http.StatusBadRequest {
		return "invalid submission"
	}
	return http.StatusText(status)
}

func newJob(statusID int) result {
	job := result{Created: time.Now()}
	job.Status.ID = statusID
	return job
}

func (s *service) tryAcquireSlot() bool {
	select {
	case s.slots <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *service) storeJob(token string, job result) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeExpiredJobs()
	if len(s.jobs) >= maxJobs {
		return false
	}
	s.jobs[token] = job
	return true
}

func (s *service) removeExpiredJobs() {
	for id, job := range s.jobs {
		if time.Since(job.Created) > jobRetention {
			delete(s.jobs, id)
		}
	}
}

func (s *service) execute(token string, language int, source, stdin []byte) {
	defer func() { <-s.slots }()
	job := newJob(statusProcessing)
	defer func() {
		s.mu.Lock()
		s.jobs[token] = job
		s.mu.Unlock()
	}()

	dir, err := os.MkdirTemp("", "amara-code-")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)
	path := dir + "/source"
	if os.WriteFile(path, source, 0600) != nil {
		return
	}
	// #nosec G302 -- Only this generated file inside a private 0700 directory is mounted read-only for the container UID.
	if os.Chmod(path, 0444) != nil {
		return
	}

	config := languageConfigs[language]
	name := "amara-run-" + token
	defer removeContainer(name)

	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()
	args := dockerArgs(name, config, path)
	// #nosec G204 -- Fixed Docker binary and allowlisted language config; name/path are server-generated, and candidate code is read-only file data, never a shell argument.
	cmd := exec.CommandContext(ctx, "docker", args...)
	reader, writer := io.Pipe()
	cmd.Stdin = reader
	go func() { _, _ = writer.Write(stdin); _ = writer.Close() }()
	defer reader.Close()
	var out, stderr output
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	job.Stdout = base64.StdEncoding.EncodeToString(out.data)
	job.Stderr = base64.StdEncoding.EncodeToString(stderr.data)
	if ctx.Err() != nil {
		job.Status.ID = statusTimeLimit
		return
	}
	if err == nil {
		job.Status.ID = statusAccepted
		return
	}
	var exited *exec.ExitError
	if errors.As(err, &exited) && (exited.ExitCode() == 124 || exited.ExitCode() == 137) {
		job.Status.ID = statusTimeLimit
		return
	}
	if errors.As(err, &exited) && exited.ExitCode() != 125 {
		job.Status.ID = statusRuntimeError
	}
}

func dockerArgs(name string, config languageConfig, sourcePath string) []string {
	program := append([]string{"timeout", "-k", "1s", "8s"}, config.program...)
	return append([]string{
		"run", "--rm", "--pull=never", "--name", name,
		"--label", "amara.sandbox=true", "--network", "none", "--read-only",
		"--cap-drop", "ALL", "--security-opt", "no-new-privileges",
		"--pids-limit", "32", "--memory", "128m", "--memory-swap", "128m",
		"--cpus", "1", "--user", "65534:65534",
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=16m",
		"--volume", sourcePath + ":" + config.filename + ":ro,Z", "-i", config.image,
	}, program...)
}

func removeContainer(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), containerCleanupTimeout)
	defer cancel()
	// #nosec G204 -- Fixed Docker command; name is a server-generated random token.
	_ = exec.CommandContext(ctx, "docker", "rm", "-f", name).Run()
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (s *service) submission(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	job, ok := s.jobs[r.PathValue("token")]
	s.mu.Unlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func languages(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		{ID: languagePython, Name: "Python"},
		{ID: languageJavaScript, Name: "JavaScript"},
		{ID: languageTypeScript, Name: "TypeScript (native Node type stripping)"},
	})
}

func main() {
	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal("Docker CLI is required")
	}
	s := &service{jobs: map[string]result{}, slots: make(chan struct{}, 2)}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /submissions", s.submit)
	mux.HandleFunc("GET /submissions/{token}", s.submission)
	mux.HandleFunc("GET /languages", languages)
	server := http.Server{Addr: "127.0.0.1:2359", Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second}
	log.Print("Development sandbox listening on 127.0.0.1:2359")
	log.Fatal(server.ListenAndServe())
}
