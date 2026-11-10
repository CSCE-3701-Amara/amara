import { Suspense } from "react";
import { Interviews } from "@/features/management/components/interviews";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <Interviews />
    </Suspense>
  );
}
