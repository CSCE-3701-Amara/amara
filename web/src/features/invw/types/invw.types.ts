export interface Amara {
  name: string;
  crest: string;
  domain: string;
  description: string;
}
export interface Criterion {
  description: string;
  amara: string;
}
export interface Question {
  prompt: string;
  format: "code" | "multiple_choice" | "short_text" | "long_text" | "video";
  options?: string[];
  amara?: string;
  difficulty?: string;
  minimum_cases_passed?: number;
  test_cases?: { input: string; expected_output: string; weight: number }[];
}
export interface ProblemReference {
  module_id: string;
  position: number;
}
export interface ModuleDefinition {
  name: string;
  description?: string;
  type: "coding" | "programming_language" | "behavioral" | "recorded_video";
  behavioral?: boolean;
  minimum_tests_passed?: number;
  active_minutes?: number;
  threshold?: number;
  questions?: Question[];
  amarat?: Amara[];
  rubric?: Criterion[];
  library?: {
    problems?: ProblemReference[];
    difficulty_mix?: Record<string, number>;
    amara: string;
  };
}
export interface AssessmentModule extends ModuleDefinition {
  id: string;
  company_id: string;
  copied_from_module_id?: string;
  published: boolean;
}
export interface LibraryProblem extends ProblemReference {
  difficulty: string;
  prompt: string;
}
export interface PipelineDraft {
  id: string;
  postingId: string;
  versionId: string;
  versionNumber: number;
  published: boolean;
}
export interface PipelineStage {
  advancementRule?: "manual" | "completion" | "threshold";
  scoreThreshold?: number;
  belowThresholdAction?: "hold" | "reject";
  id?: string;
  name: string;
  type: "assessment" | "live_interview";
  scheduledStartAt: string;
  scheduledEndAt: string;
  moduleIds?: string[];
  interview?: {
    durationMinutes: number;
    interviewerIds: string[];
    threshold: number | null;
    rubric: Criterion[];
    amarat: Amara[];
  };
}
export type ModuleSummary = Pick<
  AssessmentModule,
  "id" | "name" | "description" | "type" | "published"
>;
