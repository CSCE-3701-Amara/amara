/** Error metadata is optional for failures outside the HTTP response path. */
export type RequestError = Error & {
  readonly status?: number;
  readonly code?: string;
  readonly retryAfter?: string | null;
  readonly requestId?: string | null;
};

/** Normalize thrown values without losing an Error's message, subtype or metadata. */
export function toError(
  value: unknown,
  fallback = "The action couldn’t be completed. Try again.",
): Error {
  return value instanceof Error ? value : new Error(fallback);
}
