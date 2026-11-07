import { apiRequest, queryString, resourceId } from "@/shared/lib/api-client";
import type {
  PostingInput,
  Posting,
  PostingList,
  CompanyPostingPage,
  ApplicationSubmission,
} from "../types/job.types";
export const listPostings = (signal?: AbortSignal) =>
  apiRequest<PostingList>("/v1/jobs", { signal });
export const listCompanyPostings = (
  companyId: string,
  filter: {
    title?: string;
    status?: string;
    cursor?: string;
    limit?: number;
  } = {},
  signal?: AbortSignal,
) =>
  apiRequest<CompanyPostingPage>(
    `/v1/company/${resourceId(companyId)}/postings${queryString(filter)}`,
    { signal },
  );
export const getPosting = (id: string, signal?: AbortSignal) =>
  apiRequest<Posting>(`/v1/jobs/${resourceId(id)}`, { signal });
export const createPosting = (body: PostingInput) =>
  apiRequest<Posting>("/v1/jobs", { method: "POST", body });
export const updatePosting = (id: string, body: PostingInput) =>
  apiRequest<Posting>(`/v1/jobs/${resourceId(id)}`, { method: "PUT", body });
export const transitionPosting = (
  id: string,
  action: "publish" | "republish" | "archive",
) =>
  apiRequest<Posting>(`/v1/jobs/${resourceId(id)}/${action}`, {
    method: "POST",
  });
export const deletePosting = (id: string) =>
  apiRequest<void>(`/v1/jobs/${resourceId(id)}`, { method: "DELETE" });
export const applyToPosting = (id: string) =>
  apiRequest<ApplicationSubmission>(`/v1/jobs/${resourceId(id)}/apply`, {
    method: "POST",
  });
