export type EmploymentType =
  | "on-site full-time"
  | "on-site part-time"
  | "remote full-time"
  | "remote part-time";
export interface PostingInput {
  title: string;
  description?: string | null;
  language?: string | null;
  currency?: string | null;
  seniority_level: string;
  employment_type: EmploymentType;
  location?: string | null;
  salary_min?: number | null;
  salary_max?: number | null;
  salary_visible?: boolean | null;
  open_positions?: number | null;
  application_opening?: string | null;
  application_deadline?: string | null;
  max_applications?: number | null;
  requirements: {
    id?: string;
    text: string;
    requirement_type: "must_have" | "nice_to_have";
  }[];
  amarat: { id?: string; title: string; description?: string | null }[];
}
export interface Posting extends PostingInput {
  id: string;
  company_id: string;
  status: "draft" | "published" | "archived";
  created_at: string;
  updated_at: string;
}
export interface PostingList {
  items: Posting[] | null;
}
export interface CompanyPostingPage {
  items: {
    id: string;
    title: string;
    status: string;
    published_at: string | null;
  }[];
  next_cursor?: string;
}
export interface ApplicationSubmission {
  id: string;
  reference_number: number;
  posting_id: string;
  status: "admitted" | "screened_out";
  submitted_at: string;
  screened_out: boolean;
  failed_amarat: string[] | null;
}
