"use client";
import Link from "next/link";
import {
  useRef,
  useState,
  type PointerEvent,
  type KeyboardEvent,
} from "react";
import {
  usePathname,
  useRouter,
  useSearchParams,
} from "next/navigation";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import type { StageTransition } from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import { createPipeline } from "../api/invw.api";
import type {
  PipelineDraft,
  PipelineStage,
  ModuleSummary,
} from "../types/invw.types";
import e from "@/shared/ui/editor.module.css";
import s from "./pipeline-studio.module.css";
type GraphNode = { id: string; stage: PipelineStage; x: number; y: number };
type Edge = { from: string; to: string };
type PipelineDialog = "publish" | "migrate" | "discard";

const stageOptions: readonly [PipelineStage["type"], string][] = [
  ["assessment", "Question assessment"],
  ["live_interview", "Live interview"],
];

function localDate(value: string) {
  const date = new Date(value);
  return Number.isNaN(date.getTime())
    ? ""
    : new Date(date.getTime() - date.getTimezoneOffset() * 60000)
        .toISOString()
        .slice(0, 16);
}

/** A real DAG model: branches and joins stay explicit, independent of stage ordering. */
function GraphEditor({
  initial,
  transitions,
}: {
  initial: PipelineStage[];
  transitions?: StageTransition[];
}) {
  const [nodes, setNodes] = useState<GraphNode[]>(() =>
    initial.map((stage, index) => ({
      id: stage.id || `stage-${index}`,
      stage,
      x: 70 + (index % 2) * 300,
      y: 120 + Math.floor(index / 2) * 220,
    })),
  );
  const [edges, setEdges] = useState<Edge[]>(
    () =>
      transitions?.map((edge) => ({
        from: edge.fromStageId,
        to: edge.toStageId,
      })) || [],
  );
  const [selected, setSelected] = useState(nodes[0]?.id || "");
  const [target, setTarget] = useState("");
  const [zoom, setZoom] = useState(1);
  const [message, setMessage] = useState("");
  const drag = useRef<
    | {
        id: string;
        x: number;
        y: number;
        left: number;
        top: number;
      }
    | undefined
  >(undefined);
  const sequence = useRef(0);
  const modules = useApiResource<ModuleSummary[] | null>("/v1/modules");
  const active = nodes.find((node) => node.id === selected);
  const width = Math.max(700, ...nodes.map((n) => n.x + 250));
  const height = Math.max(710, ...nodes.map((n) => n.y + 230));
  function patch(change: Partial<PipelineStage>) {
    setNodes((items) =>
      items.map((n) =>
        n.id === selected ? { ...n, stage: { ...n.stage, ...change } } : n,
      ),
    );
    setMessage("Unsaved local graph changes.");
  }
  function add(type: PipelineStage["type"]) {
    sequence.current += 1;
    const id = `local-${sequence.current}`;
    setNodes((items) => [
      ...items,
      {
        id,
        x: 80 + (items.length % 2) * 300,
        y: 120 + Math.floor(items.length / 2) * 220,
        stage: {
          name: type === "assessment" ? "New assessment" : "Live interview",
          type,
          scheduledStartAt: "",
          scheduledEndAt: "",
          moduleIds: [],
        },
      },
    ]);
    setSelected(id);
    setMessage("Unsaved local graph changes.");
  }
  function connect() {
    if (!selected || !target || selected === target) return;
    if (edges.some((edge) => edge.from === selected && edge.to === target))
      return;
    const visited = new Set<string>();
    const stack = [target];
    while (stack.length) {
      const current = stack.pop()!;
      if (current === selected) {
        setMessage(
          "This connection would create a cycle. Choose another stage.",
        );
        return;
      }
      if (visited.has(current)) continue;
      visited.add(current);
      stack.push(
        ...edges.filter((edge) => edge.from === current).map((edge) => edge.to),
      );
    }
    setEdges((items) => [...items, { from: selected, to: target }]);
    setTarget("");
    setMessage(
      "Connection added locally. Every incoming path must pass before a joined stage unlocks.",
    );
  }
  function move(id: string, event: KeyboardEvent<HTMLButtonElement>) {
    const shifts: Record<string, [number, number]> = {
      ArrowLeft: [-20, 0],
      ArrowRight: [20, 0],
      ArrowUp: [0, -20],
      ArrowDown: [0, 20],
    };
    if (!event.altKey || !shifts[event.key]) return;
    event.preventDefault();
    const [x, y] = shifts[event.key];
    setNodes((items) =>
      items.map((n) =>
        n.id === id
          ? { ...n, x: Math.max(10, n.x + x), y: Math.max(90, n.y + y) }
          : n,
      ),
    );
  }
  function start(node: GraphNode, event: PointerEvent<HTMLButtonElement>) {
    setSelected(node.id);
    drag.current = {
      id: node.id,
      x: event.clientX,
      y: event.clientY,
      left: node.x,
      top: node.y,
    };
    event.currentTarget.setPointerCapture(event.pointerId);
  }
  function pointerMove(event: PointerEvent<HTMLButtonElement>) {
    const d = drag.current;
    if (!d) return;
    setNodes((items) =>
      items.map((n) =>
        n.id === d.id
          ? {
              ...n,
              x: Math.max(10, d.left + (event.clientX - d.x) / zoom),
              y: Math.max(90, d.top + (event.clientY - d.y) / zoom),
            }
          : n,
      ),
    );
  }
  const roots = nodes.filter((n) => !edges.some((edge) => edge.to === n.id));
  const leaves = nodes.filter((n) => !edges.some((edge) => edge.from === n.id));
  return (
    <>
      <UnavailableNotice>
        Graph connections, publication, advancement rules, and migration aren’t
        supported by the running API yet. Graph edits stay local and are never
        flattened into an ordered stage list.
      </UnavailableNotice>
      <p className={e.note}>
        Select a stage to configure it. Drag a stage to arrange it, or focus it
        and use Alt + arrow keys. Scroll the canvas on smaller screens.
      </p>
      <div className={s.editor}>
        <aside className={s.palette} aria-label="Stage palette">
          <h2>Add a stage</h2>
          {stageOptions.map(([type, label]) => (
            <button
              key={type}
              onClick={() => add(type)}
            >
              {label}
            </button>
          ))}
          <Link href="/employer/modules/new?type=recorded_video">
            Video module
          </Link>
          <Link href="/employer/modules/library">Module library</Link>
          <h2>Connections</h2>
          <p>Parallel paths and joins are supported in this local draft.</p>
        </aside>
        <section className={s.workspace} aria-label="Pipeline graph">
          <div className={s.toolbar}>
            <strong>DAG workspace</strong>
            <button
              aria-label="Zoom out"
              onClick={() =>
                setZoom((z) => Math.max(0.5, +(z - 0.1).toFixed(1)))
              }
            >
              −
            </button>
            <output aria-live="polite">{Math.round(zoom * 100)}%</output>
            <button
              aria-label="Zoom in"
              onClick={() =>
                setZoom((z) => Math.min(1.5, +(z + 0.1).toFixed(1)))
              }
            >
              +
            </button>
            <button onClick={() => setZoom(1)}>Reset zoom</button>
          </div>
          <div
            className={s.viewport}
            tabIndex={0}
            aria-label="Scrollable pipeline canvas"
          >
            <div style={{ width: width * zoom, height: (height + 160) * zoom }}>
              <div
                className={s.graph}
                style={{
                  width,
                  height: height + 160,
                  transform: `scale(${zoom})`,
                }}
              >
                <svg
                  className={s.edges}
                  width={width}
                  height={height + 160}
                  aria-hidden="true"
                >
                  <defs>
                    <marker
                      id="dag-arrow"
                      markerWidth="8"
                      markerHeight="8"
                      refX="7"
                      refY="4"
                      orient="auto"
                    >
                      <path d="M0 0L8 4L0 8" fill="#918575" />
                    </marker>
                  </defs>
                  {edges.map((edge) => {
                    const a = nodes.find((n) => n.id === edge.from),
                      b = nodes.find((n) => n.id === edge.to);
                    return a && b ? (
                      <path
                        key={`${edge.from}-${edge.to}`}
                        className={s.edge}
                        markerEnd="url(#dag-arrow)"
                        d={`M${a.x + 105},${a.y + 160} C${a.x + 105},${a.y + 200} ${b.x + 105},${b.y - 40} ${b.x + 105},${b.y}`}
                      />
                    ) : null;
                  })}
                  {roots.map((n) => (
                    <path
                      key={`root-${n.id}`}
                      className={s.edge}
                      markerEnd="url(#dag-arrow)"
                      d={`M350,60 C350,80 ${n.x + 105},80 ${n.x + 105},${n.y}`}
                    />
                  ))}
                  {leaves.map((n) => (
                    <path
                      key={`leaf-${n.id}`}
                      className={s.edge}
                      markerEnd="url(#dag-arrow)"
                      d={`M${n.x + 105},${n.y + 160} C${n.x + 105},${height - 30} 350,${height - 30} 350,${height + 30}`}
                    />
                  ))}
                </svg>
                <div className={s.terminal} style={{ left: 295, top: 12 }}>
                  Entry · All paths
                </div>
                {nodes.map((n) => (
                  <button
                    type="button"
                    className={s.node}
                    key={n.id}
                    style={{ left: n.x, top: n.y }}
                    aria-pressed={n.id === selected}
                    aria-label={`${n.stage.name}, ${edges.filter((edge) => edge.to === n.id).length} incoming connections`}
                    onClick={() => setSelected(n.id)}
                    onKeyDown={(event) => move(n.id, event)}
                    onPointerDown={(event) => start(n, event)}
                    onPointerMove={pointerMove}
                    onPointerUp={() => {
                      drag.current = undefined;
                    }}
                    onPointerCancel={() => {
                      drag.current = undefined;
                    }}
                  >
                    <small>
                      {n.stage.type === "assessment"
                        ? "ASSESSMENT"
                        : "LIVE INTERVIEW"}
                    </small>
                    <h3>{n.stage.name}</h3>
                    <p>
                      {n.stage.type === "assessment"
                        ? `${n.stage.moduleIds?.length || 0} modules`
                        : n.stage.interview
                          ? `${n.stage.interview.durationMinutes} min interview`
                          : "Interview setup"}
                    </p>
                    <small>
                      {edges.filter((edge) => edge.to === n.id).length > 1
                        ? "Join · All incoming stages must pass"
                        : "Select to configure"}
                    </small>
                  </button>
                ))}
                <div
                  className={s.terminal}
                  style={{ left: 235, top: height + 30 }}
                >
                  Finalist · All required paths
                  <br />
                  <small>Hiring is a separate decision</small>
                </div>
              </div>
            </div>
          </div>
        </section>
        <aside
          className={s.inspector}
          aria-label="Selected stage configuration"
        >
          {active ? (
            <>
              <h2>{active.stage.name}</h2>
              <label className={e.field}>
                Unique stage name
                <input
                  value={active.stage.name}
                  onChange={(event) => patch({ name: event.target.value })}
                />
              </label>
              {active.stage.type === "assessment" ? (
                <fieldset>
                  <legend>Assessment modules</legend>
                  <ApiFeedback error={modules.error} retry={modules.reload} />
                  {modules.data?.length ? (
                    modules.data.map((module) => (
                      <label className={e.check} key={module.id}>
                        <input
                          type="checkbox"
                          checked={
                            active.stage.moduleIds?.includes(module.id) || false
                          }
                          onChange={(event) =>
                            patch({
                              moduleIds: event.target.checked
                                ? [...(active.stage.moduleIds || []), module.id]
                                : active.stage.moduleIds?.filter(
                                    (id) => id !== module.id,
                                  ),
                            })
                          }
                        />
                        {module.name}
                      </label>
                    ))
                  ) : (
                    !modules.loading &&
                    !modules.error && (
                      <p>No assessment modules are available yet.</p>
                    )
                  )}
                </fieldset>
              ) : (
                <fieldset>
                  <legend>Interview setup</legend>
                  <p>
                    Interview duration, assignments, and rubric can’t be
                    changed because the API doesn’t support graph-stage
                    configuration yet.
                  </p>
                </fieldset>
              )}
              <fieldset>
                <legend>Schedule & timing</legend>
                {(["scheduledStartAt", "scheduledEndAt"] as const).map(
                  (key, index) => (
                    <label className={e.field} key={key}>
                      {index ? "Stage ends" : "Stage starts"}
                      <input
                        type="datetime-local"
                        value={localDate(active.stage[key])}
                        onChange={(event) =>
                          patch({
                            [key]: event.target.value
                              ? new Date(event.target.value).toISOString()
                              : "",
                          })
                        }
                      />
                    </label>
                  ),
                )}
                <p className={e.note}>
                  Attempt limits and candidate completion windows require
                  backend support.
                </p>
              </fieldset>
              <fieldset>
                <legend>Connections & joins</legend>
                <label className={e.field}>
                  Connect to stage
                  <select
                    value={target}
                    onChange={(event) => setTarget(event.target.value)}
                  >
                    <option value="">Choose a stage</option>
                    {nodes
                      .filter((n) => n.id !== selected)
                      .map((n) => (
                        <option value={n.id} key={n.id}>
                          {n.stage.name}
                        </option>
                      ))}
                  </select>
                </label>
                <button disabled={!target} onClick={connect}>
                  Add connection
                </button>
                {edges
                  .filter((edge) => edge.from === selected)
                  .map((edge) => (
                    <p key={edge.to}>
                      → {nodes.find((n) => n.id === edge.to)?.stage.name}{" "}
                      <button
                        aria-label={`Remove connection to ${nodes.find((n) => n.id === edge.to)?.stage.name}`}
                        onClick={() =>
                          setEdges((items) =>
                            items.filter((item) => item !== edge),
                          )
                        }
                      >
                        Remove
                      </button>
                    </p>
                  ))}
              </fieldset>
              <fieldset>
                <legend>Advancement rule</legend>
                <p>
                  Manual decisions, score thresholds, and completion rules
                  require the planned stage configuration API.
                </p>
              </fieldset>
              <button
                onClick={() => {
                  setNodes((items) => items.filter((n) => n.id !== selected));
                  setEdges((items) =>
                    items.filter(
                      (edge) => edge.from !== selected && edge.to !== selected,
                    ),
                  );
                  setSelected("");
                }}
              >
                Remove stage
              </button>
            </>
          ) : (
            <p>Add or select a stage to configure the graph.</p>
          )}
        </aside>
      </div>
      <p role="status">{message}</p>
    </>
  );
}
function Draft({ pipeline }: { pipeline: PipelineDraft }) {
  const stages = useApiResource<PipelineStage[] | null>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/stages`,
  );
  const transitions = useApiResource<StageTransition[]>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/transitions`,
  );
  const transitionsUnavailable =
    transitions.error && transitions.error.status !== 404;
  return (
    <>
      <BusyState loading={stages.loading || transitions.loading} />
      <ApiFeedback error={stages.error} retry={stages.reload} />
      {transitionsUnavailable && (
        <ApiFeedback error={transitions.error} retry={transitions.reload} />
      )}
      {transitions.error?.status === 404 && (
        <UnavailableNotice>
          Saved graph connections are not available from the current API. The
          displayed stages have no server-provided edges.
        </UnavailableNotice>
      )}
      <StatusBadge>
        Draft {pipeline.versionNumber}
      </StatusBadge>
      {stages.data &&
        !stages.loading &&
        !transitions.loading &&
        !transitionsUnavailable && (
          <GraphEditor
            key={pipeline.versionId}
            initial={stages.data}
            transitions={transitions.data}
          />
        )}
    </>
  );
}
export function PipelineStudio({ postingId }: { postingId: string }) {
  const params = useSearchParams();
  const pathname = usePathname();
  const router = useRouter();
  const state = params.get("state");
  const pipeline = useApiResource<PipelineDraft>(
    state === "published"
      ? null
      : `/v1/jobs/${resourceId(postingId)}/pipeline`,
  );
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState<Error>();
  const createPending = useRef(false);
  function updateDialog(value: PipelineDialog | null) {
    const query = new URLSearchParams(params.toString());
    if (value) query.set("dialog", value);
    else query.delete("dialog");
    const suffix = query.toString();
    router.replace(suffix ? `${pathname}?${suffix}` : pathname, {
      scroll: false,
    });
  }
  function updateState(value: string) {
    const query = new URLSearchParams(params.toString());
    query.set("state", value);
    query.delete("dialog");
    router.replace(`${pathname}?${query.toString()}`, { scroll: false });
  }
  async function createDraft() {
    if (createPending.current) return;
    createPending.current = true;
    setCreating(true);
    setCreateError(undefined);
    try {
      await createPipeline(postingId);
      pipeline.reload();
    } catch (cause) {
      setCreateError(
        cause instanceof Error
          ? cause
          : new Error("The pipeline draft could not be created."),
      );
    } finally {
      createPending.current = false;
      setCreating(false);
    }
  }
  return (
    <>
      <ProductHeading
        title={
          state === "published"
            ? "Pipeline · Published version"
            : state?.startsWith("migration")
              ? "Pipeline · Migration review"
              : "Pipeline designer"
        }
        description={
          state === "published"
            ? "Published pipeline version"
            : "Draft changes are private until explicitly published."
        }
        compact
        artwork={{
          src: "/figma/hiring-dag.svg",
          width: 200.16,
          height: 138.96,
        }}
      />
      <nav className={e.actions} aria-label="Pipeline views">
        <Link href={`/employer/pipelines/${resourceId(postingId)}`}>
          Pipeline
        </Link>
        <Link href={`/employer/pipelines/${resourceId(postingId)}/candidates`}>
          Candidates / decisions
        </Link>
        <Link href="/employer/interviews">Live interviews</Link>
        <Link
          href={`/employer/pipelines/${resourceId(postingId)}?state=migration-review`}
        >
          Versions / migration
        </Link>
        <Link href="/employer/assessments">Assessment instances</Link>
        <Link href={`/employer/pipelines/${resourceId(postingId)}/dashboard`}>
          Dashboard
        </Link>
      </nav>
      {state !== "published" && (
        <>
          <div className={e.actions}>
            <button
              type="button"
              disabled
              aria-describedby="pipeline-save-unavailable"
              title="The current API does not save graph connections."
            >
              Save graph draft
            </button>
            <button type="button" onClick={() => updateState("validation")}>
              Validate
            </button>
            <button
              type="button"
              className={e.primary}
              onClick={() => updateDialog("publish")}
            >
              Review publication
            </button>
            <button type="button" onClick={() => updateDialog("discard")}>
              Discard draft
            </button>
          </div>
          <p id="pipeline-save-unavailable" className={e.note}>
            Saving graph connections is not supported by the current API.
            Local graph edits are not persisted.
          </p>
        </>
      )}
      {state === "validation" && (
        <section className={s.state}>
          <h2>Graph validation</h2>
          <p>
            Validation is unavailable. The server cannot verify reachability,
            branch joins, or publication constraints yet.
          </p>
        </section>
      )}
      {state?.startsWith("migration") && (
        <section className={s.state}>
          <h2>
            {state === "migration-blocked"
              ? "Migration blocked"
              : "Migration review"}
          </h2>
          <p>
            Active candidate placements must be checked against the new graph
            before publication. No migration preview can be loaded yet.
          </p>
          <button type="button" onClick={() => updateDialog("migrate")}>
            Review migration
          </button>
        </section>
      )}
      <BusyState loading={pipeline.loading} />
      {pipeline.error?.status !== 404 && (
        <ApiFeedback error={pipeline.error} retry={pipeline.reload} />
      )}
      {state === "published" && (
        <section className={s.state} aria-labelledby="published-pipeline-title">
          <h2 id="published-pipeline-title">Published pipeline version</h2>
          <UnavailableNotice>
            The current API doesn’t expose the published version for this job.
            No draft data is shown as a published pipeline.
          </UnavailableNotice>
        </section>
      )}
      {pipeline.error?.status === 404 && !pipeline.loading && (
        <section className={s.emptyState} aria-labelledby="pipeline-empty-title">
          <h2 id="pipeline-empty-title">Start a pipeline draft</h2>
          <p>
            Create a draft for this job before adding stages to its pipeline.
          </p>
          <ApiFeedback error={createError} />
          <button
            type="button"
            className={e.primary}
            disabled={creating}
            onClick={createDraft}
          >
            {creating ? "Creating draft…" : "Create pipeline draft"}
          </button>
        </section>
      )}
      {state !== "published" && pipeline.data && (
        <Draft pipeline={pipeline.data} />
      )}
      <WorkspaceFooter />
    </>
  );
}
