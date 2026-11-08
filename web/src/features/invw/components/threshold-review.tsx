"use client";
import { toError } from "@/shared/lib/errors";

import { useEffect, useRef, useState, type FormEvent } from "react";
import { ApiError, apiRequest, resourceId } from "@/shared/lib/api-client";
import { ApiFeedback } from "@/shared/ui/workspace";
import e from "@/shared/ui/editor.module.css";

type Outcome = {
  candidate: string;
  status: "held" | "passed" | "rejected";
  score: number;
};

/** Preview affected candidates before applying a threshold change. */
export function ThresholdReview({
  stageId,
  current,
  onSaved,
}: {
  stageId: string;
  current: number;
  onSaved: () => void;
}) {
  const [value, setValue] = useState(String(current));
  const [preview, setPreview] = useState<Outcome[] | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [status, setStatus] = useState("");
  const request = useRef<AbortController | null>(null);
  useEffect(() => () => request.current?.abort(), []);

  async function change(apply: boolean) {
    if (request.current || (apply && preview === null)) return;
    const threshold = Number(value);
    if (
      !value.trim() ||
      !Number.isInteger(threshold) ||
      threshold < 0 ||
      threshold > 100
    ) {
      setError(new Error("Enter a whole-number threshold between 0 and 100."));
      return;
    }
    const controller = new AbortController();
    request.current = controller;
    setBusy(true);
    setError(undefined);
    setStatus("");
    try {
      const outcomes = await apiRequest<Outcome[]>(
        `/v1/pipeline-stages/${resourceId(stageId)}/threshold${apply ? "" : "/preview"}`,
        {
          method: apply ? "PUT" : "POST",
          body: {
            threshold,
            expectedThreshold: current,
            ...(apply ? { preview } : {}),
          },
          signal: controller.signal,
        },
      );
      if (controller.signal.aborted) return;
      if (apply) {
        setPreview(null);
        setStatus("Threshold saved. Refreshing candidate outcomes.");
        onSaved();
      } else {
        setPreview(outcomes);
        setStatus(`${outcomes.length} candidate outcomes ready to review.`);
      }
    } catch (cause) {
      if (controller.signal.aborted) return;
      setError(toError(cause));
      if (cause instanceof ApiError && cause.status === 409) {
        setPreview(null);
        onSaved();
      }
    } finally {
      if (request.current === controller) {
        request.current = null;
        setBusy(false);
      }
    }
  }

  return (
    <section className={e.panel} aria-label="Threshold review" aria-busy={busy}>
      <h2>Review threshold change</h2>
      <p>
        Current pass threshold: {current} / 100. Review the affected candidates
        before applying a new threshold.
      </p>
      <form
        onSubmit={(event: FormEvent) => {
          event.preventDefault();
          void change(false);
        }}
      >
        <label className={e.field}>
          Proposed pass threshold
          <input
            type="number"
            min={0}
            max={100}
            step={1}
            required
            value={value}
            disabled={busy}
            onChange={(event) => {
              setValue(event.target.value);
              setPreview(null);
              setStatus("");
            }}
          />
        </label>
        <ApiFeedback error={error} />
        <p role="status">{status}</p>
        {preview !== null && (
          <div className={e.tableWrap}>
            <table className={e.table}>
              <caption>Candidate outcomes at the proposed threshold</caption>
              <thead>
                <tr>
                  <th scope="col">Candidate</th>
                  <th scope="col">Score</th>
                  <th scope="col">Outcome</th>
                </tr>
              </thead>
              <tbody>
                {preview.map((row) => (
                  <tr key={row.candidate}>
                    <td>{row.candidate}</td>
                    <td>{row.score}</td>
                    <td>{row.status}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            {preview.length === 0 && <p>No candidates are affected.</p>}
          </div>
        )}
        <div className={e.actions}>
          <button type="submit" disabled={busy}>
            Preview outcomes
          </button>
          <button
            type="button"
            className={e.primary}
            disabled={busy || preview === null}
            onClick={() => void change(true)}
          >
            Apply reviewed threshold
          </button>
        </div>
      </form>
    </section>
  );
}
