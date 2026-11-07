"use client";
import Link from "next/link";
import { useId, useState, type FormEvent } from "react";
import { queryString, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import type { FeedFilters, FeedJob, FeedPage } from "../types/cand.types";
import styles from "./candidate.module.css";
export function salaryLabel(job: FeedJob) {
  if (
    !job.salary_visible ||
    (job.salary_min === null && job.salary_max === null)
  )
    return "Salary not disclosed";
  const format = new Intl.NumberFormat("en", {
    notation: "compact",
    maximumFractionDigits: 1,
  });
  return `${job.currency || ""} ${job.salary_min === null ? "—" : format.format(job.salary_min)}–${job.salary_max === null ? "—" : format.format(job.salary_max)}`.trim();
}
export function employmentLabel(value: string | null) {
  if (!value) return "Employment type not specified";
  return (
    value
      ?.replace("on-site ", "")
      .replace("remote ", "")
      .replace("full-time", "Full-time")
      .replace("part-time", "Part-time") +
      (value?.startsWith("remote")
        ? " · Remote"
        : value?.startsWith("on-site")
          ? " · On-site"
          : "") || "Employment type not specified"
  );
}
function JobCard({ job }: { job: FeedJob }) {
  return (
    <article className={styles.jobCard} data-figma-node="363:584">
      <header>
        <h3>{job.title}</h3>
        <span>{job.location}</span>
      </header>
      <p className={styles.companyName}>
        {job.company_name || "Company profile unavailable"}
      </p>
      <p className={styles.salary}>{salaryLabel(job)}</p>
      <div className={styles.jobTags}>
        <span>{employmentLabel(job.employment_type)}</span>
        {job.seniority_level && <span>{job.seniority_level}</span>}
        {job.language && <span>{job.language}</span>}
      </div>
      <Link
        href={`/candidate/jobs/${resourceId(job.id)}${queryString({ title: job.title })}`}
        className={s.button}
      >
        View opportunity
      </Link>
    </article>
  );
}
/** Candidate feed uses only the registered candidate endpoint and its cursor. */
export function JobFeed() {
  const searchId = useId();
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
    <section data-figma-node="363:495">
      <header className={styles.feedHeader}>
        <h1>Find your next opportunity</h1>
        <p>
          Explore open roles and find a position that matches your experience.
        </p>
      </header>
      <form onSubmit={search} aria-label="Filter roles">
        <label htmlFor={searchId} className={styles.searchLabel}>
          Search
        </label>
        <div className={styles.searchRow}>
          <input
            id={searchId}
            name="title"
            type="search"
            aria-label="Search job titles"
            placeholder="Search job titles…"
          />
          <button className={s.primary} type="submit">
            Search jobs
          </button>
        </div>
        <div className={styles.feedLayout}>
          <aside className={styles.filterPanel} aria-label="Job filters">
            <h2>Filters</h2>
            <label>
              Language
              <input name="language" placeholder="Any language" />
            </label>
            <label>
              Currency
              <input name="currency" placeholder="Any currency" />
            </label>
            <label>
              Seniority level
              <input name="seniority_level" placeholder="Any level" />
            </label>
            <label>
              Employment type
              <select name="employment_type">
                <option value="">Any type</option>
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
            <label>
              Location
              <input name="location" placeholder="City or remote" />
            </label>
            <label>
              Minimum salary
              <input
                name="salary_min"
                type="number"
                min="0"
                max="2147483647"
                placeholder="No minimum"
              />
            </label>
            <label>
              Maximum salary
              <input
                name="salary_max"
                type="number"
                min="0"
                max="2147483647"
                placeholder="No maximum"
              />
            </label>
          </aside>
          <div className={styles.results}>
            <BusyState loading={feed.loading} />
            <ApiFeedback error={feed.error} retry={feed.reload} />
            {feed.data && (
              <>
                <h2>
                  {feed.data.items.length}{" "}
                  {feed.data.items.length === 1
                    ? "opportunity"
                    : "opportunities"}
                </h2>
                <p>Explore roles currently available to apply for.</p>
                <div className={styles.jobs}>
                  {feed.data.items.map((job) => (
                    <JobCard key={job.id} job={job} />
                  ))}
                </div>
                {!feed.data.items.length && (
                  <p role="status">
                    No roles match these filters. Try a different search.
                  </p>
                )}
                <div className={styles.pagination}>
                  {history.length > 0 && (
                    <button
                      type="button"
                      onClick={() => {
                        const previous = history.at(-1);
                        setHistory(history.slice(0, -1));
                        setFilter({ ...filter, cursor: previous });
                      }}
                    >
                      Previous page
                    </button>
                  )}
                  <button
                    type="button"
                    disabled={!feed.data.next_cursor}
                    onClick={() => {
                      setHistory([...history, filter.cursor]);
                      setFilter({ ...filter, cursor: feed.data?.next_cursor });
                    }}
                  >
                    Load more opportunities
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </form>
    </section>
  );
}
