import { Suspense } from "react";
import { InterviewAvailability } from "@/features/management/components/interview-availability";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <InterviewAvailability />
    </Suspense>
  );
}
