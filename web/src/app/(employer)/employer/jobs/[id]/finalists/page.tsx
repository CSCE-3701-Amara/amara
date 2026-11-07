import { Suspense } from "react";
import { HiringApplications } from "@/features/management/components/hiring";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <HiringApplications id={id} view="finalists" />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading hiring workspace…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
