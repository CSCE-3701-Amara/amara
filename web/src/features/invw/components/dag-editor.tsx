"use client";
import { useEffect, useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { ApiFeedback } from "@/shared/ui/workspace";
import type { PipelineDraft, ModuleSummary } from "../types/invw.types";
import styles from "./dag-editor.module.css";
import { DagCanvas } from "./dag-canvas";
import { DagConnections } from "./dag-connections";
import { PipelineModulePicker } from "./pipeline-module-picker";
import {
  addConnection,
  normalizeGraph,
  removeConnection,
  removeStage,
  wouldCreateCycle,
  type DagEdge,
  type DagNode,
  type PipelineGraph,
} from "./dag-graph";
export function DagEditor({ postingId }: { postingId: string }) {
  const pipeline = useApiResource<PipelineDraft>(
    `/v1/jobs/${resourceId(postingId)}/pipeline`,
  );
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  async function create() {
    setBusy(true);
    setError(undefined);
    try {
      await apiRequest(`/v1/jobs/${resourceId(postingId)}/pipeline`, {
        method: "POST",
      });
      pipeline.reload();
    } catch (e) {
      setError(e as Error);
    } finally {
      setBusy(false);
    }
  }
  return (
    <>
      {pipeline.error?.status === 404 && (
        <section aria-label="Create pipeline graph">
          <h2>Design the pipeline</h2>
          <button disabled={busy} onClick={() => void create()}>
            {busy ? "Creating…" : "Create graph pipeline"}
          </button>
        </section>
      )}
      <ApiFeedback
        error={
          error || (pipeline.error?.status === 404 ? undefined : pipeline.error)
        }
      />
      {pipeline.data && (
        <Canvas key={pipeline.data.versionId} pipeline={pipeline.data} />
      )}
    </>
  );
}
function Canvas({ pipeline }: { pipeline: PipelineDraft }) {
  const path = `/v1/pipeline-versions/${resourceId(pipeline.versionId)}`;
  const [history, setHistory] = useState<{
    past: PipelineGraph[];
    present: PipelineGraph;
    future: PipelineGraph[];
  }>({ past: [], present: { stages: [], links: [] }, future: [] });
  const graph = history.present;
  const [selected, select] = useState("");
  const [target, setTarget] = useState("");
  const [zoom, setZoom] = useState(100);
  const [message, setMessage] = useState("");
  const [error, setError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const [supported, setSupported] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const [published, setPublished] = useState(false);
  const team = useApiResource<{ data: { id: string; name: string }[] }>(
    "/v1/companies/me/interviewers?limit=100",
  );
  const modules = useApiResource<ModuleSummary[]>("/v1/modules");
  const library = useApiResource<ModuleSummary[]>("/v1/modules/library");
  useEffect(() => {
    const controller = new AbortController();
    async function load() {
      try {
        let g: PipelineGraph;
        try {
          g = await apiRequest<PipelineGraph>(path + "/graph", {
            signal: controller.signal,
          });
          if (!g.stages) throw new Error("Invalid graph");
          setSupported(true);
        } catch {
          if (controller.signal.aborted) return;
          const stages = await apiRequest<DagNode[]>(path + "/stages", {
            signal: controller.signal,
          });
          let links: DagEdge[] = [];
          try {
            links = await apiRequest<DagEdge[]>(path + "/transitions", {
              signal: controller.signal,
            });
          } catch {}
          g = { stages, links };
        }
        if (!controller.signal.aborted) resetGraph(normalizeGraph(g));
      } catch (e) {
        if (!controller.signal.aborted) setError(e as Error);
      } finally {
        if (!controller.signal.aborted) setLoaded(true);
      }
    }
    void load();
    return () => controller.abort();
  }, [path]);
  const node = graph.stages.find((s) => s.id === selected);
  const locked = pipeline.published || published;
  function resetGraph(present: PipelineGraph) {
    setHistory({ past: [], present, future: [] });
  }
  function setGraph(change: (current: PipelineGraph) => PipelineGraph) {
    setHistory((current) => {
      const next = change(current.present);
      if (next === current.present) return current;
      return {
        past: [...current.past.slice(-49), current.present],
        present: next,
        future: [],
      };
    });
  }
  function restore(redo: boolean) {
    setHistory((current) => {
      if (redo) {
        if (!current.future.length) return current;
        return {
          past: [...current.past.slice(-49), current.present],
          present: current.future[0],
          future: current.future.slice(1),
        };
      }
      if (!current.past.length) return current;
      return {
        past: current.past.slice(0, -1),
        present: current.past[current.past.length - 1],
        future: [current.present, ...current.future],
      };
    });
    select("");
    setMessage(redo ? "Graph change redone locally." : "Graph change undone locally.");
  }
  function update(id: string, patch: Partial<DagNode>) {
    setGraph((g) => ({
      ...g,
      stages: g.stages.map((s) => (s.id === id ? { ...s, ...patch } : s)),
    }));
  }
  function connect() {
    if (!node || !target) return;
    if (wouldCreateCycle(graph, selected, target)) {
      setMessage("This connection would create a cycle.");
      return;
    }
    setGraph((current) => addConnection(current, selected, target));
    setMessage("");
  }
  async function save() {
    setBusy(true);
    setError(undefined);
    try {
      resetGraph(
        await apiRequest<PipelineGraph>(path + "/graph", {
          method: "PUT",
          body: graph,
        }),
      );
      select("");
      setMessage("Graph saved.");
    } catch (e) {
      setError(e as Error);
    } finally {
      setBusy(false);
    }
  }
  return (
    <section aria-label="Pipeline graph" className={styles.designer}>
      <div className={styles.heading}>
        <div>
          <h2>Design the pipeline</h2>
          <p>
            Build the path candidates follow, from their first stage to their
            last.
          </p>
        </div>
        <span className={styles.status}>{locked ? "Published" : "Draft"}</span>
      </div>
      <div className={styles.guide}>
        <span>
          <b>1</b> Add a stage
        </span>
        <span>
          <b>2</b> Select and set it up
        </span>
        <span>
          <b>3</b> Connect stages
        </span>
      </div>
      <div className={styles.toolbar}>
        <button
          type="button"
          aria-label="Undo graph change"
          disabled={locked || busy || !history.past.length}
          onClick={() => restore(false)}
        >
          Undo
        </button>
        <button
          type="button"
          aria-label="Redo graph change"
          disabled={locked || busy || !history.future.length}
          onClick={() => restore(true)}
        >
          Redo
        </button>
        <button
          disabled={locked || busy}
          onClick={() => {
            const id = crypto.randomUUID();
            setGraph((g) => ({
              ...g,
              stages: [
                ...g.stages,
                {
                  id,
                  name: "New assessment",
                  type: "assessment",
                  scheduledStartAt: "",
                  scheduledEndAt: "",
                  moduleIds: [],
                  x: 40 + (g.stages.length % 3) * 260,
                  y: 40 + Math.floor(g.stages.length / 3) * 200,
                },
              ],
            }));
            select(id);
          }}
        >
          <span aria-hidden="true">＋ </span>Add stage
        </button>
        <div className={styles.zoom} aria-label="Canvas zoom">
          <button
            aria-label="Zoom out"
            onClick={() => setZoom((z) => Math.max(30, z - 10))}
          >
            −
          </button>
          <output>{zoom}%</output>
          <button
            aria-label="Zoom in"
            onClick={() => setZoom((z) => Math.min(150, z + 10))}
          >
            +
          </button>
        </div>
        <span className={styles.toolbarSpacer} />
        <button
          disabled={locked || busy || !supported}
          onClick={() => void save()}
        >
          Save graph draft
        </button>
        <button
          className={styles.publish}
          disabled={locked || busy || !supported}
          onClick={async () => {
            setBusy(true);
            try {
              const saved = await apiRequest<PipelineGraph>(path + "/graph", {
                method: "PUT",
                body: graph,
              });
              resetGraph(saved);
              await apiRequest<void>(path + "/publish", { method: "POST" });
              setPublished(true);
              setMessage("Pipeline published.");
            } catch (e) {
              setError(e as Error);
            } finally {
              setBusy(false);
            }
          }}
        >
          Publish graph
        </button>
      </div>
      {locked && (
        <p className={styles.notice}>
          This pipeline is published and cannot be edited.
        </p>
      )}
      {loaded && !error && !supported && !locked && (
        <p className={styles.notice}>
          Graph editing is unavailable for this pipeline. You can view its
          stages below.
        </p>
      )}
      <div className={styles.canvasHeading}>
        <div>
          <strong>Candidate journey</strong>
          <span>
            {graph.stages.length} stages · {graph.links.length} connections
          </span>
        </div>
        <p>Drag stages to arrange them. Select a stage to edit it.</p>
      </div>
      <div className={styles.layout}>
        <DagCanvas
          graph={graph}
          selectedStageId={selected}
          zoom={zoom}
          disabled={locked || busy}
          onSelect={select}
          onUpdate={update}
        />
        <aside className={styles.inspector} aria-label="Stage editor">
          {node ? (
            <fieldset disabled={locked || busy}>
              <legend>Stage settings</legend>
              <p className={styles.inspectorIntro}>
                Editing <strong>{node.name}</strong>. Changes are saved with the
                graph draft.
              </p>
              <label>
                Advancement rule
                <select
                  value={node.advancementRule || "completion"}
                  onChange={(e) =>
                    update(node.id, {
                      advancementRule: e.target
                        .value as DagNode["advancementRule"],
                      scoreThreshold:
                        e.target.value === "threshold"
                          ? (node.scoreThreshold ?? 70)
                          : undefined,
                    })
                  }
                >
                  <option value="completion">On completion</option>
                  <option value="manual">Hiring manager decision</option>
                  <option value="threshold">Score threshold</option>
                </select>
              </label>
              {node.advancementRule === "threshold" && (
                <>
                  <label>
                    Stage score threshold
                    <input
                      type="number"
                      min="0"
                      max="100"
                      value={node.scoreThreshold ?? 70}
                      onChange={(e) =>
                        update(node.id, {
                          scoreThreshold: Number(e.target.value),
                        })
                      }
                    />
                  </label>
                  <label>
                    Below threshold
                    <select
                      value={node.belowThresholdAction || "hold"}
                      onChange={(e) =>
                        update(node.id, {
                          belowThresholdAction: e.target.value as
                            "hold" | "reject",
                        })
                      }
                    >
                      <option value="hold">Hold for review</option>
                      <option value="reject">Reject application</option>
                    </select>
                  </label>
                </>
              )}

              <label>
                Stage type
                <select
                  value={node.type}
                  onChange={(e) =>
                    update(
                      node.id,
                      e.target.value === "assessment"
                        ? {
                            type: "assessment",
                            interview: undefined,
                            moduleIds: [],
                          }
                        : {
                            type: "live_interview",
                            moduleIds: [],
                            interview: {
                              durationMinutes: 30,
                              interviewerIds: [],
                              threshold: 70,
                              rubric: [
                                {
                                  description: "Explain reasoning",
                                  amara: "Communication",
                                },
                              ],
                              amarat: [
                                {
                                  name: "Communication",
                                  crest: "speech",
                                  domain: "Communication",
                                  description: "Clear reasoning",
                                },
                              ],
                            },
                          },
                    )
                  }
                >
                  <option value="assessment">Assessment</option>
                  <option value="live_interview">Live interview</option>
                </select>
              </label>
              {node.interview && (
                <>
                  <label>
                    Interview duration (minutes)
                    <input
                      type="number"
                      min="1"
                      value={node.interview.durationMinutes}
                      onChange={(e) =>
                        update(node.id, {
                          interview: {
                            ...node.interview!,
                            durationMinutes: Number(e.target.value),
                          },
                        })
                      }
                    />
                  </label>
                  <label>
                    Interview pass threshold
                    <input
                      type="number"
                      min="0"
                      max="100"
                      value={node.interview.threshold ?? ""}
                      onChange={(e) =>
                        update(node.id, {
                          interview: {
                            ...node.interview!,
                            threshold: Number(e.target.value),
                          },
                        })
                      }
                    />
                  </label>
                  {team.data?.data.map((person) => (
                    <label key={person.id}>
                      <input
                        type="checkbox"
                        checked={node.interview!.interviewerIds.includes(
                          person.id,
                        )}
                        onChange={(e) =>
                          update(node.id, {
                            interview: {
                              ...node.interview!,
                              interviewerIds: e.target.checked
                                ? [...node.interview!.interviewerIds, person.id]
                                : node.interview!.interviewerIds.filter(
                                    (id) => id !== person.id,
                                  ),
                            },
                          })
                        }
                      />
                      {person.name}
                    </label>
                  ))}
                  <label>
                    Rubric criterion
                    <input
                      value={node.interview.rubric[0]?.description || ""}
                      onChange={(e) =>
                        update(node.id, {
                          interview: {
                            ...node.interview!,
                            rubric: node.interview!.rubric.map((r, i) =>
                              i === 0
                                ? { ...r, description: e.target.value }
                                : r,
                            ),
                          },
                        })
                      }
                    />
                  </label>
                </>
              )}
              <label>
                Name
                <input
                  value={node.name}
                  onChange={(e) => update(node.id, { name: e.target.value })}
                />
              </label>
              {(["scheduledStartAt", "scheduledEndAt"] as const).map(
                (field) => (
                  <label key={field}>
                    {field === "scheduledStartAt" ? "Opens at" : "Deadline"}
                    <input
                      type="datetime-local"
                      value={
                        node[field]
                          ? new Date(
                              Date.parse(node[field]) -
                                new Date(node[field]).getTimezoneOffset() *
                                  60000,
                            )
                              .toISOString()
                              .slice(0, 16)
                          : ""
                      }
                      onChange={(e) =>
                        update(node.id, {
                          [field]: e.target.value
                            ? new Date(e.target.value).toISOString()
                            : "",
                        })
                      }
                    />
                  </label>
                ),
              )}
              {node.type === "assessment" && (
                <PipelineModulePicker
                  key={node.id}
                  modules={modules.data ?? []}
                  library={library.data ?? []}
                  selected={node.moduleIds ?? []}
                  disabled={locked || busy}
                  onToggle={(id, checked) =>
                    update(node.id, {
                      moduleIds: checked
                        ? [...(node.moduleIds ?? []), id]
                        : node.moduleIds?.filter((value) => value !== id),
                    })
                  }
                  onCopyingChange={setBusy}
                  onCopy={(module) => {
                    update(node.id, {
                      moduleIds: [...(node.moduleIds ?? []), module.id],
                    });
                    modules.reload();
                  }}
                />
              )}
              <DagConnections
                graph={graph}
                node={node}
                target={target}
                onTargetChange={setTarget}
                onConnect={connect}
                onRemoveConnection={(edge) =>
                  setGraph((current) => removeConnection(current, edge))
                }
                onRemoveStage={() => {
                  setGraph((current) => removeStage(current, node.id));
                  select("");
                }}
              />
            </fieldset>
          ) : (
            <div className={styles.emptyInspector}>
              <span aria-hidden="true">◇</span>
              <h3>Select a stage</h3>
              <p>
                Choose a stage on the canvas to edit its details and connect it
                to the next step.
              </p>
            </div>
          )}
        </aside>
      </div>
      {message && <p role="status">{message}</p>}
      <ApiFeedback error={error} />
    </section>
  );
}
