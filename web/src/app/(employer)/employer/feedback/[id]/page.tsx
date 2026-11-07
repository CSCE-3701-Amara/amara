import { Suspense } from "react";
import { FeedbackAuthor } from "@/features/management/components/hiring";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <FeedbackAuthor id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
