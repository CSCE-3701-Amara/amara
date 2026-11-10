"use client";
import { workspaceStyles as s } from "@/shared/ui/workspace";
import type { CandidateQuestion } from "../types/product.types";
import styles from "./candidate.module.css";
import { VideoAnswer } from "./video-answer";

export function QuestionInput({
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
    return <VideoAnswer key={question.id} disabled={disabled} />;
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
