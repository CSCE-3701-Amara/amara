import { Suspense } from "react";
import { StageDecisions } from "@/features/management/components/evaluations";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <StageDecisions id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
