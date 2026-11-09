"use client";
import { useEffect, useState } from "react";
import { apiRequest, resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { ApiFeedback } from "@/shared/ui/workspace";
import type {
  PipelineDraft,
  PipelineStage,
  ModuleSummary,
} from "../types/invw.types";
import styles from "./dag-editor.module.css";
import { PipelineModulePicker } from "./pipeline-module-picker";
type Node = PipelineStage & { id: string; x: number; y: number };
type Edge = { fromStageId: string; toStageId: string };
type Graph = { stages: Node[]; links: Edge[] };
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
  const [graph, setGraph] = useState<Graph>({ stages: [], links: [] });
  const [selected, select] = useState("");
  const [target, setTarget] = useState("");
  const [zoom, setZoom] = useState(100);
  const [message, setMessage] = useState("");
  const [error, setError] = useState<Error>();
  const [busy, setBusy] = useState(false);
  const [supported, setSupported] = useState(false);
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
        let g: Graph;
        try {
          g = await apiRequest<Graph>(path + "/graph", {
            signal: controller.signal,
          });
          if (!g.stages) throw new Error("Invalid graph");
          setSupported(true);
        } catch {
          if (controller.signal.aborted) return;
          const stages = await apiRequest<Node[]>(path + "/stages", {
            signal: controller.signal,
          });
          let links: Edge[] = [];
          try {
            links = await apiRequest<Edge[]>(path + "/transitions", {
              signal: controller.signal,
            });
          } catch {}
          g = { stages, links };
        }
        if (!controller.signal.aborted)
          setGraph({
            stages: g.stages.map((s, i) => ({
              ...s,
              x: s.x ?? 40 + (i % 3) * 260,
              y: s.y ?? 40 + Math.floor(i / 3) * 200,
            })),
            links: g.links || [],
          });
      } catch (e) {
        if (!controller.signal.aborted) setError(e as Error);
      }
    }
    void load();
    return () => controller.abort();
  }, [path]);
  const node = graph.stages.find((s) => s.id === selected);
  const locked = pipeline.published || published;
  const canvasWidth = Math.max(1000, ...graph.stages.map((s) => s.x + 240));
  const canvasHeight = Math.max(800, ...graph.stages.map((s) => s.y + 180));
  function update(id: string, patch: Partial<Node>) {
    setGraph((g) => ({
      ...g,
      stages: g.stages.map((s) => (s.id === id ? { ...s, ...patch } : s)),
    }));
  }
  function connect() {
    if (!node || !target) return;
    const reaches = (id: string, seen = new Set<string>()): boolean => {
      if (id === selected) return true;
      if (seen.has(id)) return false;
      seen.add(id);
      return graph.links
        .filter((e) => e.fromStageId === id)
        .some((e) => reaches(e.toStageId, seen));
    };
    if (reaches(target)) {
      setMessage("This connection would create a cycle.");
      return;
    }
    if (
      !graph.links.some(
        (e) => e.fromStageId === selected && e.toStageId === target,
      )
    )
      setGraph((g) => ({
        ...g,
        links: [...g.links, { fromStageId: selected, toStageId: target }],
      }));
    setMessage("");
  }
  async function save() {
    setBusy(true);
    setError(undefined);
    try {
      setGraph(
        await apiRequest<Graph>(path + "/graph", {
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
    <section aria-label="Pipeline graph">
      <h2>Design the pipeline</h2>
      <p>
        Connect assessments to create branches. A joined stage opens after every
        incoming stage passes.
      </p>
      <div className={styles.toolbar}>
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
                  x: 40,
                  y: 40,
                },
              ],
            }));
            select(id);
          }}
        >
          Add stage
        </button>
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
        <button
          disabled={locked || busy || !supported}
          onClick={() => void save()}
        >
          Save graph draft
        </button>
      </div>
      <button
        disabled={locked || busy || !supported}
        onClick={async () => {
          setBusy(true);
          try {
            const saved = await apiRequest<Graph>(path + "/graph", {
              method: "PUT",
              body: graph,
            });
            setGraph(saved);
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
      <div className={styles.layout}>
        <div className={styles.viewport}>
          <div
            className={styles.canvas}
            style={{
              transform: `scale(${zoom / 100})`,
              width: canvasWidth,
              height: canvasHeight,
            }}
          >
            <svg
              className={styles.edges}
              width={canvasWidth}
              height={canvasHeight}
              aria-label="Stage connections"
            >
              {graph.links.map((e) => {
                const a = graph.stages.find((s) => s.id === e.fromStageId),
                  b = graph.stages.find((s) => s.id === e.toStageId);
                return a && b ? (
                  <g key={a.id + b.id}>
                    <path
                      d={`M${a.x + 200},${a.y + 65} C${a.x + 240},${a.y + 65} ${b.x - 40},${b.y + 65} ${b.x},${b.y + 65}`}
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                    />
                    <text x={(a.x + b.x + 200) / 2} y={(a.y + b.y) / 2 + 55}>
                      →
                    </text>
                  </g>
                ) : null;
              })}
            </svg>
            {graph.stages.map((s) => (
              <button
                key={s.id}
                className={styles.node}
                style={{ left: s.x, top: s.y }}
                aria-pressed={selected === s.id}
                aria-label={`${s.name}, ${graph.links.filter((e) => e.toStageId === s.id).length} incoming connections`}
                onClick={() => select(s.id)}
                onKeyDown={(e) => {
                  if (locked || busy || !e.altKey || !e.key.startsWith("Arrow"))
                    return;
                  e.preventDefault();
                  update(s.id, {
                    x: Math.max(
                      0,
                      s.x +
                        (e.key === "ArrowRight"
                          ? 20
                          : e.key === "ArrowLeft"
                            ? -20
                            : 0),
                    ),
                    y: Math.max(
                      0,
                      s.y +
                        (e.key === "ArrowDown"
                          ? 20
                          : e.key === "ArrowUp"
                            ? -20
                            : 0),
                    ),
                  });
                }}
                draggable={!locked && !busy}
                onDragEnd={(e) => {
                  if (locked || busy) return;
                  const rect =
                    e.currentTarget.parentElement!.getBoundingClientRect();
                  update(s.id, {
                    x: Math.max(
                      0,
                      ((e.clientX - rect.left) * 100) / zoom - 100,
                    ),
                    y: Math.max(0, ((e.clientY - rect.top) * 100) / zoom - 40),
                  });
                }}
              >
                <strong>{s.name}</strong>
                <span>
                  {s.type === "assessment"
                    ? `${s.moduleIds?.length || 0} assessment modules`
                    : "Live interview"}
                </span>
                <small>Alt + arrow keys to move</small>
              </button>
            ))}
          </div>
        </div>
        <aside className={styles.inspector}>
          {node ? (
            <fieldset disabled={locked || busy}>
              <legend>Stage settings</legend>
              <label>
                Advancement rule
                <select
                  value={node.advancementRule || "completion"}
                  onChange={(e) =>
                    update(node.id, {
                      advancementRule: e.target
                        .value as Node["advancementRule"],
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
              <label>
                Connect to stage
                <select
                  value={target}
                  onChange={(e) => setTarget(e.target.value)}
                >
                  <option value="">Choose a stage</option>
                  {graph.stages.map((s) => (
                    <option key={s.id} value={s.id}>
                      {s.name}
                    </option>
                  ))}
                </select>
              </label>
              <button onClick={connect}>Add connection</button>
              {graph.links
                .filter((e) => e.fromStageId === node.id)
                .map((e) => (
                  <button
                    key={e.toStageId}
                    onClick={() =>
                      setGraph((g) => ({
                        ...g,
                        links: g.links.filter((x) => x !== e),
                      }))
                    }
                  >
                    Remove connection to{" "}
                    {graph.stages.find((s) => s.id === e.toStageId)?.name}
                  </button>
                ))}
              <button
                onClick={() => {
                  setGraph((g) => ({
                    stages: g.stages.filter((s) => s.id !== node.id),
                    links: g.links.filter(
                      (e) =>
                        e.fromStageId !== node.id && e.toStageId !== node.id,
                    ),
                  }));
                  select("");
                }}
              >
                Remove stage
              </button>
            </fieldset>
          ) : (
            <p>Select a stage to edit its schedule, modules and connections.</p>
          )}
        </aside>
      </div>
      {message && <p role="status">{message}</p>}
      <ApiFeedback error={error} />
    </section>
  );
}
