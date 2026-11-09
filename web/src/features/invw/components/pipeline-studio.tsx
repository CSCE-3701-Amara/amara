import Link from "next/link";
import { resourceId } from "@/shared/lib/api-client";
import { ProductHeading, WorkspaceFooter } from "@/shared/ui/product-ui";
import styles from "@/shared/ui/editor.module.css";
import { PipelineEditor } from "./pipeline-editor";

export function PipelineStudio({ postingId }: { postingId: string }) {
  const pipelinePath = `/employer/pipelines/${resourceId(postingId)}`;

  return (
    <>
      <ProductHeading
        title="Pipeline designer"
        description="Build an ordered sequence of assessments and interviews for this job."
        compact
      />
      <nav className={styles.actions} aria-label="Pipeline views">
        <Link href={pipelinePath}>Pipeline</Link>
        <Link href={`${pipelinePath}/candidates`}>Candidates / decisions</Link>
        <Link href="/employer/interviews">Live interviews</Link>
        <Link href="/employer/assessments">Assessment instances</Link>
        <Link href={`${pipelinePath}/dashboard`}>Dashboard</Link>
      </nav>
      <PipelineEditor postingId={postingId} />
      <WorkspaceFooter />
    </>
  );
}
