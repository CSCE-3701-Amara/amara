"use client";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { Suspense, useRef, useState, type ReactNode } from "react";
import { deleteCurrentSession } from "@/features/auth/api/auth.api";
import { ApiError } from "@/shared/lib/api-client";
import { Brand } from "./brand";
import styles from "./workspace.module.css";
export { styles as workspaceStyles };

export function ApiFeedback({
  error,
  retry,
}: {
  error?: Error;
  retry?: () => void;
}) {
  if (!error) return null;
  return (
    <div className={styles.error} role="alert">
      <p>{error.message}</p>
      <div className={styles.actions}>
        {retry && (
          <button type="button" onClick={retry}>
            Try again
          </button>
        )}
        {error instanceof ApiError && error.status === 401 && (
          <Link href="/login">Sign in again</Link>
        )}
      </div>
    </div>
  );
}
export function Workspace({
  company = false,
  children,
}: {
  company?: boolean;
  children: ReactNode;
}) {
  const router = useRouter();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const pending = useRef(false);
  async function logout() {
    if (pending.current) return;
    pending.current = true;
    setBusy(true);
    setError(undefined);
    try {
      await deleteCurrentSession();
      router.replace(company ? "/login?account=company" : "/login");
      router.refresh();
    } catch (cause) {
      if (cause instanceof ApiError && cause.status === 401) {
        router.replace("/login?state=expired");
        router.refresh();
      } else setError(cause as Error);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return (
    <div className={styles.shell}>
      <header className={styles.header}>
        <Brand />
        <Suspense fallback={<span>Workspace navigation</span>}>
          <WorkspaceNavigation company={company} />
        </Suspense>
        <button type="button" disabled={busy} onClick={logout}>
          {busy ? "Signing out…" : "Sign out"}
        </button>
      </header>
      <main className={styles.main}>
        <p className={styles.eyebrow}>
          {company ? "Hiring workspace" : "Candidate workspace"}
        </p>
        <ApiFeedback error={error} />
        {children}
      </main>
    </div>
  );
}
export function BusyState({ loading }: { loading: boolean }) {
  return loading ? <p role="status">Loading your workspace…</p> : null;
}

function WorkspaceNavigation({ company }: { company: boolean }) {
  const path = usePathname();
  const links = company
    ? [
        ["/employer/jobs", "Job postings"],
        ["/employer/pipelines", "Assessments & pipelines"],
      ]
    : [
        ["/candidate/dashboard", "Find roles"],
        ["/candidate/applications", "My applications"],
      ];
  return (
    <nav aria-label="Workspace navigation">
      {links.map(([href, label]) => (
        <Link
          href={href}
          key={href}
          aria-current={path === href ? "page" : undefined}
        >
          {label}
        </Link>
      ))}
    </nav>
  );
}
