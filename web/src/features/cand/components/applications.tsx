"use client";
import Link from "next/link";
import { Suspense, useState } from "react";
import { useSearchParams } from "next/navigation";
import { queryString, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import { EmptyState, Stats, StatusBadge } from "@/shared/ui/product-ui";
import type { ApplicationPage } from "../types/cand.types";
import styles from "./candidate.module.css";
function ApplicationList() {
  const params = useSearchParams();
  const [dismissed, setDismissed] = useState(false);
  const presentation = dismissed ? null : params.get("state");
  const [offset, setOffset] = useState(0);
  const [group, setGroup] = useState("all");
  const page = useApiResource<ApplicationPage>(
    `/v1/candidates/me/applications${queryString({ offset, limit: 20, status: group === "progress" ? "in_progress" : undefined })}`,
  );
  const inProgress = useApiResource<ApplicationPage>(
    "/v1/candidates/me/applications?limit=1&status=in_progress",
  );
  const allCount = useApiResource<ApplicationPage>(
    "/v1/candidates/me/applications?limit=1",
  );
  const applications =
    presentation === "empty"
      ? []
      : (page.data?.data ?? []).filter(
          (item) => group !== "completed" || item.closed,
        );
  const error =
    presentation === "error"
      ? new Error("We couldn’t load your applications")
      : page.error;
  if (error)
    return (
      <section data-figma-node="89:371">
        <EmptyState
          level={1}
          title="We couldn’t load your applications"
          description="Try again to retrieve your latest progress. If your session has expired, log in again."
        >
          <ApiFeedback
            error={error}
            retry={() => {
              setDismissed(true);
              page.reload();
            }}
          />
        </EmptyState>
      </section>
    );
  if (!page.loading && applications.length === 0 && group === "all")
    return (
      <section data-figma-node="89:298">
        <EmptyState
          level={1}
          title="No applications yet"
          description="A new chapter starts with an application. Return here after applying to track assessments, interviews, and offers."
        >
          <Link
            className={`${s.button} ${s.primary}`}
            href="/candidate/dashboard"
          >
            Explore opportunities
          </Link>
        </EmptyState>
      </section>
    );
  return (
    <section data-figma-node="34:7">
      <h1>Dashboard</h1>
      <Stats
        items={[
          {
            label: "Applications",
            value:
              allCount.data?.pagination.total ??
              page.data?.pagination.total ??
              "—",
          },
          {
            label: "In progress",
            value: inProgress.data?.pagination.total ?? "—",
          },
          { label: "Offers to review", value: "—" },
        ]}
      />
      <section className={styles.applications} aria-label="Your applications">
        <h2>Your applications</h2>
        <div className={styles.tabs} aria-label="Application filters">
          {(
            [
              ["all", "All statuses"],
              ["progress", "In progress"],
              ["completed", "Completed on this page"],
            ] as const
          ).map(([value, label]) => (
            <button
              key={value}
              type="button"
              aria-pressed={group === value}
              onClick={() => {
                setGroup(value);
                setOffset(0);
              }}
            >
              {label}
            </button>
          ))}
        </div>
        <BusyState loading={page.loading} />
        {applications.map((application) => (
          <article className={styles.applicationRow} key={application.id}>
            <div>
              <h3>{application.jobTitle}</h3>
              <p>{application.companyName}</p>
            </div>
            <div className={styles.applicationState}>
              <StatusBadge>
                {application.status.replaceAll("_", " ")}
              </StatusBadge>
              <small>
                {application.updatedAt
                  ? `Updated ${new Date(application.updatedAt).toLocaleDateString()}`
                  : `Reference ${application.referenceNumber}`}
              </small>
            </div>
            <Link
              className={`${s.button} ${s.primary}`}
              href={`/candidate/applications/${resourceId(application.id)}`}
            >
              {application.closed ? "View status" : "Continue"}
            </Link>
          </article>
        ))}
        {!page.loading && !applications.length && (
          <p>No applications match this filter.</p>
        )}
        {page.data && (
          <p className={styles.listSummary}>
            {page.data.pagination.total} applications · Page{" "}
            {Math.floor(offset / 20) + 1} of{" "}
            {Math.max(1, Math.ceil(page.data.pagination.total / 20))}
          </p>
        )}
        {page.data && page.data.pagination.total > 20 && (
          <div className={s.actions}>
            <button
              type="button"
              disabled={offset === 0}
              onClick={() => setOffset(Math.max(0, offset - 20))}
            >
              Previous page
            </button>
            <button
              type="button"
              disabled={offset + 20 >= page.data.pagination.total}
              onClick={() => setOffset(offset + 20)}
            >
              Next page
            </button>
          </div>
        )}
      </section>
    </section>
  );
}
/** URL state variants share the same authenticated history route. */
export function Applications() {
  return (
    <Suspense fallback={<p role="status">Loading applications…</p>}>
      <ApplicationList />
    </Suspense>
  );
}
