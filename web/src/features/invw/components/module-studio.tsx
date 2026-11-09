"use client";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useId, useRef, useState, type FormEvent } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  EmptyState,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import e from "@/shared/ui/editor.module.css";
import {
  createModule,
  updateModule,
  publishModule,
  copyModule,
} from "../api/invw.api";
import type {
  Amara,
  AssessmentModule,
  ModuleDefinition,
  ModuleSummary,
  Question,
  LibraryProblem,
} from "../types/invw.types";
const kinds: Record<ModuleDefinition["type"], string> = {
  coding: "Coding assessment",
  programming_language: "Question assessment",
  behavioral: "Behavioral assessment",
  recorded_video: "Video response",
};
export function ModuleLibrary({
  library = false,
  problems = false,
}: {
  library?: boolean;
  problems?: boolean;
}) {
  const params = useSearchParams();
  const router = useRouter();
  const [selected, setSelected] = useState(params.get("module") || "");
  const [confirm, setConfirm] = useState(false);
  const [error, setError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const pending = useRef(false);
  const modules = useApiResource<ModuleSummary[] | null>(
    !problems ? (library ? "/v1/modules/library" : "/v1/modules") : null,
  );
  const bank = useApiResource<LibraryProblem[] | null>(
    problems ? "/v1/modules/library/problems" : null,
  );
  async function copy() {
    if (!selected || pending.current) return;
    pending.current = true;
    setBusy(true);
    setError(undefined);
    try {
      const result = await copyModule(selected);
      router.push(`/employer/modules/${resourceId(result.id)}/questions`);
    } catch (cause) {
      setError(cause as Error);
      throw cause;
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  const chosen = modules.data?.find((x) => x.id === selected);
  return (
    <>
      <ProductHeading
        title={
          problems
            ? "Coding problem library"
            : library
              ? "Copy module to company"
              : "Company assessment modules"
        }
        description={
          problems
            ? "Browse approved problems and reuse them in company modules"
            : "Assessment module library"
        }
        compact
        artwork={
          !problems
            ? { src: "/figma/hiring-dag.svg", width: 200.16, height: 138.96 }
            : undefined
        }
      />
      <BusyState loading={modules.loading || bank.loading} />
      <ApiFeedback
        error={modules.error || bank.error || error}
        retry={() => {
          modules.reload();
          bank.reload();
        }}
      />
      {problems ? (
        <>
          {bank.data?.map((problem) => (
            <article
              className={e.panel}
              key={`${problem.module_id}-${problem.position}`}
            >
              <label className={e.check}>
                <input
                  type="radio"
                  name="problem"
                  checked={
                    selected === `${problem.module_id}:${problem.position}`
                  }
                  onChange={() =>
                    setSelected(`${problem.module_id}:${problem.position}`)
                  }
                />
                <strong>{problem.prompt}</strong>
              </label>
              <p>{problem.difficulty}</p>
              <Link
                href={`/employer/modules/${resourceId(problem.module_id)}/questions`}
              >
                Open problem
              </Link>
            </article>
          ))}
          <UnavailableNotice>
            Individual problem copying isn’t supported yet. Open the source
            module to inspect it, or copy a published module from the library.
          </UnavailableNotice>
        </>
      ) : (
        <div className={e.tableWrap}>
          <table className={e.table}>
            <thead>
              <tr>
                <th scope="col">Module</th>
                <th scope="col">Type</th>
                <th scope="col">Status</th>
              </tr>
            </thead>
            <tbody>
              {modules.data?.map((module) => (
                <tr key={module.id}>
                  <td>
                    <button
                      onClick={() => setSelected(module.id)}
                      aria-pressed={selected === module.id}
                    >
                      {module.name}
                    </button>
                  </td>
                  <td>{kinds[module.type]}</td>
                  <td>
                    <StatusBadge>
                      {module.published ? "Published" : "Draft"}
                    </StatusBadge>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {!modules.loading &&
        !bank.loading &&
        (problems ? !bank.data?.length : !modules.data?.length) && (
          <EmptyState
            title={problems ? "No coding problems" : "No modules yet"}
            description={
              library
                ? "Published modules will appear here when available."
                : "Create a module to start building reusable assessments."
            }
          />
        )}
      {library && chosen ? (
        <>
          <section className={e.panel}>
            <h2>Selected module</h2>
            <p>
              {chosen.name} · {kinds[chosen.type]}
            </p>
            <h3>Destination</h3>
            <p>Your company workspace</p>
          </section>
          <section className={e.panel}>
            <h2>Copy behavior</h2>
            <p>
              The source stays unchanged. Review the copied module before using
              it in a pipeline.
            </p>
            <p>
              Copied definitions are locked by the current server; an editable
              copy requires additional backend support.
            </p>
          </section>
        </>
      ) : (
        !problems && (
          <section className={e.section}>
            <h2>Published modules</h2>
            <p>
              Create a new version to edit a published definition. Published and
              copied definitions remain locked by the server.
            </p>
          </section>
        )
      )}
      <div className={e.actions}>
        {selected && !problems && (
          <Link href={`/employer/modules/${resourceId(selected)}/questions`}>
            Open {chosen?.published ? "module" : "draft"}
          </Link>
        )}
        {library && (
          <button
            className={e.primary}
            disabled={!selected || busy}
            onClick={() => setConfirm(true)}
          >
            Copy module
          </button>
        )}
        <Link href="/employer/modules/new" className={e.primary}>
          Create module
        </Link>
        <Link
          href={library ? "/employer/modules" : "/employer/modules/library"}
        >
          {library ? "Company modules" : "Published library"}
        </Link>
        <Link href="/employer/modules/library/problems">
          Coding problem library
        </Link>
      </div>
      <WorkspaceFooter />
      {confirm && (
        <ConfirmDialog
          title="Copy module to company?"
          description={`Copy ${chosen?.name || "the selected module"} into your company workspace. The original remains unchanged.`}
          confirmLabel="Copy module"
          intent="success"
          onConfirm={copy}
          onClose={() => setConfirm(false)}
        />
      )}
    </>
  );
}
function DefinitionEditor({
  module,
  video = false,
  onSaved,
}: {
  module?: AssessmentModule;
  video?: boolean;
  onSaved: (module: AssessmentModule) => void;
}) {
  const [amarat, setAmarat] = useState<Amara[]>(() => module?.amarat ?? []);
  const amaraOptionsId = useId();
  const [questions, setQuestions] = useState<Question[]>(
    () =>
      module?.questions || [
        { prompt: "", format: video ? "video" : "short_text" },
      ],
  );
  const [selected, setSelected] = useState(0);
  const [type, setType] = useState<ModuleDefinition["type"]>(
    module?.type || (video ? "recorded_video" : "behavioral"),
  );
  const [error, setError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const [publish, setPublish] = useState(false);
  const formRef = useRef<HTMLFormElement>(null);
  const pending = useRef(false);
  const locked = !!(module?.published || module?.copied_from_module_id);
  const codeModule = type === "coding" || type === "programming_language";
  const defaultAmara =
    codeModule && amarat.length === 1 ? amarat[0].name.trim() : "";
  function update(patch: Partial<Question>) {
    setQuestions((items) =>
      items.map((item, index) =>
        index === selected ? { ...item, ...patch } : item,
      ),
    );
  }
  function addQuestionAmara(awards: Amara[], value: string) {
    const name = value.trim();
    if (!name || findAmara(awards, name)) {
      return awards;
    }
    return [...awards, { name, crest: "", domain: "", description: "" }];
  }
  function findAmara(awards: Amara[], name: string) {
    return awards.find(
      (award) => award.name.trim().toLowerCase() === name.trim().toLowerCase(),
    );
  }
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending.current || locked) return;
    try {
      onSaved(await saveDraft(event.currentTarget));
    } catch (cause) {
      setError(cause as Error);
    }
  }
  async function saveDraft(form: HTMLFormElement) {
    const data = new FormData(form);
    const awards = questions.reduce(
      (items, question) => addQuestionAmara(items, question.amara ?? ""),
      amarat,
    ).map((amara) => ({
      ...amara,
      name: amara.name.trim(),
    }));
    const definition: ModuleDefinition = {
      amarat: awards,
      rubric: module?.rubric,
      library: module?.library,
      name: String(data.get("name") || "").trim(),
      description: String(data.get("description") || ""),
      type,
      behavioral: type === "recorded_video" ? module?.behavioral : undefined,
      minimum_tests_passed: codeModule ? module?.minimum_tests_passed : undefined,
      active_minutes:
        type === "recorded_video"
          ? undefined
          : Number(data.get("active_minutes")) || undefined,
      threshold:
        type === "behavioral" || (type === "recorded_video" && module?.behavioral) || data.get("threshold") === ""
          ? undefined
          : Number(data.get("threshold")),
      questions: questions.map((question) => ({
        ...question,
        amara: findAmara(awards, question.amara ?? "")?.name ||
          defaultAmara || undefined,
      })),
    };
    if (
      !questions.length ||
      questions.some(
        (q) =>
          !q.prompt.trim() ||
          (q.format === "multiple_choice" &&
            (!q.options?.length || q.options.some((x) => !x.trim()))),
      )
    ) {
      throw new Error(
        "Every question needs a prompt. Multiple choice questions need non-empty options.",
      );
    }
    pending.current = true;
    setBusy(true);
    setError(undefined);
    try {
      if (module) return await updateModule(module.id, definition);
      return await createModule(definition);
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  const question = questions[selected];
  return (
    <>
      <form ref={formRef} onSubmit={submit} aria-busy={busy}>
        <section className={e.panel}>
          <h2>Module settings</h2>
          {locked && (
            <p>
              This definition is locked. Its original questions and
              configuration are preserved.
            </p>
          )}
          <div className={e.grid}>
            <label className={e.field}>
              Module name (required)
              <input
                name="name"
                required
                defaultValue={module?.name}
                readOnly={locked}
              />
            </label>
            <label className={e.field}>
              Type
              <select
                value={type}
                disabled={locked}
                onChange={(event) =>
                  setType(event.target.value as ModuleDefinition["type"])
                }
              >
                {Object.entries(kinds).map(([value, label]) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </label>
            <label className={`${e.field} ${e.full}`}>
              Description
              <textarea
                name="description"
                defaultValue={module?.description}
                readOnly={locked}
              />
            </label>
            {type !== "recorded_video" && (
              <label className={e.field}>
                Active attempt limit (minutes)
                <input
                  type="number"
                  min={1}
                  max={2147483647}
                  name="active_minutes"
                  defaultValue={module?.active_minutes}
                  readOnly={locked}
                />
              </label>
            )}
            <label className={e.field}>
              Pass threshold (%)
              <input
                name="threshold"
                type="number"
                min={0}
                max={100}
                defaultValue={module?.threshold}
                readOnly={locked}
              />
            </label>
          </div>
        </section>
        <section className={e.panel}>
          <h2>Amarat</h2>
          <p>Type a new Amara on a question to add it here, then fill in its crest, domain, and description before publishing.</p>
          {amarat.map((amara, index) => (
            <fieldset key={index} className={e.grid} disabled={locked}>
              <legend>Amara {index + 1}</legend>
              {(["name", "crest", "domain", "description"] as const).map((field) => (
                <label key={field} className={e.field}>
                  {field[0].toUpperCase() + field.slice(1)}
                  <input
                    value={amara[field]}
                    onChange={(event) => {
                      const value = event.target.value;
                      setAmarat((items) => items.map((item, i) =>
                        i === index ? { ...item, [field]: value } : item,
                      ));
                    }}
                  />
                </label>
              ))}
              {!locked && (
                <button type="button" onClick={() => {
                  setAmarat((items) => items.filter((_, i) => i !== index));
                  setQuestions((items) => items.map((item) =>
                    item.amara && findAmara([amara], item.amara)
                      ? { ...item, amara: "" }
                      : item,
                  ));
                }}>
                  Remove Amara
                </button>
              )}
            </fieldset>
          ))}
          {!locked && (
            <button type="button" onClick={() => setAmarat((items) => [
              ...items, { name: "", crest: "", domain: "", description: "" },
            ])}>
              Add Amara
            </button>
          )}
        </section>
        {video && (
          <UnavailableNotice>
            Prompt recording limits, retakes, human grader assignment, and stage
            scheduling need backend support. Video questions can be authored and
            saved; publication must follow the server’s checks.
          </UnavailableNotice>
        )}
        {questions.map((q, index) => (
          <article className={e.panel} key={index}>
            <button
              type="button"
              aria-pressed={index === selected}
              onClick={() => setSelected(index)}
            >
              <strong>
                Question {String(index + 1).padStart(2, "0")} ·{" "}
                {q.format.replaceAll("_", " ")}
              </strong>
            </button>
            <p>{q.prompt || "Enter a question prompt"}</p>
            <p>
              {q.amara || defaultAmara || "Skill not assigned"} ·{" "}
              {q.difficulty || "Difficulty not assigned"}
            </p>
          </article>
        ))}
        {question && (
          <section className={e.panel}>
            <h2>Selected question configuration</h2>
            <div className={e.grid}>
              <label className={`${e.field} ${e.full}`}>
                Prompt (required)
                <textarea
                  value={question.prompt}
                  required
                  readOnly={locked}
                  onChange={(event) => update({ prompt: event.target.value })}
                />
              </label>
              <label className={e.field}>
                Response format
                <select
                  value={question.format}
                  disabled={locked}
                  onChange={(event) =>
                    update({ format: event.target.value as Question["format"] })
                  }
                >
                  {[
                    "short_text",
                    "long_text",
                    "multiple_choice",
                    "code",
                    "video",
                  ].map((format) => (
                    <option value={format} key={format}>
                      {format.replaceAll("_", " ")}
                    </option>
                  ))}
                </select>
              </label>
              <label className={e.field}>
                amara skill
                <input
                  list={amaraOptionsId}
                  placeholder={defaultAmara || "Choose or type an Amara"}
                  value={question.amara ?? ""}
                  readOnly={locked}
                  onChange={(event) => update({ amara: event.target.value })}
                  onBlur={(event) => {
                    const value = event.target.value;
                    if (!locked) setAmarat((items) => addQuestionAmara(items, value));
                  }}
                />
                <datalist id={amaraOptionsId}>
                  {amarat.map((amara, index) => amara.name.trim() && (
                    <option key={index} value={amara.name.trim()} />
                  ))}
                </datalist>
              </label>
              <label className={e.field}>
                Difficulty
                <select
                  value={question.difficulty || ""}
                  disabled={locked}
                  onChange={(event) =>
                    update({ difficulty: event.target.value })
                  }
                >
                  <option value="">Choose difficulty</option>
                  {["easy", "medium", "hard"].map((d) => (
                    <option key={d}>{d}</option>
                  ))}
                </select>
              </label>
              {question.format === "multiple_choice" && (
                <label className={`${e.field} ${e.full}`}>
                  Answer options (one per line)
                  <textarea
                    value={question.options?.join("\n") || ""}
                    readOnly={locked}
                    onChange={(event) =>
                      update({ options: event.target.value.split("\n") })
                    }
                  />
                </label>
              )}
              {question.format === "code" && (
                <div className={e.full}>
                  <h3>Weighted test cases</h3>
                  {question.test_cases?.map((test, index) => (
                    <div className={e.grid} key={index}>
                      {(["input", "expected_output", "weight"] as const).map(
                        (key) => (
                          <label className={e.field} key={key}>
                            {key.replaceAll("_", " ")}
                            <input
                              value={test[key]}
                              type={key === "weight" ? "number" : "text"}
                              min={key === "weight" ? 0 : undefined}
                              readOnly={locked}
                              onChange={(event) =>
                                update({
                                  test_cases: question.test_cases!.map(
                                    (row, i) =>
                                      i === index
                                        ? {
                                            ...row,
                                            [key]:
                                              key === "weight"
                                                ? Number(event.target.value)
                                                : event.target.value,
                                          }
                                        : row,
                                  ),
                                })
                              }
                            />
                          </label>
                        ),
                      )}
                    </div>
                  ))}
                  <button
                    type="button"
                    disabled={locked}
                    onClick={() =>
                      update({
                        test_cases: [
                          ...(question.test_cases || []),
                          { input: "", expected_output: "", weight: 1 },
                        ],
                      })
                    }
                  >
                    Add test case
                  </button>
                </div>
              )}
            </div>
          </section>
        )}
        <ApiFeedback error={error} />
        <div className={e.actions}>
          <button
            type="button"
            disabled={locked}
            onClick={() => {
              setQuestions((items) => [
                ...items,
                { prompt: "", format: video ? "video" : "short_text" },
              ]);
              setSelected(questions.length);
            }}
          >
            Add {video ? "prompt" : "question"}
          </button>
          <button
            type="button"
            disabled={locked || !question}
            onClick={() => {
              setQuestions((items) => [...items, structuredClone(question)]);
              setSelected(questions.length);
            }}
          >
            Duplicate
          </button>
          <button
            type="button"
            disabled={locked || !question}
            onClick={() => {
              setQuestions((items) =>
                items.filter((_, index) => index !== selected),
              );
              setSelected(0);
            }}
          >
            Remove question
          </button>
          <Link href="/employer/modules/library">Select from library</Link>
          <button className={e.primary} type="submit" disabled={locked || busy}>
            {busy ? "Saving…" : "Save draft"}
          </button>
          <button
            type="button"
            className={e.primary}
            disabled={locked || !module || busy}
            onClick={() => {
              if (formRef.current?.reportValidity()) setPublish(true);
            }}
          >
            Publish module
          </button>
        </div>
      </form>
      <section className={e.panel}>
        <h2>Generate question drafts</h2>
        <p>
          Automated generation and question review workflows are not available
          yet. You can author questions above.
        </p>
        <div className={e.actions}>
          <button disabled>Generate drafts</button>
          <button disabled>Approve question</button>
        </div>
      </section>
      {publish && module && (
        <ConfirmDialog
          title="Publish this module?"
          description="Your current changes will be saved before publishing. Publishing locks the definition."
          confirmLabel="Publish module"
          intent="success"
          onClose={() => setPublish(false)}
          onConfirm={async () => {
            const form = formRef.current;
            if (!form) throw new Error("The module form is unavailable.");
            const saved = await saveDraft(form);
            const result = await publishModule(saved.id);
            onSaved(result);
            setPublish(false);
          }}
        />
      )}
    </>
  );
}
export function ModuleStudio({ id }: { id?: string }) {
  const params = useSearchParams();
  const router = useRouter();
  const resource = useApiResource<AssessmentModule>(
    id ? `/v1/modules/${resourceId(id)}` : null,
  );
  const video =
    params.get("type") === "recorded_video" ||
    resource.data?.type === "recorded_video";
  return (
    <>
      <ProductHeading
        title={
          video
            ? "Behavioral response · Video module"
            : id
              ? "Module question editor"
              : "Question module setup"
        }
        description={
          resource.data
            ? `${resource.data.name} · ${resource.data.published ? "Published" : "Draft"}`
            : "Build a reusable assessment definition."
        }
        compact
        artwork={{
          src: "/figma/hiring-dag.svg",
          width: 200.16,
          height: 138.96,
        }}
      />
      <BusyState loading={resource.loading} />
      <ApiFeedback error={resource.error} retry={resource.reload} />
      {(!id || resource.data) && (
        <DefinitionEditor
          key={resource.data?.id || "new"}
          module={resource.data}
          video={video}
          onSaved={(saved) => {
            if (!id)
              router.push(
                `/employer/modules/${resourceId(saved.id)}/questions`,
              );
            else resource.reload();
          }}
        />
      )}
      <WorkspaceFooter />
    </>
  );
}
