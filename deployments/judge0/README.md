# Local code sandbox

Use a disposable Linux VM for the Judge0 sandbox host. Judge0's official image uses privileged containers and requires a compatible cgroup configuration; Docker containers are not themselves VMs. The Amara API never receives the Docker socket.

From the repository root:

```sh
./scripts/start-judge0.sh
curl http://127.0.0.1:2358/languages
```

The script generates ignored, mode-600 database/Redis credentials once and starts the pinned Judge0 CE 1.13.1 image, PostgreSQL and Redis. The HTTP listener is bound to loopback. Submission networking is disabled both in server configuration and in Amara's requests. Keep this local listener private; use authenticated HTTPS when deploying Judge0 on another host.

Configure **both API and worker** when running them on the same development host:

```sh
export APP_ENV=development
export JUDGE0_URL=http://127.0.0.1:2358
export JUDGE0_ALLOW_LOCAL_HTTP=true
# Start the normal API and worker with the existing database, SMTP and MFA configuration.
```

For containerized Amara use a reachable authenticated HTTPS Judge0 endpoint; a container's loopback is not the host. `JUDGE0_AUTH_TOKEN` uses `X-Auth-Token`. Match `JUDGE0_LANGUAGE_IDS` to `/languages` on your installation. Defaults are Python 71, JavaScript 63, TypeScript 74, Go 60 and C++ 54.

Candidates can run code with their own stdin without changing saved answers. Submission queues hidden test cases for worker grading; provider outages retain pending work. Legacy percentage modules use test weights; minimum-count modules retain their existing question/case requirements. Historical percentage attempts without weighted scores are eligible for trusted regrading; terminal hiring decisions are not reversed.

Stop the local services without deleting data:

```sh
docker compose -f deployments/judge0/compose.yml down
```

Official deployment prerequisites and security release notes: https://github.com/judge0/judge0/releases/tag/v1.13.1

## Lightweight development runner

For Python, JavaScript and TypeScript without the full Judge0 compiler image:

```sh
docker pull python:3.12-slim
docker pull node:24-alpine
go run ./cmd/sandbox
```

Point both Amara processes at `JUDGE0_URL=http://127.0.0.1:2359` with `APP_ENV=development`, `JUDGE0_ALLOW_LOCAL_HTTP=true` and `JUDGE0_LANGUAGE_IDS='{"python":71,"javascript":63,"typescript":74}'`.

This separate local process implements the Judge0 submission/polling protocol using real Docker containers. Each run has no network, a read-only root, dropped capabilities, no privilege escalation, a non-root UID, 128 MiB memory, 32 processes, one CPU and an eight-second in-container deadline with a one-second forced-kill grace period, plus a ten-second coordinator deadline. It limits concurrency to two, caps output and retained receipts, and removes its container and temporary source on completion. It uses only pre-pulled, allowlisted images; input cannot select a Docker image or flags. It never mounts the Docker socket into candidate containers. Stop it with Ctrl-C. TypeScript uses Node's native type stripping; full TypeScript compilation and Go/C++ require Judge0. This is a development service, not a replacement for a separately secured production sandbox host.

Exercise actual Docker execution through Amara's Judge0 client:

```sh
AMARA_SANDBOX_URL=http://127.0.0.1:2359 go test ./internal/invw/adapters/judge0 -run TestDockerDevelopmentSandbox -v
```
