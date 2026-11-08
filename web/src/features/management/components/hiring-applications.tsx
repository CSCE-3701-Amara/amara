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
} from "@/shared/lib/contract.types";
import type { RuntimePipelineDashboard } from "@/features/invw/types/runtime.types";
import type { PipelineDraft } from "@/features/invw/types/invw.types";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  Stats,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import { PlannedRead } from "./planned-read";

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
    view === "comparison" ? `/v1/jobs/${resourceId(id)}/pipeline` : null,
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
  const dashboard = useApiResource<RuntimePipelineDashboard>(
    view === "dashboard"
      ? `/v1/jobs/${resourceId(id)}/pipeline/dashboard`
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

        className={
          view === "applications" || view === "finalists"
            ? e.titleBanner
            : undefined
        }
      />
      <UnavailableNotice>
        {view === "dashboard"
          ? "Offer counts are not available yet. Select a stage to review its candidates."
          : "This hiring view is not available yet. Recorded data can’t be changed. Identity stays hidden until the server authorizes its release."}
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
                <tr key={stage.id}>
                  <td>
                    <Link
                      href={`/employer/pipelines/${resourceId(id)}/candidates?stage=${resourceId(stage.id)}`}
                    >
                      {stage.name}
                    </Link>
                  </td>
                  <td>{stage.active}</td>
                  <td>{stage.passed}</td>
                  <td>{stage.rejected}</td>
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
                value: dashboard.data?.active ?? "—",
              },
              {
                label: "Pending offers",
                value: "Unavailable",
              },
              {
                label: "Accepted offers",
                value: "Unavailable",
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
