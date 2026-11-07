import type { Metadata } from "next";
import { Suspense } from "react";
import { AuthRoute } from "@/features/auth/components/auth-screen";

export const metadata: Metadata = { title: "Create an account" };
export default function Page() {
  return (
    <Suspense fallback={<main aria-busy="true">Loading registration…</main>}>
      <AuthRoute registration />
    </Suspense>
  );
}
