"use client";

import { useEffect, useState, type ReactNode } from "react";
import { useRouter } from "next/navigation";
import { getCurrentSession } from "../api/auth.api";
import { sessionHome } from "../lib/session-home";
import styles from "./session-entry.module.css";

/** Reuse the validated cookie without keeping tokens or identity in browser storage. */
export function SessionEntry({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [revision, setRevision] = useState(0);
  const [result, setResult] = useState<{
    revision: number;
    guest?: boolean;
    error?: Error;
  }>();

  useEffect(() => {
    const request = new AbortController();
    void getCurrentSession(request.signal).then(
      (session) => {
        if (request.signal.aborted) return;
        if (session) router.replace(sessionHome(session.user.role)!);
        else setResult({ revision, guest: true });
      },
      (error: Error) => {
        if (!request.signal.aborted) setResult({ revision, error });
      },
    );
    return () => request.abort();
  }, [revision, router]);

  const current = result?.revision === revision ? result : undefined;
  if (current?.guest) return children;
  if (current?.error) {
    return (
      <main className={styles.status}>
        <h1>We couldn’t check your session</h1>
        <p role="alert">{current.error.message}</p>
        <button type="button" onClick={() => setRevision(value => value + 1)}>
          Try again
        </button>
      </main>
    );
  }
  return (
    <main className={styles.status} aria-busy="true">
      <p role="status">Checking your session…</p>
    </main>
  );
}
