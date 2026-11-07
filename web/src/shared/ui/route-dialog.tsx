"use client";
import { useRef } from "react";
import { submitAssessmentAttempt } from "@/features/cand/api/assessment.api";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { ConfirmDialog } from "./product-ui";
import { dialogRoutes } from "./dialog-routes";
import { dialogSpecs } from "./dialog-specs";
import { performAccountAction } from "@/features/auth/api/auth.api";
import { transitionPosting, deletePosting } from "@/features/job/api/job.api";
import { publishModule } from "@/features/invw/api/invw.api";
/** Direct links expose the same modal states without duplicating page implementations. */
export function RouteDialog() {
  const path = usePathname();
  const params = useSearchParams();
  const router = useRouter();
  const navigated = useRef<string | null>(null);
  const navigationKey = path + "?" + params.toString();
  const action = params.get("dialog");
  if (!action) return null;
  const route = dialogRoutes.find(
    (row) =>
      row.dialog === action &&
      new RegExp("^" + row.path.replace("[id]", "([^/]+)") + "$").test(path),
  );
  if (!route) return null;
  const spec = dialogSpecs[route.id];
  const id = path.match(/\/(?:jobs|modules)\/([^/]+)/)?.[1];
  const account =
    path === "/candidate/profile/edit" &&
    (action === "deactivate" || action === "delete");
  const posting =
    path.startsWith("/employer/jobs/") &&
    ["publish", "archive", "republish", "delete"].includes(action);
  const moduleAction =
    path.startsWith("/employer/modules/") && action === "publish";
  const assessment =
    path.startsWith("/candidate/assessments/") && action === "submit";
  function close() {
    if (navigated.current === navigationKey) return;
    const next = new URLSearchParams(params.toString());
    next.delete("dialog");
    router.replace(path + (next.size ? "?" + next : ""), { scroll: false });
  }
  async function confirm() {
    if (account) {
      await performAccountAction(action as "delete" | "deactivate");
      navigated.current = navigationKey;
      router.replace("/login");
      router.refresh();
    } else if (posting && id) {
      if (action === "delete") {
        await deletePosting(id);
        navigated.current = navigationKey;
        router.replace("/employer/jobs");
      } else {
        await transitionPosting(
          id,
          action as "publish" | "archive" | "republish",
        );
        close();
        router.refresh();
      }
    } else if (assessment) {
      const instanceId = path.split("/").at(-1)!;
      const attemptId = params.get("attempt") || instanceId;
      await submitAssessmentAttempt(attemptId);
      const next = new URLSearchParams(params.toString());
      next.delete("dialog");
      next.set("state", "submitted");
      next.set("attempt", attemptId);
      navigated.current = navigationKey;
      router.replace(path + "?" + next, { scroll: false });
    } else if (moduleAction && id) {
      await publishModule(id);
      close();
      router.refresh();
    }
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
      unavailable={
        account || posting || moduleAction || assessment
          ? undefined
          : "This action is not available yet. No changes have been saved."
      }
      onConfirm={confirm}
      onClose={close}
    />
  );
}
