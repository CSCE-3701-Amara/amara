"use client";
import Link from "next/link";
import { Calendar, localDateKey } from "@/shared/ui/calendar";
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
import type { Interview, TimeSlot } from "../types/product.types";
import styles from "./candidate.module.css";
export function CandidateInterview({ id }: { id: string }) {
  const [dialog, setDialog] = useState<string | null>(null);
  const [selected, setSelected] = useState<TimeSlot>();
  const interview = useApiResource<Interview>(
    `/v1/interviews/${resourceId(id)}`,
  );
  const availability = useApiResource<TimeSlot[]>(
    `/v1/interviews/${resourceId(id)}/availability`,
  );
  const error = interview.error || availability.error;
  const slots = availability.data || [];
  const dates = [
    ...new Set(slots.map((slot) => localDateKey(new Date(slot.startAt)))),
  ];
  const [day, setDay] = useState("");
  const displayedDay = day || dates[0] || "";
  return (
    <section data-figma-node="375:1077">
      <ProductHeading
        title="Schedule your interview"
        description="Choose a date and time to meet the hiring team."
      />
      <UnavailableNotice>
        Interview scheduling isn’t available yet. Existing confirmations and
        availability can’t be changed.
      </UnavailableNotice>
      <section className={styles.panel}>
        <div className={s.actions}>
          <h2>Meet the hiring team</h2>
          <StatusBadge>
            {interview.data?.status || "Scheduling unavailable"}
          </StatusBadge>
        </div>
        <p>
          Meeting details are shared with authorized participants after
          confirmation.
        </p>
        {interview.data?.joinUrl &&
          /^https?:\/\//.test(interview.data.joinUrl) && (
            <a
              href={interview.data.joinUrl}
              target="_blank"
              rel="noopener noreferrer"
            >
              Join interview
            </a>
          )}
      </section>
      <BusyState loading={interview.loading || availability.loading} />
      {error?.status !== 404 && (
        <ApiFeedback
          error={error}
          retry={() => {
            interview.reload();
            availability.reload();
          }}
        />
      )}
      <div className={styles.schedule}>
        <section className={styles.panel}>
          <h2>Choose a date</h2>
          <p>All times are shown in your local timezone.</p>
          <Calendar
            dates={dates}
            selected={displayedDay}
            onSelect={(date) => {
              setDay(date);
              setSelected(undefined);
            }}
          />
          {!slots.length && <p>No interview slots can be loaded yet.</p>}
        </section>
        <section className={styles.panel}>
          <h2>Available times</h2>
          <div className={styles.slots}>
            {slots
              .filter(
                (slot) => localDateKey(new Date(slot.startAt)) === displayedDay,
              )
              .map((slot) => (
                <button
                  key={slot.startAt}
                  type="button"
                  aria-pressed={selected?.startAt === slot.startAt}
                  onClick={() => setSelected(slot)}
                >
                  {new Date(slot.startAt).toLocaleTimeString(undefined, {
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                  <br />
                  {Math.round(
                    (Date.parse(slot.endAt) - Date.parse(slot.startAt)) / 60000,
                  )}{" "}
                  minutes
                </button>
              ))}
          </div>
          <div className={s.actions}>
            <button
              type="button"
              className={s.primary}
              disabled={!selected}
              onClick={() => setDialog("confirm")}
            >
              Confirm interview time
            </button>
            <button
              type="button"
              disabled={!interview.data}
              onClick={() => setDialog("reschedule")}
            >
              Request a different time
            </button>
          </div>
          <p>
            You can request a different time later if your availability changes.
          </p>
        </section>
      </div>
      <Link href="/candidate/applications">Back to applications</Link>
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "reschedule"
              ? "Reschedule interview?"
              : "Confirm interview time?"
          }
          description={
            selected
              ? `You selected ${new Date(selected.startAt).toLocaleString()}. Confirming reserves this interview slot.`
              : "Review your interview time before continuing."
          }
          confirmLabel={dialog === "reschedule" ? "Reschedule" : "Confirm time"}
          intent="success"
          unavailable="Interview confirmations and rescheduling aren’t available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </section>
  );
}
