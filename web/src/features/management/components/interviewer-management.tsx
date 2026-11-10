"use client";
import Link from "next/link";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { InterviewerPage } from "@/shared/lib/contract.types";
import {
  EmptyState,
  ProductHeading,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import { PlannedRead } from "./planned-read";

export function InterviewerManagement() {
  const team = useApiResource<InterviewerPage>(
    "/v1/companies/me/interviewers?limit=20&offset=0",
  );
  return (
    <>
      <ProductHeading
        title="Interviewer management"
        description="Manage company interviewers, availability, and assigned evaluations"
      />
      <UnavailableNotice>
        Interviewer invitations and assignment changes aren’t available yet.
      </UnavailableNotice>
      <PlannedRead
        loading={team.loading}
        error={team.error}
        retry={team.reload}
      />
      {team.data?.data.map((person) => (
        <article className={e.panel} key={person.id}>
          <h2>{person.name}</h2>
          <p>{person.accountStatus}</p>
        </article>
      ))}
      {!team.loading && !team.data?.data.length && (
        <EmptyState
          title="No interviewers available"
          description="Your company interviewers will appear here when available."
        />
      )}
      <div className={e.actions}>
        <button disabled>Invite interviewer</button>
        <button disabled>Assign interviews</button>
        <Link href="/employer/interviews">Manage interviews</Link>
      </div>
      <WorkspaceFooter />
    </>
  );
}
