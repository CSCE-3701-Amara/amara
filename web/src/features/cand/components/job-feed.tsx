"use client";
import Link from "next/link";
import { useRef, useState, type FormEvent } from "react";
import { queryString } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import { applyToPosting } from "@/features/job/api/job.api";
import type { ApplicationSubmission } from "@/features/job/types/job.types";
import type { FeedFilters, FeedJob, FeedPage } from "../types/cand.types";

function JobCard({ job }: { job: FeedJob }) {
  const [result, setResult] = useState<ApplicationSubmission>();
  const [error, setError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const pending = useRef(false);
  async function apply() {
    if (pending.current) return;
    pending.current = true;
    setBusy(true);
    setError(undefined);
    try {
      setResult(await applyToPosting(job.id));
    } catch (cause) {
      setError(cause as Error);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return (
    <article className={s.card}>
      <h2>{job.title}</h2>
      <p className={s.muted}>
        {[job.location, job.employment_type, job.seniority_level]
          .filter(Boolean)
          .join(" · ")}
      </p>
      {job.description && <p>{job.description}</p>}
      {job.salary_visible &&
        (job.salary_min !== null || job.salary_max !== null) && (
          <p>
            Salary: {job.salary_min ?? "Not specified"} –{" "}
            {job.salary_max ?? "Not specified"} {job.currency}
          </p>
        )}
      {job.application_deadline && (
        <p>
          Apply by {new Date(job.application_deadline).toLocaleDateString()}
        </p>
      )}
      <ApiFeedback error={error} />
      {result ? (
        <div className={s.success} role="status">
          <p>
            {result.screened_out
              ? "Your application was submitted but did not pass initial screening."
              : "Application submitted."}{" "}
            Reference: {result.reference_number}.
          </p>
          {!!result.failed_amarat?.length && (
            <p>Required skills: {result.failed_amarat.join(", ")}</p>
          )}
          <Link href="/candidate/applications">View my applications</Link>
        </div>
      ) : (
        <button
          type="button"
          className={s.primary}
          disabled={busy}
          onClick={apply}
        >
          {busy ? "Submitting…" : "Apply for this role"}
        </button>
      )}
    </article>
  );
}
export function JobFeed() {
  const [filter, setFilter] = useState<FeedFilters>({ limit: 20 });
  const [history, setHistory] = useState<(string | undefined)[]>([]);
  const feed = useApiResource<FeedPage>(`/v1/jobs/feed${queryString(filter)}`);
  function search(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const values = new FormData(event.currentTarget);
    setHistory([]);
    setFilter({
      limit: 20,
      title: String(values.get("title") || "").trim(),
      location: String(values.get("location") || "").trim(),
      employment_type: String(
        values.get("employment_type") || "",
      ) as FeedFilters["employment_type"],
      language: String(values.get("language") || "").trim(),
      currency: String(values.get("currency") || "").trim(),
      seniority_level: String(values.get("seniority_level") || "").trim(),
      ...Object.fromEntries(
        ["salary_min", "salary_max"]
          .filter((name) => values.get(name) !== "")
          .map((name) => [name, Number(values.get(name))]),
      ),
    });
  }
  return (
    <>
      <h1>Find your next opportunity</h1>
      <p className={s.muted}>Explore roles that value what you can do.</p>
      <form onSubmit={search} className={s.card} aria-label="Filter roles">
        <div className={s.filters}>
          <label className={s.field}>
            Job title
            <input name="title" type="search" placeholder="Search roles" />
          </label>
          <label className={s.field}>
            Location
            <input name="location" placeholder="City or region" />
          </label>
          <label className={s.field}>
            Employment type
            <select name="employment_type">
              <option value="">All types</option>
              {[
                "on-site full-time",
                "on-site part-time",
                "remote full-time",
                "remote part-time",
              ].map((type) => (
                <option key={type}>{type}</option>
              ))}
            </select>
          </label>
          <button className={s.primary} type="submit">
            Search roles
          </button>
        </div>
        <details>
          <summary>More filters</summary>
          <div className={s.filters}>
            <label className={s.field}>
              Language
              <input name="language" />
            </label>
            <label className={s.field}>
              Currency
              <input name="currency" />
            </label>
            <label className={s.field}>
              Seniority
              <input name="seniority_level" />
            </label>
            <label className={s.field}>
              Minimum salary
              <input name="salary_min" type="number" min="0" max="2147483647" />
            </label>
            <label className={s.field}>
              Maximum salary
              <input name="salary_max" type="number" min="0" max="2147483647" />
            </label>
          </div>
        </details>
      </form>
      <BusyState loading={feed.loading} />
      <ApiFeedback error={feed.error} retry={feed.reload} />
      {feed.data && (
        <>
          <div className={s.grid}>
            {feed.data.items.map((job) => (
              <JobCard key={job.id} job={job} />
            ))}
          </div>
          {!feed.data.items.length && (
            <p role="status">
              No roles match these filters. Try a different search.
            </p>
          )}
          <div className={s.actions}>
            <button
              disabled={!history.length}
              onClick={() => {
                const previous = history.at(-1);
                setHistory(history.slice(0, -1));
                setFilter({ ...filter, cursor: previous });
              }}
            >
              Previous page
            </button>
            <button
              disabled={!feed.data.next_cursor}
              onClick={() => {
                setHistory([...history, filter.cursor]);
                setFilter({ ...filter, cursor: feed.data?.next_cursor });
              }}
            >
              Next page
            </button>
          </div>
        </>
      )}
    </>
  );
}
