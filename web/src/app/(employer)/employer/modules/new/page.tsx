import { Suspense } from "react";
import { ModuleStudio } from "@/features/invw/components/module-studio";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading modules…</p>}>
      <ModuleStudio />
    </Suspense>
  );
}
