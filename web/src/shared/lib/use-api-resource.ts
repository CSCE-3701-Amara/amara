"use client";
import { useEffect, useState } from "react";
import { apiRequest, ApiError } from "./api-client";

/** Cancel stale reads, including hidden Next.js routes; keep private data uncached. */
export function useApiResource<T>(path: string | null) {
  const [revision, setRevision] = useState(0);
  const key = `${path}:${revision}`;
  const [result, setResult] = useState<{
    key: string;
    data?: T;
    error?: ApiError;
  }>();
  useEffect(() => {
    if (!path) return;
    const request = new AbortController();
    void apiRequest<T>(path, { signal: request.signal }).then(
      (data) => {
        if (!request.signal.aborted) setResult({ key, data });
      },
      (error) => {
        if (!request.signal.aborted) setResult({ key, error });
      },
    );
    return () => request.abort();
  }, [path, key]);
  const current = result?.key === key ? result : undefined;
  return {
    data: current?.data,
    error: current?.error,
    loading: !!path && !current,
    reload: () => setRevision((value) => value + 1),
  };
}
