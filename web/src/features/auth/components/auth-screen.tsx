"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useId, useRef, useState, type FormEvent } from "react";
import { AuthError, createSession, registerAccount } from "../api/auth.api";
import type { AccountKind, AuthState } from "../types/auth.types";
import { AuthShell } from "./auth-shell";
import { FormField } from "./form-field";
import styles from "./auth.module.css";

const states: Record<
  Exclude<AuthState, "form">,
  { title: string; description: string }
> = {
  mfa: {
    title: "Verify your sign-in",
    description: "Enter the multi-factor code for your account.",
  },
  "code-rejected": {
    title: "Try your code again",
    description:
      "That code couldn’t be verified. Check the code and try again.",
  },
  "credentials-rejected": {
    title: "Check your sign-in details",
    description: "The email or password didn’t match. Try again.",
  },
  "rate-limited": {
    title: "Too many sign-in attempts",
    description: "Wait before trying to sign in again.",
  },
  expired: {
    title: "Sign in again",
    description: "Your session ended. Sign in to continue.",
  },
  "verification-pending": {
    title: "Verification pending",
    description:
      "Your registration was received. Your account is awaiting verification.",
  },
  "review-pending": {
    title: "Registration under review",
    description:
      "Your company registration was received. Sign-in is available after approval.",
  },
};

/** Reads only allowlisted URL state; passwords stay in memory throughout MFA. */
export function AuthRoute({
  registration = false,
}: {
  registration?: boolean;
}) {
  const params = useSearchParams();
  const kind = params.get("account") === "company" ? "company" : "candidate";
  const requestedState = params.get("state") || "form";
  const initialState: AuthState = Object.hasOwn(states, requestedState)
    ? (requestedState as AuthState)
    : "form";
  return (
    <AuthScreen
      key={`${initialState}:${registration}`}
      kind={kind}
      registration={registration}
      initialState={initialState}
    />
  );
}

