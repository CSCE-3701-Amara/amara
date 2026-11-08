import { Suspense } from "react";
import { HiringOffer } from "@/features/management/components/hiring-offer";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <HiringOffer />
    </Suspense>
  );
}
