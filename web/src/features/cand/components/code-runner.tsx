"use client";
import { useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import styles from "./candidate.module.css";
export function CodeRunner({
  attemptId,
  questionId,
  source,
  language,
  disabled,
}: {
  attemptId: string;
  questionId: string;
  source: string;
  language: string;
  disabled: boolean;
}) {
  const [stdin, input] = useState("");
  const [output, setOutput] = useState("");
  const [busy, setBusy] = useState(false);
  async function run() {
    if (busy) return;
    setBusy(true);
    try {
      const r = await apiRequest<{ stdout: string; succeeded: boolean }>(
        `/v1/assessment-attempts/${resourceId(attemptId)}/run`,
        { method: "POST", body: { questionId, source, language, stdin } },
      );
      setOutput(
        r.stdout +
          (r.succeeded ? "" : "\nProgram did not complete successfully."),
      );
    } catch (e) {
      setOutput(e instanceof Error ? e.message : "Execution failed");
    } finally {
      setBusy(false);
    }
  }
  return (
    <section className={styles.codeRunner} aria-label="Code execution">
      <label>
        Standard input
        <textarea
          value={stdin}
          onChange={(e) => input(e.target.value)}
          maxLength={16384}
        />
      </label>
      <button disabled={disabled || busy} onClick={() => void run()}>
        {busy ? "Running…" : "Run code"}
      </button>
      <pre role="status" aria-live="polite">
        {output}
      </pre>
      <p>Runs use your input. Saving and submitting are separate actions.</p>
    </section>
  );
}
