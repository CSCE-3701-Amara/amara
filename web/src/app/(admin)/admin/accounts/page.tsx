import { Suspense } from "react";
import { Administration } from "@/features/management/components/admin";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <Administration accounts />
    </Suspense>
  );
}
