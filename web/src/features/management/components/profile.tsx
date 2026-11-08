"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useProfileUpdate } from "@/shared/lib/use-profile-update";
import p from "./profile.module.css";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type {
  CompanyProfile,
  InterviewerProfile,
  InterviewerPage,
} from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
export function TeamProfile({
  interviewer = false,
}: {
  interviewer?: boolean;
}) {
  const params = useSearchParams();
  const editing = interviewer || params.get("edit") === "true";
  const update = useProfileUpdate<CompanyProfile | InterviewerProfile>(
    interviewer ? "/v1/interviewers/me" : "/v1/companies/me",
    !interviewer,
  );
  const company = useApiResource<CompanyProfile>(
    !interviewer ? "/v1/companies/me" : null,
  );
  const profile = useApiResource<InterviewerProfile>(
    interviewer ? "/v1/interviewers/me" : null,
  );
  const team = useApiResource<InterviewerPage>(
    !interviewer ? "/v1/companies/me/interviewers?limit=20&offset=0" : null,
  );
  const error = company.error || profile.error;
  const companyData =
    !interviewer && update.saved
      ? (update.saved as CompanyProfile)
      : company.data;
  const profileData =
    interviewer && update.saved
      ? (update.saved as InterviewerProfile)
      : profile.data;
  const fields: readonly (readonly [string, string, string])[] = [
    ["legalName", "Legal name", "text"],
    ["displayName", "Display name", "text"],
    ["websiteLink", "Website", "url"],
    ["logoUrl", "Logo link", "url"],
    ["contactEmail", "Public contact email", "email"],
    ["contactPhone", "Contact phone", "tel"],
    ["linkedinProfile", "LinkedIn profile", "url"],
  ];
  return (
    <>
      <ProductHeading
        title={
          interviewer
            ? "Edit interviewer profile"
            : editing
              ? "Edit company profile"
              : "Company profile"
        }
        description={interviewer ? "Interviewer workspace" : undefined}
        className={p.banner}
      >
        {!editing && (
          <Link className={p.editLink} href="/employer/profile?edit=true">
            Edit profile
          </Link>
        )}
      </ProductHeading>
      {!companyData && !profileData && (
        <UnavailableNotice>
          Current profile details can’t be loaded yet. You can save changes to
          individual fields; untouched fields are kept.
        </UnavailableNotice>
      )}
      <BusyState loading={company.loading || profile.loading} />
      {error?.status !== 404 && (
        <ApiFeedback
          error={error}
          retry={() => {
            company.reload();
            profile.reload();
          }}
        />
      )}
      {!interviewer && editing && (
        <section className={e.section}>
          <h2>Your interviewers</h2>
          {team.data?.data.map((x) => (
            <p key={x.id}>
              {x.name} · {x.accountStatus}
            </p>
          ))}
          <p>
            Interviewers submit assigned scorecards. Hiring decisions remain
            with the hiring manager.
          </p>
          <Link href="/employer/interviewers">Manage interviewers</Link>
        </section>
      )}
      {!editing && (
        <section className={e.panel} data-figma-node="471:877">
          <h2>Company overview</h2>
          {companyData && (
            <StatusBadge>
              {companyData.approved ? "Approved" : "Awaiting approval"}
            </StatusBadge>
          )}
          <dl className={p.overview}>
            {fields
              .filter(([name]) => name !== "logoUrl")
              .map(([name, label]) => (
                <div key={name}>
                  <dt>{label}</dt>
                  <dd>
                    {String(
                      companyData?.[name as keyof CompanyProfile] ||
                        "Unavailable",
                    )}
                  </dd>
                </div>
              ))}
            <div className={p.description}>
              <dt>Description</dt>
              <dd>
                {companyData?.description ||
                  "Company details aren’t available yet."}
              </dd>
            </div>
          </dl>
        </section>
      )}
      {editing && !company.loading && !profile.loading && (
        <form
          className={e.panel}
          onSubmit={update.onSubmit}
          onChange={update.onChange}
          aria-busy={update.busy}
          key={profileData?.id || companyData?.id || "draft"}
        >
          <h2>{interviewer ? "Profile details" : "Company details"}</h2>
          {companyData && (
            <StatusBadge>
              {companyData.approved ? "Approved" : "Awaiting approval"}
            </StatusBadge>
          )}
          <fieldset className={p.fields} disabled={update.busy}>
            <div className={e.grid}>
              {interviewer ? (
                <label className={e.field}>
                  Name (required)
                  <input
                    name="name"
                    required
                    maxLength={100}
                    autoComplete="name"
                    defaultValue={profileData?.name}
                  />
                </label>
              ) : (
                <>
                  {fields.map(([name, label, type]) => (
                    <label className={e.field} key={name}>
                      {label}
                      {name === "legalName" ? " (required)" : ""}
                      <input
                        name={name}
                        type={type}
                        autoComplete={
                          name === "contactEmail"
                            ? "email"
                            : name === "contactPhone"
                              ? "tel"
                              : name === "legalName"
                                ? "organization"
                                : "off"
                        }
                        spellCheck={!["email", "url"].includes(type)}
                        required={name === "legalName" && !!companyData}
                        maxLength={
                          name === "contactPhone"
                            ? 50
                            : name === "legalName" || name === "displayName"
                              ? 255
                              : undefined
                        }
                        defaultValue={String(
                          companyData?.[name as keyof CompanyProfile] || "",
                        )}
                      />
                    </label>
                  ))}
                  <label className={`${e.field} ${e.full}`}>
                    Description
                    <textarea
                      name="description"
                      maxLength={5000}
                      defaultValue={companyData?.description || ""}
                    />
                  </label>
                </>
              )}
            </div>
          </fieldset>
          <ApiFeedback error={update.error} />
          {update.message && <p role="status">{update.message}</p>}
          <div className={e.actions}>
            <Link
              onNavigate={(event) => {
                if (
                  update.dirty &&
                  !window.confirm("Discard your unsaved profile changes?")
                )
                  event.preventDefault();
              }}
              href={
                interviewer ? "/interviewer/interviews" : "/employer/profile"
              }
            >
              Cancel
            </Link>
            <button className={p.save} type="submit" disabled={update.busy}>
              {update.busy
                ? "Saving…"
                : interviewer
                  ? "Save profile"
                  : "Save company profile"}
            </button>
          </div>
        </form>
      )}
      {!interviewer && !editing && (
        <section className={e.section}>
          <h2>Hiring team</h2>
          {team.data?.data.map((x) => (
            <p key={x.id}>
              {x.name} · {x.accountStatus}
            </p>
          ))}
          <p>
            Interviewers submit assigned scorecards. Hiring decisions remain
            with the hiring manager.
          </p>
          <Link href="/employer/interviewers">Manage interviewers</Link>
        </section>
      )}
    </>
  );
}
