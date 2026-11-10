const labels: Record<string, string> = {
  "/candidate/dashboard": "opportunities",
  "/candidate/applications": "applications",
  "/candidate/offers": "offers",
  "/candidate/profile": "profile",
  "/employer/jobs": "job postings",
  "/employer/pipelines": "pipelines",
  "/employer/modules": "modules",
  "/employer/modules/library": "module library",
  "/employer/grading": "grading",
  "/employer/offers": "offers",
  "/interviewer/interviews": "interviews",
  "/admin/dashboard": "company registrations",
};

/** Return to an existing parent screen, including for directly opened URLs. */
export function workspaceBackLink(path: string, home: string, query = "") {
  if (path === home) return null;

  let href: string;
  if (/^\/candidate\/jobs\//.test(path)) {
    href = "/candidate/dashboard";
  } else if (/^\/candidate\/(assessments|interviews)\/[^/]+$/.test(path)) {
    href = "/candidate/applications";
  } else if (/^\/employer\/applications\//.test(path)) {
    href = "/employer/jobs";
  } else if (/^\/employer\/feedback\//.test(path)) {
    href = "/employer/grading";
  } else if (/^\/interviewer\/interviews\//.test(path)) {
    href = "/interviewer/interviews";
  } else {
    const segments = path.split("/").filter(Boolean);
    href = segments.length > 2 ? "/" + segments.slice(0, -1).join("/") : home;
  }

  const parentSection = href.split("/")[2];
  const detailLabels: Record<string, string> = {
    applications: "application",
    assessments: "assessment",
    jobs: "job posting",
    modules: "module",
    offers: "offer",
    pipelines: "pipeline",
  };
  const label = labels[href] || detailLabels[parentSection] || "workspace";
  if (href.startsWith("/candidate/assessments/")) {
    // Variant screens must return to the same attempt, rather than treating
    // the assessment-instance ID as an attempt ID. Do not reopen dialogs.
    const current = new URLSearchParams(query);
    const context = new URLSearchParams();
    for (const key of ["attempt", "application", "state"]) {
      const value = current.get(key);
      if (value) context.set(key, value);
    }
    if (context.size) href += "?" + context.toString();
  }
  return { href, label: `Back to ${label}` };
}
