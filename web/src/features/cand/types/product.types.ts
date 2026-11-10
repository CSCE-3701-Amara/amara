/** Planned contract projections. Writes remain unavailable until their Go handlers exist. */
export interface Offer {
  id: string;
  applicationId: string;
  jobTitle: string;
  compensationSummary: string;
  proposedStartDate: string;
  responseDeadline: string;
  message?: string | null;
  status: "pending" | "accepted" | "declined" | "expired" | "withdrawn";
  createdAt: string;
  response?: "accepted" | "declined" | null;
  respondedAt?: string | null;
}
export interface OfferPage {
  data: Offer[];
  pagination: { limit: number; offset: number; total: number };
}
export interface CandidateQuestion {
  id: string;
  prompt: string;
  responseFormat:
    "multiple_choice" | "short_text" | "long_text" | "code" | "video";
  options?: string[] | null;
  position: number;
}
export interface AssessmentAttempt {
  savedResponses?: Record<string, CandidateSubmission["response"]>;
  id: string;
  instanceId: string;
  stageExecutionId: string;
  status: string;
  deadline: string;
  startedAt?: string | null;
  submittedAt?: string | null;
  autoSubmitted: boolean;
  questions: CandidateQuestion[];
}
export interface AssessmentInstance {
  id: string;
  applicationId: string;
  stageId: string;
  status: string;
  startedAt?: string | null;
  deadline?: string | null;
}
export interface Interview {
  id: string;
  stageExecutionId: string;
  scheduledStartAt: string;
  scheduledEndAt: string;
  status: string;
  rescheduleCount: number;
  joinUrl?: string | null;
  provider?: string | null;
}
export interface TimeSlot {
  startAt: string;
  endAt: string;
}
export interface CandidateSubmission {
  id: string;
  attemptId: string;
  questionId: string;
  response?:
    | string
    | number[]
    | { language: string; code: string }
    | { videoUrl: string }
    | null;
  savedAt: string;
}
export interface AssessmentResult {
  id: string;
  instanceId: string;
  attemptId?: string | null;
  passed: boolean;
  score?: number | null;
  feedback?: string | null;
  feedbackIsAiGenerated: boolean;
  disclosedAt: string;
  scoredAt?: string | null;
}
