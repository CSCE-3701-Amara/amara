import { Suspense } from "react";
import { ApplicationDetail } from "@/features/cand/components/application-detail";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <ApplicationDetail id={id} />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading application…</p>}>
      <Content params={params} />
    </Suspense>
  );
}
