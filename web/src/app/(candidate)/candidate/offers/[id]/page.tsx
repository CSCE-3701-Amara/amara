import { Suspense } from "react";
import { CandidateOffers } from "@/features/cand/components/offers";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <CandidateOffers id={id} />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading offer…</p>}>
      <Content params={params} />
    </Suspense>
  );
}
