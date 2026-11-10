import { expect, test } from "@playwright/test";
import {
  FALLBACK_EXECUTION_LANGUAGES,
  findSavedResponse,
  responseLanguage,
  responseValue,
} from "../../src/features/cand/lib/assessment-responses";
import type { CandidateSubmission } from "../../src/features/cand/types/product.types";

const submission = (
  response: CandidateSubmission["response"],
): CandidateSubmission => ({
  id: "submission",
  attemptId: "attempt",
  questionId: "question",
  response,
  savedAt: "2026-11-09T00:00:00Z",
});

test("newly saved response wins over the attempt snapshot", () => {
  expect(
    findSavedResponse("question", [submission("new")], { question: "old" }),
  ).toBe("new");
});

test("null saved response falls back exactly like the existing nullish contract", () => {
  expect(
    findSavedResponse("question", [submission(null)], { question: "snapshot" }),
  ).toBe("snapshot");
  expect(findSavedResponse("missing", [], undefined)).toBeUndefined();
});

test("extracts editable values from text, selections, and code", () => {
  expect(responseValue("answer")).toBe("answer");
  expect(responseValue([1, 3])).toEqual([1, 3]);
  expect(responseValue({ language: "go", code: "package main" })).toBe(
    "package main",
  );
  expect(responseValue({ videoUrl: "https://example.test/video" })).toBe("");
  expect(responseValue(null)).toBe("");
});

test("extracts languages only from code responses", () => {
  expect(responseLanguage({ language: "python", code: "print(1)" })).toBe(
    "python",
  );
  expect(responseLanguage("python")).toBe("");
  expect(responseLanguage({ videoUrl: "https://example.test/video" })).toBe("");
  expect(FALLBACK_EXECUTION_LANGUAGES).toEqual([
    "javascript",
    "typescript",
    "python",
    "go",
    "cpp",
  ]);
});
