"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { InterviewPage } from "@/shared/lib/contract.types";
import type { TimeSlot } from "@/features/cand/types/product.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

export function InterviewAvailability() {
  const params = useSearchParams();
  const [id, setId] = useState(params.get("interview") || "");
  const interviews = useApiResource<InterviewPage>(
    "/v1/interviews?limit=100&offset=0",
  );
  const slots = useApiResource<TimeSlot[]>(
    id ? `/v1/interviews/${resourceId(id)}/availability` : null,
  );
  const [drafts, setDrafts] = useState<{ startAt: string; endAt: string }[]>([
    { startAt: "", endAt: "" },
  ]);
  return (
    <>
      <ProductHeading
        title="Your available times"
        description="Interview availability"
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Availability updates aren’t available yet. Entered times remain on this
        screen.
      </UnavailableNotice>
      <label className={e.field}>
        Interview
        <select value={id} onChange={(event) => setId(event.target.value)}>
          <option value="">Choose an interview</option>
          {interviews.data?.data.map((row) => (
            <option key={row.id} value={row.id}>
              {new Date(row.scheduledStartAt).toLocaleString()} · {row.id}
            </option>
          ))}
        </select>
      </label>
      <BusyState loading={slots.loading} />
      {slots.error?.status !== 404 && (
        <ApiFeedback error={slots.error} retry={slots.reload} />
      )}
      <section className={e.section}>
        <h2>Your available times</h2>
        {slots.data?.map((slot) => (
          <p key={slot.startAt}>
            {new Date(slot.startAt).toLocaleString()} –{" "}
            {new Date(slot.endAt).toLocaleString()}
          </p>
        ))}
        {drafts.map((draft, index) => (
          <div className={e.grid} key={index}>
            {(["startAt", "endAt"] as const).map((key, i) => (
              <label className={e.field} key={key}>
                {i ? "Ends" : "Starts"}
                <input
                  type="datetime-local"
                  value={draft[key]}
                  onChange={(event) =>
                    setDrafts((rows) =>
                      rows.map((row, j) =>
                        j === index
                          ? { ...row, [key]: event.target.value }
                          : row,
                      ),
                    )
                  }
                />
              </label>
            ))}
          </div>
        ))}
        <div className={e.actions}>
          <button
            onClick={() =>
              setDrafts((rows) => [...rows, { startAt: "", endAt: "" }])
            }
          >
            Add time interval
          </button>
        </div>
      </section>
      <section className={e.section}>
        <h2>Scheduling</h2>
        <p>
          Candidates can choose a slot that matches the assigned interviewers’
          common availability.
        </p>
      </section>
      <div className={e.actions}>
        <Link href="/interviewer/interviews">Cancel</Link>
        <button className={e.primary} disabled>
          Save availability
        </button>
      </div>
      <WorkspaceFooter />
    </>
  );
}
