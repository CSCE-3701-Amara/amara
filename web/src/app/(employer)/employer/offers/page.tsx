import { Suspense } from "react";
import { OffersOverview } from "@/features/management/components/offers-overview";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <OffersOverview />
    </Suspense>
  );
}
