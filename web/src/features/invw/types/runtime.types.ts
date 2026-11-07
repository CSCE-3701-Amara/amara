/** Runtime projection returned by GET /v1/jobs/{postingId}/pipeline/dashboard. */
export interface RuntimePipelineDashboard {
  finalistTarget: number | null;
  active: number;
  links: { from: string; to: string }[];
  stages: {
    id: string;
    name: string;
    stage_order: number;
    rule: "manual" | "threshold" | "completion";
    threshold: number | null;
    active: number;
    passed: number;
    rejected: number;
    candidates: {
      id: string;
      candidate: string;
      status: string;
      active: boolean;
      score: number | null;
      niceToHaveScore: number | null;
      percentile: number | null;
    }[];
  }[];
}
