"use client";
import { useState } from "react";
import { workspaceStyles as s } from "@/shared/ui/workspace";
import { StatusBadge } from "@/shared/ui/product-ui";
import type {
  AssessmentAttempt,
  CandidateQuestion,
} from "../types/product.types";
import styles from "./candidate.module.css";
import { AttemptClock } from "./attempt-clock";
import { QuestionInput } from "./assessment-input";

export function AdvancedAssessment({
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
                      "javascript"
                    }
                    onChange={(event) => {
                      onLanguageChange(question.id, event.target.value);
                      setLanguages((values) => ({
                        ...values,
                        [question.id]: event.target.value,
                      }));
                    }}
                  >
                    {["javascript", "typescript", "python", "go", "cpp"].map(
                      (language) => (
                        <option key={language}>{language}</option>
                      ),
                    )}
                  </select>
                </label>
              )}
              <QuestionInput
                key={question.id}
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
                  "javascript",
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
