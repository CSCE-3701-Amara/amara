"use client";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { GradingSubmissionPage } from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

const art = { src: "/figma/hiring-dag.svg", width: 200.16, height: 138.96 };
export function AssignedGrading() {
  const [dialog, setDialog] = useState<string | null>(null);
  const [selected, setSelected] = useState("");
  const [offset, setOffset] = useState(0);
  const resource = useApiResource<GradingSubmissionPage>(
    `/v1/grading/submissions?limit=20&offset=${offset}`,
  );
  const active = resource.data?.data.find((x) => x.id === selected);
  return (
    <>
      <ProductHeading
        title="Interviewer workspace"
        description="Assigned grading"
        compact
        artwork={art}
      />
      <UnavailableNotice>
        Grades cannot be saved yet. Only assigned, authorized submissions belong
        in this workspace.
      </UnavailableNotice>
      <BusyState loading={resource.loading} />
      {resource.error?.status !== 404 && (
        <ApiFeedback error={resource.error} retry={resource.reload} />
      )}
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">Submission</th>
              <th scope="col">Assessment</th>
              <th scope="col">Status</th>
            </tr>
          </thead>
          <tbody>
            {resource.data?.data.map((row) => (
              <tr key={row.id}>
                <td>
                  <button
                    aria-pressed={selected === row.id}
                    onClick={() => setSelected(row.id)}
                  >
                    Candidate {row.anonymousIdentifier}
                  </button>
                </td>
                <td>{row.prompt}</td>
                <td>{row.gradingStatus.replaceAll("_", " ")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!resource.loading && !resource.data?.data.length && (
        <EmptyState
          title="No submissions assigned"
          description="Your assigned answers will appear here when grading becomes available."
        />
      )}
      <section className={e.section}>
        <h2>Selected answer</h2>
        <p>Question: {active?.prompt || "Select a submission"}</p>
        <div className={e.panel}>
          {typeof active?.response === "string" ? (
            active.response
          ) : active?.response && "code" in active.response ? (
            <pre>{active.response.code}</pre>
          ) : Array.isArray(active?.response) ? (
            active.response.join(", ")
          ) : (
            "No written response available"
          )}
        </div>
        <label className={e.field}>
          Earned points{active ? ` (out of ${active.maxPoints})` : ""}
          <input
            key={active?.id || "empty"}
            type="number"
            min={0}
            max={active?.maxPoints}
            step="any"
            defaultValue={active?.earnedPoints ?? ""}
            disabled={!active}
          />
        </label>
        <div className={e.actions}>
          <button
            className={e.primary}
            disabled={!active}
            onClick={() => setDialog("submit")}
          >
            Review grade
          </button>
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
          title="Submit grading results?"
          description="Review earned points and the candidate response before submitting. Scores remain private until released through the supported hiring workflow."
          confirmLabel="Submit grades"
          unavailable="Saving grades is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
