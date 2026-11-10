import styles from "./dag-editor.module.css";
import { graphDimensions, type DagNode, type PipelineGraph } from "./dag-graph";

type DagCanvasProps = {
  graph: PipelineGraph;
  selectedStageId: string;
  zoom: number;
  disabled: boolean;
  onSelect: (stageId: string) => void;
  onUpdate: (stageId: string, patch: Partial<DagNode>) => void;
};

export function DagCanvas({
  graph,
  selectedStageId,
  zoom,
  disabled,
  onSelect,
  onUpdate,
}: DagCanvasProps) {
  const { width, height } = graphDimensions(graph);

  return (
    <div className={styles.viewport}>
      {graph.stages.length === 0 && (
        <div className={styles.emptyCanvas}>
          <strong>Start with your first stage</strong>
          <p>
            Add an assessment or interview to begin designing the candidate
            journey.
          </p>
        </div>
      )}
      <div
        className={styles.canvas}
        style={{ transform: `scale(${zoom / 100})`, width, height }}
      >
        <svg
          className={styles.edges}
          width={width}
          height={height}
          aria-label="Stage connections"
        >
          {graph.links.map((edge) => {
            const from = graph.stages.find(
              (stage) => stage.id === edge.fromStageId,
            );
            const to = graph.stages.find(
              (stage) => stage.id === edge.toStageId,
            );
            return from && to ? (
              <g key={from.id + to.id}>
                <path
                  d={`M${from.x + 200},${from.y + 65} C${from.x + 240},${from.y + 65} ${to.x - 40},${to.y + 65} ${to.x},${to.y + 65}`}
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                />
                <text
                  x={(from.x + to.x + 200) / 2}
                  y={(from.y + to.y) / 2 + 55}
                >
                  →
                </text>
              </g>
            ) : null;
          })}
        </svg>
        {graph.stages.map((stage) => (
          <button
            key={stage.id}
            className={styles.node}
            style={{ left: stage.x, top: stage.y }}
            aria-pressed={selectedStageId === stage.id}
            aria-label={`${stage.name}, ${graph.links.filter((edge) => edge.toStageId === stage.id).length} incoming connections`}
            onClick={() => onSelect(stage.id)}
            onKeyDown={(event) => {
              if (disabled || !event.altKey || !event.key.startsWith("Arrow"))
                return;
              event.preventDefault();
              onUpdate(stage.id, {
                x: Math.max(
                  0,
                  stage.x +
                    (event.key === "ArrowRight"
                      ? 20
                      : event.key === "ArrowLeft"
                        ? -20
                        : 0),
                ),
                y: Math.max(
                  0,
                  stage.y +
                    (event.key === "ArrowDown"
                      ? 20
                      : event.key === "ArrowUp"
                        ? -20
                        : 0),
                ),
              });
            }}
            draggable={!disabled}
            onDragEnd={(event) => {
              if (disabled) return;
              const parent = event.currentTarget.parentElement;
              if (!parent) return;
              const rect = parent.getBoundingClientRect();
              onUpdate(stage.id, {
                x: Math.max(
                  0,
                  ((event.clientX - rect.left) * 100) / zoom - 100,
                ),
                y: Math.max(0, ((event.clientY - rect.top) * 100) / zoom - 40),
              });
            }}
          >
            <strong>{stage.name}</strong>
            <span>
              {stage.type === "assessment"
                ? `${stage.moduleIds?.length || 0} assessment modules`
                : "Live interview"}
            </span>
            <small>Alt + arrow keys to move</small>
          </button>
        ))}
      </div>
    </div>
  );
}
