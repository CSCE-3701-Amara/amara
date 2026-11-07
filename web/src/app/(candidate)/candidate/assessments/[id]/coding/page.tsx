import { Suspense } from "react";
import { CandidateAssessment } from "@/features/cand/components/assessment";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <CandidateAssessment id={id} variant="coding" />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading assessment…</p>}>
      <Content params={params} />
    </Suspense>
  );
}
