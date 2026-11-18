"use client";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { ConfirmDialog } from "./product-ui";
import { dialogRoutes } from "./dialog-routes";
import { dialogSpecs } from "./dialog-specs";

/**
 * Direct links preview the designed modal states. They never perform the
 * action: a URL can be crafted by anyone, so a link must not be able to put a
 * genuine "Delete account" or "Submit assessment" confirmation in front of a
 * user. Real actions start only from controls on the page itself.
 */
export function RouteDialog() {
  const path = usePathname();
  const params = useSearchParams();
  const router = useRouter();
  const action = params.get("dialog");
  if (!action) return null;
  const route = dialogRoutes.find(
    (row) => row.dialog === action && matchesRoute(row.path, path),
  );
  if (!route) return null;
  const spec = dialogSpecs[route.id];
  function close() {
    const next = new URLSearchParams(params.toString());
    next.delete("dialog");
    router.replace(path + (next.size ? "?" + next : ""), { scroll: false });
  }
  const description = spec.description
    .replace(/Acme Corp/g, "the company shown in the offer")
    .replace(
      /Thursday, 12 November at 1:00 PM \(UTC\+02:00\)/g,
      "the time shown in your interview details",
    );
  return (
    <ConfirmDialog
      key={route.id}
      title={spec.title}
      description={description}
      confirmLabel={spec.confirmLabel}
      intent={spec.intent}
      reason={spec.requiresReason}
      unavailable="This preview can't make changes. Close it and use the button on this page to continue."
      onClose={close}
    />
  );
}

/** Compare a route pattern such as `/candidate/offers/[id]` with a path. */
function matchesRoute(pattern: string, path: string) {
  const expected = pattern.split("/");
  const actual = path.split("/");
  return (
    expected.length === actual.length &&
    expected.every(
      (segment, i) =>
        segment === actual[i] || (segment === "[id]" && actual[i] !== ""),
    )
  );
}
