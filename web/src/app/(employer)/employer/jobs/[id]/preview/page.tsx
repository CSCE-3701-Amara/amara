import { Suspense } from "react";
import { PostingPreview } from "@/features/job/components/posting-preview";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <PostingPreview id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading preview…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
