import { Suspense } from "react";
import { CandidateInterview } from "@/features/cand/components/interview";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <CandidateInterview id={id} />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading interview…</p>}>
      <Content params={params} />
    </Suspense>
  );
}
