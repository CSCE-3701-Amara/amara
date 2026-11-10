"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { OfferPage } from "@/features/cand/types/product.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  EmptyState,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import s from "./offers-overview.module.css";

const filters = [
  "all",
  "pending",
  "accepted",
  "declined",
  "withdrawn",
  "expired",
] as const;
type Filter = (typeof filters)[number];

/** The offers destination is a list; offer creation has its own existing route. */
export function OffersOverview() {
  const params = useSearchParams();
  const status: Filter = filters.includes(params.get("status") as Filter)
    ? (params.get("status") as Filter)
    : "all";
  const requestedOffset = Number(params.get("offset"));
  const offset =
    Number.isSafeInteger(requestedOffset) && requestedOffset >= 0
      ? requestedOffset
      : 0;
  const resource = useApiResource<OfferPage>(
    `/v1/offers?limit=20&offset=${offset}${status === "all" ? "" : `&status=${status}`}`,
  );
  const rows = resource.data?.data.filter(
    (offer) => status === "all" || offer.status === status,
  );
  return (
    <>
      <ProductHeading
        title="Offers"
        description="Track finalist offers, response deadlines, and outcomes across your postings."
        className={e.largeTitle}
      />
      <nav className={e.actions} aria-label="Offer status filters">
        {filters.map((filter) => (
          <Link
            key={filter}
            href={`/employer/offers${filter === "all" ? "" : `?status=${filter}`}`}
            aria-current={status === filter ? "page" : undefined}
            className={status === filter ? s.selected : undefined}
          >
            {filter === "all"
              ? "All offers"
              : filter[0].toUpperCase() + filter.slice(1)}
          </Link>
        ))}
      </nav>
      <UnavailableNotice>
        Offer listing and changes are not available from the server yet.
      </UnavailableNotice>
      <BusyState loading={resource.loading} />
      {resource.error?.status !== 404 && (
        <ApiFeedback error={resource.error} retry={resource.reload} />
      )}
      <div className={s.list}>
        {rows?.map((offer) => (
          <article className={s.card} key={offer.id}>
            <div>
              <h2>{offer.jobTitle}</h2>
              <p>
                {offer.compensationSummary} · Proposed start{" "}
                {new Date(
                  `${offer.proposedStartDate}T00:00:00Z`,
                ).toLocaleDateString(undefined, {
                  dateStyle: "medium",
                  timeZone: "UTC",
                })}{" "}
                · Respond by {new Date(offer.responseDeadline).toLocaleString()}
              </p>
            </div>
            <StatusBadge
              tone={
                offer.status === "accepted"
                  ? "success"
                  : offer.status === "pending"
                    ? "warning"
                    : "neutral"
              }
            >
              {offer.status[0].toUpperCase() + offer.status.slice(1)}
            </StatusBadge>
            <Link
              className={e.pill}
              href={`/employer/offers/${resourceId(offer.id)}`}
            >
              {offer.status === "pending" ? "Manage offer" : "View offer"}
            </Link>
          </article>
        ))}
      </div>
      {!resource.loading && !rows?.length && (
        <EmptyState
          title="No offers available"
          description="Finalist offers will appear here when the service is available."
        />
      )}
      {resource.data && (
        <>
          <p className={s.count}>
            {resource.data.pagination.total} offers · Page{" "}
            {Math.floor(offset / 20) + 1}
          </p>
          <nav className={e.actions} aria-label="Offer pages">
            {offset > 0 && (
              <Link
                href={`?status=${status}&offset=${Math.max(0, offset - 20)}`}
              >
                Previous page
              </Link>
            )}
            {offset + 20 < resource.data.pagination.total && (
              <Link href={`?status=${status}&offset=${offset + 20}`}>
                Next page
              </Link>
            )}
          </nav>
        </>
      )}
      <div className={e.actions}>
        <Link href="/employer/jobs">Back to postings</Link>
        <Link className={e.primary} href="/employer/offers/new">
          Create finalist offer
        </Link>
      </div>
    </>
  );
}
