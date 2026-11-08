"use client";
import Link from "next/link";
import Image from "next/image";
import {
  useMemo,
  useRef,
  useState,
  type PointerEvent,
  type KeyboardEvent,
} from "react";
import { useSearchParams } from "next/navigation";
import { resourceId } from "@/shared/lib/api-client";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { useUnsavedChanges } from "@/shared/lib/use-unsaved-changes";
import type {
  StageTransition,
  PipelineMigrationPreview,
} from "@/shared/lib/contract.types";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import {
  ConfirmDialog,
  ProductHeading,
  StatusBadge,
  UnavailableNotice,
  WorkspaceFooter,
} from "@/shared/ui/product-ui";
import { PipelineEditor } from "./pipeline-editor";
import type {
  PipelineDraft,
  PipelineStage,
  ModuleSummary,
} from "../types/invw.types";
import e from "@/shared/ui/editor.module.css";
import s from "./pipeline-studio.module.css";
type StageKind = "question" | "coding" | "video" | "interview";
type GraphNode = {
  id: string;
  stage: PipelineStage;
  x: number;
  y: number;
  kind?: StageKind;
};
type Edge = { from: string; to: string };
const stageKinds = [
  {
    kind: "question",
    label: "Question-based",
    detail: "Questions & formats",
    icon: "list-checks",
  },
  {
    kind: "coding",
    label: "Coding challenge",
    detail: "Code & test cases",
    icon: "code-xml",
  },
  {
    kind: "video",
    label: "Video-recorded",
    detail: "Prompts & graders",
    icon: "video",
  },
  {
    kind: "interview",
    label: "Live interview",
    detail: "Slots & scorecards",
    icon: "calendar-days",
  },
] as const;
function StudioIcon({ name, size = 15.36 }: { name: string; size?: number }) {
  return (
    <Image
      src={`/figma/pipeline/${name}.svg`}
      width={size}
      height={size}
      alt=""
      unoptimized
    />
  );
}
/** Lay out actual dependencies; stage ordering never creates graph connections. */
function arrange(
  stages: PipelineStage[],
  transitions: StageTransition[] = [],
): GraphNode[] {
  const depths = new Map(
    stages.map((stage, index) => [stage.id || `stage-${index}`, 0]),
  );
  for (let pass = 0; pass < stages.length; pass++) {
    let changed = false;
    for (const edge of transitions) {
      if (!depths.has(edge.fromStageId) || !depths.has(edge.toStageId))
        continue;
      const depth = Math.min(
        stages.length - 1,
        depths.get(edge.fromStageId)! + 1,
      );
      if (depth > depths.get(edge.toStageId)!) {
        depths.set(edge.toStageId, depth);
        changed = true;
      }
    }
    if (!changed) break;
  }
  const groups = new Map<number, string[]>();
  for (const [id, depth] of depths)
    groups.set(depth, [...(groups.get(depth) || []), id]);
  return stages.map((stage, index) => {
    const id = stage.id || `stage-${index}`,
      depth = depths.get(id)!,
      group = groups.get(depth)!;
    return {
      id,
      stage,
      x: 176 + depth * 290,
      y:
        Math.max(245, (group.length - 1) * 152 + 90) +
        (group.indexOf(id) - (group.length - 1) / 2) * 304,
    };
  });
}
/** A real DAG model: branches and joins stay explicit, independent of stage ordering. */
function GraphEditor({
  initial,
  transitions,
  published,
}: {
  initial: PipelineStage[];
  transitions?: StageTransition[];
  published: boolean;
}) {
  const [nodes, setNodes] = useState<GraphNode[]>(() =>
    arrange(initial, transitions),
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
  const [dirty, setDirty] = useState(false);
  const history = useRef<{ nodes: GraphNode[]; edges: Edge[] }[]>([]);
  const future = useRef<{ nodes: GraphNode[]; edges: Edge[] }[]>([]);
  const [historySize, setHistorySize] = useState({ past: 0, future: 0 });
  const viewport = useRef<HTMLDivElement>(null);
  const connection = useRef<HTMLSelectElement>(null);
  useUnsavedChanges(dirty);
  function remember() {
    history.current = [...history.current.slice(-49), { nodes, edges }];
    future.current = [];
    setHistorySize({ past: history.current.length, future: 0 });
    setDirty(true);
  }
  function restore(redo: boolean) {
    const source = redo ? future : history,
      destination = redo ? history : future;
    const snapshot = source.current.pop();
    if (!snapshot) return;
    destination.current.push({ nodes, edges });
    setNodes(snapshot.nodes);
    setEdges(snapshot.edges);
    setHistorySize({
      past: history.current.length,
      future: future.current.length,
    });
    setDirty(true);
    setMessage(
      redo ? "Graph change redone locally." : "Graph change undone locally.",
    );
  }
  const drag = useRef<
    | {
        id: string;
        x: number;
        y: number;
        left: number;
        top: number;
        moved?: boolean;
      }
    | undefined
  >(undefined);
  const sequence = useRef(0);
  const modules = useApiResource<ModuleSummary[] | null>("/v1/modules");
  const nodeById = new Map(nodes.map((node) => [node.id, node]));
  const active = nodeById.get(selected);
  const { incoming, outgoing } = useMemo(() => {
    const incoming = new Map<string, number>(),
      outgoing = new Map<string, number>();
    for (const edge of edges) {
      incoming.set(edge.to, (incoming.get(edge.to) || 0) + 1);
      outgoing.set(edge.from, (outgoing.get(edge.from) || 0) + 1);
    }
    return { incoming, outgoing };
  }, [edges]);
  const width = Math.max(700, ...nodes.map((n) => n.x + 250));
  const height = Math.max(710, ...nodes.map((n) => n.y + 230));
  function patch(change: Partial<PipelineStage>) {
    remember();
    setNodes((items) =>
      items.map((n) =>
        n.id === selected ? { ...n, stage: { ...n.stage, ...change } } : n,
      ),
    );
    setMessage("Unsaved local graph changes.");
  }
  function add(kind: StageKind) {
    remember();
    const type = kind === "interview" ? "live_interview" : "assessment";
    sequence.current += 1;
    const id = `local-${sequence.current}`;
    setNodes((items) => [
      ...items,
      {
        id,
        kind,
        x: 80 + (items.length % 2) * 300,
        y: 120 + Math.floor(items.length / 2) * 220,
        stage: {
          name: stageKinds.find((item) => item.kind === kind)!.label,
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
    remember();
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
    if (!event.altKey || !shifts[event.key] || published) return;
    event.preventDefault();
    remember();
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
    if (published) return;
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
    if (!d.moved) {
      if (Math.hypot(event.clientX - d.x, event.clientY - d.y) < 3) return;
      remember();
      d.moved = true;
    }
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
  const roots = nodes.filter((n) => !incoming.has(n.id));
  const leaves = nodes.filter((n) => !outgoing.has(n.id));
  function removeSelected() {
    if (!active || published) return;
    remember();
    setNodes((items) => items.filter((n) => n.id !== selected));
    setEdges((items) =>
      items.filter((edge) => edge.from !== selected && edge.to !== selected),
    );
    setSelected("");
    setMessage("Stage and its connections removed locally.");
  }
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
          <h2>Add stage</h2>
          {stageKinds.map(({ kind, label, detail, icon }) => (
            <button key={kind} disabled={published} onClick={() => add(kind)}>
              <span className={s.paletteIcons}>
                <StudioIcon name={icon} size={19.2} />
                <StudioIcon name="grip-vertical" size={13.44} />
              </span>
              <strong>{label}</strong>
              <span>{detail}</span>
            </button>
          ))}
          <div className={s.paletteTerminal}>
            <StudioIcon name="flag" /> Finalist
            <p>End of assessment. Not an offer or hire.</p>
          </div>
          <Link href="/employer/modules/library">Module library</Link>
        </aside>
        <section className={s.workspace} aria-label="Pipeline graph">
          <div className={s.toolbar}>
            <div className={s.toolGroup}>
              <button onClick={() => viewport.current?.focus()}>
                <StudioIcon name="mouse-pointer2" /> Select / move
              </button>
              <button
                disabled={!active || published}
                onClick={() => connection.current?.focus()}
              >
                <StudioIcon name="workflow" /> Connect
              </button>
              <button disabled={!active || published} onClick={removeSelected}>
                <StudioIcon name="trash2" /> Delete
              </button>
            </div>
            <div className={s.toolGroup}>
              <button
                aria-label="Undo graph change"
                disabled={!historySize.past || published}
                onClick={() => restore(false)}
              >
                <StudioIcon name="undo2" />
              </button>
              <button
                aria-label="Redo graph change"
                disabled={!historySize.future || published}
                onClick={() => restore(true)}
              >
                <StudioIcon name="redo2" />
              </button>
              <button
                aria-label="Zoom out"
                onClick={() =>
                  setZoom((z) => Math.max(0.5, +(z - 0.1).toFixed(1)))
                }
              >
                <StudioIcon name="minus" />
              </button>
              <output aria-live="polite">{Math.round(zoom * 100)}%</output>
              <button
                aria-label="Zoom in"
                onClick={() =>
                  setZoom((z) => Math.min(1.5, +(z + 0.1).toFixed(1)))
                }
              >
                <StudioIcon name="plus" />
              </button>
              <button
                aria-label="Fit graph to viewport"
                onClick={() => {
                  setZoom(
                    Math.max(
                      0.3,
                      Math.min(
                        1,
                        (viewport.current?.clientWidth || width) / width,
                      ),
                    ),
                  );
                  viewport.current?.scrollTo({ top: 0, left: 0 });
                }}
              >
                <StudioIcon name="maximize" /> Fit
              </button>
              <button className={s.resetZoom} onClick={() => setZoom(1)}>
                Reset zoom
              </button>
            </div>
          </div>
          <div
            className={s.viewport}
            ref={viewport}
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
                    const a = nodeById.get(edge.from),
                      b = nodeById.get(edge.to);
                    return a && b ? (
                      <path
                        key={`${edge.from}-${edge.to}`}
                        className={s.edge}
                        markerEnd="url(#dag-arrow)"
                        d={`M${a.x + 234},${a.y + 85} C${a.x + 265},${a.y + 85} ${b.x - 30},${b.y + 85} ${b.x},${b.y + 85}`}
                      />
                    ) : null;
                  })}
                  {roots.map((n) => (
                    <path
                      key={`root-${n.id}`}
                      className={s.edge}
                      markerEnd="url(#dag-arrow)"
                      d={`M129,330 C150,330 ${n.x - 25},${n.y + 85} ${n.x},${n.y + 85}`}
                    />
                  ))}
                  {leaves.map((n) => (
                    <path
                      key={`leaf-${n.id}`}
                      className={s.edge}
                      markerEnd="url(#dag-arrow)"
                      d={`M${n.x + 117},${n.y + 171} C${n.x + 117},${height - 30} 350,${height - 30} 350,${height + 30}`}
                    />
                  ))}
                </svg>
                <p className={s.graphContext}>
                  {nodes.length} stages · {edges.length} directed links · All
                  dates in utc
                </p>
                <div className={s.terminal} style={{ left: 28, top: 302 }}>
                  Entry<small>Start here</small>
                </div>
                {nodes.map((n) => (
                  <button
                    type="button"
                    className={s.node}
                    key={n.id}
                    style={{ left: n.x, top: n.y }}
                    aria-pressed={n.id === selected}
                    aria-label={`${n.stage.name}, ${incoming.get(n.id) || 0} incoming connections`}
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
                    <span className={s.incomingPort}>
                      <StudioIcon name="incoming-port" size={7.68} />
                    </span>
                    <span className={s.outgoingPort}>
                      <StudioIcon name="incoming-port" size={7.68} />
                    </span>
                    <span className={s.nodeType}>
                      <small>
                        {n.stage.type === "assessment"
                          ? stageKinds.find((kind) => kind.kind === n.kind)
                              ?.label || "Assessment"
                          : "Live interview"}
                      </small>
                      <StudioIcon name="grip-vertical" size={13.44} />
                    </span>
                    <h3>{n.stage.name}</h3>
                    <p>{n.stage.moduleIds?.length || 0} modules</p>
                    <p>
                      {n.stage.scheduledStartAt
                        ? new Date(n.stage.scheduledStartAt).toLocaleDateString(
                            undefined,
                            { day: "2-digit", month: "short", timeZone: "UTC" },
                          )
                        : "Schedule needed"}
                      {n.stage.scheduledEndAt
                        ? ` – ${new Date(n.stage.scheduledEndAt).toLocaleDateString(undefined, { day: "2-digit", month: "short", timeZone: "UTC" })}`
                        : ""}
                    </p>
                    <small className={s.nodeState}>
                      <StudioIcon name="circle-check" size={11.52} />
                      {(incoming.get(n.id) || 0) > 1
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
          <div className={s.validation}>
            <strong>Local graph draft</strong>
            <span>{dirty ? "Unsaved changes" : "Loaded stages"}</span>
            <p>
              Directed links and joins remain explicit. Publication checks and
              migration review are unavailable.
            </p>
            <Link href="?state=migration-review">Migration review →</Link>
          </div>
        </section>
        <aside
          className={s.inspector}
          aria-label="Selected stage configuration"
        >
          {active ? (
            <>
              <div className={s.inspectorHeading}>
                <h2>{active.stage.name}</h2>
                <span>Manager only</span>
              </div>
              <label className={e.field}>
                Unique stage name
                <input
                  value={active.stage.name}
                  readOnly={published}
                  onChange={(event) => patch({ name: event.target.value })}
                />
              </label>
              <fieldset>
                <legend>Assessment modules</legend>
                <ApiFeedback error={modules.error} retry={modules.reload} />
                {modules.data?.map((module) => (
                  <label className={e.check} key={module.id}>
                    <input
                      type="checkbox"
                      disabled={published}
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
                ))}
              </fieldset>
              <fieldset>
                <legend>Schedule & timing</legend>
                <div className={s.schedule}>
                  {(["scheduledStartAt", "scheduledEndAt"] as const).map(
                    (key, index) => (
                      <label className={e.field} key={key}>
                        {index ? "Stage ends" : "Stage starts"}
                        <input
                          type="datetime-local"
                          disabled={published}
                          value={
                            active.stage[key]
                              ? active.stage[key].slice(0, 16)
                              : ""
                          }
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
                </div>
                <div className={s.schedule}>
                  <label className={e.field}>
                    Active attempt limit
                    <input disabled placeholder="Unavailable" />
                  </label>
                  <label className={e.field}>
                    Candidate completion window
                    <input disabled placeholder="Unavailable" />
                  </label>
                </div>
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
                    ref={connection}
                    value={target}
                    disabled={published}
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
                <button disabled={!target || published} onClick={connect}>
                  Add connection
                </button>
                {edges
                  .filter((edge) => edge.from === selected)
                  .map((edge) => (
                    <p key={edge.to}>
                      → {nodeById.get(edge.to)?.stage.name}{" "}
                      <button
                        disabled={published}
                        aria-label={`Remove connection to ${nodeById.get(edge.to)?.stage.name}`}
                        onClick={() => {
                          remember();
                          setEdges((items) =>
                            items.filter((item) => item !== edge),
                          );
                        }}
                      >
                        Remove
                      </button>
                    </p>
                  ))}
              </fieldset>
              <fieldset>
                <legend>Advancement rule</legend>
                <div
                  className={s.ruleChoices}
                  aria-label="Advancement rules unavailable"
                >
                  <button disabled>Manual</button>
                  <button disabled>Threshold</button>
                  <button disabled>On completion</button>
                </div>
                <p>
                  Manual decisions, score thresholds, and completion rules
                  require the planned stage configuration API.
                </p>
              </fieldset>
              <button disabled={published} onClick={removeSelected}>
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
function Draft({
  pipeline,
  published,
}: {
  pipeline: PipelineDraft;
  published: boolean;
}) {
  const stages = useApiResource<PipelineStage[] | null>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/stages`,
  );
  const transitions = useApiResource<StageTransition[]>(
    `/v1/pipeline-versions/${resourceId(pipeline.versionId)}/transitions`,
  );
  return (
    <>
      <BusyState loading={stages.loading || transitions.loading} />
      <ApiFeedback error={stages.error} retry={stages.reload} />
      {transitions.error?.status !== 404 && (
        <ApiFeedback error={transitions.error} retry={transitions.reload} />
      )}
      <StatusBadge>
        {published ? "Published version" : "Draft"} {pipeline.versionNumber}
      </StatusBadge>
      {!stages.loading && !stages.error && !transitions.loading && (
        <GraphEditor
          key={pipeline.versionId}
          initial={stages.data || []}
          transitions={transitions.data}
          published={published}
        />
      )}
    </>
  );
}
export function PipelineStudio({ postingId }: { postingId: string }) {
  const params = useSearchParams();
  const state = params.get("state");
  const [dialog, setDialog] = useState<string | null>(null);
  const pipeline = useApiResource<PipelineDraft>(
    `/v1/jobs/${resourceId(postingId)}/pipeline`,
  );
  const migration = useApiResource<PipelineMigrationPreview>(null);
  return (
    <>
      <section
        className={s.header}
        aria-label="Pipeline navigation and version"
      >
        <ProductHeading
          title={
            state === "published"
              ? "Pipeline · Published version"
              : state?.startsWith("migration")
                ? "Pipeline · Migration review"
                : "Pipeline designer"
          }
          description="Draft changes are private until explicitly published."
          compact
          className={s.heading}
        />
        <nav className={e.actions} aria-label="Pipeline views">
          <Link href={`/employer/pipelines/${resourceId(postingId)}`}>
            Pipeline
          </Link>
          <Link
            href={`/employer/pipelines/${resourceId(postingId)}/candidates`}
          >
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
        <div className={e.actions}>
          <button disabled>Save graph draft</button>
          <button onClick={() => setDialog("validate")}>Validate</button>
          <button className={e.primary} onClick={() => setDialog("publish")}>
            Review publication
          </button>
          <button onClick={() => setDialog("discard")}>Discard draft</button>
        </div>
        <p className={s.versionNote}>
          {pipeline.data
            ? `Version ${pipeline.data.versionNumber}`
            : "Pipeline draft"}{" "}
          · {state === "published" ? "Read only" : "Private draft"}
        </p>
      </section>
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
          {migration.data?.affectedExecutions}
          <button disabled>Confirm migration</button>
        </section>
      )}
      <BusyState loading={pipeline.loading} />
      {pipeline.error?.status !== 404 && (
        <ApiFeedback error={pipeline.error} retry={pipeline.reload} />
      )}{" "}
      {pipeline.data ? (
        <Draft pipeline={pipeline.data} published={state === "published"} />
      ) : (
        !pipeline.loading && <GraphEditor initial={[]} published={false} />
      )}
      <section className={e.section} aria-label="Privacy and feedback policy">
        <h2>Privacy and feedback</h2>
        <div className={e.grid}>
          <label className={e.field}>
            Candidate anonymity
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Identity reveal point
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Rejection feedback
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
          <label className={e.field}>
            Selected stage feedback
            <select disabled>
              <option>Policy unavailable</option>
            </select>
          </label>
        </div>
        <p>
          Identity must be revealed before incompatible video and interview
          stages on every path. Policy loading and saving are not available yet.
        </p>
        <div className={e.actions}>
          <button disabled>Save policy</button>
        </div>
      </section>
      <details>
        <summary>Server-supported stage configuration</summary>
        <p>
          This existing editor saves stage definitions through the current API.
          It does not save graph connections.
        </p>
        <PipelineEditor postingId={postingId} />
      </details>
      <WorkspaceFooter />
      {dialog && (
        <ConfirmDialog
          title={
            dialog === "discard"
              ? "Discard pipeline draft?"
              : dialog === "migrate"
                ? "Confirm pipeline migration?"
                : dialog === "validate"
                  ? "Validate pipeline graph?"
                  : "Publish pipeline version?"
          }
          description="Review the graph, schedules, candidate placements, and required stage configuration before proceeding."
          confirmLabel={
            dialog === "discard"
              ? "Discard draft"
              : dialog === "validate"
                ? "Validate"
                : "Confirm publication"
          }
          intent={dialog === "discard" ? "danger" : "warning"}
          unavailable="The required graph operation is not available yet. Your local graph has not been published or discarded."
          onClose={() => setDialog(null)}
        />
      )}
    </>
  );
}
