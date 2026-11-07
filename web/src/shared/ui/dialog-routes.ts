/** Figma dialog deep links, separate from page routes and API operations. */
export const dialogRoutes = [
  {
    id: "212:329",
    path: "/candidate/applications/[id]",
    dialog: "withdraw-review",
  },
  {
    id: "376:555",
    path: "/candidate/interviews/[id]",
    dialog: "confirm",
  },
  {
    id: "376:564",
    path: "/candidate/interviews/[id]",
    dialog: "reschedule",
  },
  {
    id: "376:573",
    path: "/candidate/offers/[id]",
    dialog: "accept",
  },
  {
    id: "376:582",
    path: "/candidate/offers/[id]",
    dialog: "decline",
  },
  {
    id: "376:591",
    path: "/candidate/profile/edit",
    dialog: "deactivate",
  },
  {
    id: "376:600",
    path: "/candidate/profile/edit",
    dialog: "delete",
  },
  {
    id: "377:553",
    path: "/candidate/assessments/[id]",
    dialog: "submit",
  },
  {
    id: "378:571",
    path: "/candidate/applications/[id]",
    dialog: "withdraw",
  },
  {
    id: "376:1131",
    path: "/employer/jobs/[id]",
    dialog: "publish",
  },
  {
    id: "376:1140",
    path: "/employer/jobs/[id]",
    dialog: "archive",
  },
  {
    id: "376:1149",
    path: "/employer/jobs/[id]",
    dialog: "republish",
  },
  {
    id: "376:1158",
    path: "/employer/jobs/[id]",
    dialog: "delete",
  },
  {
    id: "376:1167",
    path: "/employer/modules/[id]",
    dialog: "publish",
  },
  {
    id: "376:1176",
    path: "/employer/modules/[id]",
    dialog: "delete",
  },
  {
    id: "376:1185",
    path: "/employer/feedback/[id]",
    dialog: "release",
  },
  {
    id: "376:1194",
    path: "/employer/pipelines/[id]/candidates",
    dialog: "decision",
  },
  {
    id: "376:1203",
    path: "/employer/pipelines/[id]/candidates",
    dialog: "override",
  },
  {
    id: "376:1212",
    path: "/employer/offers/new",
    dialog: "capacity",
  },
  {
    id: "376:1221",
    path: "/employer/offers/[id]",
    dialog: "withdraw",
  },
  {
    id: "376:1230",
    path: "/employer/offers/[id]",
    dialog: "extend",
  },
  {
    id: "377:1082",
    path: "/employer/pipelines/[id]/candidates",
    dialog: "pass",
  },
  {
    id: "377:1091",
    path: "/employer/pipelines/[id]/candidates",
    dialog: "reject",
  },
  {
    id: "376:1995",
    path: "/employer/pipelines/[id]",
    dialog: "publish",
  },
  {
    id: "376:2004",
    path: "/employer/pipelines/[id]",
    dialog: "migrate",
  },
  {
    id: "376:2013",
    path: "/employer/pipelines/[id]",
    dialog: "discard",
  },
  {
    id: "377:1856",
    path: "/interviewer/interviews/[id]/scorecard",
    dialog: "submit",
  },
  {
    id: "377:1865",
    path: "/interviewer/grading",
    dialog: "submit",
  },
  {
    id: "376:2109",
    path: "/admin/dashboard",
    dialog: "approve",
  },
  {
    id: "376:2121",
    path: "/admin/dashboard",
    dialog: "reject",
  },
  {
    id: "376:2133",
    path: "/admin/accounts",
    dialog: "ban",
  },
  {
    id: "376:2145",
    path: "/admin/accounts",
    dialog: "unban",
  },
  {
    id: "376:2157",
    path: "/admin/accounts",
    dialog: "deactivate",
  },
  {
    id: "376:2169",
    path: "/admin/accounts",
    dialog: "reactivate",
  },
  {
    id: "376:2181",
    path: "/admin/accounts",
    dialog: "delete",
  },
] as const;
