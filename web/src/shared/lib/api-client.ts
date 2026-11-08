/** Same-origin API calls preserve the backend's HttpOnly session cookie. */
export class ApiError extends Error {
  constructor(
    message: string,
    readonly code: string,
    readonly status: number,
    readonly retryAfter?: string | null,
    readonly requestId?: string | null,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/** Both registered error formats are untrusted JSON at the transport boundary. */
function errorFields(data: unknown) {
  const problem = isRecord(data) ? data : {};
  const legacy = isRecord(problem.error) ? problem.error : {};
  return {
    message: [problem.detail, legacy.message].find(
      (value): value is string => typeof value === "string" && value.length > 0,
    ),
    code: [problem.code, legacy.code].find(
      (value): value is string => typeof value === "string" && value.length > 0,
    ),
    requestId:
      typeof problem.request_id === "string" ? problem.request_id : undefined,
  };
}

export async function apiRequest<T>(
  path: string,
  options: {
    method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
    body?: unknown;
    signal?: AbortSignal;
  } = {},
): Promise<T> {
  // Never send session-bearing requests to URLs supplied by API response data.
  if (!path.startsWith("/v1/") || path.startsWith("//"))
    throw new Error("Invalid API path");
  let response: Response;
  try {
    response = await fetch(path, {
      method: options.method || "GET",
      credentials: "include",
      cache: "no-store",
      signal: options.signal,
      headers: {
        Accept: "application/json, application/problem+json",
        ...(options.body === undefined
          ? {}
          : { "Content-Type": "application/json" }),
      },
      ...(options.body === undefined
        ? {}
        : { body: JSON.stringify(options.body) }),
    });
  } catch (cause) {
    if (options.signal?.aborted) throw cause;
    throw new ApiError(
      "We couldn’t reach amara. Check your connection and try again.",
      "network_error",
      0,
    );
  }
  if (response.status === 204 && response.ok) return undefined as T;
  const contentType = response.headers.get("Content-Type") || "";
  const data: unknown = contentType.includes("json")
    ? await response.json().catch(() => null)
    : null;
  if (!response.ok) {
    const fields = errorFields(data);
    throw new ApiError(
      fields.message || "We couldn’t complete your request. Please try again.",
      fields.code || "request_failed",
      response.status,
      response.headers.get("Retry-After"),
      response.headers.get("X-Request-ID") || fields.requestId,
    );
  }
  if (data === null)
    throw new ApiError(
      "The server returned an unexpected response. Please try again.",
      "invalid_response",
      response.status,
    );
  return data as T;
}

export function queryString(values: object) {
  const params = new URLSearchParams();
  for (const [name, value] of Object.entries(values)) {
    if (value !== undefined && value !== null && value !== "")
      params.set(name, String(value));
  }
  const query = params.toString();
  return query ? `?${query}` : "";
}
export function resourceId(value: string) {
  return encodeURIComponent(value);
}
