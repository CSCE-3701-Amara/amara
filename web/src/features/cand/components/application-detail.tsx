"use client";
import Link from "next/link";
import Image from "next/image";
import type { CandidateStageProgressView } from "@/shared/lib/contract.types";
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
  const stages = useApiResource<CandidateStageProgressView[]>(
    progress ? `/v1/applications/${resourceId(id)}/stage-progress` : null,
  );
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
          <div className={styles.detailGrid}>
            <section className={styles.panel}>
              <h2>{detail.data.jobTitle}</h2>
              <p>{detail.data.companyName}</p>
              {progress ? (
                <>
                  <h2>Your hiring journey</h2>
                  <UnavailableNotice>
                    Stage-by-stage progress isn’t available yet. Your recorded
                    application status is shown below.
                  </UnavailableNotice>
                  <p>
                    {stages.data
                      ? `${stages.data.filter((stage) => stage.status === "passed").length} of ${stages.data.length} stages completed`
                      : "Stage details unavailable"}
                  </p>
                  <div className={styles.journey}>
                    {stages.data?.map((stage, index) => (
                      <div className={styles.journeyStep} key={stage.id}>
                        <article className={styles.panel}>
                          <StatusBadge
                            tone={
                              stage.status === "passed" ? "success" : "neutral"
                            }
                          >
                            {stage.status.replaceAll("_", " ")}
                          </StatusBadge>
                          <h3>{stage.name}</h3>
                          {stage.deadline && (
                            <p>
                              Due {new Date(stage.deadline).toLocaleString()}
                            </p>
                          )}
                          {stage.score !== null &&
                            stage.score !== undefined && (
                              <p>Score: {stage.score} / 100</p>
                            )}
                          {stage.feedback && <p>{stage.feedback}</p>}
                          {stage.amaraAwarded && (
                            <p>Earned amara: {stage.amaraAwarded}</p>
                          )}
                        </article>
                        {index < (stages.data?.length || 0) - 1 && (
                          <Image
                            src="/figma/journey-connector.svg"
                            alt=""
                            width={74}
                            height={14}
                          />
                        )}
                      </div>
                    ))}
                  </div>
                  {stages.error?.status !== 404 && (
                    <ApiFeedback error={stages.error} retry={stages.reload} />
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
