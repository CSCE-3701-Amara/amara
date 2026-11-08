"use client";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { FeedbackView } from "@/shared/lib/contract.types";
import {
  ConfirmDialog,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import { PlannedRead } from "./planned-read";

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
