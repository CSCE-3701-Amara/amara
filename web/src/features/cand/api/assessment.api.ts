import { apiRequest, resourceId } from "@/shared/lib/api-client";
import type {
  AssessmentAttempt,
  CandidateSubmission,
} from "../types/product.types";

/** Registered assessment handlers return the attempt directly, without an envelope. */
export const startAssessment = (instanceId: string) =>
  apiRequest<AssessmentAttempt>(
    `/v1/assessment-instances/${resourceId(instanceId)}/start`,
    { method: "POST" },
  );
export const saveAssessmentAnswer = (
  attemptId: string,
  questionId: string,
  response: CandidateSubmission["response"],
) =>
  apiRequest<CandidateSubmission>(
    `/v1/assessment-attempts/${resourceId(attemptId)}/submissions/${resourceId(questionId)}`,
    { method: "PUT", body: { response } },
  );
export const submitAssessmentAttempt = (attemptId: string) =>
  apiRequest<AssessmentAttempt>(
    `/v1/assessment-attempts/${resourceId(attemptId)}/submission`,
    { method: "POST" },
  );
