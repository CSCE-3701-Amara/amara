"use client";

import { useState } from "react";
import { DagEditor } from "./dag-editor";
import Link from "next/link";
import { resourceId } from "@/shared/lib/api-client";
import { ProductHeading, WorkspaceFooter } from "@/shared/ui/product-ui";
import styles from "@/shared/ui/editor.module.css";
import { PipelineEditor } from "./pipeline-editor";

export function PipelineStudio({ postingId }: { postingId: string }) {
  const pipelinePath = `/employer/pipelines/${resourceId(postingId)}`;
  const [stageEditorVisited, setStageEditorVisited] = useState(false);

  return (
    <>
      <ProductHeading
        title="Pipeline designer"
        description="Design branches and joins between assessments and interviews for this job."
        compact
      />
      <nav className={styles.actions} aria-label="Pipeline views">
        <Link href={pipelinePath}>Pipeline</Link>
        <Link href={`${pipelinePath}/candidates`}>Candidates / decisions</Link>
        <Link href="/employer/interviews">Live interviews</Link>
        <Link href="/employer/assessments">Assessment instances</Link>
        <Link href={`${pipelinePath}/dashboard`}>Dashboard</Link>
      </nav>
      <DagEditor postingId={postingId} />
      <details
        onToggle={(event) => {
          if (event.currentTarget.open) setStageEditorVisited(true);
        }}
      >
        <summary>Legacy stage editor</summary>
        <p>
          Use this editor for pipelines created with the earlier stage format.
        </p>
        {stageEditorVisited && <PipelineEditor postingId={postingId} />}
      </details>
      <WorkspaceFooter />
    </>
  );
}
