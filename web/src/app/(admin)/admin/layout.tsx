import { Workspace } from "@/shared/ui/workspace";
export default function Layout({ children }: { children: React.ReactNode }) {
  return <Workspace role="admin">{children}</Workspace>;
}
