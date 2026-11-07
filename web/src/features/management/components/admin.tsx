"use client";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type {
  CompanyRegistrationPage,
  ManagedAccountPage,
} from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
/** Administrator projections only; no credentials or private candidate records are inferred. */
export function Administration({ accounts = false }: { accounts?: boolean }) {
  const [dialog, setDialog] = useState<string | null>(null);
  const [selected, setSelected] = useState("");
  const [reason, setReason] = useState("");
  const [offset, setOffset] = useState(0);
  const requests = useApiResource<CompanyRegistrationPage>(
    !accounts
      ? `/v1/admin/company-registration-requests?limit=20&offset=${offset}`
      : null,
  );
  const managed = useApiResource<ManagedAccountPage>(
    accounts ? `/v1/admin/accounts?limit=20&offset=${offset}` : null,
  );
  const rows = accounts
    ? managed.data?.data.map((x) => ({
        id: x.id,
        name: x.name,
        type: x.accountType,
        status: x.banned
          ? "Banned"
          : x.deactivatedAt
            ? "Deactivated"
            : x.accountStatus || "Unknown",
      }))
    : requests.data?.data.map((x) => ({
        id: x.id,
        name: x.companyName,
        type: "Registration",
        status: x.status,
      }));
  const active = rows?.find((x) => x.id === selected);
  const error = managed.error || requests.error;
  const page = accounts ? managed.data?.pagination : requests.data?.pagination;
  const titles: Record<string, string> = {
    approve: "Approve company registration?",
    reject: "Reject company registration?",
    ban: "Ban this account?",
    unban: "Unban this account?",
    deactivate: "Deactivate this account?",
    reactivate: "Reactivate this account?",
    delete: "Delete this account?",
  };
  return (
    <>
      <ProductHeading
        title={accounts ? "Manage accounts" : "Company registration review"}
        description="System administration"
        gradient
        compact
        artwork={{ src: "/figma/admin.svg", width: 200.16, height: 138.96 }}
      />
      <UnavailableNotice>
        Administration is not available yet. Decisions and account changes
        cannot be saved.
      </UnavailableNotice>
      <BusyState loading={requests.loading || managed.loading} />
      {error?.status !== 404 && (
        <ApiFeedback
          error={error}
          retry={() => {
            requests.reload();
            managed.reload();
          }}
        />
      )}
      <div className={e.tableWrap}>
        <table className={e.table}>
          <thead>
            <tr>
              <th scope="col">{accounts ? "Account" : "Request"}</th>
              <th scope="col">{accounts ? "Type" : "Company"}</th>
              <th scope="col">Status</th>
            </tr>
          </thead>
          <tbody>
            {rows?.map((row) => (
              <tr key={row.id}>
                <td>
                  <button
                    onClick={() => setSelected(row.id)}
                    aria-pressed={selected === row.id}
                  >
                    {accounts ? row.name : row.id}
                  </button>
                </td>
                <td>{accounts ? row.type : row.name}</td>
                <td>{row.status.replaceAll("_", " ")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {!rows?.length && !requests.loading && !managed.loading && (
        <EmptyState
          title="No records available"
          description="Records will appear here when administration becomes available."
        />
      )}
      <section className={e.section}>
        <h2>{accounts ? "Selected account" : "Selected request"}</h2>
        <p>
          {active
            ? `${active.name} · ${active.status}`
            : "Select a record to review its details."}
        </p>
        <label className={e.field}>
          {accounts ? "Action reason" : "Decision reason"}
          <textarea
            value={reason}
            onChange={(event) => setReason(event.target.value)}
            maxLength={2000}
          />
        </label>
        {accounts && (
          <>
            <h3>Ban effect</h3>
            <p>Banning a company archives its published postings.</p>
          </>
        )}
        <div className={e.actions}>
          {(accounts
            ? ["ban", "unban", "deactivate", "reactivate"]
            : ["reject", "approve"]
          ).map((action) => (
            <button
              key={action}
              disabled={!active}
              onClick={() => setDialog(action)}
            >
              {action[0].toUpperCase() + action.slice(1)}{" "}
              {accounts ? "account" : "request"}
            </button>
          ))}
        </div>
      </section>
      {page && (
        <div className={e.actions}>
          <button
            disabled={!offset}
            onClick={() => setOffset((x) => Math.max(0, x - 20))}
          >
            Previous page
          </button>
          <span>
            {offset + 1}–{Math.min(offset + 20, page.total)} of {page.total}
          </span>
          <button
            disabled={offset + 20 >= page.total}
            onClick={() => setOffset((x) => x + 20)}
          >
            Next page
          </button>
        </div>
      )}
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={titles[dialog] || "Review this action?"}
          description={
            active
              ? `This changes the status of ${active.name}. Review the reason before continuing.`
              : "Select and review a record before making this decision."
          }
          confirmLabel={dialog[0].toUpperCase() + dialog.slice(1)}
          intent={
            ["reject", "ban", "delete"].includes(dialog) ? "danger" : "success"
          }
          reason
          unavailable="This administration action is not available yet."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
