"use client";
import Link from "next/link";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import type { PostingList } from "../types/job.types";
export function Postings() {
  const jobs = useApiResource<PostingList>("/v1/jobs");
  const [title, setTitle] = useState("");
  const [status, setStatus] = useState("");
  const visible = (jobs.data?.items ?? []).filter(
    (job) =>
      job.title.toLowerCase().includes(title.toLowerCase()) &&
      (!status || job.status === status),
  );
  return (
    <>
      <h1>Your job postings</h1>
      <p className={s.muted}>Manage the roles your team is hiring for.</p>
      <Link href="/employer/jobs/new" className={`${s.button} ${s.primary}`}>
        Create job posting
      </Link>
      <div className={s.filters}>
        <label className={s.field}>
          Search postings
          <input
            type="search"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
          />
        </label>
        <label className={s.field}>
          Posting status
          <select
            value={status}
            onChange={(event) => setStatus(event.target.value)}
          >
            <option value="">All statuses</option>
            {["draft", "published", "archived"].map((value) => (
              <option key={value}>{value}</option>
            ))}
          </select>
        </label>
      </div>
      <BusyState loading={jobs.loading} />
      <ApiFeedback error={jobs.error} retry={jobs.reload} />
      {jobs.data && (
        <>
          <div className={s.grid}>
            {visible.map((job) => (
              <article className={s.card} key={job.id}>
                <h2>{job.title}</h2>
                <span className={s.badge}>{job.status}</span>
                <p>
                  {job.location || "Location not set"} · {job.employment_type}
                </p>
                <Link
                  href={`/employer/jobs/${encodeURIComponent(job.id)}`}
                  className={s.button}
                >
                  Manage posting
                </Link>
              </article>
            ))}
          </div>
          {!visible.length && <p>No postings match these filters.</p>}
        </>
      )}
    </>
  );
}
