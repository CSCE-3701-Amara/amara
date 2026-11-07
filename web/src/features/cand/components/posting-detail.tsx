"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { ApiError, apiRequest, queryString } from "@/shared/lib/api-client";
import { applyToPosting } from "@/features/job/api/job.api";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import { ProductHeading, StatusBadge } from "@/shared/ui/product-ui";
import type { ApplicationSubmission } from "@/features/job/types/job.types";
import type { FeedJob, FeedPage } from "../types/cand.types";
import { employmentLabel, salaryLabel } from "./job-feed";
import styles from "./candidate.module.css";
/** Resolve candidate details through the authenticated feed, never the company-only posting read. */
export function PostingDetail({ id }: { id: string }) {
  const params = useSearchParams();
  const title = params.get("title") || "";
  const [revision, setRevision] = useState(0);
  const key = `${id}:${title}:${revision}`;
  const [result, setResult] = useState<{
    key: string;
    job?: FeedJob;
    error?: Error;
  }>();
  const [application, setApplication] = useState<ApplicationSubmission>();
  const [applyError, setApplyError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const pending = useRef(false);
  useEffect(() => {
    const request = new AbortController();
    async function load() {
      try {
        let cursor: string | undefined;
        const seen = new Set<string>();
        do {
          const page = await apiRequest<FeedPage>(
            `/v1/jobs/feed${queryString({ title, limit: 100, cursor })}`,
            { signal: request.signal },
          );
          const job = page.items.find((item) => item.id === id);
          if (job) {
            if (!request.signal.aborted) setResult({ key, job });
            return;
          }
          cursor = page.next_cursor;
          if (cursor && seen.has(cursor))
            throw new Error("The job feed returned an invalid page cursor.");
          if (cursor) seen.add(cursor);
        } while (cursor && !request.signal.aborted);
        throw new ApiError(
          "This opportunity is no longer available in the job feed.",
          "posting_not_found",
          404,
        );
      } catch (error) {
        if (!request.signal.aborted) setResult({ key, error: error as Error });
      }
    }
    void load();
    return () => request.abort();
  }, [id, title, key]);
  const current = result?.key === key ? result : undefined;
  const job = current?.job;
  async function apply() {
    if (pending.current) return;
    pending.current = true;
    setBusy(true);
    setApplyError(undefined);
    try {
      setApplication(await applyToPosting(id));
    } catch (error) {
      setApplyError(error as Error);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return (
    <section data-figma-node="384:641">
      <Link href="/candidate/dashboard">‹ Back to opportunities</Link>
      {!job && <h1>Opportunity details</h1>}
      <BusyState loading={!current} />
      <ApiFeedback
        error={current?.error}
        retry={() => setRevision((value) => value + 1)}
      />
      {job && (
        <>
          <ProductHeading
            title={job.title}
            description={[job.company_name, job.location]
              .filter(Boolean)
              .join(" · ")}
          />
          <p className={s.muted}>
            {[
              employmentLabel(job.employment_type),
              job.seniority_level,
              job.language,
            ]
              .filter(Boolean)
              .join(" · ")}
          </p>
          <div className={styles.detailGrid}>
            <section className={styles.panel}>
              <h2>About the role</h2>
              <p>
                {job.description ||
                  "The hiring team hasn’t provided a role description."}
              </p>
              <h2>What we are looking for</h2>
              <p>
                Review the role description before applying. Required skills are
                checked when your application is submitted.
              </p>
              <h2>Benefits and hiring process</h2>
              <p>Further details will be shared by the hiring team.</p>
            </section>
            <aside className={styles.statusPanel}>
              <h2>Salary range</h2>
              <p>{salaryLabel(job)}</p>
              <p>
                {job.application_deadline
                  ? `Apply by ${new Date(job.application_deadline).toLocaleString()}`
                  : "Applications are reviewed on a rolling basis."}
              </p>
              <ApiFeedback error={applyError} />
              {application ? (
                <div className={s.success} role="status">
                  <StatusBadge
                    tone={application.screened_out ? "warning" : "success"}
                  >
                    {application.screened_out
                      ? "Screening complete"
                      : "Application submitted"}
                  </StatusBadge>
                  <p>
                    {application.screened_out
                      ? "Your application was submitted but did not pass initial screening."
                      : "Your application has been received."}{" "}
                    Reference: {application.reference_number}.
                  </p>
                  {!!application.failed_amarat?.length && (
                    <p>
                      Required skills: {application.failed_amarat.join(", ")}
                    </p>
                  )}
                  <Link href="/candidate/applications">
                    View my applications
                  </Link>
                </div>
              ) : (
                <button
                  className={s.primary}
                  type="button"
                  disabled={busy}
                  onClick={apply}
                >
                  {busy ? "Submitting…" : "Apply now"}
                </button>
              )}
            </aside>
          </div>
        </>
      )}
    </section>
  );
}
