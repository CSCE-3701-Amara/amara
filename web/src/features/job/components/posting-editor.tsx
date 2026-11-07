"use client";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useRef, useState, type FormEvent } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import { PipelineEditor } from "@/features/invw/components/pipeline-editor";
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
function PostingForm({ initial }: { initial?: Posting }) {
  const router = useRouter();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [message, setMessage] = useState("");
  const pending = useRef(false);
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
    setMessage("");
    if (
      body.salary_min !== null &&
      body.salary_max !== null &&
      body.salary_min! > body.salary_max!
    ) {
      setError(new Error("Minimum salary cannot exceed maximum salary."));
      return;
    }
    if (
      body.application_opening &&
      body.application_deadline &&
      body.application_deadline <= body.application_opening
    ) {
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
      else setMessage("Posting saved.");
    } catch (cause) {
      setError(cause as Error);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return (
    <form onSubmit={submit} className={`${s.card} ${s.form}`} aria-busy={busy}>
      <p>
        Title, seniority level, employment type, and required skills are
        required.
      </p>
      <label className={s.field}>
        Job title
        <input name="title" required defaultValue={initial?.title} />
      </label>
      <label className={s.field}>
        Description
        <textarea
          name="description"
          defaultValue={initial?.description ?? ""}
        />
      </label>
      <div className={s.grid}>
        <label className={s.field}>
          Seniority level
          <input
            name="seniority_level"
            required
            defaultValue={initial?.seniority_level}
          />
        </label>
        <label className={s.field}>
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
        <label className={s.field}>
          Location
          <input name="location" defaultValue={initial?.location ?? ""} />
        </label>
        <label className={s.field}>
          Language
          <input name="language" defaultValue={initial?.language ?? ""} />
        </label>
        <label className={s.field}>
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
          <label key={name} className={s.field}>
            {label}
            <input
              name={name}
              type="number"
              min={min}
              max="2147483647"
              defaultValue={initial?.[name] ?? ""}
            />
          </label>
        ))}
        <label className={s.field}>
          Applications open
          <input
            name="application_opening"
            type="datetime-local"
            defaultValue={localDate(initial?.application_opening)}
          />
        </label>
        <label className={s.field}>
          Application deadline
          <input
            name="application_deadline"
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
      <label className={s.field}>
        Required skills (Amarat), one per line
        <textarea
          name="amarat"
          required
          defaultValue={initial?.amarat.map((item) => item.title).join("\n")}
          placeholder={"React\nTypeScript"}
        />
      </label>
      <div className={s.grid}>
        <label className={s.field}>
          Must-have requirements, one per line
          <textarea
            name="must_have"
            defaultValue={initial?.requirements
              .filter((item) => item.requirement_type === "must_have")
              .map((item) => item.text)
              .join("\n")}
          />
        </label>
        <label className={s.field}>
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
      <ApiFeedback error={error} />
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
  const router = useRouter();
  async function act(action: "publish" | "republish" | "archive" | "delete") {
    if (busy) return;
    if (
      action === "delete" &&
      !window.confirm(`Delete “${posting.title}”? This cannot be undone.`)
    )
      return;
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
              act(posting.status === "archived" ? "republish" : "publish")
            }
          >
            {posting.status === "archived"
              ? "Republish posting"
              : "Publish posting"}
          </button>
        )}
        {posting.status === "published" && (
          <button disabled={busy} onClick={() => act("archive")}>
            Archive posting
          </button>
        )}
        <button
          className={s.danger}
          disabled={busy}
          onClick={() => act("delete")}
        >
          Delete posting
        </button>
      </div>
      <ApiFeedback error={error} />
    </>
  );
}
export function PostingEditor({ id }: { id?: string }) {
  const posting = useApiResource<Posting>(
    id ? `/v1/jobs/${resourceId(id)}` : null,
  );
  return (
    <>
      <Link href="/employer/jobs">Back to postings</Link>
      <h1>{id ? "Manage job posting" : "Create a job posting"}</h1>
      <BusyState loading={posting.loading} />
      <ApiFeedback error={posting.error} retry={posting.reload} />
      {(!id || posting.data) && (
        <PostingForm key={posting.data?.id || "new"} initial={posting.data} />
      )}
      {posting.data && (
        <PostingActions posting={posting.data} reload={posting.reload} />
      )}
      {id && <PipelineEditor postingId={id} />}
    </>
  );
}
