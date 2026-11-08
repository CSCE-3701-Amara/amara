"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { Offer, OfferPage } from "@/features/cand/types/product.types";
import {
  ConfirmDialog,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import { PlannedRead } from "./planned-read";

const art = { src: "/figma/hiring-dag.svg", width: 200.16, height: 138.96 };
export function HiringOffer({ id }: { id?: string }) {
  const params = useSearchParams();
  const [dialog, setDialog] = useState<string | null>(null);
  const detail = useApiResource<Offer>(
    id ? `/v1/offers/${resourceId(id)}` : null,
  );
  const list = useApiResource<OfferPage>(
    !id && !params.has("application") ? "/v1/offers?limit=20&offset=0" : null,
  );
  const create = !id;
  const capacity = params.get("state") === "capacity";
  return (
    <>
      <ProductHeading
        title={create ? "Create finalist offer" : "Manage offer"}
        description={detail.data?.jobTitle}
        compact
        artwork={art}
      />
      <UnavailableNotice>
        Offers cannot be sent, withdrawn, or extended yet. Entered details
        remain on this screen.
      </UnavailableNotice>
      <PlannedRead
        loading={detail.loading || list.loading}
        error={detail.error || list.error}
        retry={() => {
          detail.reload();
          list.reload();
        }}
      />
      {list.data?.data.map((offer) => (
        <article key={offer.id} className={e.panel}>
          <h2>{offer.jobTitle}</h2>
          <p>
            {offer.status} · {offer.compensationSummary}
          </p>
          <Link href={`/employer/offers/${resourceId(offer.id)}`}>
            Manage offer
          </Link>
        </article>
      ))}
      <form
        onSubmit={(event) => {
          event.preventDefault();
          setDialog(create ? "send" : "extend");
        }}
      >
        <section className={e.section}>
          <h2>
            {create
              ? "Offer details"
              : detail.data?.status || "Pending response"}
          </h2>
          {create ? (
            <div className={e.grid}>
              <label className={e.field}>
                Compensation summary (required)
                <textarea name="compensationSummary" required />
              </label>
              <label className={e.field}>
                Proposed start date (required)
                <input name="proposedStartDate" type="date" required />
              </label>
              <label className={e.field}>
                Response deadline (required)
                <input name="responseDeadline" type="datetime-local" required />
              </label>
              <label className={e.field}>
                Message
                <textarea name="message" />
              </label>
            </div>
          ) : (
            <>
              <p>
                {detail.data?.compensationSummary ||
                  "No offer details available"}
              </p>
              <p>Proposed start: {detail.data?.proposedStartDate || "—"}</p>
              <p>
                Response deadline:{" "}
                {detail.data?.responseDeadline
                  ? new Date(detail.data.responseDeadline).toLocaleString()
                  : "—"}
              </p>
              <label className={e.field}>
                New response deadline
                <input name="responseDeadline" type="datetime-local" required />
              </label>
            </>
          )}
        </section>
        <section className={e.section}>
          <h2>{create ? "Capacity" : "Deadline extension"}</h2>
          <p>
            {create
              ? "If the posting limit is exceeded, confirmation is required before sending."
              : "A new deadline does not change the other terms of this offer."}
          </p>
          {capacity && (
            <p className={e.error} role="alert">
              Capacity must be verified before sending. This offer has not been
              sent.
            </p>
          )}
        </section>
        <div className={e.actions}>
          <Link href="/employer/jobs">Cancel</Link>
          {id && (
            <button type="button" onClick={() => setDialog("withdraw")}>
              Review withdrawal
            </button>
          )}
          <button type="submit" className={e.primary}>
            {create ? "Review offer" : "Review deadline extension"}
          </button>
        </div>
      </form>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "withdraw"
              ? "Withdraw this offer?"
              : dialog === "extend"
                ? "Extend offer deadline?"
                : capacity
                  ? "Send offer above posting capacity?"
                  : "Send this offer?"
          }
          description="Review the saved offer terms and recipient before continuing. No offer has been sent or changed."
          confirmLabel={
            dialog === "withdraw"
              ? "Withdraw offer"
              : dialog === "extend"
                ? "Extend deadline"
                : "Send offer"
          }
          intent={dialog === "withdraw" ? "danger" : "warning"}
          unavailable="This offer operation is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
