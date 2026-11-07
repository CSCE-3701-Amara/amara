export type AccountKind = "candidate" | "company";
export type AuthState =
  | "form"
  | "mfa"
  | "code-rejected"
  | "credentials-rejected"
  | "rate-limited"
  | "expired"
  | "verification-pending"
  | "review-pending";
export interface LoginRequest {
  email: string;
  password: string;
  mfaCode?: string;
}
export interface LoginResponse {
  user: {
    id: string;
    email: string;
    role: "candidate" | "interviewer" | "company_representative";
  };
  session: { expiresAt: string };
}
export interface RegistrationRequest {
  email: string;
  password: string;
  firstName?: string;
  lastName?: string;
  companyName?: string;
  registrationNumber?: string;
}
export interface RegistrationResponse {
  id: string;
  accountType: "candidate" | "company" | "interviewer";
  email: string;
  status: "pending_verification" | "pending_review" | "active";
  createdAt: string;
}
