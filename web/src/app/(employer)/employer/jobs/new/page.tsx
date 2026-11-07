import { Suspense } from "react";
import { PostingEditor } from "@/features/job/components/posting-editor";
export default function Page() {
  return (
    <Suspense fallback={<p>Loading editor…</p>}>
      <PostingEditor />
    </Suspense>
  );
}
