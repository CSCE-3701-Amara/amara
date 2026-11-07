import { Suspense } from "react";
import { CandidateOffers } from "@/features/cand/components/offers";
export default function Page() {
  return (
    <Suspense fallback={<p role="status">Loading offers…</p>}>
      <CandidateOffers />
    </Suspense>
  );
}