function AuthScreen({
  kind,
  registration,
  initialState,
}: {
  kind: AccountKind;
  registration: boolean;
  initialState: AuthState;
}) {
  const stateTitleId = useId();
  const company = kind === "company";
  const [state, setState] = useState<AuthState>(initialState);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const [retryAt, setRetryAt] = useState(0);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const previousState = useRef(initialState);
  const requestRef = useRef<AbortController | null>(null);
  const inFlight = useRef(false);
  const screenRef = useRef<HTMLElement>(null);
  const formRef = useRef<HTMLFormElement>(null);
  const router = useRouter();
  const isMfa = state === "mfa" || state === "code-rejected";
  const loginHref = company ? "/login?account=company" : "/login";
  const registerHref = company ? "/register?account=company" : "/register";

  useEffect(() => {
    return () => {
      requestRef.current?.abort();
      inFlight.current = false;
    };
  }, [kind]);

  useEffect(() => {
    if (previousState.current === state) return;
    previousState.current = state;
    screenRef.current
      ?.querySelector<HTMLElement>(
        state === "form"
          ? '[name="email"]'
          : isMfa
            ? '[name="mfaCode"]'
            : "[data-state-title]",
      )
      ?.focus();
  }, [state, isMfa]);

  function backToForm() {
    setState("form");
    setError("");
    setCode("");
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (inFlight.current) return;
    setError("");
    setFieldErrors({});
    if (Date.now() < retryAt) {
      setState("rate-limited");
      return;
    }
    if (!email.trim() || !password.trim()) {
      const id = !email.trim() ? "email" : "password";
      setFieldErrors({ [id]: `Enter your ${id} to continue.` });
      screenRef.current?.querySelector<HTMLElement>(`[name="${id}"]`)?.focus();
      if (isMfa) setState("form");
      return;
    }
    const data = new FormData(event.currentTarget);
    const request = new AbortController();
    requestRef.current = request;
    inFlight.current = true;
    setBusy(true);
    try {
      if (registration) {
        const name = company
          ? {
              companyName: String(data.get("companyName") || "").trim(),
              registrationNumber: String(
                data.get("registrationNumber") || "",
              ).trim(),
            }
          : {
              firstName: String(data.get("firstName") || "").trim(),
              lastName: String(data.get("lastName") || "").trim(),
            };
        if (Object.values(name).some((value) => !value)) {
          const errors = Object.fromEntries(
            Object.entries(name)
              .filter(([, value]) => !value)
              .map(([key]) => [
                key,
                "Enter a value; spaces alone aren’t accepted.",
              ]),
          );
          setFieldErrors(errors);
          screenRef.current
            ?.querySelector<HTMLElement>(`[name="${Object.keys(errors)[0]}"]`)
            ?.focus();
          return;
        }
        await registerAccount(
          kind,
          { email: email.trim(), password, ...name },
          request.signal,
        );
        if (request.signal.aborted) return;
        setPassword("");
        setState(company ? "review-pending" : "verification-pending");
      } else {
        const result = await createSession(
          {
            email: email.trim(),
            password,
            ...(isMfa ? { mfaCode: code.trim() } : {}),
          },
          request.signal,
        );
        if (request.signal.aborted) return;
        if (
          !["candidate", "company_representative", "interviewer"].includes(
            result.user?.role,
          ) ||
          !Number.isFinite(Date.parse(result.session?.expiresAt))
        )
          throw new Error("Unexpected session response");
        setPassword("");
        router.replace(
          result.user.role === "company_representative"
            ? "/employer/jobs"
            : result.user.role === "candidate"
              ? "/candidate/dashboard"
              : "/",
        );
        router.refresh();
      }
    } catch (cause) {
      if (request.signal.aborted) return;
      if (cause instanceof AuthError) {
        if (cause.code === "mfa_required") setState("mfa");
        else if (cause.code === "invalid_mfa_code") {
          setCode("");
          setState("code-rejected");
          screenRef.current
            ?.querySelector<HTMLElement>('[name="mfaCode"]')
            ?.focus();
        } else if (cause.code === "invalid_credentials")
          setState("credentials-rejected");
        else if (cause.status === 429) {
          const seconds = Number(cause.retryAfter);
          setRetryAt(
            Number.isFinite(seconds) && seconds > 0
              ? Date.now() + seconds * 1000
              : Date.parse(cause.retryAfter || "") || Date.now() + 60000,
          );
          setState("rate-limited");
        } else setError(cause.message);
      } else
        setError(
          "We couldn’t reach amara. Check your connection and try again.",
        );
    } finally {
      if (requestRef.current === request) {
        inFlight.current = false;
        setBusy(false);
      }
    }
  }

  return (
    <AuthShell
      contentRef={screenRef}
      company={company}
      registration={registration && state === "form"}
    >
      {state === "form" ? (
        <>
          <header className={styles.heading} key={kind}>
            <h1 className={registration ? styles.registrationTitle : undefined}>
              {registration
                ? `Create a ${company ? "company" : "candidate"} account`
                : company
                  ? "Log in to your company"
                  : "Log in to your account"}
            </h1>
            <p>
              {registration
                ? company
                  ? "Set up your company workspace to post roles and assess talent."
                  : "Create an account to discover roles and show what you can do."
                : company
                  ? "Access your hiring workspace, postings, and assessments."
                  : "Pick up where you left off with jobs, applications, and assessments."}
            </p>
          </header>
          <form
            ref={formRef}
            onSubmit={submit}
            className={styles.card}
            aria-busy={busy}
          >
            <nav
              className={styles.selector}
              aria-label="Account type"
              data-account={kind}
            >
              <Link
                href={`${registration ? "/register" : "/login"}?account=company`}
                aria-disabled={busy || undefined}
                onClick={(event) => {
                  if (inFlight.current) event.preventDefault();
                }}
                aria-current={company ? "page" : undefined}
              >
                <span
                  className={`${styles.optionIcon} ${styles.hiringIcon}`}
                  aria-hidden="true"
                />
                Hiring team
              </Link>
              <Link
                href={registration ? "/register" : "/login"}
                aria-disabled={busy || undefined}
                onClick={(event) => {
                  if (inFlight.current) event.preventDefault();
                }}
                aria-current={!company ? "page" : undefined}
              >
                <span
                  className={`${styles.optionIcon} ${styles.candidateIcon}`}
                  aria-hidden="true"
                />
                Candidate
              </Link>
            </nav>
            <p className={styles.requiredNote}>All fields are required.</p>
            <div
              className={registration && company ? styles.row : styles.fields}
            >
              <FormField
                id="email"
                error={fieldErrors.email}
                name="email"
                label="Email"
                type="email"
                autoComplete="email"
                required
                placeholder={
                  company && registration
                    ? "you@company.com"
                    : "you@example.com"
                }
                icon="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <FormField
                id="password"
                error={fieldErrors.password}
                name="password"
                label="Password"
                type="password"
                autoComplete={
                  registration ? "new-password" : "current-password"
                }
                required
                minLength={registration ? 8 : undefined}
                maxLength={registration ? 128 : undefined}
                hint={registration ? "Use 8–128 characters." : undefined}
                placeholder={
                  registration ? "Create password" : "Enter password"
                }
                icon="lock"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            {registration && (
              <div className={styles.row} key={kind}>
                {company ? (
                  <>
                    <FormField
                      id="companyName"
                      error={fieldErrors.companyName}
                      name="companyName"
                      label="Company name"
                      required
                      maxLength={200}
                      autoComplete="organization"
                      placeholder="Enter company name"
                    />
                    <FormField
                      id="registrationNumber"
                      error={fieldErrors.registrationNumber}
                      name="registrationNumber"
                      label="Tax registration number"
                      required
                      maxLength={100}
                      placeholder="Enter tax number"
                    />
                  </>
                ) : (
                  <>
                    <FormField
                      id="firstName"
                      error={fieldErrors.firstName}
                      name="firstName"
                      label="First name"
                      required
                      maxLength={100}
                      autoComplete="given-name"
                      placeholder="First name"
                    />
                    <FormField
                      id="lastName"
                      error={fieldErrors.lastName}
                      name="lastName"
                      label="Last name"
                      required
                      maxLength={100}
                      autoComplete="family-name"
                      placeholder="Last name"
                    />
                  </>
                )}
              </div>
            )}
            {error && (
              <p role="alert" className={styles.error}>
                {error}
              </p>
            )}
            <div className={styles.actions}>
              <button className={styles.primary} disabled={busy} type="submit">
                {busy
                  ? "Please wait…"
                  : registration
                    ? "Create account"
                    : "Log in to amara"}
              </button>
              {!registration && (
                <Link
                  className={styles.secondary}
                  href={registerHref}
                  aria-disabled={busy || undefined}
                  onClick={(event) => {
                    if (inFlight.current) event.preventDefault();
                  }}
                >
                  Create an account
                </Link>
              )}
            </div>
          </form>
          {registration ? (
            <Link
              className={styles.loginLink}
              href={loginHref}
              aria-disabled={busy || undefined}
              onClick={(event) => {
                if (inFlight.current) event.preventDefault();
              }}
            >
              Already have an account? Log in
            </Link>
          ) : (
            <button
              disabled={busy}
              className={styles.mfaLink}
              onClick={() => {
                if (formRef.current?.reportValidity()) {
                  setError("");
                  setState("mfa");
                }
              }}
            >
              Enter an MFA code
            </button>
          )}
        </>
      ) : (
        <section
          className={styles.stateCard}
          aria-labelledby={stateTitleId}
          aria-live="polite"
        >
          <h1 id={stateTitleId} data-state-title tabIndex={-1}>
            {states[state].title}
          </h1>
          <p>{states[state].description}</p>
          {isMfa ? (
            <form onSubmit={submit} aria-busy={busy}>
              <FormField
                id="mfaCode"
                hint="Enter the six-digit code from your authenticator."
                error={
                  state === "code-rejected"
                    ? "That code couldn’t be verified. Try a new code."
                    : undefined
                }
                name="mfaCode"
                label="Multi-factor code"
                autoComplete="one-time-code"
                inputMode="numeric"
                pattern="[0-9]{6}"
                maxLength={6}
                minLength={6}
                required
                placeholder="Enter code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
              />
              {error && (
                <p role="alert" className={styles.error}>
                  {error}
                </p>
              )}
              <button
                className={styles.stateAction}
                type="submit"
                disabled={busy}
              >
                {busy ? "Verifying…" : "Verify and sign in"}
              </button>
              <button
                disabled={busy}
                className={styles.stateBack}
                type="button"
                onClick={backToForm}
              >
                Back to sign-in
              </button>
            </form>
          ) : (
            <button
              className={styles.stateAction}
              onClick={() => {
                if (registration) router.push(loginHref);
                else backToForm();
              }}
            >
              {state === "review-pending"
                ? "Back to company sign-in"
                : "Back to sign-in"}
            </button>
          )}
        </section>
      )}
    </AuthShell>
  );
}
