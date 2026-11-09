#!/usr/bin/env bash
set -euo pipefail
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
