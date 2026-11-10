"use client";
import Link from "next/link";
import {
  CandidatePipeline,
  type CandidatePipelineData,
} from "./candidate-pipeline";
import type { AssessmentInstance } from "../types/product.types";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
} from "@/shared/ui/product-ui";
import type { ApplicationDetail as ApplicationDetailData } from "../types/cand.types";
import styles from "./candidate.module.css";
/** Application detail uses the existing ownership-checked candidate read handler. */
export function ApplicationDetail({
  id,
  progress = false,
}: {
  id: string;
  progress?: boolean;
}) {
  const params = useSearchParams();
  const [dialog, setDialog] = useState<string | null>(null);
  const detail = useApiResource<ApplicationDetailData>(
    `/v1/applications/${resourceId(id)}`,
  );
  const stageProgress = useApiResource<CandidatePipelineData>(
    progress ? `/v1/applications/${resourceId(id)}/pipeline-progress` : null,
  );
  const assessments = useApiResource<AssessmentInstance[]>(
    `/v1/applications/${resourceId(id)}/assessment-instances`,
  );
  function assessmentLink(instance: AssessmentInstance) {
    return (
      <Link
        href={`/candidate/assessments/${resourceId(instance.id)}?state=availability&application=${resourceId(id)}`}
      >
        {instance.status === "in_progress"
          ? "Continue assessment"
          : instance.status === "available"
            ? "Open assessment"
            : "Review assessment"}
      </Link>
    );
  }
  const withdrawn =
    params.get("state") === "withdrawn" || detail.data?.status === "withdrawn";
  if (withdrawn)
    return (
      <section data-figma-node="224:137">
        <EmptyState
          level={1}
          title="Application withdrawn"
          description="This application is no longer active. You can still review its status and explore other opportunities."
        >
          <Link href="/candidate/applications" className={s.button}>
            Back to applications
          </Link>
          <Link
            href="/candidate/dashboard"
            className={`${s.button} ${s.primary}`}
          >
            Explore opportunities
          </Link>
        </EmptyState>
      </section>
    );
  return (
    <section data-figma-node={progress ? "125:51985" : "375:958"}>
      <ProductHeading
        title={progress ? "Application progress" : "Application details"}
        description={
          detail.data
            ? `${detail.data.jobTitle} · ${detail.data.companyName}${detail.data.submittedAt ? ` · Submitted ${new Date(detail.data.submittedAt).toLocaleDateString()}` : ""}`
            : undefined
        }
      />
      <BusyState loading={detail.loading} />
      <ApiFeedback error={detail.error} retry={detail.reload} />
      {detail.data && (
        <>
          <div className={progress ? styles.progressLayout : styles.detailGrid}>
            <section className={styles.panel}>
              <h2>{detail.data.jobTitle}</h2>
              <p>{detail.data.companyName}</p>
              {progress ? (
                <>
                  <h2>Your hiring journey</h2>
                  <BusyState loading={stageProgress.loading} />
                  {stageProgress.error?.status === 404 && (
                    <UnavailableNotice>
                      The hiring pipeline isn’t available for this application
                      yet.
                    </UnavailableNotice>
                  )}
                  {stageProgress.data && (
                    <p>
                      {
                        stageProgress.data.stages.filter(
                          (stage) => stage.status === "completed",
                        ).length
                      }{" "}
                      of {stageProgress.data.stages.length} stages completed
                    </p>
                  )}
                  {stageProgress.data && (
                    <CandidatePipeline
                      data={stageProgress.data}
                      assessments={assessments.data ?? []}
                      applicationId={id}
                    />
                  )}
                  {stageProgress.error?.status !== 404 && (
                    <ApiFeedback
                      error={stageProgress.error}
                      retry={stageProgress.reload}
                    />
                  )}
                  <StatusBadge>
                    {detail.data.status.replaceAll("_", " ")}
                  </StatusBadge>
                  <div className={styles.actions}>
                    <Link
                      href={`/candidate/applications/${resourceId(id)}`}
                      className={s.button}
                    >
                      View application details
                    </Link>
                  </div>
                </>
              ) : (
                <>
                  <h2>Application timeline</h2>
                  <ol className={styles.timeline}>
                    {detail.data.submittedAt && (
                      <li>
                        <time dateTime={detail.data.submittedAt}>
                          {new Date(
                            detail.data.submittedAt,
                          ).toLocaleDateString()}
                        </time>
                        <h3>Application submitted</h3>
                        <p>Your application was received.</p>
                      </li>
                    )}
                    {detail.data.screeningCompletedAt && (
                      <li>
                        <time dateTime={detail.data.screeningCompletedAt}>
                          {new Date(
                            detail.data.screeningCompletedAt,
                          ).toLocaleDateString()}
                        </time>
                        <h3>Initial review</h3>
                        <p>
                          {detail.data.screenedOut
                            ? "Your application did not pass initial screening."
                            : "Your application passed initial screening."}
                        </p>
                      </li>
                    )}
                    {detail.data.closed && detail.data.updatedAt && (
                      <li>
                        <time dateTime={detail.data.updatedAt}>
                          {new Date(detail.data.updatedAt).toLocaleDateString()}
                        </time>
                        <h3>Decision released</h3>
                        <p>The hiring team has closed this application.</p>
                      </li>
                    )}
                  </ol>
                  {!detail.data.screenedOut && (
                    <Link
                      href={`/candidate/applications/${resourceId(id)}/progress`}
                      className={s.button}
                    >
                      View hiring pipeline
                    </Link>
                  )}
                </>
              )}
            </section>
            <aside className={styles.statusPanel}>
              <StatusBadge>
                {detail.data.status.replaceAll("_", " ")}
              </StatusBadge>
              <h2>
                {detail.data.closed
                  ? "Application closed"
                  : "Your application is in progress"}
              </h2>
              <h3>Assessments</h3>
              <BusyState loading={assessments.loading} />
              <ApiFeedback
                error={assessments.error}
                retry={assessments.reload}
              />
              {assessments.data?.map((instance) => (
                <p key={instance.id}>
                  {assessmentLink(instance)} ·{" "}
                  {instance.status.replaceAll("_", " ")}
                </p>
              ))}
              {assessments.data?.length === 0 && (
                <p>No assessments assigned.</p>
              )}
              {detail.data.finalFeedback && (
                <>
                  <h3>Hiring team feedback</h3>
                  <p>{detail.data.finalFeedback}</p>
                </>
              )}
              <h3>Submitted information</h3>
              <p>Reference: {detail.data.referenceNumber}</p>
              <p>
                {detail.data.submittedAt
                  ? new Date(detail.data.submittedAt).toLocaleString()
                  : "Submission time unavailable"}
              </p>
              <h3>Next steps</h3>
              <p>
                {detail.data.closed
                  ? "Your feedback remains available here for future reference."
                  : "Check back here for updates from the hiring team."}
              </p>
            </aside>
          </div>
          <div className={styles.actions}>
            <Link href="/candidate/applications" className={s.button}>
              Back to applications
            </Link>
            <Link
              href="/candidate/dashboard"
              className={`${s.button} ${s.primary}`}
            >
              Explore opportunities
            </Link>
            {!detail.data.closed && (
              <button type="button" onClick={() => setDialog("withdraw")}>
                Withdraw application
              </button>
            )}
          </div>
        </>
      )}
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "withdraw-review"
              ? "Review before withdrawing"
              : "Withdraw this application?"
          }
          description="Withdrawing closes your application and removes it from the active hiring process. Review your decision before continuing."
          confirmLabel="Withdraw application"
          intent="danger"
          unavailable="Application withdrawal isn’t available yet. Your application will remain unchanged."
          onClose={() => setDialog(null)}
        />
      )}
    </section>
  );
}
