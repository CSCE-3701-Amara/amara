import { apiRequest } from "@/shared/lib/api-client";
import type {
  AccountKind,
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
export function deleteCurrentSession() {
  return apiRequest<void>("/v1/auth/sessions/current", { method: "DELETE" });
}
export function performAccountAction(action: "deactivate" | "delete") {
  return apiRequest<void>("/v1/auth/account-actions", {
    method: "POST",
    body: { action },
  });
}
