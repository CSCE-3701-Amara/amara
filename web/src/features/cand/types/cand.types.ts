import type { EmploymentType } from "@/features/job/types/job.types";
export interface FeedFilters {
  title?: string;
  language?: string;
  currency?: string;
  seniority_level?: string;
  employment_type?: EmploymentType | "";
  location?: string;
  salary_min?: number;
  salary_max?: number;
  limit?: number;
  cursor?: string;
}
export interface FeedJob {
  id: string;
  company_id: string;
  title: string;
  description: string | null;
  language: string | null;
  currency: string | null;
  seniority_level: string | null;
  employment_type: string | null;
  location: string | null;
  salary_min: number | null;
  salary_max: number | null;
  salary_visible: boolean | null;
  open_positions: number | null;
  application_opening: string | null;
  application_deadline: string | null;
  max_applications: number | null;
  status: string;
  published_at: string;
  application_url: string;
}
export interface FeedPage {
  items: FeedJob[];
  next_cursor?: string;
}
export interface ApplicationSummary {
  id: string;
  referenceNumber: number;
  postingId: string;
  jobTitle: string;
  companyName: string;
  createdAt: string | null;
  updatedAt: string | null;
  submittedAt: string | null;
  withdrawnAt: string | null;
  screenedOut: boolean | null;
  screeningCompletedAt: string | null;
  status: string;
  closed: boolean;
  identityHidden: boolean;
}
export interface ApplicationDetail extends ApplicationSummary {
  pipelineExecutionId: string | null;
  finalDecision: boolean | null;
  finalFeedback?: string;
  rejectionFeedback?: unknown;
}
export interface ApplicationPage {
  data: ApplicationSummary[];
  pagination: { limit: number; offset: number; total: number };
}
