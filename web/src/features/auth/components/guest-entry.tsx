import { cookies } from "next/headers";
import type { ReactNode } from "react";
import { SessionEntry } from "./session-entry";

/** Cookie presence triggers validation; presence alone never authenticates a visitor. */
export async function GuestEntry({ children }: { children: ReactNode }) {
  const cookieStore = await cookies();
  if (!cookieStore.get("amara_session")?.value) return children;
  return <SessionEntry>{children}</SessionEntry>;
}
