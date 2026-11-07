import { Suspense } from "react";
import { PostingEditor } from "@/features/job/components/posting-editor";
async function PostingPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <PostingEditor id={id} />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading posting…</p>}>
      <PostingPage params={params} />
    </Suspense>
  );
}
