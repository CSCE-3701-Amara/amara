import { NextResponse, type NextRequest } from "next/server";

/** The browser's same-origin request becomes server-to-server at the rewrite. */
export function proxy(request: NextRequest) {
  const origin = request.headers.get("origin");
  // Next normalizes loopback hostnames in nextUrl; Host retains the browser's
  // actual origin, including its port. Forwarded-host headers are not trusted.
  const frontendOrigin = `${request.nextUrl.protocol}//${request.headers.get("host")}`;
  if (origin !== frontendOrigin) return NextResponse.next();

  // Next rewrites change Host to the API host. Retaining the frontend Origin
  // would make Go mistake this validated same-origin request for a CORS call.
  const headers = new Headers(request.headers);
  headers.delete("origin");
  return NextResponse.next({ request: { headers } });
}

export const config = {
  matcher: ["/v1/:path*", "/health/:path*"],
};
