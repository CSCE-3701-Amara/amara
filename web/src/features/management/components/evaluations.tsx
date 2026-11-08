"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { RuntimePipelineDashboard } from "@/features/invw/types/runtime.types";
import { ThresholdReview } from "@/features/invw/components/threshold-review";
import type {
  GradingSubmissionPage,
  InterviewPage,
  InterviewFeedback,
} from "@/shared/lib/contract.types";
import type {
  AssessmentInstance,
  TimeSlot,
} from "@/features/cand/types/product.types";
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
const art = { src: "/figma/hiring-dag.svg", width: 200.16, height: 138.96 };
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
export function InterviewScorecard({ id }: { id: string }) {
  const [dialog, setDialog] = useState<string | null>(null);
  const scorecard = useApiResource<InterviewFeedback>(
    `/v1/interviews/${resourceId(id)}/feedback/me`,
  );
  const criteria = Object.keys(scorecard.data?.ratings || {});
  return (
    <>
      <ProductHeading
        title="Interview scorecard"
        description="Assigned interviewer · Private hiring workspace"
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Scorecard drafts, notes, and submission aren’t available yet. Private
        notes are never shown to candidates.
      </UnavailableNotice>
      <BusyState loading={scorecard.loading} />
      {scorecard.error?.status !== 404 && (
        <ApiFeedback error={scorecard.error} retry={scorecard.reload} />
      )}
      <div className={e.split}>
        <section>
          <h2>Mapped evaluation rubric</h2>
          <p>
            Rate every criterion · 1 = Limited evidence · 5 = Strong evidence
          </p>
          {!criteria.length && (
            <p>
              The assigned rubric cannot be loaded yet. No criteria or ratings
              have been invented.
            </p>
          )}
          {criteria.map((name) => (
            <fieldset key={name} className={e.panel}>
              <legend>{name}</legend>
              <div className={e.actions}>
                {[1, 2, 3, 4, 5].map((rating) => (
                  <label className={e.check} key={rating}>
                    <input
                      type="radio"
                      name={name}
                      value={rating}
                      defaultChecked={scorecard.data?.ratings[name] === rating}
                    />
                    {rating}
                  </label>
                ))}
              </div>
            </fieldset>
          ))}
          <label className={e.field}>
            Written recommendation (required)
            <textarea
              required
              maxLength={5000}
              defaultValue={scorecard.data?.recommendation || ""}
            />
          </label>
          <div className={e.actions}>
            <button disabled>Save scorecard draft</button>
            <button className={e.primary} onClick={() => setDialog("submit")}>
              Review submission
            </button>
          </div>
        </section>
        <aside className={e.panel}>
          <h2>Private interview notes</h2>
          <StatusBadge>Hiring manager & assigned interviewer</StatusBadge>
          <p>Never shared with the candidate.</p>
          <label className={e.field}>
            New note
            <textarea
              maxLength={5000}
              defaultValue={scorecard.data?.notes || ""}
            />
          </label>
          <button disabled>Save notes</button>
          <h3>Submission</h3>
          <p>
            Submitting locks your scorecard for review. Only an authorized
            hiring manager can pass, reject, or advance a candidate.
          </p>
        </aside>
      </div>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title="Submit interview scorecard?"
          description="Complete every rubric rating and your written recommendation before submitting. Your scorecard will be locked for hiring-team review."
          confirmLabel="Submit scorecard"
          unavailable="Scorecard submission is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
export function InterviewAvailability() {
  const params = useSearchParams();
  const [id, setId] = useState(params.get("interview") || "");
  const interviews = useApiResource<InterviewPage>(
    "/v1/interviews?limit=100&offset=0",
  );
  const slots = useApiResource<TimeSlot[]>(
    id ? `/v1/interviews/${resourceId(id)}/availability` : null,
  );
  const [drafts, setDrafts] = useState<{ startAt: string; endAt: string }[]>([
    { startAt: "", endAt: "" },
  ]);
  return (
    <>
      <ProductHeading
        title="Your available times"
        description="Interview availability"
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Availability updates aren’t available yet. Entered times remain on this
        screen.
      </UnavailableNotice>
      <label className={e.field}>
        Interview
        <select value={id} onChange={(event) => setId(event.target.value)}>
          <option value="">Choose an interview</option>
          {interviews.data?.data.map((row) => (
            <option key={row.id} value={row.id}>
              {new Date(row.scheduledStartAt).toLocaleString()} · {row.id}
            </option>
          ))}
        </select>
      </label>
      <BusyState loading={slots.loading} />
      {slots.error?.status !== 404 && (
        <ApiFeedback error={slots.error} retry={slots.reload} />
      )}
      <section className={e.section}>
        <h2>Your available times</h2>
        {slots.data?.map((slot) => (
          <p key={slot.startAt}>
            {new Date(slot.startAt).toLocaleString()} –{" "}
            {new Date(slot.endAt).toLocaleString()}
          </p>
        ))}
        {drafts.map((draft, index) => (
          <div className={e.grid} key={index}>
            {(["startAt", "endAt"] as const).map((key, i) => (
              <label className={e.field} key={key}>
                {i ? "Ends" : "Starts"}
                <input
                  type="datetime-local"
                  value={draft[key]}
                  onChange={(event) =>
                    setDrafts((rows) =>
                      rows.map((row, j) =>
                        j === index
                          ? { ...row, [key]: event.target.value }
                          : row,
                      ),
                    )
                  }
                />
              </label>
            ))}
          </div>
        ))}
        <div className={e.actions}>
          <button
            onClick={() =>
              setDrafts((rows) => [...rows, { startAt: "", endAt: "" }])
            }
          >
            Add time interval
          </button>
        </div>
      </section>
      <section className={e.section}>
        <h2>Scheduling</h2>
        <p>
          Candidates can choose a slot that matches the assigned interviewers’
          common availability.
        </p>
      </section>
      <div className={e.actions}>
        <Link href="/interviewer/interviews">Cancel</Link>
        <button className={e.primary} disabled>
          Save availability
        </button>
      </div>
      <WorkspaceFooter />
    </>
  );
}
export function AssessmentInstances() {
  const params = useSearchParams();
  const [application, setApplication] = useState(
    params.get("application") || "",
  );
  const [query, setQuery] = useState(application);
  const [selected, setSelected] = useState("");
  const resource = useApiResource<AssessmentInstance[]>(null);
  const active = resource.data?.find((row) => row.id === selected);
  return (
    <>
      <ProductHeading
        title="Assessment instances & extensions"
        description="Individual candidate + stage"
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Hiring-team assessment timing and extensions require backend support.
        The contract’s instance list is candidate scoped, so this workspace does
        not call it with a hiring-manager session.
      </UnavailableNotice>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          setApplication(query.trim());
        }}
        className={e.actions}
      >
        <label className={e.field}>
          Application ID
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </label>
        <button type="submit" disabled>
          Search instances
        </button>
      </form>
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">Candidate / instance</th>
              <th scope="col">Stage</th>
              <th scope="col">Instance state</th>
              <th scope="col">Final score</th>
            </tr>
          </thead>
          <tbody>
            {resource.data?.map((row) => (
              <tr key={row.id}>
                <td>
                  <button onClick={() => setSelected(row.id)}>{row.id}</button>
                </td>
                <td>{row.stageId}</td>
                <td>{row.status}</td>
                <td>Unavailable</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className={e.split}>
        <section className={e.panel}>
          <h2>Selected instance</h2>
          <p>
            {active?.id ||
              "Select an instance when authorized data becomes available."}
          </p>
          <p>
            Active time and completion windows are separate. Provider failures
            and pending grading never become a zero score automatically.
          </p>
        </section>
        <section className={e.panel}>
          <h2>Give additional time</h2>
          <div className={e.grid}>
            <label className={e.field}>
              Additional active time (minutes)
              <input type="number" min={1} />
            </label>
            <label className={e.field}>
              Individual completion deadline
              <input type="datetime-local" />
            </label>
            <label className={`${e.field} ${e.full}`}>
              Required reason
              <textarea required />
            </label>
          </div>
          <div className={e.actions}>
            <button disabled>Save extension</button>
          </div>
        </section>
      </div>
      <WorkspaceFooter />
    </>
  );
}
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
