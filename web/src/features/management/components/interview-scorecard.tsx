"use client";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { InterviewFeedback } from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

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
