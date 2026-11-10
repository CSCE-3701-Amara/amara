import Link from "next/link";
import { useId } from "react";
import { resourceId } from "@/shared/lib/api-client";
import { StatusBadge } from "@/shared/ui/product-ui";
import type { AssessmentInstance } from "../types/product.types";
import styles from "./candidate-pipeline.module.css";

export type CandidatePipelineData = {
  stages: {
    id: string;
    name: string;
    status: string;
    passed?: boolean;
    completed_at?: string;
  }[];
  links?: { from_stage_id: string; to_stage_id: string }[];
};

export function CandidatePipeline({
  data,
  assessments,
  applicationId,
}: {
  data: CandidatePipelineData;
  assessments: AssessmentInstance[];
  applicationId: string;
}) {
  const arrowId = useId();
  const names = new Map(data.stages.map((stage) => [stage.id, stage.name]));
  const incoming = new Map(
    data.stages.map((stage) => [stage.id, [] as string[]]),
  );
  for (const link of data.links ?? []) {
    if (names.has(link.from_stage_id) && names.has(link.to_stage_id))
      incoming.get(link.to_stage_id)!.push(link.from_stage_id);
  }
  const remaining = new Set(names.keys());
  const groups: (typeof data.stages)[] = [];
  while (remaining.size) {
    const ready = data.stages.filter(
      (stage) =>
        remaining.has(stage.id) &&
        incoming.get(stage.id)!.every((id) => !remaining.has(id)),
    );
    // Never hide stages if a malformed provider graph contains a cycle.
    if (!ready.length) {
      groups.push(data.stages.filter((stage) => remaining.has(stage.id)));
      break;
    }
    groups.push(ready);
    ready.forEach((stage) => remaining.delete(stage.id));
  }
  const byStage = new Map<string, AssessmentInstance[]>();
  for (const instance of assessments)
    byStage.set(instance.stageId, [
      ...(byStage.get(instance.stageId) ?? []),
      instance,
    ]);

  const nodeWidth = 300;
  const nodeHeight = 340;
  const columnGap = 140;
  const rowGap = 48;
  const padding = 36;
  const headingHeight = 56;
  const rowCount = Math.max(1, ...groups.map((group) => group.length));
  const contentHeight = rowCount * nodeHeight + (rowCount - 1) * rowGap;
  const width =
    padding * 2 +
    groups.length * nodeWidth +
    Math.max(0, groups.length - 1) * columnGap;
  const height = padding * 2 + headingHeight + contentHeight;
  const positions = new Map<string, { x: number; y: number }>();
  groups.forEach((group, column) => {
    const groupHeight = group.length * nodeHeight + (group.length - 1) * rowGap;
    group.forEach((stage, row) =>
      positions.set(stage.id, {
        x: padding + column * (nodeWidth + columnGap),
        y:
          padding +
          headingHeight +
          (contentHeight - groupHeight) / 2 +
          row * (nodeHeight + rowGap),
      }),
    );
  });

  return (
    <div className={styles.journey}>
      <p className={styles.intro}>
        Follow your path from left to right. Separate branches run in parallel;
        connecting arrows show what must pass before the next stage opens.
      </p>
      <div className={styles.legend}>
        <span>
          <i className={styles.activeDot} />
          Current stages
        </span>
        <span>
          <i className={styles.lockedDot} />
          Upcoming stages
        </span>
        <span>Scroll the map to explore your journey</span>
      </div>
      <div
        className={styles.viewport}
        role="region"
        aria-label="Hiring journey graph"
        tabIndex={0}
      >
        <div className={styles.graph} style={{ width, height }}>
          <svg
            className={styles.connections}
            width={width}
            height={height}
            role="img"
            aria-label="Stage dependency arrows"
          >
            <defs>
              <marker
                id={arrowId}
                markerWidth="10"
                markerHeight="10"
                refX="8"
                refY="4"
                orient="auto"
                markerUnits="userSpaceOnUse"
              >
                <path
                  d="M0,0 L8,4 L0,8"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                />
              </marker>
            </defs>
            {(data.links ?? []).map((link) => {
              const from = positions.get(link.from_stage_id);
              const to = positions.get(link.to_stage_id);
              if (!from || !to) return null;
              const startX = from.x + nodeWidth;
              const startY = from.y + nodeHeight / 2;
              const endY = to.y + nodeHeight / 2;
              const middleX = (startX + to.x) / 2;
              return (
                <path
                  key={`${link.from_stage_id}:${link.to_stage_id}`}
                  data-from={link.from_stage_id}
                  data-to={link.to_stage_id}
                  d={`M ${startX} ${startY} C ${middleX} ${startY}, ${middleX} ${endY}, ${to.x - 8} ${endY}`}
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  markerEnd={`url(#${arrowId})`}
                />
              );
            })}
          </svg>
          {groups.map((group, index) => (
            <section
              className={styles.group}
              key={group.map((stage) => stage.id).join(":")}
              aria-label={`Pipeline step ${index + 1}`}
              style={{
                left: padding + index * (nodeWidth + columnGap),
                top: padding,
                width: nodeWidth,
              }}
            >
              <div className={styles.groupHeading}>
                <h3>Step {index + 1}</h3>
                <span>
                  {group.length > 1
                    ? `${group.length} parallel stages`
                    : "1 stage"}
                </span>
              </div>
              <div className={styles.stages}>
                {group.map((stage) => {
                  const prerequisites = incoming.get(stage.id)!;
                  const instances = byStage.get(stage.id) ?? [];
                  const label =
                    stage.status === "active"
                      ? "Available or in progress"
                      : stage.status === "completed"
                        ? stage.passed
                          ? "Passed"
                          : "Completed"
                        : "Not open yet";
                  return (
                    <article
                      key={stage.id}
                      className={styles.stage}
                      data-state={stage.status}
                      style={{
                        top: positions.get(stage.id)!.y - padding,
                        width: nodeWidth,
                        height: nodeHeight,
                      }}
                    >
                      <StatusBadge tone={stage.passed ? "success" : "neutral"}>
                        {label}
                      </StatusBadge>
                      <h4>{stage.name}</h4>
                      <p className={styles.prerequisite}>
                        {prerequisites.length
                          ? `Opens after ${prerequisites.map((id) => names.get(id)).join(" and ")} pass.`
                          : "Entry stage · no earlier stages required."}
                      </p>
                      {instances.map((instance) => (
                        <div className={styles.assessment} key={instance.id}>
                          <span>{instance.status.replaceAll("_", " ")}</span>
                          {instance.deadline && (
                            <p>
                              Due {new Date(instance.deadline).toLocaleString()}
                            </p>
                          )}
                          <Link
                            href={`/candidate/assessments/${resourceId(instance.id)}?state=availability&application=${resourceId(applicationId)}`}
                          >
                            {instance.status === "in_progress"
                              ? "Continue assessment"
                              : instance.status === "available"
                                ? "Open assessment"
                                : "Review assessment"}
                          </Link>
                        </div>
                      ))}
                      {!instances.length && stage.status !== "completed" && (
                        <p>
                          {stage.status === "active"
                            ? "The hiring team will share the next action here."
                            : "Your next action will appear when this stage opens."}
                        </p>
                      )}
                    </article>
                  );
                })}
              </div>
            </section>
          ))}
        </div>
      </div>
      {!data.stages.length && (
        <p>No stages have been published for this application yet.</p>
      )}
    </div>
  );
}
