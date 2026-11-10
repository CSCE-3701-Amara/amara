import { expect, test } from "@playwright/test";
import { NextRequest } from "next/server";
import { proxy } from "../../src/proxy";

for (const origin of ["https://amara.example", "http://127.0.0.1:3000"]) {
  test(`proxy normalizes a same-origin mutation from ${origin}`, () => {
    const request = new NextRequest(`${origin}/v1/auth/sessions`, {
      method: "POST",
      headers: { origin, host: new URL(origin).host, cookie: "amara_session=example", "content-type": "application/json" },
      body: "{}",
    });
    const response = proxy(request);
    const overridden = response.headers.get("x-middleware-override-headers")?.split(",");
    expect(overridden).toBeDefined();
    expect(overridden).not.toContain("origin");
    expect(response.headers.get("x-middleware-request-cookie")).toBe("amara_session=example");
    expect(response.headers.get("x-middleware-request-content-type")).toBe("application/json");
    expect(request.headers.get("origin")).toBe(origin);
  });
}
for (const origin of ["https://untrusted.example", "https://amara.example.evil.test", "null", ""]) {
  test(`proxy preserves backend CORS checks for origin ${origin || "absent"}`, () => {
    const request = new NextRequest("https://amara.example/v1/auth/sessions", {
      method: "POST",
      headers: { host: "amara.example", ...(origin ? { origin } : {}), "x-forwarded-host": "untrusted.example" },
    });
    const response = proxy(request);
    expect(response.headers.has("x-middleware-override-headers")).toBe(false);
    expect(request.headers.get("origin")).toBe(origin || null);
  });
}
