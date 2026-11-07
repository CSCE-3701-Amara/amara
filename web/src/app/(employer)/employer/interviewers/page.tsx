import { Suspense } from "react";
import { InterviewerManagement } from "@/features/management/components/hiring";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <InterviewerManagement />
    </Suspense>
  );
}
