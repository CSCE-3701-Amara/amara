import { Suspense } from "react";
import { ModuleLibrary } from "@/features/invw/components/module-studio";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading modules…</p>}>
      <ModuleLibrary />
    </Suspense>
  );
}
