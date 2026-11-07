"use client";
import Link from "next/link";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import {
  ProductHeading,
  Stats,
  StatusBadge,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
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
      <ProductHeading
        title="Job postings"
        description="Manage the roles your team is hiring for."
        artwork={{ src: "/figma/postings.svg", width: 345.6, height: 238.8 }}
      />
      <section className={e.list}>
        <h2>Your job postings</h2>
        <div className={e.actions}>
          {["", "draft", "published", "archived"].map((value) => (
            <button
              type="button"
              key={value}
              aria-pressed={status === value}
              onClick={() => setStatus(value)}
            >
              {value ? value[0].toUpperCase() + value.slice(1) : "All states"}
            </button>
          ))}
          <Link href="/employer/jobs/new" className={e.primary}>
            Create posting
          </Link>
        </div>
        <label className={s.field}>
          Search postings
          <input
            type="search"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
          />
        </label>
        <BusyState loading={jobs.loading} />
        <ApiFeedback error={jobs.error} retry={jobs.reload} />
        {jobs.data && (
          <>
            <div className={e.list}>
              {visible.map((job) => (
                <article className={e.row} key={job.id}>
                  <div>
                    <h3>{job.title}</h3>
                    <p>
                      {job.employment_type} · {job.seniority_level}
                    </p>
                  </div>
                  <StatusBadge>{job.status}</StatusBadge>
                  <Link
                    href={`/employer/jobs/${encodeURIComponent(job.id)}`}
                    className={s.button}
                  >
                    Edit
                  </Link>
                  <Link
                    href={`/employer/jobs/${encodeURIComponent(job.id)}/preview`}
                    className={s.button}
                  >
                    Preview
                  </Link>
                </article>
              ))}
            </div>
            {!visible.length && <p>No postings match these filters.</p>}
          </>
        )}
      </section>
      <Stats
        items={["draft", "published", "archived"].map((value) => ({
          label: value[0].toUpperCase() + value.slice(1),
          value:
            jobs.data?.items?.filter((job) => job.status === value).length ??
            "—",
        }))}
      />
      <WorkspaceFooter />
    </>
  );
}
