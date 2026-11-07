import { Suspense } from "react";
import { AssessmentInstances } from "@/features/management/components/evaluations";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <AssessmentInstances />
    </Suspense>
  );
}
