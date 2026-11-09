import { apiRequest, resourceId } from "@/shared/lib/api-client";
import type {
  ModuleDefinition,
  AssessmentModule,
  ModuleSummary,
  LibraryProblem,
  PipelineDraft,
  PipelineStage,
} from "../types/invw.types";
export const listModules = (signal?: AbortSignal) =>
  apiRequest<ModuleSummary[] | null>("/v1/modules", { signal });
export const listLibraryModules = (signal?: AbortSignal) =>
  apiRequest<ModuleSummary[] | null>("/v1/modules/library", { signal });
export const listLibraryProblems = (signal?: AbortSignal) =>
  apiRequest<LibraryProblem[] | null>("/v1/modules/library/problems", {
    signal,
  });
export const getModule = (id: string, signal?: AbortSignal) =>
  apiRequest<AssessmentModule>(`/v1/modules/${resourceId(id)}`, { signal });
export const createModule = (body: ModuleDefinition) =>
  apiRequest<AssessmentModule>("/v1/modules", { method: "POST", body });
export const updateModule = (id: string, body: ModuleDefinition) =>
  apiRequest<AssessmentModule>(`/v1/modules/${resourceId(id)}`, {
    method: "PUT",
    body,
  });
export const publishModule = (id: string) =>
  apiRequest<AssessmentModule>(`/v1/modules/${resourceId(id)}/publish`, {
    method: "POST",
  });
export const copyModule = (id: string) =>
  apiRequest<AssessmentModule>(`/v1/modules/${resourceId(id)}/copy`, {
    method: "POST",
  });
export const getPipeline = (postingId: string, signal?: AbortSignal) =>
  apiRequest<PipelineDraft>(`/v1/jobs/${resourceId(postingId)}/pipeline`, {
    signal,
  });
export const createPipeline = (postingId: string) =>
  apiRequest<PipelineDraft>(`/v1/jobs/${resourceId(postingId)}/pipeline`, {
    method: "POST",
  });
export const publishPipeline = (versionId: string) =>
  apiRequest<void>(`/v1/pipeline-versions/${resourceId(versionId)}/publish`, {
    method: "POST",
  });
export const listPipelineStages = (versionId: string, signal?: AbortSignal) =>
  apiRequest<PipelineStage[] | null>(
    `/v1/pipeline-versions/${resourceId(versionId)}/stages`,
    { signal },
  );
export const replacePipelineStages = (
  versionId: string,
  stages: PipelineStage[],
) =>
  apiRequest<PipelineStage[]>(
    `/v1/pipeline-versions/${resourceId(versionId)}/stages`,
    { method: "PUT", body: stages },
  );
