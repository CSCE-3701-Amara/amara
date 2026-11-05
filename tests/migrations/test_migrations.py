"""Replay raw Supabase migrations in a disposable local PostgreSQL cluster.

Run with: python3 tests/migrations/test_migrations.py
Requires initdb, pg_ctl, and psql on PATH. No Supabase project is contacted.
"""

import json
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[2]
POSTGRES_TOOLS = ("initdb", "pg_ctl", "psql")


@unittest.skipUnless(
    all(shutil.which(tool) for tool in POSTGRES_TOOLS),
    "Requires initdb, pg_ctl, and psql on PATH",
)
class MigrationReplayTest(unittest.TestCase):
    def test_fresh_install_preserves_schema_and_existing_data(self):
        expected = {}
        diagram = (ROOT / "docs/sdd/database/erd_diagrams/main.mmd").read_text()
        for table, body in re.findall(
            r"^    ([A-Z_]+) \{\n(.*?)^    \}", diagram, re.MULTILINE | re.DOTALL
        ):
            expected[table.lower()] = sorted(
                re.findall(r"^        \w+ (\w+)", body, re.MULTILINE)
            )
        self.assertTrue(expected, "ER diagram must define tables")

        def run(*args):
            result = subprocess.run(args, text=True, capture_output=True)
            self.assertEqual(
                result.returncode,
                0,
                f"Command failed: {args}\n{result.stdout}\n{result.stderr}",
            )
            return result.stdout.strip()

        with tempfile.TemporaryDirectory(prefix="amara-migrations-") as temp:
            data = str(Path(temp) / "data")
            run("initdb", "-D", data, "-A", "trust", "-U", "migration_test",
                "--no-locale", "--encoding=UTF8")
            run("pg_ctl", "-D", data, "-l", str(Path(temp) / "postgres.log"),
                "-o", f"-F -h '' -k {temp} -p 55432", "-w", "start")
            try:
                psql = (
                    "psql", "-X", "-h", temp, "-p", "55432",
                    "-U", "migration_test", "-d", "postgres",
                    "-v", "ON_ERROR_STOP=1", "-At",
                )
                for migration in sorted((ROOT / "supabase/migrations").glob("*.sql")):
                    with self.subTest(migration=migration.name):
                        # Execute the whole file, as Supabase does: no Goose filtering.
                        run(*psql, "--single-transaction", "-f", str(migration))
                        if migration.name == "20261029000007_create_applications_table.sql":
                            run(*psql, "-c", """
                                INSERT INTO candidates
                                    (id, name, email, password_hash, salt, created_at)
                                VALUES ('00000000-0000-0000-0000-000000000001',
                                    'Migration fixture', 'migration@example.invalid',
                                    decode('00', 'hex'), decode('00', 'hex'), now());
                            """)

                actual = json.loads(run(*psql, "-c", """
                    SELECT jsonb_object_agg(table_name, columns)
                    FROM (
                        SELECT table_name, jsonb_agg(column_name ORDER BY column_name) AS columns
                        FROM information_schema.columns
                        WHERE table_schema = 'public'
                        GROUP BY table_name
                    ) AS tables;
                """))
                self.assertEqual(actual, expected)
                self.assertEqual(run(*psql, "-c", """
                    SELECT name || ':' || account_status FROM candidates
                    WHERE id = '00000000-0000-0000-0000-000000000001';
                """), "Migration fixture:pending_verification")
            finally:
                run("pg_ctl", "-D", data, "-m", "fast", "-w", "stop")


if __name__ == "__main__":
    unittest.main()
