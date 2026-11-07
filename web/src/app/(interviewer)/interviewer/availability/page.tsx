import { Suspense } from "react";
import { InterviewAvailability } from "@/features/management/components/evaluations";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <InterviewAvailability />
    </Suspense>
  );
}
