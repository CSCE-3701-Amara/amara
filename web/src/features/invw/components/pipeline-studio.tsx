"use client";
import Link from "next/link";
import dynamic from "next/dynamic";
import { useState } from "react";
import { useSearchParams } from "next/navigation";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type {
  StageTransition,
  PipelineMigrationPreview,
} from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import type { PipelineDraft, PipelineStage } from "../types/invw.types";
import e from "@/shared/ui/editor.module.css";
import s from "./pipeline-studio.module.css";
import { PipelineGraph } from "./pipeline-graph";

const PipelineEditor = dynamic(() =>
  import("./pipeline-editor").then((module) => module.PipelineEditor),
);

function Draft({
  pipeline,
  published,
}: {
  pipeline: PipelineDraft;
  published: boolean;
}) {
  const stages = useApiResource<PipelineStage[] | null>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/stages`,
  );
  const transitions = useApiResource<StageTransition[]>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/transitions`,
  );
  return (
    <>
      <BusyState loading={stages.loading || transitions.loading} />
      <ApiFeedback error={stages.error} retry={stages.reload} />
      {transitions.error?.status !== 404 && (
        <ApiFeedback error={transitions.error} retry={transitions.reload} />
      )}
      <StatusBadge>
        {published ? "Published version" : "Draft"} {pipeline.versionNumber}
      </StatusBadge>
      {!stages.loading && !stages.error && !transitions.loading && (
        <PipelineGraph
          key={pipeline.versionId}
          initial={stages.data || []}
          transitions={transitions.data}
          published={published}
        />
      )}
    </>
  );
}

export function PipelineStudio({ postingId }: { postingId: string }) {
  const params = useSearchParams();
  const state = params.get("state");
  const [dialog, setDialog] = useState<string | null>(null);
  const [stageEditorVisited, setStageEditorVisited] = useState(false);
  const pipeline = useApiResource<PipelineDraft>(
    `/v1/jobs/${resourceId(postingId)}/pipeline`,
  );
  const migration = useApiResource<PipelineMigrationPreview>(null);
  return (
    <>
      <section
        className={s.header}
        aria-label="Pipeline navigation and version"
      >
        <ProductHeading
          title={
            state === "published"
              ? "Pipeline · Published version"
              : state?.startsWith("migration")
                ? "Pipeline · Migration review"
                : "Pipeline designer"
          }
          description="Draft changes are private until explicitly published."
          compact
          className={s.heading}
        />
        <nav className={e.actions} aria-label="Pipeline views">
          <Link href={`/employer/pipelines/${resourceId(postingId)}`}>
            Pipeline
          </Link>
          <Link
            href={`/employer/pipelines/${resourceId(postingId)}/candidates`}
          >
            Candidates / decisions
          </Link>
          <Link href="/employer/interviews">Live interviews</Link>
          <Link
            href={`/employer/pipelines/${resourceId(postingId)}?state=migration-review`}
          >
            Versions / migration
          </Link>
          <Link href="/employer/assessments">Assessment instances</Link>
          <Link href={`/employer/pipelines/${resourceId(postingId)}/dashboard`}>
            Dashboard
          </Link>
        </nav>
        <div className={e.actions}>
          <button disabled>Save graph draft</button>
          <button onClick={() => setDialog("validate")}>Validate</button>
          <button className={e.primary} onClick={() => setDialog("publish")}>
            Review publication
          </button>
          <button onClick={() => setDialog("discard")}>Discard draft</button>
        </div>
        <p className={s.versionNote}>
          {pipeline.data
            ? `Version ${pipeline.data.versionNumber}`
            : "Pipeline draft"}{" "}
          · {state === "published" ? "Read only" : "Private draft"}
        </p>
      </section>
      {state === "validation" && (
        <section className={s.state}>
          <h2>Graph validation</h2>
          <p>
            Validation is unavailable. The server cannot verify reachability,
            branch joins, or publication constraints yet.
          </p>
        </section>
      )}
      {state?.startsWith("migration") && (
        <section className={s.state}>
          <h2>
            {state === "migration-blocked"
              ? "Migration blocked"
              : "Migration review"}
          </h2>
          <p>
            Active candidate placements must be checked against the new graph
            before publication. No migration preview can be loaded yet.
          </p>
          {migration.data?.affectedExecutions}
          <button disabled>Confirm migration</button>
        </section>
      )}
      <BusyState loading={pipeline.loading} />
      {pipeline.error?.status !== 404 && (
        <ApiFeedback error={pipeline.error} retry={pipeline.reload} />
      )}{" "}
      {pipeline.data ? (
        <Draft pipeline={pipeline.data} published={state === "published"} />
      ) : (
        !pipeline.loading && <PipelineGraph initial={[]} published={false} />
      )}
      <section className={e.section} aria-label="Privacy and feedback policy">
        <h2>Privacy and feedback</h2>
        <div className={e.grid}>
          <label className={e.field}>
            Candidate anonymity
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Identity reveal point
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Rejection feedback
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Selected stage feedback
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
        </div>
        <p>
          Identity must be revealed before incompatible video and interview
          stages on every path. Policy loading and saving are not available yet.
        </p>
        <div className={e.actions}>
          <button disabled>Save policy</button>
        </div>
      </section>
      <details
        onToggle={(event) => {
          if (event.currentTarget.open) setStageEditorVisited(true);
        }}
      >
        <summary>Server-supported stage configuration</summary>
        <p>
          This existing editor saves stage definitions through the current API.
          It does not save graph connections.
        </p>
        {stageEditorVisited && <PipelineEditor postingId={postingId} />}
      </details>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "discard"
              ? "Discard pipeline draft?"
              : dialog === "migrate"
                ? "Confirm pipeline migration?"
                : dialog === "validate"
                  ? "Validate pipeline graph?"
                  : "Publish pipeline version?"
          }
          description="Review the graph, schedules, candidate placements, and required stage configuration before proceeding."
          confirmLabel={
            dialog === "discard"
              ? "Discard draft"
              : dialog === "validate"
                ? "Validate"
                : "Confirm publication"
          }
          intent={dialog === "discard" ? "danger" : "warning"}
          unavailable="The required graph operation is not available yet. Your local graph has not been published or discarded."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
