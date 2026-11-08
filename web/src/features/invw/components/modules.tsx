"use client";
import Link from "next/link";
import { useState, type FormEvent } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import {
  createModule,
  updateModule,
  publishModule,
  copyModule,
} from "../api/invw.api";
import type {
  ModuleDefinition,
  ModuleSummary,
  AssessmentModule,
} from "../types/invw.types";
function ModuleForm({
  module,
  saved,
}: {
  module?: AssessmentModule;
  saved: () => void;
}) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy) return;
    const data = new FormData(event.currentTarget);
    const form = event.currentTarget;
    const name = String(data.get("name") || "").trim();
    const prompt = String(data.get("prompt") || "").trim();
    if (!name || ((!module || module.type === "behavioral") && !prompt)) {
      setError(new Error("Enter a module name and a question."));
      return;
    }
    const definition: ModuleDefinition = {
      name,
      description: String(data.get("description") || "").trim(),
      type: module?.type || "behavioral",
      ...(data.get("active_minutes")
        ? { active_minutes: Number(data.get("active_minutes")) }
        : {}),
      ...(module
        ? {
            threshold: module.threshold,
            behavioral: module.behavioral,
            minimum_tests_passed: module.minimum_tests_passed,
            amarat: module.amarat,
            rubric: module.rubric,
            library: module.library,
            questions: module.questions,
          }
        : {}),
    };
    if (!module || module.type === "behavioral")
      definition.questions = [
        {
          ...(module?.questions?.[0] || { format: "short_text" as const }),
          prompt,
        },
        ...(module?.questions?.slice(1) || []),
      ];
    setBusy(true);
    setError(undefined);
    try {
      if (module) await updateModule(module.id, definition);
      else await createModule(definition);
      form.reset();
      saved();
    } catch (cause) {
      setError(cause as Error);
    } finally {
      setBusy(false);
    }
  }
  return (
    <form
      onSubmit={submit}
      className={s.form}
      aria-label={module ? "Edit assessment" : "Create assessment"}
      aria-busy={busy}
    >
      <p>
        {module
          ? "Update this draft’s details."
          : "Create a behavioral assessment with a written-response question."}
      </p>
      <label className={s.field}>
        Assessment name
        <input name="name" required defaultValue={module?.name} />
      </label>
      <label className={s.field}>
        Assessment description
        <textarea name="description" defaultValue={module?.description} />
      </label>
      {module?.type !== "recorded_video" && (
        <label className={s.field}>
          Time limit in minutes (optional)
          <input
            name="active_minutes"
            type="number"
            min="1"
            max="2147483647"
            defaultValue={module?.active_minutes}
          />
        </label>
      )}
      {(!module || module.type === "behavioral") && (
        <label className={s.field}>
          Question
          <textarea
            name="prompt"
            required
            defaultValue={module?.questions?.[0]?.prompt}
          />
        </label>
      )}
      <ApiFeedback error={error} />
      <button type="submit" className={s.primary} disabled={busy}>
        {busy
          ? "Saving…"
          : module
            ? "Save assessment"
            : "Create assessment draft"}
      </button>
    </form>
  );
}
function ModuleDetails({ id, saved }: { id: string; saved: () => void }) {
  const resource = useApiResource<AssessmentModule>(
    `/v1/modules/${resourceId(id)}`,
  );
  return (
    <>
      <BusyState loading={resource.loading} />
      <ApiFeedback error={resource.error} retry={resource.reload} />
      {resource.data &&
        (resource.data.published || resource.data.copied_from_module_id ? (
          <>
            <p>
              This assessment is locked. Its questions are preserved for use in
              your pipeline.
            </p>
            {resource.data.questions?.map((question, index) => (
              <p key={index}>
                {index + 1}. {question.prompt}
              </p>
            ))}
          </>
        ) : (
          <ModuleForm module={resource.data} saved={saved} />
        ))}
    </>
  );
}
export function Modules() {
  const owned = useApiResource<ModuleSummary[] | null>("/v1/modules");
  const library = useApiResource<ModuleSummary[] | null>("/v1/modules/library");
  const [selected, setSelected] = useState<string>();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [message, setMessage] = useState("");
  function saved() {
    owned.reload();
    library.reload();
    setSelected(undefined);
    setMessage("Assessment saved.");
  }
  async function action(id: string, type: "publish" | "copy") {
    if (busy) return;
    setBusy(true);
    setError(undefined);
    setMessage("");
    try {
      if (type === "copy") await copyModule(id);
      else await publishModule(id);
      owned.reload();
      library.reload();
      setSelected(undefined);
      setMessage(
        type === "copy"
          ? "Assessment copied into your workspace."
          : "Assessment published.",
      );
    } catch (cause) {
      setError(cause as Error);
    } finally {
      setBusy(false);
    }
  }
  return (
    <>
      <h1>Assessments & pipelines</h1>
      <p className={s.muted}>
        Build reusable assessments. Open a job posting to configure its ordered
        pipeline stages.
      </p>
      <Link href="/employer/jobs">Choose a job posting</Link>
      <section className={s.card}>
        <h2>Create an assessment</h2>
        <ModuleForm saved={saved} />
      </section>
      <ApiFeedback error={error} />
      {message && (
        <p role="status" className={s.success}>
          {message}
        </p>
      )}
      <h2>Your assessments</h2>
      <BusyState loading={owned.loading} />
      <ApiFeedback error={owned.error} retry={owned.reload} />
      {owned.data?.length === 0 && (
        <p>No assessments yet. Create your first assessment above.</p>
      )}
      <div className={s.grid}>
        {owned.data?.map((module) => (
          <article className={s.card} key={module.id}>
            <h3>{module.name}</h3>
            <span className={s.badge}>
              {module.published ? "Published" : "Draft"}
            </span>
            <p>{module.description}</p>
            <div className={s.actions}>
              <button
                onClick={() =>
                  setSelected(selected === module.id ? undefined : module.id)
                }
                aria-expanded={selected === module.id}
              >
                View assessment
              </button>
              {!module.published && (
                <button
                  disabled={busy}
                  onClick={() => action(module.id, "publish")}
                >
                  Publish assessment
                </button>
              )}
            </div>
            {selected === module.id && (
              <ModuleDetails id={module.id} saved={saved} />
            )}
          </article>
        ))}
      </div>
      <h2>Published assessment library</h2>
      <BusyState loading={library.loading} />
      <ApiFeedback error={library.error} retry={library.reload} />
      {library.data?.length === 0 && <p>The assessment library is empty.</p>}
      <div className={s.grid}>
        {library.data?.map((module) => (
          <article className={s.card} key={module.id}>
            <h3>{module.name}</h3>
            <p>{module.description}</p>
            <button disabled={busy} onClick={() => action(module.id, "copy")}>
              Copy to workspace
            </button>
          </article>
        ))}
      </div>
    </>
  );
}
