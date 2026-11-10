"use client";
import Link from "next/link";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { resourceId } from "@/shared/lib/api-client";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import { ConfirmDialog, StatusBadge } from "@/shared/ui/product-ui";
import { transitionPosting, deletePosting } from "../api/job.api";
import type { Posting, PostingList } from "../types/job.types";
import styles from "./postings.module.css";

type Action = "publish" | "republish" | "archive" | "delete";
const labels: Record<Action, string> = {
  publish: "Publish",
  republish: "Republish",
  archive: "Archive",
  delete: "Delete",
};
/** Unified Figma posting feed with status-specific lifecycle confirmations. */
export function Postings() {
  const jobs = useApiResource<PostingList>("/v1/jobs");
  const [title, setTitle] = useState("");
  const [status, setStatus] = useState("");
  const [confirmation, setConfirmation] = useState<{
    job: Posting;
    action: Action;
  } | null>(null);
  const items = jobs.data?.items ?? [];
  const visible = items.filter(
    (job) =>
      job.title.toLowerCase().includes(title.toLowerCase()) &&
      (!status || job.status === status),
  );
  if (jobs.error)
    return (
      <section className={styles.errorLayout} data-figma-node="109:14385">
        <header>
          <h1>We couldn’t load your postings.</h1>
          <p>Your saved postings have not changed.</p>
        </header>
        <div>
          <h2>Please try again</h2>
          <p>The latest posting list is temporarily unavailable.</p>
          <ApiFeedback error={jobs.error} retry={jobs.reload} />
        </div>
      </section>
    );
  return (
    <section data-figma-node="109:13232">
      <header className={styles.heading}>
        <div>
          <h1>Job postings</h1>
          <p>Create, publish, and manage every role from one feed.</p>
        </div>
        <Link className={styles.primary} href="/employer/jobs/new">
          Create posting
        </Link>
      </header>
      <section className={styles.panel}>
        <h2>All postings</h2>
        <p className={styles.muted}>
          Search every posting and use status-specific actions to manage it.
        </p>
        <div className={styles.filters}>
          <label>
            <span className={styles.srOnly}>Search postings</span>
            <input
              type="search"
              placeholder="Search by job title…"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
            />
          </label>
          <div className={styles.tabs} aria-label="Posting status">
            {["", "draft", "published", "archived"].map((value) => (
              <button
                type="button"
                key={value}
                aria-pressed={status === value}
                onClick={() => setStatus(value)}
              >
                {value ? value[0].toUpperCase() + value.slice(1) : "All"} (
                {jobs.data
                  ? items.filter((job) => !value || job.status === value).length
                  : "—"}
                )
              </button>
            ))}
          </div>
        </div>
        <BusyState loading={jobs.loading} />
        <ApiFeedback error={jobs.error} retry={jobs.reload} />
        {jobs.data && (
          <>
            <div className={styles.grid}>
              {visible.map((job) => (
                <article className={styles.card} key={job.id}>
                  <header>
                    <div>
                      <h3>{job.title}</h3>
                      <p>
                        {job.employment_type} · {job.seniority_level}
                      </p>
                    </div>
                    <StatusBadge
                      tone={
                        job.status === "published"
                          ? "success"
                          : job.status === "draft"
                            ? "warning"
                            : "neutral"
                      }
                    >
                      {job.status[0].toUpperCase() + job.status.slice(1)}
                    </StatusBadge>
                  </header>
                  <div className={styles.actions}>
                    <Link
                      href={`/employer/jobs/${resourceId(job.id)}${job.status === "draft" ? "" : "/preview"}`}
                    >
                      {job.status === "draft" ? "Edit" : "Details"}
                    </Link>
                    {job.status === "published" && (
                      <Link
                        href={`/employer/jobs/${resourceId(job.id)}/applications`}
                      >
                        Applicants
                      </Link>
                    )}
                    <button
                      type="button"
                      className={
                        job.status === "published" ? undefined : styles.primary
                      }
                      onClick={() =>
                        setConfirmation({
                          job,
                          action:
                            job.status === "draft"
                              ? "publish"
                              : job.status === "published"
                                ? "archive"
                                : "republish",
                        })
                      }
                    >
                      {job.status === "draft"
                        ? "Publish"
                        : job.status === "published"
                          ? "Archive"
                          : "Republish"}
                    </button>
                    <button
                      type="button"
                      className={styles.danger}
                      onClick={() => setConfirmation({ job, action: "delete" })}
                    >
                      Delete
                    </button>
                  </div>
                </article>
              ))}
            </div>
            {!visible.length && (
              <div className={styles.empty}>
                <h3>
                  {items.length
                    ? "No matching postings"
                    : "No job postings yet"}
                </h3>
                <p>
                  {items.length
                    ? "Try another title or status filter."
                    : "Create your first posting to start hiring with amara."}
                </p>
                {!items.length && (
                  <Link className={styles.primary} href="/employer/jobs/new">
                    Create posting
                  </Link>
                )}
              </div>
            )}
            <p className={styles.muted}>
              Showing {visible.length}{" "}
              {visible.length === 1 ? "posting" : "postings"}
            </p>
          </>
        )}
      </section>
      {confirmation && (
        <ConfirmDialog
          title={`${labels[confirmation.action]} this job?`}
          description={`${labels[confirmation.action]} “${confirmation.job.title}”? ${confirmation.action === "delete" ? "This permanently removes the posting and cannot be undone." : "The server checks your permissions and the posting’s current state before saving."}`}
          confirmLabel={`${labels[confirmation.action]} posting`}
          intent={confirmation.action === "delete" ? "danger" : "warning"}
          onClose={() => setConfirmation(null)}
          onConfirm={async () => {
            if (confirmation.action === "delete")
              await deletePosting(confirmation.job.id);
            else
              await transitionPosting(confirmation.job.id, confirmation.action);
            jobs.reload();
          }}
        />
      )}
    </section>
  );
}
