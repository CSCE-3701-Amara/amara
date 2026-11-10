import { workspaceStyles as s } from "@/shared/ui/workspace";
import { StatusBadge } from "@/shared/ui/product-ui";
import type {
  AssessmentAttempt,
  AssessmentInstance,
  AssessmentResult,
  CandidateQuestion,
} from "../types/product.types";
import styles from "./candidate.module.css";

export function AssessmentResultView({
  result,
}: {
  result?: AssessmentResult;
}) {
  return (
    <>
      <section className={result ? s.success : styles.panel}>
        <StatusBadge tone={result ? "success" : "neutral"}>
          {result ? "Result available" : "Result unavailable"}
        </StatusBadge>
        <h2>
          {result ? "Assessment completed" : "Your result hasn’t been released"}
        </h2>
        <p>
          {result
            ? "Your score and reviewer feedback are now available."
            : "Scores and reviewer feedback appear only after disclosure."}
        </p>
      </section>
      {result && (
        <>
          <div className={styles.detailGrid}>
            <section className={styles.panel}>
              <p>Overall score</p>
              <strong className={styles.clock}>
                {result.score !== null && result.score !== undefined
                  ? `${result.score} / 100`
                  : "Not scored"}
              </strong>
            </section>
            <section className={styles.panel}>
              <p>Assessment outcome</p>
              <h2>
                {result.passed
                  ? "Threshold achieved"
                  : "Threshold not achieved"}
              </h2>
              <p>
                This outcome does not automatically determine whether you
                progress to the next hiring stage.
              </p>
            </section>
          </div>
          <section className={styles.panel}>
            <h2>Reviewer feedback</h2>
            <p>{result.feedback || "No feedback was released."}</p>
            <p>
              {result.feedbackIsAiGenerated ? "AI-generated feedback · " : ""}
              Released {new Date(result.disclosedAt).toLocaleString()}
            </p>
          </section>
        </>
      )}
    </>
  );
}

export function AssessmentAvailabilityView({
  instance,
  busy,
  onOpen,
}: {
  instance?: AssessmentInstance;
  busy: boolean;
  onOpen: () => Promise<void>;
}) {
  return (
    <>
      <section className={styles.panel}>
        <StatusBadge>
          {instance?.status || "Availability unavailable"}
        </StatusBadge>
        <h2>Technical skills assessment</h2>
        <p>
          {instance?.deadline
            ? `Available until ${new Date(instance.deadline).toLocaleString()}`
            : "Assessment availability will appear here when it is released."}
        </p>
      </section>
      <div className={styles.detailGrid}>
        <section className={styles.panel}>
          <h2>Continue where you left off</h2>
          <p>Your saved answers can be restored while the attempt is active.</p>
          <button
            type="button"
            disabled={busy || !instance || instance.status !== "in_progress"}
            onClick={() => void onOpen()}
          >
            {busy ? "Opening…" : "Resume assessment"}
          </button>
        </section>
        <section className={styles.panel}>
          <h2>Start a new attempt</h2>
          <p>
            Check your connection, available time, and the completion deadline
            before starting.
          </p>
          <button
            type="button"
            className={s.primary}
            disabled={busy || (!!instance && instance.status !== "available")}
            onClick={() => void onOpen()}
          >
            {busy ? "Opening…" : "Start assessment"}
          </button>
        </section>
      </div>
      <section className={styles.panel}>
        <h2>Countdown and auto-submit</h2>
        <p>
          Starting begins the active timer. The server enforces the completion
          deadline and locks closed attempts. Opening this overview does not
          start an attempt.
        </p>
      </section>
    </>
  );
}

export function AssessmentLockedView({
  locked,
  attempt,
  questions,
}: {
  locked: boolean;
  attempt?: AssessmentAttempt;
  questions: CandidateQuestion[];
}) {
  return (
    <div className={styles.detailGrid}>
      <div>
        <section className={styles.panel}>
          <h2>
            {locked
              ? "Your answers are saved and locked"
              : "Submission status unavailable"}
          </h2>
          <StatusBadge>{attempt?.status || "Not confirmed"}</StatusBadge>
          <p>
            {locked
              ? "This attempt is closed. Its saved responses can no longer be edited."
              : "A submission can only be confirmed by the assessment service."}
          </p>
          {attempt?.submittedAt && (
            <p>Submitted {new Date(attempt.submittedAt).toLocaleString()}</p>
          )}
        </section>
        <section className={styles.panel}>
          <h2>Submission receipt</h2>
          {questions.map((question) => (
            <p key={question.id}>
              {question.position + 1} · {question.prompt} · Locked
            </p>
          ))}
        </section>
      </div>
      <aside className={styles.panel}>
        <h2>What happens next</h2>
        <p>
          Assessment evaluation and stage progression are separate. A pending
          result is not a rejection or hiring decision.
        </p>
      </aside>
    </div>
  );
}
