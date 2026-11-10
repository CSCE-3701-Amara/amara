import type { DagEdge, DagNode, PipelineGraph } from "./dag-graph";
import styles from "./dag-editor.module.css";

export function DagConnections({
  graph,
  node,
  target,
  onTargetChange,
  onConnect,
  onRemoveConnection,
  onRemoveStage,
}: {
  graph: PipelineGraph;
  node: DagNode;
  target: string;
  onTargetChange: (stageId: string) => void;
  onConnect: () => void;
  onRemoveConnection: (edge: DagEdge) => void;
  onRemoveStage: () => void;
}) {
  return (
    <>
      <div className={styles.connectionHeading}>
        <h3>What comes next?</h3>
        <p>
          Connect this stage to the next step. A stage with multiple incoming
          connections opens after all of them pass.
        </p>
      </div>
      <label>
        Connect to stage
        <select
          value={target}
          onChange={(event) => onTargetChange(event.target.value)}
        >
          <option value="">Choose a stage</option>
          {graph.stages.map((stage) => (
            <option key={stage.id} value={stage.id}>
              {stage.name}
            </option>
          ))}
        </select>
      </label>
      <button onClick={onConnect}>Add connection</button>
      {graph.links
        .filter((edge) => edge.fromStageId === node.id)
        .map((edge) => (
          <button key={edge.toStageId} onClick={() => onRemoveConnection(edge)}>
            Remove connection to{" "}
            {graph.stages.find((stage) => stage.id === edge.toStageId)?.name}
          </button>
        ))}
      <button onClick={onRemoveStage}>Remove stage</button>
    </>
  );
}
