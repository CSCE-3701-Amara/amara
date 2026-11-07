import { Suspense } from "react";
import { PostingDetail } from "@/features/cand/components/posting-detail";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <PostingDetail id={id} />;
}
export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p role="status">Loading opportunity…</p>}>
      <Content params={params} />
    </Suspense>
  );
}
