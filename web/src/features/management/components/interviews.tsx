"use client";
import Link from "next/link";
import { useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { InterviewPage } from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

export function Interviews({ interviewer = false }: { interviewer?: boolean }) {
  const [dialog, setDialog] = useState<string | null>(null);
  const [selected, setSelected] = useState("");
  const [offset, setOffset] = useState(0);
  const [outcome, setOutcome] = useState<
    "held" | "candidate_absent" | "interviewer_absent"
  >("held");
  const resource = useApiResource<InterviewPage>(
    `/v1/interviews?limit=20&offset=${offset}`,
  );
  const active = resource.data?.data.find((x) => x.id === selected);
  return (
    <>
      <ProductHeading
        title={
          interviewer
            ? "Your interviews"
            : "Live interview · Scheduling & outcomes"
        }
        description={
          interviewer
            ? "Assigned interviews and scorecards"
            : "Manage scheduling, meeting confirmations, and recorded outcomes."
        }
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Interview discovery and scheduling are not available yet. An authorized
        participant can record an outcome when an interview is available.
      </UnavailableNotice>
      <BusyState loading={resource.loading} />
      {resource.error?.status !== 404 && (
        <ApiFeedback error={resource.error} retry={resource.reload} />
      )}
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">Interview</th>
              <th scope="col">Time</th>
              <th scope="col">Outcome / meeting</th>
              <th scope="col">Action</th>
            </tr>
          </thead>
          <tbody>
            {resource.data?.data.map((row) => (
              <tr key={row.id}>
                <td>
                  <button
                    aria-pressed={row.id === selected}
                    onClick={() => setSelected(row.id)}
                  >
                    {row.id}
                  </button>
                </td>
                <td>
                  {new Date(row.scheduledStartAt).toLocaleString()} –{" "}
                  {new Date(row.scheduledEndAt).toLocaleTimeString()}
                </td>
                <td>{row.outcome || row.status}</td>
                <td>
                  {interviewer ? (
                    <Link
                      href={`/interviewer/interviews/${resourceId(row.id)}/scorecard`}
                    >
                      Open scorecard
                    </Link>
                  ) : (
                    <button
                      onClick={() => {
                        setSelected(row.id);
                        setDialog("outcome");
                      }}
                    >
                      Review outcome
                    </button>
                  )}
                  {row.status === "confirmed" &&
                    row.joinUrl &&
                    /^https?:\/\//.test(row.joinUrl) && (
                      <a
                        href={row.joinUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        Join meeting
                      </a>
                    )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!resource.loading && !resource.data?.data.length && (
        <EmptyState
          title="No interviews available"
          description="Assigned interviews will appear here when the service becomes available."
        />
      )}
      <section className={e.panel}>
        <h2>Selected interview</h2>
        <p>
          {active
            ? `${active.status} · ${new Date(active.scheduledStartAt).toLocaleString()}`
            : "Select an interview to review its scheduling details."}
        </p>
        <p>
          Absence holds the stage and does not record a failing score. A meeting
          is joinable only after provider confirmation.
        </p>
        <label className={e.field}>
          Interview outcome
          <select
            value={outcome}
            disabled={!active}
            onChange={(event) =>
              setOutcome(event.target.value as typeof outcome)
            }
          >
            <option value="held">Held</option>
            <option value="candidate_absent">Candidate absent</option>
            <option value="interviewer_absent">Interviewer absent</option>
          </select>
        </label>
        <div className={e.actions}>
          <Link
            href={`/interviewer/availability${selected ? "?interview=" + resourceId(selected) : ""}`}
          >
            View availability
          </Link>
          {!interviewer && (
            <>
              <button disabled={!active} onClick={() => setDialog("waiver")}>
                Record scorecard waiver
              </button>
              <button disabled={!active} onClick={() => setDialog("outcome")}>
                Review outcome
              </button>
            </>
          )}
        </div>
      </section>
      {resource.data && (
        <div className={e.actions}>
          <button
            disabled={!offset}
            onClick={() => setOffset((n) => Math.max(0, n - 20))}
          >
            Previous page
          </button>
          <button
            disabled={offset + 20 >= resource.data!.pagination.total}
            onClick={() => setOffset((n) => n + 20)}
          >
            Next page
          </button>
        </div>
      )}
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "waiver"
              ? "Waive missing scorecard?"
              : "Record interview outcome?"
          }
          description={
            dialog === "waiver"
              ? "A waiver cannot create a score when no scorecards were submitted."
              : `Record this interview as ${outcome.replaceAll("_", " ")}. Absence holds the stage without recording a failing score.`
          }
          confirmLabel={
            dialog === "waiver" ? "Record waiver" : "Record outcome"
          }
          reason={dialog === "waiver"}
          unavailable={
            dialog === "waiver"
              ? "Scorecard waivers are not available yet."
              : !active
                ? "Select an authorized interview first."
                : undefined
          }
          onConfirm={async () => {
            if (!active)
              throw new Error("Select an authorized interview first.");
            await apiRequest(
              `/v1/interviews/${resourceId(active.id)}/outcome`,
              { method: "POST", body: { outcome } },
            );
            resource.reload();
          }}
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
