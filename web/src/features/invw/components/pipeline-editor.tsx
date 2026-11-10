"use client";
import { toDateTimeLocal } from "@/shared/lib/date-time";
import { toError } from "@/shared/lib/errors";
import { useState, type FormEvent } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import {
  ApiFeedback,
  BusyState,
  workspaceStyles as s,
} from "@/shared/ui/workspace";
import { ConfirmDialog } from "@/shared/ui/product-ui";
import { createPipeline, publishPipeline, replacePipelineStages } from "../api/invw.api";
import { PipelineModulePicker } from "./pipeline-module-picker";
import type {
  PipelineDraft,
  PipelineStage,
  ModuleSummary,
} from "../types/invw.types";
function StageForm({
  initial,
  versionId,
}: {
  initial: PipelineStage[];
  versionId: string;
}) {
  const [stages, setStages] = useState(initial);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [message, setMessage] = useState("");
  const modules = useApiResource<ModuleSummary[] | null>("/v1/modules");
  const library = useApiResource<ModuleSummary[] | null>("/v1/modules/library");
  function update(index: number, patch: Partial<PipelineStage> | ((stage: PipelineStage) => Partial<PipelineStage>)) {
    setStages((items) => items.map((stage, i) =>
      i === index ? { ...stage, ...(typeof patch === "function" ? patch(stage) : patch) } : stage,
    ));
    setMessage("");
  }
  function move(index: number, by: number) {
    const next = [...stages];
    [next[index], next[index + by]] = [next[index + by], next[index]];
    setStages(next);
    setMessage("");
  }
  async function save(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy) return;
    setError(undefined);
    setMessage("");
    const names = new Set<string>();
    for (const stage of stages) {
      if (
        !stage.name.trim() ||
        names.has(stage.name.trim()) ||
        !stage.scheduledStartAt ||
        !stage.scheduledEndAt ||
        new Date(stage.scheduledEndAt) <= new Date(stage.scheduledStartAt) ||
        (stage.type === "assessment" && !stage.moduleIds?.length)
      ) {
        setError(
          new Error(
            "Each stage needs a unique name, an end time after its start time, and at least one assessment.",
          ),
        );
        return;
      }
      names.add(stage.name.trim());
    }
    setBusy(true);
    try {
      setStages((await replacePipelineStages(versionId, stages)) || []);
      setMessage("Pipeline stages saved.");
    } catch (cause) {
      setError(toError(cause));
    } finally {
      setBusy(false);
    }
  }
  return (
    <form onSubmit={save} aria-label="Pipeline stages" aria-busy={busy}>
      <p>
        Stages run in the order shown. Saving replaces this draft’s stage list.
      </p>
      <BusyState loading={modules.loading || library.loading} />
      <ApiFeedback error={modules.error} retry={modules.reload} />
      <ApiFeedback error={library.error} retry={library.reload} />
      {stages.map((stage, index) => (
        <fieldset
          className={`${s.card} ${s.form}`}
          key={stage.id || `new-${index}`}
        >
          <legend>Stage {index + 1}</legend>
          <label className={s.field}>
            Stage name
            <input
              value={stage.name}
              required
              onChange={(event) => update(index, { name: event.target.value })}
            />
          </label>
          <div className={s.grid}>
            <label className={s.field}>
              Stage opens
              <input
                type="datetime-local"
                required
                value={toDateTimeLocal(stage.scheduledStartAt)}
                onChange={(event) =>
                  update(index, {
                    scheduledStartAt: event.target.value
                      ? new Date(event.target.value).toISOString()
                      : "",
                  })
                }
              />
            </label>
            <label className={s.field}>
              Stage deadline
              <input
                type="datetime-local"
                required
                value={toDateTimeLocal(stage.scheduledEndAt)}
                onChange={(event) =>
                  update(index, {
                    scheduledEndAt: event.target.value
                      ? new Date(event.target.value).toISOString()
                      : "",
                  })
                }
              />
            </label>
          </div>
          {stage.type === "assessment" ? (
            <PipelineModulePicker
              modules={modules.data ?? []}
              library={library.data ?? []}
              selected={stage.moduleIds ?? []}
              disabled={busy}
              onToggle={(id, checked) =>
                update(index, {
                  moduleIds: checked
                    ? [...(stage.moduleIds ?? []), id]
                    : stage.moduleIds?.filter((item) => item !== id),
                })
              }
              onCopyingChange={setBusy}
              onCopy={(module) => {
                update(index, (current) => ({
                  moduleIds: [...(current.moduleIds ?? []), module.id],
                }));
                modules.reload();
              }}
            />
          ) : (
            <p>
              Live interview · {stage.interview?.durationMinutes} minutes.
              Existing interviewer assignments and rubric will be preserved.
            </p>
          )}
          <div className={s.actions}>
            <button
              type="button"
              disabled={index === 0 || busy}
              onClick={() => move(index, -1)}
            >
              Move earlier
            </button>
            <button
              type="button"
              disabled={index === stages.length - 1 || busy}
              onClick={() => move(index, 1)}
            >
              Move later
            </button>
            <button
              type="button"
              disabled={busy}
              onClick={() => {
                setStages(stages.filter((_, i) => i !== index));
                setMessage("");
              }}
            >
              Remove stage
            </button>
          </div>
        </fieldset>
      ))}
      <ApiFeedback error={error} />
      {message && (
        <p role="status" className={s.success}>
          {message}
        </p>
      )}
      <div className={s.actions}>
        <button
          type="button"
          disabled={busy}
          onClick={() => {
            setStages([
              ...stages,
              {
                name: "",
                type: "assessment",
                scheduledStartAt: "",
                scheduledEndAt: "",
                moduleIds: [],
              },
            ]);
            setMessage("");
          }}
        >
          Add assessment stage
        </button>
        <button className={s.primary} type="submit" disabled={busy}>
          {busy ? "Saving…" : "Save pipeline stages"}
        </button>
      </div>
    </form>
  );
}
function Stages({ pipeline }: { pipeline: PipelineDraft }) {
  const stages = useApiResource<PipelineStage[] | null>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/stages`,
  );
  return (
    <>
      <BusyState loading={stages.loading} />
      <ApiFeedback error={stages.error} retry={stages.reload} />
      {!stages.loading && !stages.error && stages.data !== undefined && (
        pipeline.published ? (
          <ol>
            {stages.data?.map((stage) => <li key={stage.id}>{stage.name}</li>)}
          </ol>
        ) : (
          <StageForm initial={stages.data || []} versionId={pipeline.versionId} />
        )
      )}
    </>
  );
}
export function PipelineEditor({ postingId }: { postingId: string }) {
  const pipeline = useApiResource<PipelineDraft>(
    `/v1/jobs/${resourceId(postingId)}/pipeline`,
  );
  const [busy, setBusy] = useState(false);
  const [confirmPublish, setConfirmPublish] = useState(false);
  const [error, setError] = useState<Error>();
  async function create() {
    if (busy) return;
    setBusy(true);
    setError(undefined);
    try {
      await createPipeline(postingId);
      pipeline.reload();
    } catch (cause) {
      setError(toError(cause));
    } finally {
      setBusy(false);
    }
  }
  return (
    <section className={s.card} aria-label="Hiring pipeline">
      <h2>Hiring pipeline</h2>
      <BusyState loading={pipeline.loading} />
      {pipeline.error?.status === 404 ? (
        <>
          <p>This posting has no pipeline yet.</p>
          <button disabled={busy} onClick={create}>
            {busy ? "Creating…" : "Create pipeline draft"}
          </button>
        </>
      ) : (
        <ApiFeedback error={pipeline.error} retry={pipeline.reload} />
      )}
      <ApiFeedback error={error} />
      {pipeline.data && (
        <>
          <p>{pipeline.data.published ? "Published" : "Draft"}</p>
          <Stages key={pipeline.data.versionId} pipeline={pipeline.data} />
          {!pipeline.data.published && (
            <button type="button" onClick={() => setConfirmPublish(true)}>
              Publish saved stages
            </button>
          )}
        </>
      )}
      {confirmPublish && pipeline.data && (
        <ConfirmDialog
          title="Publish this pipeline?"
          description="Publishing locks the saved stages and starts the pipeline for eligible applicants."
          confirmLabel="Publish pipeline"
          onClose={() => setConfirmPublish(false)}
          onConfirm={async () => {
            await publishPipeline(pipeline.data!.versionId);
            pipeline.reload();
          }}
        />
      )}
    </section>
  );
}
