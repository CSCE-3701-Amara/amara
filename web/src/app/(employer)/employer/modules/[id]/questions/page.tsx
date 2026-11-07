import { Suspense } from "react";
import { ModuleStudio } from "@/features/invw/components/module-studio";
async function Content({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <ModuleStudio id={id} />;
}
export default function Page(props: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<p>Loading module…</p>}>
      <Content {...props} />
    </Suspense>
  );
}
