"use client";
import Link from "next/link";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
} from "@/shared/ui/product-ui";
import type { Offer, OfferPage } from "../types/product.types";
import styles from "./candidate.module.css";
export function CandidateOffers({ id }: { id?: string }) {
  const [dialog, setDialog] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const list = useApiResource<OfferPage>(
    id ? null : `/v1/offers?limit=20&offset=${offset}`,
  );
  const detail = useApiResource<Offer>(
    id ? `/v1/offers/${resourceId(id)}` : null,
  );
  const [selected, setSelected] = useState<Offer>();
  const offer = detail.data || selected;
  const error = detail.error || list.error;
  return (
    <section data-figma-node={id ? "375:1232" : "212:71"}>
      <ProductHeading
        title={id ? "Your offer" : "Offers"}
        description={
          offer
            ? `${offer.jobTitle} · Received ${new Date(offer.createdAt).toLocaleDateString()}`
            : undefined
        }
      />
      <UnavailableNotice>
        Offer review and responses aren’t available yet. Accepting or declining
        an offer is currently unavailable.
      </UnavailableNotice>
      <BusyState loading={list.loading || detail.loading} />
      {error?.status !== 404 && (
        <ApiFeedback error={error} retry={id ? detail.reload : list.reload} />
      )}
      {!id && list.data && (
        <>
          {list.data.data.map((item) => (
            <article className={styles.panel} key={item.id}>
              <div className={s.actions}>
                <h2>{item.jobTitle}</h2>
                <StatusBadge>{item.status}</StatusBadge>
              </div>
              <p>{item.compensationSummary}</p>
              {item.message && <p>{item.message}</p>}
              <p>
                Proposed · {item.proposedStartDate} · Deadline ·{" "}
                {new Date(item.responseDeadline).toLocaleString()}
              </p>
              <div className={s.actions}>
                <button
                  type="button"
                  className={s.primary}
                  disabled={item.status !== "pending"}
                  onClick={() => {
                    setSelected(item);
                    setDialog("accept");
                  }}
                >
                  Accept offer
                </button>
                <Link
                  className={s.button}
                  href={`/candidate/offers/${resourceId(item.id)}`}
                >
                  Details
                </Link>
                <button
                  type="button"
                  disabled={item.status !== "pending"}
                  onClick={() => {
                    setSelected(item);
                    setDialog("decline");
                  }}
                >
                  Decline offer
                </button>
              </div>
            </article>
          ))}
          {list.data.data.length === 0 && (
            <p>No offers are available to review.</p>
          )}
          <div className={s.actions}>
            <button
              type="button"
              disabled={offset === 0}
              onClick={() => setOffset(Math.max(0, offset - 20))}
            >
              Previous page
            </button>
            <button
              type="button"
              disabled={offset + 20 >= list.data.pagination.total}
              onClick={() => setOffset(offset + 20)}
            >
              Next page
            </button>
          </div>
        </>
      )}
      {id && (
        <div className={styles.detailGrid}>
          <section className={styles.panel}>
            <h2>{offer?.jobTitle || "Offer details"}</h2>
            <StatusBadge>{offer?.status || "Offer unavailable"}</StatusBadge>
            <h2>Offer details</h2>
            <dl>
              <dt>Compensation</dt>
              <dd>{offer?.compensationSummary || "Not provided"}</dd>
              <dt>Start date</dt>
              <dd>{offer?.proposedStartDate || "Not provided"}</dd>
              <dt>Response deadline</dt>
              <dd>
                {offer?.responseDeadline
                  ? new Date(offer.responseDeadline).toLocaleString()
                  : "Not provided"}
              </dd>
            </dl>
            <h2>A message from the hiring team</h2>
            <p>
              {offer?.message ||
                "The offer message will appear here when it becomes available."}
            </p>
          </section>
          <aside className={styles.statusPanel}>
            <h2>Ready to respond?</h2>
            <p>
              Accepting confirms your intention to join. You can also decline
              the offer.
            </p>
            <div className={s.actions}>
              <button
                type="button"
                disabled={!offer || offer.status !== "pending"}
                onClick={() => setDialog("decline")}
              >
                Decline offer
              </button>
              <button
                type="button"
                className={s.primary}
                disabled={!offer || offer.status !== "pending"}
                onClick={() => setDialog("accept")}
              >
                Accept offer
              </button>
            </div>
            <Link href="/candidate/offers">Back to offers</Link>
          </aside>
        </div>
      )}
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "decline" ? "Decline this offer?" : "Accept this offer?"
          }
          description={
            dialog === "decline"
              ? "Declining tells the hiring team that you will not be accepting this offer."
              : "Accepting confirms your intention to join under the terms shown in this offer."
          }
          confirmLabel={dialog === "decline" ? "Decline offer" : "Accept offer"}
          intent={dialog === "decline" ? "danger" : "success"}
          unavailable="Offer responses aren’t available yet. Your offer will remain unchanged."
          onClose={() => setDialog(null)}
        />
      )}
    </section>
  );
}
