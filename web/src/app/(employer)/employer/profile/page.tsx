import { Suspense } from "react";
import { TeamProfile } from "@/features/management/components/profile";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <TeamProfile />
    </Suspense>
  );
}
