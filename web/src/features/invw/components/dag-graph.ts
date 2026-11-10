import type { PipelineStage } from "../types/invw.types";

export type DagNode = PipelineStage & {
  id: string;
  x: number;
  y: number;
};

export type DagEdge = {
  fromStageId: string;
  toStageId: string;
};

export type PipelineGraph = {
  stages: DagNode[];
  links: DagEdge[];
};

export function normalizeGraph(graph: PipelineGraph): PipelineGraph {
  return {
    stages: graph.stages.map((stage, index) => ({
      ...stage,
      x: stage.x ?? 40 + (index % 3) * 260,
      y: stage.y ?? 40 + Math.floor(index / 3) * 200,
    })),
    links: graph.links ?? [],
  };
}

export function graphDimensions(graph: PipelineGraph) {
  return {
    width: Math.max(1000, ...graph.stages.map((stage) => stage.x + 240)),
    height: Math.max(800, ...graph.stages.map((stage) => stage.y + 180)),
  };
}

export function wouldCreateCycle(
  graph: PipelineGraph,
  fromStageId: string,
  toStageId: string,
): boolean {
  const visited = new Set<string>();
  const pending = [toStageId];

  while (pending.length > 0) {
    const stageId = pending.pop()!;
    if (stageId === fromStageId) return true;
    if (visited.has(stageId)) continue;
    visited.add(stageId);
    graph.links
      .filter((edge) => edge.fromStageId === stageId)
      .forEach((edge) => pending.push(edge.toStageId));
  }

  return false;
}

export function addConnection(
  graph: PipelineGraph,
  fromStageId: string,
  toStageId: string,
): PipelineGraph {
  const exists = graph.links.some(
    (edge) => edge.fromStageId === fromStageId && edge.toStageId === toStageId,
  );
  return exists
    ? graph
    : {
        ...graph,
        links: [...graph.links, { fromStageId, toStageId }],
      };
}

export function removeConnection(
  graph: PipelineGraph,
  edge: DagEdge,
): PipelineGraph {
  return {
    ...graph,
    links: graph.links.filter((candidate) => candidate !== edge),
  };
}

export function removeStage(
  graph: PipelineGraph,
  stageId: string,
): PipelineGraph {
  return {
    stages: graph.stages.filter((stage) => stage.id !== stageId),
    links: graph.links.filter(
      (edge) => edge.fromStageId !== stageId && edge.toStageId !== stageId,
    ),
  };
}
