"use client";
import Link from "next/link";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import { transitionPosting } from "../api/job.api";
import type { Posting } from "../types/job.types";
import e from "@/shared/ui/editor.module.css";
export function PostingPreview({
  id,
  review = false,
}: {
  id: string;
  review?: boolean;
}) {
  const posting = useApiResource<Posting>(`/v1/jobs/${resourceId(id)}`);
  const [confirm, setConfirm] = useState(false);
  return (
    <>
      <ProductHeading
        title={review ? "Ready to publish?" : "Job posting preview"}
        description={posting.data?.title}
        artwork={{
          src: "/figma/publication-review.svg",
          width: 345.6,
          height: 238.8,
        }}
      />
      <BusyState loading={posting.loading} />
      <ApiFeedback error={posting.error} retry={posting.reload} />
      {posting.data && (
        <article className={e.panel}>
          <h2>{review ? "Posting summary" : posting.data.title}</h2>
          <StatusBadge>{posting.data.status}</StatusBadge>
          <h3>
            {posting.data.seniority_level} · {posting.data.employment_type} ·{" "}
            {posting.data.location}
          </h3>
          <p>{posting.data.description}</p>
          {(review || posting.data.salary_visible) && (
            <p>
              {posting.data.salary_min ?? "—"}–{posting.data.salary_max ?? "—"}{" "}
              {posting.data.currency} ·{" "}
              {posting.data.open_positions ?? "Unspecified"} openings
            </p>
          )}
          <h3>Required amarat</h3>
          <ul>
            {posting.data.amarat.map((amara, index) => (
              <li key={amara.id || index}>
                {amara.title}
                {amara.description ? ` · ${amara.description}` : ""}
              </li>
            ))}
          </ul>
          <h3>Requirements</h3>
          <ul>
            {posting.data.requirements.map((requirement, index) => (
              <li key={requirement.id || index}>
                {requirement.text} ·{" "}
                {requirement.requirement_type.replaceAll("_", " ")}
              </li>
            ))}
          </ul>
          <p>
            Applications:{" "}
            {posting.data.application_opening
              ? new Date(posting.data.application_opening).toLocaleString()
              : "Open date unspecified"}{" "}
            –{" "}
            {posting.data.application_deadline
              ? new Date(posting.data.application_deadline).toLocaleString()
              : "No deadline specified"}
          </p>
          <p className={e.note}>
            This preview uses the saved posting. Save your changes before
            reviewing publication.
          </p>
          <div className={e.actions}>
            {review && posting.data.status !== "published" && (
              <button className={e.primary} onClick={() => setConfirm(true)}>
                Confirm publication
              </button>
            )}
            <Link href={`/employer/jobs/${resourceId(id)}`}>
              Back to editor
            </Link>
            <Link href={`/employer/jobs/${resourceId(id)}/applications`}>
              Applications
            </Link>
            <Link href={`/employer/jobs/${resourceId(id)}/finalists`}>
              Finalists
            </Link>
          </div>
        </article>
      )}
      <WorkspaceFooter />
      {confirm && posting.data && (
        <ConfirmDialog
          title="Publish this posting?"
          description="Publishing makes the saved posting available to eligible candidates. The server checks all publication requirements."
          confirmLabel="Publish posting"
          onClose={() => setConfirm(false)}
          onConfirm={async () => {
            await transitionPosting(
              id,
              posting.data!.status === "archived" ? "republish" : "publish",
            );
            posting.reload();
            setConfirm(false);
          }}
        />
      )}
    </>
  );
}
