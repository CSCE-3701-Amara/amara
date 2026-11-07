import { Suspense } from "react";
import { Interviews } from "@/features/management/components/evaluations";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <Interviews />
    </Suspense>
  );
}
