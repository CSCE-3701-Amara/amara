"use client";
import { useState, type FormEvent } from "react";
import { queryString, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import type { ApplicationPage, ApplicationDetail } from "../types/cand.types";
function Details({ id }: { id: string }) {
  const detail = useApiResource<ApplicationDetail>(
    `/v1/applications/${resourceId(id)}`,
  );
  return (
    <section aria-label="Application details">
      <BusyState loading={detail.loading} />
      <ApiFeedback error={detail.error} retry={detail.reload} />
      {detail.data && (
        <>
          <p>Status: {detail.data.status.replaceAll("_", " ")}</p>
          <p>
            {detail.data.closed
              ? "This application is closed."
              : "This application is open."}
          </p>
          {detail.data.finalFeedback && <p>{detail.data.finalFeedback}</p>}
        </>
      )}
    </section>
  );
}
export function Applications() {
  const [offset, setOffset] = useState(0);
  const [status, setStatus] = useState("");
  const [selected, setSelected] = useState<string>();
  const page = useApiResource<ApplicationPage>(
    `/v1/candidates/me/applications${queryString({ offset, limit: 20, status })}`,
  );
  function filter(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStatus(
      String(new FormData(event.currentTarget).get("status") || "").trim(),
    );
    setOffset(0);
    setSelected(undefined);
  }
  return (
    <>
      <h1>My applications</h1>
      <p className={s.muted}>
        Follow your progress and read feedback when it’s available.
      </p>
      <form onSubmit={filter} className={s.filters}>
        <label className={s.field}>
          Application status
          <input name="status" placeholder="All statuses" />
        </label>
        <button type="submit">Filter applications</button>
      </form>
      <BusyState loading={page.loading} />
      <ApiFeedback error={page.error} retry={page.reload} />
      {page.data && (
        <>
          {!page.data.data.length && (
            <p role="status">No applications found.</p>
          )}
          {page.data.data.map((application) => (
            <article key={application.id} className={s.card}>
              <h2>{application.jobTitle}</h2>
              <p>{application.companyName}</p>
              <span className={s.badge}>
                {application.status.replaceAll("_", " ")}
              </span>
              <p>Reference: {application.referenceNumber}</p>
              <button
                type="button"
                aria-expanded={selected === application.id}
                onClick={() =>
                  setSelected(
                    selected === application.id ? undefined : application.id,
                  )
                }
              >
                {selected === application.id ? "Hide details" : "View details"}
              </button>
              {selected === application.id && <Details id={application.id} />}
            </article>
          ))}
          <div className={s.actions}>
            <button
              disabled={offset === 0}
              onClick={() => {
                setOffset(Math.max(0, offset - 20));
                setSelected(undefined);
              }}
            >
              Previous page
            </button>
            <p>{page.data.pagination.total} applications</p>
            <button
              disabled={offset + 20 >= page.data.pagination.total}
              onClick={() => {
                setOffset(offset + 20);
                setSelected(undefined);
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
