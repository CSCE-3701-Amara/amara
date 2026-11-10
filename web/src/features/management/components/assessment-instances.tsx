"use client";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { AssessmentInstance } from "@/features/cand/types/product.types";
import {
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";

export function AssessmentInstances() {
  const params = useSearchParams();
  const [application, setApplication] = useState(
    params.get("application") || "",
  );
  const [query, setQuery] = useState(application);
  const [selected, setSelected] = useState("");
  const resource = useApiResource<AssessmentInstance[]>(null);
  const active = resource.data?.find((row) => row.id === selected);
  return (
    <>
      <ProductHeading
        title="Assessment instances & extensions"
        description="Individual candidate + stage"
        compact
        className={e.titleBanner}
      />
      <UnavailableNotice>
        Hiring-team assessment timing and extensions require backend support.
        The contract’s instance list is candidate scoped, so this workspace does
        not call it with a hiring-manager session.
      </UnavailableNotice>
      <form
        onSubmit={(event) => {
          event.preventDefault();
          setApplication(query.trim());
        }}
        className={e.actions}
      >
        <label className={e.field}>
          Application ID
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </label>
        <button type="submit" disabled>
          Search instances
        </button>
      </form>
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">Candidate / instance</th>
              <th scope="col">Stage</th>
              <th scope="col">Instance state</th>
              <th scope="col">Final score</th>
            </tr>
          </thead>
          <tbody>
            {resource.data?.map((row) => (
              <tr key={row.id}>
                <td>
                  <button onClick={() => setSelected(row.id)}>{row.id}</button>
                </td>
                <td>{row.stageId}</td>
                <td>{row.status}</td>
                <td>Unavailable</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className={e.split}>
        <section className={e.panel}>
          <h2>Selected instance</h2>
          <p>
            {active?.id ||
              "Select an instance when authorized data becomes available."}
          </p>
          <p>
            Active time and completion windows are separate. Provider failures
            and pending grading never become a zero score automatically.
          </p>
        </section>
        <section className={e.panel}>
          <h2>Give additional time</h2>
          <div className={e.grid}>
            <label className={e.field}>
              Additional active time (minutes)
              <input type="number" min={1} />
            </label>
            <label className={e.field}>
              Individual completion deadline
              <input type="datetime-local" />
            </label>
            <label className={`${e.field} ${e.full}`}>
              Required reason
              <textarea required />
            </label>
          </div>
          <div className={e.actions}>
            <button disabled>Save extension</button>
          </div>
        </section>
      </div>
      <WorkspaceFooter />
    </>
  );
}
