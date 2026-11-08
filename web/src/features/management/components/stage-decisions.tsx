"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { RuntimePipelineDashboard } from "@/features/invw/types/runtime.types";
import { ThresholdReview } from "@/features/invw/components/threshold-review";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

export function StageDecisions({ id }: { id: string }) {
  const params = useSearchParams();
  const [dialog, setDialog] = useState<string | null>(null);
  const [stage, setStage] = useState(params.get("stage") || "");
  const [search, setSearch] = useState("");
  const [selected, setSelected] = useState<string[]>([]);
  const executions = useApiResource<RuntimePipelineDashboard>(
    `/v1/jobs/${resourceId(id)}/pipeline/dashboard`,
  );
  const stageId = stage || executions.data?.stages[0]?.id || "";
  const activeStage = executions.data?.stages.find((row) => row.id === stageId);
  const rows =
    activeStage?.candidates.filter((row) =>
      row.candidate.toLowerCase().includes(search.toLowerCase()),
    ) || [];
  const [decision, setDecision] = useState<"pass" | "reject">("pass");
  return (
    <>
      <ProductHeading
        title="Candidates · Stage decisions"
        description="Decision queue · Hiring manager only"
        compact
        className={e.titleBanner}
      />
      <nav className={e.actions}>
        <Link href={`/employer/pipelines/${resourceId(id)}`}>Pipeline</Link>
        <Link href="/employer/interviews">Live interviews</Link>
        <Link href="/employer/assessments">Assessment instances</Link>
      </nav>
      <UnavailableNotice>
        Manual pass and reject decisions are available for active candidates.
        Threshold changes require a separate preview of candidate outcomes.
        Decision overrides are not supported yet.
      </UnavailableNotice>
      {params.get("state") === "conflict" && (
        <p className={e.error} role="alert">
          The candidate state changed. Reload the current record before
          reviewing a new decision.
        </p>
      )}
      {params.get("state") === "override" && (
        <section className={e.panel}>
          <h2>Override decision</h2>
          <p>
            An override requires an authorized hiring manager and a recorded
            nonblank reason. The previous decision remains in history.
          </p>
        </section>
      )}
      <div className={e.grid}>
        <label className={e.field}>
          Stage
          <select
            aria-label="Stage"
            value={stageId}
            onChange={(event) => {
              setStage(event.target.value);
              setSelected([]);
            }}
          >
            <option value="">Choose a stage</option>
            {executions.data?.stages.map((row) => (
              <option key={row.id} value={row.id}>
                {row.name}
              </option>
            ))}
          </select>
        </label>
        <label className={e.field}>
          Search candidate
          <input
            type="search"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
      </div>
      <BusyState loading={executions.loading} />
      {executions.error?.status !== 404 && (
        <ApiFeedback error={executions.error} retry={executions.reload} />
      )}
      {activeStage?.rule === "threshold" && activeStage.threshold !== null && (
        <ThresholdReview
          key={`${activeStage.id}-${activeStage.threshold}`}
          stageId={activeStage.id}
          current={activeStage.threshold}
          onSaved={executions.reload}
        />
      )}
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">Candidate / selection</th>
              <th scope="col">Assessment</th>
              <th scope="col">Stage decision</th>
              <th scope="col">Recorded by / time</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row) => (
              <tr key={row.id}>
                <td>
                  <label className={e.check}>
                    <input
                      type="checkbox"
                      disabled={!row.active || activeStage?.rule !== "manual"}
                      checked={selected.includes(row.id)}
                      onChange={(event) =>
                        setSelected((ids) =>
                          event.target.checked
                            ? [...ids, row.id]
                            : ids.filter((x) => x !== row.id),
                        )
                      }
                    />
                    {row.candidate}
                  </label>
                </td>
                <td>
                  {row.score === null || row.score === undefined
                    ? "Pending"
                    : `${row.score} / 100`}
                </td>
                <td>{row.status}</td>
                <td>Decision history unavailable</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!rows.length && !executions.loading && (
        <EmptyState
          title="No candidate decisions available"
          description="Choose an authorized stage to inspect its recorded executions."
        />
      )}
      <div className={e.actions}>
        <StatusBadge>{selected.length} selected · Same stage</StatusBadge>
        <button
          disabled={!selected.length || activeStage?.rule !== "manual"}
          onClick={() => {
            setDecision("pass");
            setDialog("pass");
          }}
        >
          Pass
        </button>
        <button
          disabled={!selected.length || activeStage?.rule !== "manual"}
          onClick={() => {
            setDecision("reject");
            setDialog("reject");
          }}
        >
          Reject
        </button>

        <button onClick={() => setSelected([])}>Clear selection</button>
      </div>
      <section className={e.panel}>
        <h2>Finalist is not hired</h2>
        <p>
          Completing all required paths marks a candidate as a finalist. Hiring
          is a separate decision outside this pipeline.
        </p>
      </section>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "reject"
              ? "Reject this candidate?"
              : dialog === "pass"
                ? "Pass this candidate?"
                : "Apply this bulk decision?"
          }
          description="Review the selected candidates, their current stage, and the recorded reason. Candidates in different stages cannot be decided together."
          confirmLabel="Confirm decision"
          intent={dialog === "reject" ? "danger" : "success"}
          reason
          unavailable={
            !activeStage || activeStage.rule !== "manual"
              ? "Select candidates from a manual-decision stage."
              : undefined
          }
          onConfirm={async (reason) => {
            const chosen = rows.filter((row) => selected.includes(row.id));
            if (!chosen.length || chosen.some((row) => !row.active))
              throw new Error(
                "Select active candidates and refresh changed records before deciding.",
              );
            if (chosen.length === 1)
              await apiRequest(
                `/v1/stage-executions/${resourceId(chosen[0].id)}/decisions`,
                {
                  method: "POST",
                  body: {
                    decision,
                    decisionType: "manual",
                    reason,
                    expectedStatus: chosen[0].status,
                  },
                },
              );
            else
              await apiRequest(
                `/v1/pipeline-stages/${resourceId(stageId)}/decisions`,
                {
                  method: "POST",
                  body: {
                    decision,
                    decisionType: "manual",
                    reason,
                    stageExecutionIds: chosen.map((row) => row.id),
                  },
                },
              );
            setSelected([]);
            executions.reload();
          }}
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
