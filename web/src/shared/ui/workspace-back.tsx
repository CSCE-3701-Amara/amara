"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";
import { workspaceBackLink } from "@/shared/lib/workspace-back-link";
import styles from "./workspace.module.css";

export function WorkspaceBack({ home }: { home: string }) {
  const path = usePathname();
  const params = useSearchParams();
  const destination = workspaceBackLink(path, home, params.toString());
  if (!destination) return null;

  return (
    <Link href={destination.href} className={styles.backLink}>
      <Image src="/figma/chevron-left.svg" width={16} height={16} alt="" />
      {destination.label}
    </Link>
  );
}
