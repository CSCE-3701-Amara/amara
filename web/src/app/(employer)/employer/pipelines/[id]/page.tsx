import { Suspense } from "react";
import { PipelineStudio } from "@/features/invw/components/pipeline-studio";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <PipelineStudio postingId={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading pipeline…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
