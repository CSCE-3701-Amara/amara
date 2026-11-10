import { Suspense } from "react";
import { InterviewScorecard } from "@/features/management/components/interview-scorecard";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <InterviewScorecard id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading workspace…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
