import { Suspense } from "react";
import { HiringOffer } from "@/features/management/components/hiring";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <HiringOffer id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
