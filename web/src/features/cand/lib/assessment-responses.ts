import type { CandidateSubmission } from "../types/product.types";

type CandidateResponse = CandidateSubmission["response"];

export const FALLBACK_EXECUTION_LANGUAGES = [
  "javascript",
  "typescript",
  "python",
  "go",
  "cpp",
] as const;

export function findSavedResponse(
  questionId: string,
  savedAnswers: CandidateSubmission[],
  savedResponses?: Record<string, CandidateResponse>,
): CandidateResponse {
  return (
    savedAnswers.find((answer) => answer.questionId === questionId)?.response ??
    savedResponses?.[questionId]
  );
}

export function responseValue(response: CandidateResponse): string | number[] {
  if (typeof response === "string" || Array.isArray(response)) return response;
  if (response && "code" in response) return response.code;
  return "";
}

export function responseLanguage(response: CandidateResponse): string {
  return response &&
    !Array.isArray(response) &&
    typeof response === "object" &&
    "language" in response
    ? response.language
    : "";
}
