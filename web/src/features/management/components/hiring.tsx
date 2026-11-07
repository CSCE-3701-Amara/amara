"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type {
  ApplicationScorecard,
  FinalistPage,
  ManagedApplicationPage,
  PipelineDashboard,
  FeedbackView,
  InterviewerPage,
} from "@/shared/lib/contract.types";
import type { PipelineDraft } from "@/features/invw/types/invw.types";
import type { Offer, OfferPage } from "@/features/cand/types/product.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  Stats,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
const art = { src: "/figma/hiring-dag.svg", width: 200.16, height: 138.96 };
function PlannedRead({
  loading,
  error,
  retry,
}: {
  loading: boolean;
  error?: Error & { status?: number };
  retry: () => void;
}) {
  return (
    <>
      <BusyState loading={loading} />
      {error?.status !== 404 && <ApiFeedback error={error} retry={retry} />}
    </>
  );
}
function identity(row: ApplicationScorecard) {
  return row.identity?.name || `Candidate ${row.anonymousIdentifier}`;
}
export function HiringApplications({
  id,
  view = "applications",
}: {
  id: string;
  view?:
    "applications" | "finalists" | "comparison" | "scorecard" | "dashboard";
}) {
  const params = useSearchParams();
  const [selected, setSelected] = useState<string[]>([]);
  const [dialog, setDialog] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const pipeline = useApiResource<PipelineDraft>(
    view === "comparison" || view === "dashboard"
      ? `/v1/jobs/${resourceId(id)}/pipeline`
      : null,
  );
  const list = useApiResource<ManagedApplicationPage>(
    view === "applications"
      ? `/v1/jobs/${resourceId(id)}/applications?limit=20&offset=${offset}`
      : null,
  );
  const finalists = useApiResource<FinalistPage>(
    view === "finalists"
      ? `/v1/jobs/${resourceId(id)}/finalists?limit=20&offset=${offset}`
      : null,
  );
  const scorecard = useApiResource<ApplicationScorecard>(
    view === "scorecard"
      ? `/v1/applications/${resourceId(id)}/scorecard`
      : null,
  );
  const comparisons = useApiResource<ApplicationScorecard[]>(
    view === "comparison" && pipeline.data
      ? `/v1/pipelines/${resourceId(pipeline.data.id)}/comparison${selected.length ? "?" + selected.map((value) => `applicationIds=${resourceId(value)}`).join("&") : ""}`
      : null,
  );
  const dashboard = useApiResource<PipelineDashboard>(
    view === "dashboard" && pipeline.data
      ? `/v1/pipelines/${resourceId(pipeline.data.id)}/dashboard`
      : null,
  );
  const ranking = params.get("state") === "ranking";
  const cards =
    view === "finalists"
      ? finalists.data?.data
      : view === "comparison"
        ? comparisons.data
        : undefined;
  const title =
    view === "applications"
      ? "Applications"
      : view === "finalists"
        ? ranking
          ? "Finalist ranking"
          : "Finalists"
        : view === "comparison"
          ? "Compare applications"
          : view === "scorecard"
            ? "Application scorecard"
            : "Pipeline dashboard";
  const resources = [
    pipeline,
    list,
    finalists,
    scorecard,
    comparisons,
    dashboard,
  ];
  const error = resources.find((r) => r.error)?.error;
  const loading = resources.some((r) => r.loading);
  const page = list.data?.pagination || finalists.data?.pagination;
  const stageNames = [
    ...new Set(
      cards?.flatMap((row) => row.stages.map((stage) => stage.stageName)) || [],
    ),
  ];
  return (
    <>
      <ProductHeading
        title={title}
        description={
          view === "comparison"
            ? "Compare up to four shortlisted candidates"
            : view === "scorecard"
              ? scorecard.data
                ? identity(scorecard.data)
                : "Hiring manager workspace"
              : undefined
        }
        compact
        artwork={art}
      />
      <UnavailableNotice>
        This hiring view is not available yet. Recorded data can’t be changed.
        Identity stays hidden until the server authorizes its release.
      </UnavailableNotice>
      <PlannedRead
        loading={loading}
        error={error}
        retry={() => resources.forEach((r) => r.reload())}
      />
      <div className={e.tableWrap}>
        {view === "applications" && (
          <table className={e.table}>
            <thead>
              <tr>
                <th scope="col">Application</th>
                <th scope="col">Screening</th>
                <th scope="col">Pipeline</th>
              </tr>
            </thead>
            <tbody>
              {list.data?.data.map((row) => (
                <tr key={row.application.id}>
                  <td>
                    <Link
                      href={`/employer/applications/${resourceId(row.application.id)}/scorecard`}
                    >
                      {row.identity?.name ||
                        `Candidate ${row.anonymousIdentifier || row.application.referenceNumber}`}
                    </Link>
                  </td>
                  <td>{row.screeningSummary || row.application.status}</td>
                  <td>
                    <Link
                      href={`/employer/applications/${resourceId(row.application.id)}/scorecard`}
                    >
                      View scorecard
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {view === "finalists" && (
          <table className={e.table}>
            <thead>
              <tr>
                {ranking && <th scope="col">Display order</th>}
                <th scope="col">Finalist</th>
                <th scope="col">Mean score</th>
                <th scope="col">Selection</th>
                <th scope="col">Offer</th>
              </tr>
            </thead>
            <tbody>
              {cards?.map((row, index) => (
                <tr key={row.applicationId}>
                  {ranking && <td>{offset + index + 1}</td>}
                  <td>{identity(row)}</td>
                  <td>
                    {row.meanScore === null || row.meanScore === undefined
                      ? "Pending"
                      : `${row.meanScore} / 100`}
                  </td>
                  <td>
                    <label className={e.check}>
                      <input
                        type="checkbox"
                        aria-label={`Select ${identity(row)}`}
                        checked={selected.includes(row.applicationId)}
                        onChange={(event) =>
                          setSelected((ids) =>
                            event.target.checked
                              ? [...ids, row.applicationId]
                              : ids.filter((x) => x !== row.applicationId),
                          )
                        }
                      />
                      Select
                    </label>
                  </td>
                  <td>
                    <Link
                      href={`/employer/offers/new?application=${resourceId(row.applicationId)}`}
                    >
                      Create offer
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {view === "comparison" && (
          <table className={e.table}>
            <thead>
              <tr>
                <th scope="col">Stage / measure</th>
                {cards?.slice(0, 4).map((row) => (
                  <th scope="col" key={row.applicationId}>
                    {identity(row)}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {stageNames.map((name) => (
                <tr key={name}>
                  <th scope="row">{name}</th>
                  {cards?.slice(0, 4).map((row) => {
                    const stage = row.stages.find(
                      (stage) => stage.stageName === name,
                    );
                    return (
                      <td key={row.applicationId}>
                        {stage?.score === null || stage?.score === undefined
                          ? stage?.status || "Not available"
                          : `${stage.score} / 100`}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {view === "scorecard" && (
          <table className={e.table}>
            <thead>
              <tr>
                <th scope="col">Stage</th>
                <th scope="col">Outcome</th>
                <th scope="col">Score</th>
              </tr>
            </thead>
            <tbody>
              {scorecard.data?.stages.map((stage) => (
                <tr key={stage.stageId}>
                  <td>{stage.stageName}</td>
                  <td>{stage.status.replaceAll("_", " ")}</td>
                  <td>{stage.score ?? "—"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {view === "dashboard" && (
          <table className={e.table}>
            <thead>
              <tr>
                <th scope="col">Stage</th>
                <th scope="col">Active</th>
                <th scope="col">Passed</th>
                <th scope="col">Rejected</th>
              </tr>
            </thead>
            <tbody>
              {dashboard.data?.stages.map((stage) => (
                <tr key={stage.stageId}>
                  <td>{stage.name}</td>
                  <td>{stage.activeCount}</td>
                  <td>{stage.passedCount}</td>
                  <td>{stage.rejectedCount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      {!loading &&
        !list.data?.data.length &&
        !cards?.length &&
        !scorecard.data &&
        !dashboard.data && (
          <EmptyState
            title="No records available"
            description="This view will display authorized records when the service becomes available."
          />
        )}
      {view === "scorecard" && (
        <section className={e.section}>
          <h2>Application record</h2>
          <p>Identity is hidden until the configured reveal point.</p>
          {scorecard.data?.insightsAreSystemGenerated && (
            <p>These insights are system generated and support human review.</p>
          )}
          <h3>Strengths</h3>
          <ul>
            {scorecard.data?.strengths?.map((text) => (
              <li key={text}>{text}</li>
            ))}
          </ul>
          <h3>Areas for development</h3>
          <ul>
            {scorecard.data?.areasForDevelopment?.map((text) => (
              <li key={text}>{text}</li>
            ))}
          </ul>
        </section>
      )}
      {view === "dashboard" && (
        <>
          <h2>Offers</h2>
          <Stats
            items={[
              {
                label: "Active candidates",
                value: dashboard.data?.activeCandidates ?? "—",
              },
              {
                label: "Pending offers",
                value: dashboard.data?.pendingOffers ?? "—",
              },
              {
                label: "Accepted offers",
                value: dashboard.data?.acceptedOffers ?? "—",
              },
            ]}
          />
        </>
      )}
      {page && (
        <div className={e.actions}>
          <button
            disabled={!offset}
            onClick={() => setOffset((n) => Math.max(0, n - 20))}
          >
            Previous page
          </button>
          <button
            disabled={offset + 20 >= page.total}
            onClick={() => setOffset((n) => n + 20)}
          >
            Next page
          </button>
        </div>
      )}
      <div className={e.actions}>
        <Link href="/employer/jobs">Back to postings</Link>
        {view !== "scorecard" && (
          <>
            <Link href={`/employer/jobs/${resourceId(id)}/applications`}>
              View applications
            </Link>
            <Link href={`/employer/jobs/${resourceId(id)}/comparison`}>
              Compare applications
            </Link>
            <Link
              href={`/employer/jobs/${resourceId(id)}/finalists?state=ranking`}
            >
              Rank finalists
            </Link>
          </>
        )}
        {view === "finalists" && (
          <button
            disabled={!selected.length}
            onClick={() => setDialog("reject")}
          >
            Review rejection
          </button>
        )}
        {ranking && <button disabled>Save ranking</button>}
      </div>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "bulk-reject"
              ? "Reject selected finalists?"
              : "Reject this finalist?"
          }
          description="Rejection closes the selected applications. Review each scorecard before recording a hiring decision."
          confirmLabel="Reject finalist"
          intent="danger"
          reason
          unavailable="Finalist rejection is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
export function HiringOffer({ id }: { id?: string }) {
  const params = useSearchParams();
  const [dialog, setDialog] = useState<string | null>(null);
  const detail = useApiResource<Offer>(
    id ? `/v1/offers/${resourceId(id)}` : null,
  );
  const list = useApiResource<OfferPage>(
    !id && !params.has("application") ? "/v1/offers?limit=20&offset=0" : null,
  );
  const create = !id;
  const capacity = params.get("state") === "capacity";
  return (
    <>
      <ProductHeading
        title={create ? "Create finalist offer" : "Manage offer"}
        description={detail.data?.jobTitle}
        compact
        artwork={art}
      />
      <UnavailableNotice>
        Offers cannot be sent, withdrawn, or extended yet. Entered details
        remain on this screen.
      </UnavailableNotice>
      <PlannedRead
        loading={detail.loading || list.loading}
        error={detail.error || list.error}
        retry={() => {
          detail.reload();
          list.reload();
        }}
      />
      {list.data?.data.map((offer) => (
        <article key={offer.id} className={e.panel}>
          <h2>{offer.jobTitle}</h2>
          <p>
            {offer.status} · {offer.compensationSummary}
          </p>
          <Link href={`/employer/offers/${resourceId(offer.id)}`}>
            Manage offer
          </Link>
        </article>
      ))}
      <form
        onSubmit={(event) => {
          event.preventDefault();
          setDialog(create ? "send" : "extend");
        }}
      >
        <section className={e.section}>
          <h2>
            {create
              ? "Offer details"
              : detail.data?.status || "Pending response"}
          </h2>
          {create ? (
            <div className={e.grid}>
              <label className={e.field}>
                Compensation summary (required)
                <textarea name="compensationSummary" required />
              </label>
              <label className={e.field}>
                Proposed start date (required)
                <input name="proposedStartDate" type="date" required />
              </label>
              <label className={e.field}>
                Response deadline (required)
                <input name="responseDeadline" type="datetime-local" required />
              </label>
              <label className={e.field}>
                Message
                <textarea name="message" />
              </label>
            </div>
          ) : (
            <>
              <p>
                {detail.data?.compensationSummary ||
                  "No offer details available"}
              </p>
              <p>Proposed start: {detail.data?.proposedStartDate || "—"}</p>
              <p>
                Response deadline:{" "}
                {detail.data?.responseDeadline
                  ? new Date(detail.data.responseDeadline).toLocaleString()
                  : "—"}
              </p>
              <label className={e.field}>
                New response deadline
                <input name="responseDeadline" type="datetime-local" required />
              </label>
            </>
          )}
        </section>
        <section className={e.section}>
          <h2>{create ? "Capacity" : "Deadline extension"}</h2>
          <p>
            {create
              ? "If the posting limit is exceeded, confirmation is required before sending."
              : "A new deadline does not change the other terms of this offer."}
          </p>
          {capacity && (
            <p className={e.error} role="alert">
              Capacity must be verified before sending. This offer has not been
              sent.
            </p>
          )}
        </section>
        <div className={e.actions}>
          <Link href="/employer/jobs">Cancel</Link>
          {id && (
            <button type="button" onClick={() => setDialog("withdraw")}>
              Review withdrawal
            </button>
          )}
          <button type="submit" className={e.primary}>
            {create ? "Review offer" : "Review deadline extension"}
          </button>
        </div>
      </form>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "withdraw"
              ? "Withdraw this offer?"
              : dialog === "extend"
                ? "Extend offer deadline?"
                : capacity
                  ? "Send offer above posting capacity?"
                  : "Send this offer?"
          }
          description="Review the saved offer terms and recipient before continuing. No offer has been sent or changed."
          confirmLabel={
            dialog === "withdraw"
              ? "Withdraw offer"
              : dialog === "extend"
                ? "Extend deadline"
                : "Send offer"
          }
          intent={dialog === "withdraw" ? "danger" : "warning"}
          unavailable="This offer operation is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
export function FeedbackAuthor({ id }: { id: string }) {
  const params = useSearchParams();
  const [dialog, setDialog] = useState<string | null>(null);
  const [preview, setPreview] = useState(false);
  const rejection = params.get("state") === "rejection";
  const record = useApiResource<FeedbackView>(
    `/v1/${rejection ? "pipeline-executions" : "stage-executions"}/${resourceId(id)}/${rejection ? "rejection-feedback" : "feedback"}`,
  );
  return (
    <>
      <ProductHeading
        title={rejection ? "Review rejection feedback" : "Feedback authoring"}
        description="Draft, review, and release candidate feedback after a hiring decision"
      />
      <UnavailableNotice>
        Feedback cannot be saved, approved, or released yet. Private notes are
        never included in candidate feedback automatically.
      </UnavailableNotice>
      <PlannedRead
        loading={record.loading}
        error={record.error}
        retry={record.reload}
      />
      <section className={e.panel}>
        <h2>Feedback type</h2>
        <p>{rejection ? "Rejection feedback" : "Stage feedback"}</p>
      </section>
      <section className={e.panel}>
        <h2>Draft message</h2>
        <label className={e.field}>
          Candidate feedback
          <textarea
            key={record.data?.feedback || "empty"}
            name="feedback"
            defaultValue={record.data?.feedback || ""}
            maxLength={5000}
          />
        </label>
        {record.data?.feedbackIsAiGenerated && (
          <p>System generated draft · Human approval required</p>
        )}
      </section>
      <section className={e.panel}>
        <h2>Internal review</h2>
        <p>
          {record.data?.approvedAt
            ? "Approved"
            : "Hiring manager approval required"}
        </p>
      </section>
      <section className={e.panel}>
        <h2>Candidate visibility</h2>
        <p>
          {record.data?.releasedAt
            ? "Released"
            : "Hidden until approved and released. Release is recorded in the application timeline."}
        </p>
      </section>
      <div className={e.actions}>
        <button disabled>Save draft</button>
        <button onClick={() => setPreview(!preview)} aria-expanded={preview}>
          Preview candidate view
        </button>
        <button className={e.primary} onClick={() => setDialog("release")}>
          Approve and release
        </button>
      </div>
      {preview && (
        <section className={e.panel}>
          <h2>Candidate preview · Saved feedback</h2>
          <p>
            {record.data?.feedback ||
              "No saved candidate feedback is available."}
          </p>
        </section>
      )}
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title="Release candidate feedback?"
          description="Only approved candidate-facing feedback is released. Private scorecard notes remain confidential."
          confirmLabel="Approve and release"
          intent="success"
          unavailable="Feedback release is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
export function InterviewerManagement() {
  const team = useApiResource<InterviewerPage>(
    "/v1/companies/me/interviewers?limit=20&offset=0",
  );
  return (
    <>
      <ProductHeading
        title="Interviewer management"
        description="Manage company interviewers, availability, and assigned evaluations"
      />
      <UnavailableNotice>
        Interviewer invitations and assignment changes aren’t available yet.
      </UnavailableNotice>
      <PlannedRead
        loading={team.loading}
        error={team.error}
        retry={team.reload}
      />
      {team.data?.data.map((person) => (
        <article className={e.panel} key={person.id}>
          <h2>{person.name}</h2>
          <p>{person.accountStatus}</p>
        </article>
      ))}
      {!team.loading && !team.data?.data.length && (
        <EmptyState
          title="No interviewers available"
          description="Your company interviewers will appear here when available."
        />
      )}
      <div className={e.actions}>
        <button disabled>Invite interviewer</button>
        <button disabled>Assign interviews</button>
        <Link href="/employer/interviews">Manage interviews</Link>
      </div>
      <WorkspaceFooter />
    </>
  );
}
