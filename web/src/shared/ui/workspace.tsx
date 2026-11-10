"use client";
import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { Suspense, useId, useRef, useState, type ReactNode } from "react";
import { deleteCurrentSession } from "@/features/auth/api/auth.api";
import { ApiError } from "@/shared/lib/api-client";
import { sessionHome } from "@/features/auth/lib/session-home";
import { Brand } from "./brand";
import { RouteDialog } from "./route-dialog";
import styles from "./workspace.module.css";
export { styles as workspaceStyles };
export type WorkspaceRole = "candidate" | "company" | "interviewer" | "admin";
const destinations: Record<
  WorkspaceRole,
  readonly (readonly [string, string])[]
> = {
  candidate: [
    ["/candidate/applications", "Applications"],
    ["/candidate/offers", "Offers"],
    ["/candidate/profile", "Profile"],
    ["/candidate/dashboard", "Feed"],
  ],
  company: [
    ["/employer/jobs", "Job postings"],
    ["/employer/pipelines", "Pipeline studio"],
    ["/employer/modules", "Modules"],
    ["/employer/interviews", "Interviews"],
    ["/employer/grading", "Grading"],
    ["/employer/offers", "Offers"],
    ["/employer/profile", "Company profile"],
  ],
  interviewer: [
    ["/interviewer/interviews", "Interviews"],
    ["/interviewer/grading", "Assigned grading"],
    ["/interviewer/availability", "Availability"],
    ["/interviewer/profile", "Profile"],
  ],
  admin: [
    ["/admin/dashboard", "Company registrations"],
    ["/admin/accounts", "Accounts"],
  ],
};

/** Displays the API's real failure without interpreting every 401 as an expired session. */
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
function WorkspaceNavigation({
  role,
  onNavigate,
}: {
  role: WorkspaceRole;
  onNavigate: () => void;
}) {
  const path = usePathname();
  return (
    <nav aria-label="Workspace navigation" className={styles.navigation}>
      {destinations[role].map(([href, label]) => (
        <Link
          key={href}
          href={href}
          onClick={onNavigate}
          aria-current={
            path === href ||
            (href !== "/candidate/dashboard" && path.startsWith(href + "/"))
              ? "page"
              : undefined
          }
        >
          {label}
        </Link>
      ))}
    </nav>
  );
}

/** Shared workspace chrome with stable branding and collapsible navigation. */
export function Workspace({
  company = false,
  role: suppliedRole,
  children,
}: {
  company?: boolean;
  role?: WorkspaceRole;
  children: ReactNode;
}) {
  const role = suppliedRole || (company ? "company" : "candidate");
  const workspaceLabel =
    role === "company"
      ? "Hiring workspace"
      : `${role[0].toUpperCase()}${role.slice(1)} workspace`;
  const router = useRouter();
  const menuId = useId();
  const [collapsed, setCollapsed] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
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
      router.replace(role === "company" ? "/login?account=company" : "/login");
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
    <div
      className={styles.shell}
      data-collapsed={collapsed}
      data-mobile-open={mobileOpen}
    >
      <a className={styles.skipLink} href="#workspace-content">
        Skip to content
      </a>
      <aside
        className={styles.sidebar}
        aria-label={`${role === "company" ? "Hiring" : role} workspace`}
      >
        <div className={styles.brandRow}>
          <Brand variant="workspace" href={sessionHome(role)} />
          <button
            className={styles.mobileToggle}
            type="button"
            aria-expanded={mobileOpen}
            aria-controls={menuId}
            onClick={() => {
              setCollapsed(false);
              setMobileOpen((value) => !value);
            }}
          >
            {mobileOpen ? "Close menu" : "Menu"}
          </button>
        </div>
        <div className={styles.menu} id={menuId}>
          <button
            type="button"
            className={styles.collapseToggle}
            aria-expanded={!collapsed}
            aria-label={collapsed ? "Show menu" : "Hide menu"}
            onClick={() => setCollapsed((value) => !value)}
          >
            <Image
              src="/figma/chevron-left.svg"
              alt=""
              width={16}
              height={16}
            />
            {!collapsed && "Hide menu"}
          </button>
          <div className={styles.menuInner}>
            {!collapsed && (
              <p className={styles.workspaceLabel}>{workspaceLabel}</p>
            )}
            {!collapsed && (
              <Suspense fallback={<p>Loading navigation…</p>}>
                <WorkspaceNavigation
                  role={role}
                  onNavigate={() => setMobileOpen(false)}
                />
              </Suspense>
            )}
            <button
              type="button"
              className={styles.logout}
              disabled={busy}
              onClick={logout}
              aria-label={busy ? "Signing out" : "Sign out"}
            >
              {collapsed ? (
                <span aria-hidden="true">↪</span>
              ) : busy ? (
                "Signing out…"
              ) : (
                "Log out"
              )}
            </button>
          </div>
        </div>
      </aside>
      <main id="workspace-content" tabIndex={-1} className={styles.main}>
        <div className={styles.contextBar}>
          <span>{workspaceLabel}</span>
          <span>amara / Skills-first recruitment</span>
        </div>
        <ApiFeedback error={error} />
        {children}
      </main>
      <Suspense fallback={null}>
        <RouteDialog />
      </Suspense>
    </div>
  );
}
/** A polite loading announcement shared by resource-backed screens. */
export function BusyState({ loading }: { loading: boolean }) {
  return loading ? (
    <p role="status" className={styles.loading}>
      Loading your workspace…
    </p>
  ) : null;
}
