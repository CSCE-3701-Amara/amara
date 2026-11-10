#!/usr/bin/env bash
set -euo pipefail

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required to run Judge0." >&2
  exit 1
fi

cgroup_version=""
if ! cgroup_version="$(docker info --format '{{.CgroupVersion}}' 2>/dev/null)"; then
  echo "The Docker daemon must be running to start Judge0." >&2
  exit 1
fi
if [[ "$cgroup_version" != "1" ]]; then
  cat >&2 <<EOF
Judge0 CE 1.13.1 cannot execute submissions on this host because Docker uses
cgroup v${cgroup_version:-unknown}. Its bundled isolate 1.8.1 requires cgroup v1.

Use a disposable Linux VM configured with cgroup v1 for full Judge0, or run
"make sandbox" for Amara's limited Python/JavaScript/TypeScript development
runner on this host. The development runner uses Docker's cgroup v2 limits.
EOF
  exit 1
fi

cd "$(dirname "$0")/../deployments/judge0"
python3 - <<'PY'
from pathlib import Path
import secrets,os
p=Path('judge0.conf')
if not p.exists():
 config='\n'.join(['JUDGE0_TELEMETRY_ENABLE=false','REDIS_HOST=redis','REDIS_PASSWORD='+secrets.token_hex(32),'POSTGRES_HOST=db','POSTGRES_DB=judge0','POSTGRES_USER=judge0','POSTGRES_PASSWORD='+secrets.token_hex(32),'ALLOW_ENABLE_NETWORK=false','ENABLE_NETWORK=false','CPU_TIME_LIMIT=2','WALL_TIME_LIMIT=5','MEMORY_LIMIT=128000','MAX_PROCESSES_AND_OR_THREADS=32','MAX_FILE_SIZE=64','NUMBER_OF_RUNS=1','RAILS_ENV=production'])+'\n'
 fd=os.open(p,os.O_WRONLY|os.O_CREAT|os.O_EXCL,0o600)
 with os.fdopen(fd,'w') as f:f.write(config)
PY
docker compose -f compose.yml up -d db redis
docker compose -f compose.yml up -d server workers
