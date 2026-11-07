import { Suspense } from "react";
import { CandidateProfileScreen } from "@/features/cand/components/profile";
export default function Page() {
  return (
    <Suspense fallback={<p role="status">Loading profile…</p>}>
      <CandidateProfileScreen />
    </Suspense>
  );
}
