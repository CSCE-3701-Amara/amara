"use client";
import Link from "next/link";
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
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
export function TeamProfile({
  interviewer = false,
}: {
  interviewer?: boolean;
}) {
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
          interviewer ? "Edit interviewer profile" : "Edit company profile"
        }
        description={interviewer ? "Interviewer workspace" : undefined}
        compact
        artwork={{
          src: interviewer ? "/figma/hiring-dag.svg" : "/figma/company.svg",
          width: 210,
          height: 150,
        }}
      />
      <UnavailableNotice>
        Profile updates aren’t available yet. Changes entered here remain on
        this screen.
      </UnavailableNotice>
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
      {!interviewer && (
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
      <form
        className={e.panel}
        onSubmit={(event) => event.preventDefault()}
        key={profile.data?.id || company.data?.id || "draft"}
      >
        <h2>{interviewer ? "Profile details" : "Company details"}</h2>
        {company.data && (
          <StatusBadge>
            {company.data.approved ? "Approved" : "Awaiting approval"}
          </StatusBadge>
        )}
        <div className={e.grid}>
          {interviewer ? (
            <label className={e.field}>
              Name (required)
              <input name="name" required defaultValue={profile.data?.name} />
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
                    required={name === "legalName"}
                    defaultValue={String(
                      company.data?.[name as keyof CompanyProfile] || "",
                    )}
                  />
                </label>
              ))}
              <label className={`${e.field} ${e.full}`}>
                Description
                <textarea
                  name="description"
                  maxLength={5000}
                  defaultValue={company.data?.description || ""}
                />
              </label>
            </>
          )}
        </div>
        <div className={e.actions}>
          <Link
            href={interviewer ? "/interviewer/interviews" : "/employer/jobs"}
          >
            Cancel
          </Link>
          <button className={e.primary} type="submit" disabled>
            Save {interviewer ? "profile" : "company profile"}
          </button>
        </div>
      </form>
      <WorkspaceFooter />
    </>
  );
}
