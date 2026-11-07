import { apiRequest, queryString, resourceId } from "@/shared/lib/api-client";
import type {
  FeedFilters,
  FeedPage,
  ApplicationPage,
  ApplicationDetail,
} from "../types/cand.types";
export const listJobFeed = (filter: FeedFilters = {}, signal?: AbortSignal) =>
  apiRequest<FeedPage>(`/v1/jobs/feed${queryString(filter)}`, { signal });
export const listApplications = (
  filter: { status?: string; offset?: number; limit?: number } = {},
  signal?: AbortSignal,
) =>
  apiRequest<ApplicationPage>(
    `/v1/candidates/me/applications${queryString(filter)}`,
    { signal },
  );
export const getApplication = (id: string, signal?: AbortSignal) =>
  apiRequest<ApplicationDetail>(`/v1/applications/${resourceId(id)}`, {
    signal,
  });
