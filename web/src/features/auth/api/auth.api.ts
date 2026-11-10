import { apiRequest, ApiError } from "@/shared/lib/api-client";
import type {
  AccountKind,
  CurrentSessionResponse,
  LoginRequest,
  LoginResponse,
  RegistrationRequest,
  RegistrationResponse,
} from "../types/auth.types";
export { ApiError as AuthError } from "@/shared/lib/api-client";

export function createSession(input: LoginRequest, signal?: AbortSignal) {
  return apiRequest<LoginResponse>("/v1/auth/sessions", {
    method: "POST",
    body: input,
    signal,
  });
}
export function registerAccount(
  kind: AccountKind,
  input: RegistrationRequest,
  signal?: AbortSignal,
) {
  return apiRequest<RegistrationResponse>(
    kind === "company"
      ? "/v1/auth/company-registrations"
      : "/v1/auth/registrations",
    {
      method: "POST",
      body: kind === "company" ? input : { ...input, accountType: "candidate" },
      signal,
    },
  );
}
export function verifyAccount(
  input: { email: string; code: string },
  signal?: AbortSignal,
) {
  return apiRequest<RegistrationResponse>("/v1/auth/verifications", {
    method: "POST",
    body: input,
    signal,
  });
}
export function deleteCurrentSession() {
  return apiRequest<void>("/v1/auth/sessions/current", { method: "DELETE" });
}
export function performAccountAction(action: "deactivate" | "delete") {
  return apiRequest<void>("/v1/auth/account-actions", {
    method: "POST",
    body: { action },
  });
}

/** A revoked cookie means guest; service failures must not masquerade as sign-out. */
export async function getCurrentSession(
  signal?: AbortSignal,
): Promise<CurrentSessionResponse | null> {
  let current: CurrentSessionResponse;
  try {
    current = await apiRequest<CurrentSessionResponse>(
      "/v1/auth/sessions/current",
      { signal },
    );
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) return null;
    throw error;
  }
  const expiry = Date.parse(current?.session?.expiresAt ?? "");
  if (
    !["candidate", "company_representative", "interviewer"].includes(current?.user?.role) ||
    typeof current?.user?.id !== "string" ||
    !current.user.id ||
    !Number.isFinite(expiry)
  ) {
    throw new ApiError(
      "The session service returned an unexpected response. Please try again.",
      "invalid_response",
      200,
    );
  }
  return expiry > Date.now() ? current : null;
}
