# ADR: Supabase PostgreSQL as the relational data store

Owner: Ziad Eliwa | Status: Accepted | Last updated: 2026-10-28 | Jira: AMARA-1292

## Context

SRS requires a relational, tenant-aware data store; SDD §3.5 defers the engine
choice to an ADR. Hosting targets AWS/Azure, the team is student-sized, and
Supabase was already provisioned as the platform project (`amara-dev`,
West EU).

## Decision

Use **Supabase PostgreSQL** as the authoritative relational store.

- Applications connect over the **Supavisor session pooler** (`:5432`,
  `postgres.<project-ref>@aws-0-eu-west-1.pooler.supabase.com`) for IPv4
  reachability; the direct `db.<ref>.supabase.co` host is IPv6-only on this
  tier. Prepared statements require the session pooler, not the transaction
  (`:6543`) mode.
- Connections use pgx v5 pools configured in `internal/platform`.
- Schema **does not use Supabase's Auth services**; the AUTH epic owns its own
  tables. Because Supabase reserves the `auth` schema, epic schemas are
  `accounts`, `cand`, `job`, `invw`, `ana`.
- Schema migrations are plain SQL under `supabase/migrations/`, applied with
  **Supabase CLI** (`supabase migration new`, `supabase db push`); the
  migration history table (`supabase_migrations`) lives server-side.

## Alternatives considered

- Self-hosted PostgreSQL on the AWS/Azure deployment — more control, but
  operationally too heavy for the T1 team; revisit if Managed services must
  be replaced.
- golang-migrate with embedded SQL — engine-neutral and CI-friendly, but
  requires building migration state management ourselves; Supabase CLI brings
  diffing (`db diff`), reset, and a remote history table for free.
- Direct (IPv6) connection — fails from typical student IPv4 home networks.

## Consequences

- Shared-instance connection limits are small: pools start at `MaxConns=4`
  (`DB_POOL_MAX_CONNS` env), load-test ceiling is 100.
- Reserved-schema rule constrains our epic schema naming.
- CI schema pushes need a `SUPABASE_ACCESS_TOKEN` machine identity (deferred
  to the CI ticket).
- Object storage and auth providers remain separate decisions (TBD-014).
