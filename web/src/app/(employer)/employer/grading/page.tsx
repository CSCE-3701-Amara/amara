import { Suspense } from "react";
import { AssignedGrading } from "@/features/management/components/evaluations";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <AssignedGrading />
    </Suspense>
  );
}
