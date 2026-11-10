import { expect, test } from "@playwright/test";
import {
  addConnection,
  graphDimensions,
  normalizeGraph,
  removeConnection,
  removeStage,
  wouldCreateCycle,
  type PipelineGraph,
} from "../../src/features/invw/components/dag-graph";

const graph = (links: PipelineGraph["links"] = []): PipelineGraph => ({
  stages: [
    {
      id: "a",
      name: "A",
      type: "assessment",
      scheduledStartAt: "",
      scheduledEndAt: "",
      x: 0,
      y: 0,
    },
    {
      id: "b",
      name: "B",
      type: "assessment",
      scheduledStartAt: "",
      scheduledEndAt: "",
      x: 100,
      y: 100,
    },
    {
      id: "c",
      name: "C",
      type: "assessment",
      scheduledStartAt: "",
      scheduledEndAt: "",
      x: 200,
      y: 200,
    },
  ],
  links,
});

test("normalizes only missing coordinates and supplies missing links", () => {
  const normalized = normalizeGraph({
    stages: [
      {
        id: "a",
        name: "A",
        type: "assessment",
        scheduledStartAt: "",
        scheduledEndAt: "",
        x: 0,
        y: 0,
      },
      {
        id: "b",
        name: "B",
        type: "assessment",
        scheduledStartAt: "",
        scheduledEndAt: "",
        x: undefined!,
        y: undefined!,
      },
    ],
    links: undefined!,
  });
  expect(normalized.stages[0]).toMatchObject({ x: 0, y: 0 });
  expect(normalized.stages[1]).toMatchObject({ x: 300, y: 40 });
  expect(normalized.links).toEqual([]);
});

test("dimensions enforce minimums and grow with node positions", () => {
  expect(graphDimensions({ stages: [], links: [] })).toEqual({
    width: 1000,
    height: 800,
  });
  expect(graphDimensions(graph())).toEqual({ width: 1000, height: 800 });
  const large = graph();
  large.stages[2] = { ...large.stages[2], x: 900, y: 700 };
  expect(graphDimensions(large)).toEqual({ width: 1140, height: 880 });
});

test("rejects self, direct, and transitive cycles but permits a branch", () => {
  const chain = graph([
    { fromStageId: "a", toStageId: "b" },
    { fromStageId: "b", toStageId: "c" },
  ]);
  expect(wouldCreateCycle(chain, "a", "a")).toBe(true);
  expect(wouldCreateCycle(chain, "b", "a")).toBe(true);
  expect(wouldCreateCycle(chain, "c", "a")).toBe(true);
  expect(wouldCreateCycle(chain, "a", "c")).toBe(false);
});

test("connection mutations avoid duplicates and remove incident links", () => {
  const initial = graph([{ fromStageId: "a", toStageId: "b" }]);
  expect(addConnection(initial, "a", "b")).toBe(initial);
  const connected = addConnection(initial, "b", "c");
  expect(connected.links).toHaveLength(2);
  expect(removeConnection(connected, connected.links[0]).links).toEqual([
    { fromStageId: "b", toStageId: "c" },
  ]);
  expect(removeStage(connected, "b")).toEqual({
    stages: [initial.stages[0], initial.stages[2]],
    links: [],
  });
});
