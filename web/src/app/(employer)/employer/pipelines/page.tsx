"use client";

import Link from "next/link";
import type { PostingList } from "@/features/job/types/job.types";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { ProductHeading, StatusBadge } from "@/shared/ui/product-ui";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import styles from "@/shared/ui/editor.module.css";

export default function Page() {
  const postings = useApiResource<PostingList>("/v1/jobs");
  const items = postings.data?.items ?? [];
  return (
    <section>
      <ProductHeading
        title="Pipeline studio"
        description="Choose a job posting to configure its assessment and interview stages."
      />
      <BusyState loading={postings.loading} />
      <ApiFeedback error={postings.error} retry={postings.reload} />
      {postings.data && (items.length ? (
        <ul className={styles.resourceList}>
          {items.map((posting) => (
            <li key={posting.id}>
              <Link href={`/employer/pipelines/${resourceId(posting.id)}`}>
                {posting.title}
              </Link>{" "}
              <StatusBadge>{posting.status}</StatusBadge>
            </li>
          ))}
        </ul>
      ) : (
        <p>
          Create a job posting to start building its pipeline.{" "}
          <Link href="/employer/jobs/new">Create posting</Link>
        </p>
      ))}
    </section>
  );
}
