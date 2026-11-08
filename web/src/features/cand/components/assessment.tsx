"use client";
import { toError } from "@/shared/lib/errors";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useId, useRef, useState } from "react";
import { useUnsavedChanges } from "@/shared/lib/use-unsaved-changes";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import type {
  AssessmentAttempt,
  AssessmentInstance,
  AssessmentResult,
  CandidateQuestion,
  CandidateSubmission,
} from "../types/product.types";
import styles from "./candidate.module.css";
import {
  startAssessment,
  saveAssessmentAnswer,
  submitAssessmentAttempt,
} from "../api/assessment.api";
import { AttemptClock } from "./attempt-clock";
import { QuestionInput } from "./assessment-input";
import { AdvancedAssessment } from "./assessment-workspace";

/** Candidate-safe question projections only; unreleased grading evidence is never shown. */
export function CandidateAssessment({
  id,
  variant = "questions",
}: {
  id: string;
  variant?: "questions" | "coding" | "video";
}) {
  const params = useSearchParams();
  const router = useRouter();
  const state = params.get("state");
  const attemptId = params.get("attempt") || id;
  const prefix = useId();
  const [dialog, setDialog] = useState<string | null>(null);
  const pending = useRef(false);
  const answerLanguages = useRef<Record<string, string>>({});
  const dirtyQuestions = useRef(new Set<string>());
  const [dirty, setDirty] = useState(false);
  useUnsavedChanges(dirty);
  const [busy, setBusy] = useState(false);
  const [actionError, setActionError] = useState<Error>();
  const [message, setMessage] = useState("");
  const [savedAnswers, setSavedAnswers] = useState<CandidateSubmission[]>([]);
  const [submission, setSubmission] = useState<AssessmentAttempt>();
  const applicationId = params.get("application");
  const instances = useApiResource<AssessmentInstance[]>(
    state === "availability" && applicationId
      ? `/v1/applications/${resourceId(applicationId)}/assessment-instances`
      : null,
  );
  const [answers, setAnswers] = useState<Record<string, string | number[]>>({});
  const attempt = useApiResource<AssessmentAttempt>(
    state === "availability" || state === "result"
      ? null
      : `/v1/assessment-attempts/${resourceId(attemptId)}`,
    state || "",
  );
  const attemptData = submission?.id === attemptId ? submission : attempt.data;
  const instanceData = instances.data?.find((item) => item.id === id);
  const result = useApiResource<AssessmentResult>(
    state === "result"
      ? `/v1/assessment-instances/${resourceId(id)}/result`
      : null,
  );
  const error = attempt.error || instances.error || result.error;
  const loading = attempt.loading || instances.loading || result.loading;
  const locked =
    !!attemptData &&
    ["submitted", "pending_grading", "scored", "missed"].includes(
      attemptData.status,
    );
  const questions = (attemptData?.questions || []).filter((question) =>
    variant === "coding"
      ? question.responseFormat === "code"
      : variant === "video"
        ? question.responseFormat === "video"
        : true,
  );
  function answer(question: CandidateQuestion): string | number[] {
    if (Object.hasOwn(answers, question.id)) return answers[question.id];
    const response = savedAnswers?.find(
      (value) => value.questionId === question.id,
    )?.response;
    if (typeof response === "string" || Array.isArray(response))
      return response;
    if (response && "code" in response) return response.code;
    return "";
  }
  const canEdit = attemptData?.status === "in_progress";
  async function writeAnswers(only?: CandidateQuestion, language?: string) {
    const targets = only
      ? [only]
      : questions.filter((question) => dirtyQuestions.current.has(question.id));
    const results = await Promise.allSettled(
      targets
        .filter((question) => question.responseFormat !== "video")
        .map((question) => {
          const value = answer(question);
          return saveAssessmentAnswer(
            attemptId,
            question.id,
            question.responseFormat === "code"
              ? {
                  language:
                    language ||
                    answerLanguages.current[question.id] ||
                    "javascript",
                  code: typeof value === "string" ? value : "",
                }
              : value,
          );
        }),
    );
    const responses = results.flatMap((result) =>
      result.status === "fulfilled" ? [result.value] : [],
    );
    responses.forEach((response) =>
      dirtyQuestions.current.delete(response.questionId),
    );
    setDirty(dirtyQuestions.current.size > 0);
    setSavedAnswers((previous) => [
      ...previous.filter(
        (value) =>
          !responses.some((row) => row.questionId === value.questionId),
      ),
      ...responses,
    ]);
    const failure = results.find((result) => result.status === "rejected");
    if (failure?.status === "rejected") throw failure.reason;
  }
  async function saveProgress(question?: CandidateQuestion, language?: string) {
    if (pending.current || !canEdit) return;
    if (!question && dirtyQuestions.current.size === 0) {
      setMessage("No changes to save.");
      return;
    }
    pending.current = true;
    setBusy(true);
    setActionError(undefined);
    setMessage("");
    try {
      await writeAnswers(question, language);
      setMessage("Progress saved.");
    } catch (cause) {
      setActionError(toError(cause));
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  async function openAttempt() {
    if (pending.current) return;
    pending.current = true;
    setBusy(true);
    setActionError(undefined);
    try {
      const opened = await startAssessment(id);
      router.replace(
        `/candidate/assessments/${resourceId(id)}?attempt=${resourceId(opened.id)}${applicationId ? "&application=" + resourceId(applicationId) : ""}`,
      );
    } catch (cause) {
      setActionError(toError(cause));
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  async function confirmSubmission() {
    if (pending.current || !canEdit)
      throw new Error("This attempt is not open for submission.");
    pending.current = true;
    setBusy(true);
    try {
      await writeAnswers();
      setSubmission(await submitAssessmentAttempt(attemptId));
      router.replace(
        `/candidate/assessments/${resourceId(id)}?attempt=${resourceId(attemptId)}&state=submitted`,
      );
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  const title =
    state === "availability"
      ? "Your assessment"
      : state === "result"
        ? "Assessment result"
        : state === "submitted"
          ? "Coding challenge · Submission details"
          : variant === "coding"
            ? "Coding challenge"
            : variant === "video"
              ? "Behavioral response"
              : "Skills assessment";
  return (
    <section
      className={variant !== "questions" ? styles.advancedScreen : undefined}
      data-figma-node={
        state === "availability"
          ? "383:1196"
          : state === "result"
            ? "383:675"
            : state === "submitted"
              ? "125:52481"
              : variant === "coding"
                ? "125:52134"
                : variant === "video"
                  ? "125:52370"
                  : "383:553"
      }
    >
      <ProductHeading
        title={title}
        description={
          state === "result"
            ? "Your assessment result is separate from the hiring-stage decision."
            : state === "availability"
              ? "Check availability, continue an existing attempt, or start when the assessment opens."
              : "Review the questions and your saved responses for this assessment instance."
        }
      />
      <UnavailableNotice>
        Saved answers cannot yet be restored after reloading because the
        answer-read endpoint is unavailable. Save before leaving. Video uploads
        and result disclosure are not available yet.
      </UnavailableNotice>
      <ApiFeedback error={actionError} />
      {message && <p role="status">{message}</p>}
      <BusyState loading={loading} />
      {error?.status !== 404 && (
        <ApiFeedback
          error={error}
          retry={() => {
            attempt.reload();
            instances.reload();
            result.reload();
          }}
        />
      )}
      {state === "result" ? (
        <>
          <section className={result.data ? s.success : styles.panel}>
            <StatusBadge tone={result.data ? "success" : "neutral"}>
              {result.data ? "Result available" : "Result unavailable"}
            </StatusBadge>
            <h2>
              {result.data
                ? "Assessment completed"
                : "Your result hasn’t been released"}
            </h2>
            <p>
              {result.data
                ? "Your score and reviewer feedback are now available."
                : "Scores and reviewer feedback appear only after disclosure."}
            </p>
          </section>
          {result.data && (
            <>
              <div className={styles.detailGrid}>
                <section className={styles.panel}>
                  <p>Overall score</p>
                  <strong className={styles.clock}>
                    {result.data.score !== null &&
                    result.data.score !== undefined
                      ? `${result.data.score} / 100`
                      : "Not scored"}
                  </strong>
                </section>
                <section className={styles.panel}>
                  <p>Assessment outcome</p>
                  <h2>
                    {result.data.passed
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
                <p>{result.data.feedback || "No feedback was released."}</p>
                <p>
                  {result.data.feedbackIsAiGenerated
                    ? "AI-generated feedback · "
                    : ""}
                  Released {new Date(result.data.disclosedAt).toLocaleString()}
                </p>
              </section>
            </>
          )}
        </>
      ) : state === "availability" ? (
        <>
          <section className={styles.panel}>
            <StatusBadge>
              {instanceData?.status || "Availability unavailable"}
            </StatusBadge>
            <h2>Technical skills assessment</h2>
            <p>
              {instanceData?.deadline
                ? `Available until ${new Date(instanceData.deadline).toLocaleString()}`
                : "Assessment availability will appear here when it is released."}
            </p>
          </section>
          <div className={styles.detailGrid}>
            <section className={styles.panel}>
              <h2>Continue where you left off</h2>
              <p>
                Your saved answers can be restored while the attempt is active.
              </p>
              <button
                type="button"
                disabled={
                  busy || !instanceData || instanceData.status !== "in_progress"
                }
                onClick={() => void openAttempt()}
              >
                {busy ? "Opening…" : "Resume assessment"}
              </button>
            </section>
            <section className={styles.panel}>
              <h2>Start a new attempt</h2>
              <p>
                Check your connection, available time, and the completion
                deadline before starting.
              </p>
              <button
                type="button"
                className={s.primary}
                disabled={
                  busy ||
                  (!!instanceData && instanceData.status !== "available")
                }
                onClick={() => void openAttempt()}
              >
                {busy ? "Opening…" : "Start assessment"}
              </button>
            </section>
          </div>
          <section className={styles.panel}>
            <h2>Countdown and auto-submit</h2>
            <p>
              Starting begins the active timer. The server enforces the
              completion deadline and locks closed attempts. Opening this
              overview does not start an attempt.
            </p>
          </section>
        </>
      ) : state === "submitted" || locked ? (
        <div className={styles.detailGrid}>
          <div>
            <section className={styles.panel}>
              <h2>
                {locked
                  ? "Your answers are saved and locked"
                  : "Submission status unavailable"}
              </h2>
              <StatusBadge>
                {attemptData?.status || "Not confirmed"}
              </StatusBadge>
              <p>
                {locked
                  ? "This attempt is closed. Its saved responses can no longer be edited."
                  : "A submission can only be confirmed by the assessment service."}
              </p>
              {attemptData?.submittedAt && (
                <p>
                  Submitted {new Date(attemptData.submittedAt).toLocaleString()}
                </p>
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
              Assessment evaluation and stage progression are separate. A
              pending result is not a rejection or hiring decision.
            </p>
          </aside>
        </div>
      ) : variant !== "questions" ? (
        <AdvancedAssessment
          questions={questions}
          attempt={attemptData}
          variant={variant}
          answer={answer}
          onAnswer={(id, value) => {
            dirtyQuestions.current.add(id);
            setDirty(true);
            setAnswers((previous) => ({ ...previous, [id]: value }));
          }}
          onSubmit={() => setDialog("submit")}
          prefix={prefix}
          onSave={saveProgress}
          busy={busy}
          canEdit={canEdit}
          onLanguageChange={(id, language) => {
            answerLanguages.current[id] = language;
            dirtyQuestions.current.add(id);
            setDirty(true);
          }}
          savedLanguage={(question) => {
            const response = savedAnswers?.find(
              (row) => row.questionId === question.id,
            )?.response;
            return response &&
              !Array.isArray(response) &&
              typeof response === "object" &&
              "language" in response
              ? response.language
              : "";
          }}
        />
      ) : (
        <div className={styles.assessmentGrid}>
          <aside className={styles.panel}>
            <h2>Questions</h2>
            <p>{questions.length} questions</p>
            <div className={styles.problemList}>
              {questions.map((question, index) => (
                <a
                  key={question.id}
                  className={s.button}
                  href={`#${prefix}-${question.id}`}
                >
                  {index + 1} · {question.responseFormat.replaceAll("_", " ")}
                </a>
              ))}
            </div>
            <AttemptClock
              deadline={attemptData?.deadline}
              active={attemptData?.status === "in_progress"}
            />
          </aside>
          <div>
            {questions.length === 0 ? (
              <section className={styles.panel}>
                <h2>Questions aren’t available yet</h2>
                <p>
                  Your assigned questions will appear here when the assessment
                  service is available.
                </p>
              </section>
            ) : (
              questions.map((question, index) => (
                <section
                  className={`${styles.panel} ${styles.questionPanel}`}
                  id={`${prefix}-${question.id}`}
                  key={question.id}
                >
                  <p className={styles.questionMeta}>
                    Question {String(index + 1).padStart(2, "0")} ·{" "}
                    {question.responseFormat.replaceAll("_", " ")}
                  </p>
                  <h2>{question.prompt}</h2>
                  <QuestionInput
                    question={question}
                    value={answer(question)}
                    onChange={(value) => {
                      dirtyQuestions.current.add(question.id);
                      setDirty(true);
                      setAnswers((previous) => ({
                        ...previous,
                        [question.id]: value,
                      }));
                    }}
                    disabled={locked || busy || !canEdit}
                    questionId={`${prefix}-${question.id}`}
                  />
                </section>
              ))
            )}
            <div className={s.actions}>
              <button
                type="button"
                disabled={busy || !canEdit}
                onClick={() => void saveProgress()}
              >
                {busy ? "Saving…" : "Save progress"}
              </button>
              <button
                type="button"
                className={s.primary}
                disabled={busy || !canEdit}
                onClick={() => setDialog("submit")}
              >
                Submit assessment
              </button>
            </div>
            <p className={s.muted}>
              Save progress to store your edited answers. Submission saves
              current edits before locking the attempt.
            </p>
          </div>
        </div>
      )}
      <div className={s.actions}>
        <Link href="/candidate/applications" className={s.button}>
          Back to application
        </Link>
      </div>
      {variant !== "questions" && <WorkspaceFooter />}
      {dialog && (
        <ConfirmDialog
          title="Submit assessment?"
          description="Submitting sends your saved responses for evaluation and locks your answers. It is separate from a stage or hiring decision."
          confirmLabel="Submit assessment"
          intent="warning"
          unavailable={
            !canEdit ? "This attempt is not open for submission." : undefined
          }
          onConfirm={confirmSubmission}
          onClose={() => setDialog(null)}
        />
      )}
    </section>
  );
}
