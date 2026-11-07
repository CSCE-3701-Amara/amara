"use client";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { PostingPreview } from "./posting-preview";
import { useRouter } from "next/navigation";
import { useId, useRef, useState, type FormEvent } from "react";
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
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import {
  createPosting,
  updatePosting,
  transitionPosting,
  deletePosting,
} from "../api/job.api";
import type { EmploymentType, Posting, PostingInput } from "../types/job.types";

function localDate(value?: string | null) {
  if (!value) return "";
  const date = new Date(value);
  return new Date(date.getTime() - date.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16);
}
function PostingForm({
  initial,
  onSaved,
}: {
  initial?: Posting;
  onSaved?: () => void;
}) {
  const router = useRouter();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [message, setMessage] = useState("");
  const pending = useRef(false);
  const errorId = useId();
  const [invalidField, setInvalidField] = useState("");
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending.current) return;
    const data = new FormData(event.currentTarget);
    const text = (name: string) => String(data.get(name) || "").trim();
    const number = (name: string) =>
      text(name) === "" ? null : Number(text(name));
    const date = (name: string) =>
      text(name) ? new Date(text(name)).toISOString() : null;
    const lines = (name: string) =>
      text(name)
        .split("\n")
        .map((line) => line.trim())
        .filter(Boolean);
    const body: PostingInput = {
      title: text("title"),
      description: text("description") || null,
      seniority_level: text("seniority_level"),
      employment_type: text("employment_type") as EmploymentType,
      language: text("language") || null,
      location: text("location") || null,
      currency: text("currency") || null,
      salary_min: number("salary_min"),
      salary_max: number("salary_max"),
      salary_visible: data.has("salary_visible"),
      open_positions: number("open_positions"),
      max_applications: number("max_applications"),
      application_opening: date("application_opening"),
      application_deadline: date("application_deadline"),
      requirements: [
        ...lines("must_have").map((value) => ({
          text: value,
          requirement_type: "must_have" as const,
        })),
        ...lines("nice_to_have").map((value) => ({
          text: value,
          requirement_type: "nice_to_have" as const,
        })),
      ],
      amarat: lines("amarat").map((title) => ({
        title,
        description: initial?.amarat.find((item) => item.title === title)
          ?.description,
      })),
    };
    setError(undefined);
    setInvalidField("");
    setMessage("");
    if (
      body.salary_min !== null &&
      body.salary_max !== null &&
      body.salary_min! > body.salary_max!
    ) {
      setError(new Error("Minimum salary cannot exceed maximum salary."));
      setInvalidField("salary_max");
      event.currentTarget
        .querySelector<HTMLInputElement>('[name="salary_max"]')
        ?.focus();
      return;
    }
    if (
      body.application_opening &&
      body.application_deadline &&
      body.application_deadline <= body.application_opening
    ) {
      setInvalidField("application_deadline");
      event.currentTarget
        .querySelector<HTMLInputElement>('[name="application_deadline"]')
        ?.focus();
      setError(
        new Error("The application deadline must be after the opening time."),
      );
      return;
    }
    if (!body.title || !body.seniority_level || !body.amarat.length) {
      setError(
        new Error(
          "Enter a title, seniority level, and at least one required skill.",
        ),
      );
      return;
    }
    pending.current = true;
    setBusy(true);
    try {
      const saved = initial
        ? await updatePosting(initial.id, body)
        : await createPosting(body);
      if (!initial) router.push(`/employer/jobs/${resourceId(saved.id)}`);
      else {
        setMessage("Posting saved.");
        onSaved?.();
      }
    } catch (cause) {
      setError(cause as Error);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return (
    <form
      onSubmit={submit}
      className={`${e.section} ${s.form}`}
      aria-busy={busy}
    >
      <p>
        Title, seniority level, employment type, and required skills are
        required.
      </p>
      <h2>Role details</h2>
      <label className={e.field}>
        Job title
        <input name="title" required defaultValue={initial?.title} />
      </label>
      <label className={e.field}>
        Description
        <textarea
          name="description"
          defaultValue={initial?.description ?? ""}
        />
      </label>
      <div className={e.grid}>
        <label className={e.field}>
          Seniority level
          <input
            name="seniority_level"
            required
            defaultValue={initial?.seniority_level}
          />
        </label>
        <label className={e.field}>
          Employment type
          <select
            name="employment_type"
            defaultValue={initial?.employment_type || "on-site full-time"}
          >
            {[
              "on-site full-time",
              "on-site part-time",
              "remote full-time",
              "remote part-time",
            ].map((value) => (
              <option key={value}>{value}</option>
            ))}
          </select>
        </label>
        <label className={e.field}>
          Location
          <input name="location" defaultValue={initial?.location ?? ""} />
        </label>
        <label className={e.field}>
          Language
          <input name="language" defaultValue={initial?.language ?? ""} />
        </label>
        <label className={e.field}>
          Currency
          <input name="currency" defaultValue={initial?.currency ?? ""} />
        </label>
        {(
          [
            ["salary_min", "Minimum salary", 0],
            ["salary_max", "Maximum salary", 0],
            ["open_positions", "Open positions", 1],
            ["max_applications", "Application limit", 1],
          ] as const
        ).map(([name, label, min]) => (
          <label key={name} className={e.field}>
            {label}
            <input
              name={name}
              aria-invalid={invalidField === name || undefined}
              aria-describedby={invalidField === name ? errorId : undefined}
              type="number"
              min={min}
              max="2147483647"
              defaultValue={initial?.[name] ?? ""}
            />
          </label>
        ))}
        <label className={e.field}>
          Applications open
          <input
            name="application_opening"
            type="datetime-local"
            defaultValue={localDate(initial?.application_opening)}
          />
        </label>
        <label className={e.field}>
          Application deadline
          <input
            name="application_deadline"
            aria-invalid={invalidField === "application_deadline" || undefined}
            aria-describedby={
              invalidField === "application_deadline" ? errorId : undefined
            }
            type="datetime-local"
            defaultValue={localDate(initial?.application_deadline)}
          />
        </label>
      </div>
      <label className={s.check}>
        <input
          name="salary_visible"
          type="checkbox"
          defaultChecked={initial?.salary_visible ?? false}
        />
        Show salary to candidates
      </label>
      <label className={e.field}>
        Required skills (amarat), one per line
        <textarea
          name="amarat"
          required
          defaultValue={initial?.amarat.map((item) => item.title).join("\n")}
          placeholder={"React\nTypeScript"}
        />
      </label>
      <div className={e.grid}>
        <label className={e.field}>
          Must-have requirements, one per line
          <textarea
            name="must_have"
            defaultValue={initial?.requirements
              .filter((item) => item.requirement_type === "must_have")
              .map((item) => item.text)
              .join("\n")}
          />
        </label>
        <label className={e.field}>
          Nice-to-have requirements, one per line
          <textarea
            name="nice_to_have"
            defaultValue={initial?.requirements
              .filter((item) => item.requirement_type === "nice_to_have")
              .map((item) => item.text)
              .join("\n")}
          />
        </label>
      </div>
      {invalidField && error ? (
        <p id={errorId} className={s.error} role="alert">
          {error.message}
        </p>
      ) : (
        <ApiFeedback error={error} />
      )}
      {message && (
        <p className={s.success} role="status">
          {message}
        </p>
      )}
      <button type="submit" className={s.primary} disabled={busy}>
        {busy ? "Saving…" : initial ? "Save posting" : "Create draft"}
      </button>
    </form>
  );
}
function PostingActions({
  posting,
  reload,
}: {
  posting: Posting;
  reload: () => void;
}) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [confirm, setConfirm] = useState<
    "publish" | "republish" | "archive" | "delete"
  >();
  const router = useRouter();
  async function act(action: "publish" | "republish" | "archive" | "delete") {
    if (busy) return;
    setBusy(true);
    setError(undefined);
    try {
      if (action === "delete") {
        await deletePosting(posting.id);
        router.push("/employer/jobs");
      } else {
        await transitionPosting(posting.id, action);
        reload();
      }
    } catch (cause) {
      setError(cause as Error);
      throw cause;
    } finally {
      setBusy(false);
    }
  }
  return (
    <>
      <div className={s.actions}>
        {posting.status !== "published" && (
          <button
            className={s.primary}
            disabled={busy}
            onClick={() =>
              setConfirm(
                posting.status === "archived" ? "republish" : "publish",
              )
            }
          >
            {posting.status === "archived"
              ? "Republish posting"
              : "Publish posting"}
          </button>
        )}
        {posting.status === "published" && (
          <button disabled={busy} onClick={() => setConfirm("archive")}>
            Archive posting
          </button>
        )}
        <button
          className={s.danger}
          disabled={busy}
          onClick={() => setConfirm("delete")}
        >
          Delete posting
        </button>
      </div>
      <ApiFeedback error={error} />
      {confirm && (
        <ConfirmDialog
          title={`${confirm[0].toUpperCase() + confirm.slice(1)} this posting?`}
          description={`Review “${posting.title}” before continuing. This changes the posting’s availability to candidates.`}
          confirmLabel={confirm === "delete" ? "Delete posting" : "Confirm"}
          intent={confirm === "delete" ? "danger" : "warning"}
          onClose={() => setConfirm(undefined)}
          onConfirm={async () => {
            await act(confirm);
            setConfirm(undefined);
          }}
        />
      )}
    </>
  );
}
export function PostingEditor({ id }: { id?: string }) {
  const params = useSearchParams();
  const posting = useApiResource<Posting>(
    id ? `/v1/jobs/${resourceId(id)}` : null,
  );
  if (id && params.get("state") === "publication-review")
    return <PostingPreview id={id} review />;
  return (
    <>
      <Link href="/employer/jobs">Back to postings</Link>
      <ProductHeading
        title={id ? "Edit job posting" : "Create job posting"}
        description={
          posting.data
            ? `${posting.data.title} · ${posting.data.status}`
            : "Build a clear role and choose the skills you need."
        }
        artwork={{
          src: "/figma/posting-editor.svg",
          width: 261.6,
          height: 180,
        }}
      />
      <BusyState loading={posting.loading} />
      <ApiFeedback error={posting.error} retry={posting.reload} />
      {(!id || posting.data) && (
        <PostingForm key={posting.data?.id || "new"} initial={posting.data} />
      )}
      {id && (
        <div className={e.actions}>
          <Link href={`/employer/jobs/${resourceId(id)}/preview`}>Preview</Link>
          <Link
            href={`/employer/jobs/${resourceId(id)}?state=publication-review`}
          >
            Review publication
          </Link>
        </div>
      )}
      {posting.data && (
        <PostingActions posting={posting.data} reload={posting.reload} />
      )}
      <section className={e.section}>
        <h2>Hiring pipeline</h2>
        <p>Set up assessments and interviews in Pipeline studio.</p>
        <Link
          href={
            id ? `/employer/pipelines/${resourceId(id)}` : "/employer/pipelines"
          }
          className={s.button}
        >
          Open Pipeline studio
        </Link>
      </section>
      <WorkspaceFooter />
    </>
  );
}
