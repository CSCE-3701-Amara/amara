"use client";
import { CodeRunner } from "./code-runner";
import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useId, useRef, useState } from "react";
import { toError } from "@/shared/lib/errors";
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
import {
  FALLBACK_EXECUTION_LANGUAGES,
  findSavedResponse,
  responseLanguage,
  responseValue,
} from "../lib/assessment-responses";
import {
  AssessmentAvailabilityView,
  AssessmentLockedView,
  AssessmentResultView,
} from "./assessment-state-views";

function AttemptClock({
  deadline,
  active,
}: {
  deadline?: string;
  active: boolean;
}) {
  const [now, setNow] = useState(0);
  useEffect(() => {
    if (!active) return;
    const firstTick = setTimeout(() => setNow(Date.now()), 0);
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => {
      clearTimeout(firstTick);
      clearInterval(timer);
    };
  }, [active]);
  const milliseconds =
    deadline && now ? Math.max(0, Date.parse(deadline) - now) : undefined;
  return (
    <section>
      <p>Remaining active time</p>
      <strong className={styles.clock}>
        {milliseconds !== undefined && Number.isFinite(milliseconds)
          ? `${Math.floor(milliseconds / 60000)}:${String(Math.floor(milliseconds / 1000) % 60).padStart(2, "0")}`
          : "Time unavailable"}
      </strong>
      {deadline && (
        <p>Completion deadline · {new Date(deadline).toLocaleString()}</p>
      )}
    </section>
  );
}
function VideoAnswer({ disabled }: { disabled: boolean }) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const recorderRef = useRef<MediaRecorder | null>(null);
  const recordingUrl = useRef<string | null>(null);
  const [recording, setRecording] = useState(false);
  const [preview, setPreview] = useState("");
  const [error, setError] = useState("");
  useEffect(
    () => () => {
      if (recorderRef.current?.state === "recording")
        recorderRef.current.stop();
      streamRef.current?.getTracks().forEach((track) => track.stop());
      if (recordingUrl.current) URL.revokeObjectURL(recordingUrl.current);
    },
    [],
  );
  async function record() {
    setError("");
    try {
      if (
        !navigator.mediaDevices?.getUserMedia ||
        typeof MediaRecorder === "undefined"
      )
        throw new Error(
          "Recording isn’t supported in this browser. Use a supported browser with camera access.",
        );
      const stream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true,
      });
      streamRef.current = stream;
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        videoRef.current.muted = true;
        await videoRef.current.play();
      }
      const recorder = new MediaRecorder(stream);
      const chunks: BlobPart[] = [];
      recorder.ondataavailable = (event) => {
        if (event.data.size) chunks.push(event.data);
      };
      recorder.onstop = () => {
        if (recordingUrl.current) URL.revokeObjectURL(recordingUrl.current);
        const url = URL.createObjectURL(
          new Blob(chunks, { type: recorder.mimeType || "video/webm" }),
        );
        recordingUrl.current = url;
        setPreview(url);
        setRecording(false);
        stream.getTracks().forEach((track) => track.stop());
        if (videoRef.current) {
          videoRef.current.srcObject = null;
          videoRef.current.muted = false;
        }
      };
      recorderRef.current = recorder;
      setPreview("");
      recorder.start();
      setRecording(true);
    } catch (cause) {
      streamRef.current?.getTracks().forEach((track) => track.stop());
      setError(
        cause instanceof Error
          ? cause.message
          : "Camera or microphone access was denied. Update your browser permissions and try again.",
      );
    }
  }
  return (
    <>
      <div className={styles.videoPreview}>
        {!recording && !preview && (
          <Image src="/figma/video.svg" width={40.32} height={40.32} alt="" />
        )}
        <video
          ref={videoRef}
          src={preview || undefined}
          controls={!!preview}
          playsInline
          hidden={!recording && !preview}
          aria-label="Your recorded video response"
        />
      </div>
      <div className={s.actions}>
        <button
          type="button"
          className={s.primary}
          disabled={disabled || recording}
          onClick={record}
        >
          {preview ? "Retake" : "Record"}
        </button>
        <button
          type="button"
          disabled={!recording}
          onClick={() => recorderRef.current?.stop()}
        >
          Stop
        </button>
        <button type="button" disabled aria-describedby="video-upload-note">
          Save recording
        </button>
      </div>
      {error && (
        <p role="alert" className={s.error}>
          {error}
        </p>
      )}
      <p>
        Changing the prompt or leaving this screen discards unsaved recordings.
      </p>
      <p id="video-upload-note">
        Recording previews stay on this screen. Uploading and saving recordings
        aren’t available yet.
      </p>
    </>
  );
}
function QuestionInput({
  question,
  value,
  onChange,
  disabled,
  questionId,
}: {
  question: CandidateQuestion;
  value: string | number[];
  onChange: (value: string | number[]) => void;
  disabled: boolean;
  questionId: string;
}) {
  if (question.responseFormat === "video")
    return <VideoAnswer disabled={disabled} />;
  if (question.responseFormat === "multiple_choice")
    return (
      <fieldset className={styles.choiceGroup}>
        <legend className={styles.srOnly}>Your answer</legend>
        {question.options?.map((option, index) => (
          <label className={styles.choice} key={index}>
            <input
              type="radio"
              name={questionId}
              checked={Array.isArray(value) && value.includes(index)}
              onChange={() => onChange([index])}
              disabled={disabled}
            />
            <span>
              {String.fromCharCode(65 + index)}. {option}
            </span>
          </label>
        ))}
      </fieldset>
    );
  return (
    <div className={s.field}>
      <label htmlFor={`${questionId}-answer`}>Your answer</label>
      <textarea
        id={`${questionId}-answer`}
        className={
          question.responseFormat === "code" ? styles.codeArea : undefined
        }
        name={questionId}
        autoComplete="off"
        spellCheck={question.responseFormat !== "code"}
        value={typeof value === "string" ? value : ""}
        onChange={(event) => onChange(event.target.value)}
        maxLength={question.responseFormat === "code" ? undefined : 2000}
        rows={
          question.responseFormat === "short_text"
            ? 2
            : question.responseFormat === "code"
              ? 14
              : 5
        }
        disabled={disabled}
        placeholder={
          question.responseFormat === "code"
            ? "Write your solution here…"
            : question.responseFormat === "short_text"
              ? "Write a concise answer…"
              : "Explain your approach, investigation steps and trade-offs…"
        }
      />
      {question.responseFormat === "long_text" && (
        <small>
          {typeof value === "string" ? value.length : 0} / 2,000 characters
        </small>
      )}
    </div>
  );
}
function AdvancedAssessment({
  questions,
  attempt,
  variant,
  answer,
  onAnswer,
  onSubmit,
  prefix,
  savedLanguage,
  onSave,
  busy,
  canEdit,
  onLanguageChange,
}: {
  questions: CandidateQuestion[];
  attempt?: AssessmentAttempt;
  variant: "coding" | "video";
  answer: (question: CandidateQuestion) => string | number[];
  onAnswer: (id: string, value: string | number[]) => void;
  onSubmit: () => void;
  prefix: string;
  savedLanguage: (question: CandidateQuestion) => string;
  onSave: (question: CandidateQuestion, language: string) => Promise<void>;
  busy: boolean;
  canEdit: boolean;
  onLanguageChange: (id: string, language: string) => void;
}) {
  const [index, setIndex] = useState(0);
  const [languages, setLanguages] = useState<Record<string, string>>({});
  const availableLanguages = useApiResource<string[]>(
    "/v1/execution-languages",
  );
  const defaultLanguage = availableLanguages.data?.[0] || "typescript";
  const question = questions[index];
  return (
    <div className={styles.advancedAssessment}>
      <aside>
        <section className={styles.panel}>
          <h2>
            {variant === "coding" ? "Problems" : "Your prompts"} ·{" "}
            {questions.length} {variant === "coding" ? "questions" : "prompts"}
          </h2>
          <div className={styles.problemList}>
            {questions.map((item, i) => (
              <button
                type="button"
                key={item.id}
                aria-pressed={index === i}
                onClick={() => setIndex(i)}
              >
                {i + 1} · {item.prompt}
              </button>
            ))}
          </div>
          <p>
            Your assigned questions stay the same across opens of this attempt.
          </p>
        </section>
        <section className={styles.panel}>
          <h2>Attempt timing</h2>
          <AttemptClock
            deadline={attempt?.deadline}
            active={attempt?.status === "in_progress"}
          />
          <p>
            The deadline is controlled by the assessment service. Local drafts
            are not automatically submitted.
          </p>
        </section>
      </aside>
      <section className={styles.panel}>
        <h2>
          {variant === "coding" ? "Problem" : "Prompt"} {index + 1}
        </h2>
        {question ? (
          <>
            <StatusBadge>
              {variant === "coding" ? "Code response" : "Video response"}
            </StatusBadge>
            <p className={styles.publicPrompt}>{question.prompt}</p>
            <h3>Requirements</h3>
            <p>
              Follow the instructions in this public prompt. Private evaluation
              details are not available here.
            </p>
            <div className={s.actions}>
              <button
                type="button"
                disabled={index === 0}
                onClick={() => setIndex((i) => i - 1)}
              >
                Previous
              </button>
              <button
                type="button"
                disabled={index >= questions.length - 1}
                onClick={() => setIndex((i) => i + 1)}
              >
                Next {variant === "coding" ? "problem" : "prompt"}
              </button>
            </div>
          </>
        ) : (
          <p>
            Your assigned {variant === "coding" ? "problems" : "prompts"} aren’t
            available yet.
          </p>
        )}
      </section>
      <div>
        <section className={styles.panel}>
          <h2>Your answer</h2>
          {question && (
            <>
              {variant === "coding" && (
                <label className={s.field}>
                  Language
                  <select
                    disabled={busy || !canEdit}
                    value={
                      languages[question.id] ||
                      savedLanguage(question) ||
                      defaultLanguage
                    }
                    onChange={(event) => {
                      onLanguageChange(question.id, event.target.value);
                      setLanguages((values) => ({
                        ...values,
                        [question.id]: event.target.value,
                      }));
                    }}
                  >
                    {(
                      availableLanguages.data ?? FALLBACK_EXECUTION_LANGUAGES
                    ).map((language) => (
                      <option key={language}>{language}</option>
                    ))}
                  </select>
                </label>
              )}
              {variant === "coding" && attempt && (
                <CodeRunner
                  attemptId={attempt.id}
                  questionId={question.id}
                  source={
                    typeof answer(question) === "string"
                      ? (answer(question) as string)
                      : ""
                  }
                  language={
                    languages[question.id] ||
                    savedLanguage(question) ||
                    defaultLanguage
                  }
                  disabled={busy || !canEdit}
                />
              )}
              <QuestionInput
                question={question}
                value={answer(question)}
                onChange={(value) => onAnswer(question.id, value)}
                disabled={busy || !canEdit}
                questionId={`${prefix}-${question.id}`}
              />
            </>
          )}
          <p className={s.muted}>
            Save your answer before leaving this assessment.
          </p>
          <button
            type="button"
            disabled={busy || !canEdit || !question || variant === "video"}
            onClick={() =>
              question &&
              void onSave(
                question,
                languages[question.id] ||
                  savedLanguage(question) ||
                  defaultLanguage,
              )
            }
          >
            {busy
              ? "Saving…"
              : `Save ${variant === "coding" ? "answer" : "recording"}`}
          </button>
        </section>
        <section className={styles.panel}>
          <h2>Submit assessment?</h2>
          <p>
            Submission applies to this instance only. Saved responses would be
            locked for evaluation. This is separate from a hiring decision.
          </p>
          <button
            type="button"
            className={s.primary}
            disabled={busy || !canEdit}
            onClick={onSubmit}
          >
            Review submission
          </button>
        </section>
      </div>
    </div>
  );
}
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
  const [runLanguages, setRunLanguages] = useState<Record<string, string>>({});
  const availableLanguages = useApiResource<string[]>(
    "/v1/execution-languages",
  );
  const defaultLanguage = availableLanguages.data?.[0] || "typescript";
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
    return responseValue(
      findSavedResponse(question.id, savedAnswers, attemptData?.savedResponses),
    );
  }
  function savedLanguage(question: CandidateQuestion): string {
    return responseLanguage(
      findSavedResponse(question.id, savedAnswers, attemptData?.savedResponses),
    );
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
                    savedLanguage(question) ||
                    defaultLanguage,
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
        <AssessmentResultView result={result.data} />
      ) : state === "availability" ? (
        <AssessmentAvailabilityView
          instance={instanceData}
          busy={busy}
          onOpen={openAttempt}
        />
      ) : state === "submitted" || locked ? (
        <AssessmentLockedView
          locked={locked}
          attempt={attemptData}
          questions={questions}
        />
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
          savedLanguage={savedLanguage}
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
                  {question.responseFormat === "code" && attemptData && (
                    <>
                      <label className={s.field}>
                        Language
                        <select
                          disabled={busy || !canEdit}
                          value={
                            runLanguages[question.id] ||
                            savedLanguage(question) ||
                            defaultLanguage
                          }
                          onChange={(e) => {
                            answerLanguages.current[question.id] =
                              e.target.value;
                            setRunLanguages((previous) => ({
                              ...previous,
                              [question.id]: e.target.value,
                            }));
                            dirtyQuestions.current.add(question.id);
                            setDirty(true);
                          }}
                        >
                          {(
                            availableLanguages.data ??
                            FALLBACK_EXECUTION_LANGUAGES
                          ).map((l) => (
                            <option key={l}>{l}</option>
                          ))}
                        </select>
                      </label>
                      <CodeRunner
                        attemptId={attemptId}
                        questionId={question.id}
                        source={
                          typeof answer(question) === "string"
                            ? (answer(question) as string)
                            : ""
                        }
                        language={
                          runLanguages[question.id] ||
                          savedLanguage(question) ||
                          defaultLanguage
                        }
                        disabled={busy || !canEdit}
                      />
                    </>
                  )}
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
