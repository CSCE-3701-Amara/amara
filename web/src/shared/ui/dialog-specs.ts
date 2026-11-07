/** Dialog copy and intent tokens transcribed from the full-file Figma inventory. */
export const dialogSpecs = {
  "212:329": {
    id: "212:329",
    title: "Review before withdrawing",
    description:
      "Withdrawing ends your active application. Review your progress before continuing. This action cannot be undone.",
    confirmLabel: "Withdraw application",
    intent: "danger",
    requiresReason: false,
  },
  "376:555": {
    id: "376:555",
    title: "Confirm interview time?",
    description:
      "You are booking Thursday, 12 November at 1:00 PM (UTC+02:00). This will reserve the selected 45-minute slot.",
    confirmLabel: "Confirm time",
    intent: "success",
    requiresReason: false,
  },
  "376:564": {
    id: "376:564",
    title: "Reschedule interview?",
    description:
      "Your current interview time will be replaced with the newly selected slot. The hiring team will be notified of the change.",
    confirmLabel: "Reschedule",
    intent: "warning",
    requiresReason: false,
  },
  "376:573": {
    id: "376:573",
    title: "Accept this offer?",
    description:
      "By accepting, you confirm your intention to join Acme Corp under the terms shown in the offer details.",
    confirmLabel: "Accept offer",
    intent: "success",
    requiresReason: false,
  },
  "376:582": {
    id: "376:582",
    title: "Decline this offer?",
    description:
      "You are about to decline the offer from Acme Corp. This response will be recorded for the hiring team.",
    confirmLabel: "Decline offer",
    intent: "danger",
    requiresReason: false,
  },
  "376:591": {
    id: "376:591",
    title: "Deactivate your account?",
    description:
      "Your candidate account will be deactivated. You may need to contact support to restore access.",
    confirmLabel: "Deactivate",
    intent: "warning",
    requiresReason: false,
  },
  "376:600": {
    id: "376:600",
    title: "Delete your account?",
    description:
      "This permanently deletes your candidate account and cannot be undone. Make sure you have saved any information you need.",
    confirmLabel: "Delete account",
    intent: "danger",
    requiresReason: false,
  },
  "377:553": {
    id: "377:553",
    title: "Submit your assessment?",
    description:
      "You have completed the required questions. Once submitted, your answers can no longer be edited.",
    confirmLabel: "Submit assessment",
    intent: "warning",
    requiresReason: false,
  },
  "378:571": {
    id: "378:571",
    title: "Withdraw this application?",
    description:
      "Withdrawing ends your participation in this hiring process and closes the active application stages.",
    confirmLabel: "Withdraw",
    intent: "danger",
    requiresReason: false,
  },
  "376:1131": {
    id: "376:1131",
    title: "Publish this job?",
    description:
      "The job posting will become visible to candidates according to its visibility and application-window settings.",
    confirmLabel: "Publish job",
    intent: "warning",
    requiresReason: false,
  },
  "376:1140": {
    id: "376:1140",
    title: "Archive this job?",
    description:
      "{`Archiving changes the job's status and removes it from active job management. Review the posting before continuing.`}",
    confirmLabel: "Archive job",
    intent: "warning",
    requiresReason: false,
  },
  "376:1149": {
    id: "376:1149",
    title: "Republish this job?",
    description:
      "The job will return to a published state using its current saved details.",
    confirmLabel: "Republish",
    intent: "warning",
    requiresReason: false,
  },
  "376:1158": {
    id: "376:1158",
    title: "Delete this job?",
    description:
      "This will soft-delete the job posting. It will no longer appear as an active posting.",
    confirmLabel: "Delete job",
    intent: "danger",
    requiresReason: false,
  },
  "376:1167": {
    id: "376:1167",
    title: "Publish this module?",
    description:
      "Published modules become read-only. Review the module details and threshold before publishing.",
    confirmLabel: "Publish module",
    intent: "warning",
    requiresReason: false,
  },
  "376:1176": {
    id: "376:1176",
    title: "Delete this module?",
    description:
      "{`This removes the module from the company's module library. This action should only be used when the module is no longer needed.`}",
    confirmLabel: "Delete module",
    intent: "danger",
    requiresReason: false,
  },
  "376:1185": {
    id: "376:1185",
    title: "Release feedback to candidate?",
    description:
      "The approved feedback will become visible to the candidate in their submitted application details.",
    confirmLabel: "Release feedback",
    intent: "warning",
    requiresReason: false,
  },
  "376:1194": {
    id: "376:1194",
    title: "Confirm stage decision?",
    description:
      "This records the selected pass or reject decision for the candidate at the current pipeline stage.",
    confirmLabel: "Submit decision",
    intent: "warning",
    requiresReason: false,
  },
  "376:1203": {
    id: "376:1203",
    title: "Override the current decision?",
    description:
      "This replaces the existing stage decision. The override should be intentional and documented.",
    confirmLabel: "Override decision",
    intent: "warning",
    requiresReason: false,
  },
  "376:1212": {
    id: "376:1212",
    title: "Create offer beyond open positions?",
    description:
      "{`The number of offers will exceed the job's configured open positions. Confirm that you have reviewed this exception.`}",
    confirmLabel:
      "I understand this offer will exceed the configured open positions.",
    intent: "warning",
    requiresReason: false,
  },
  "376:1221": {
    id: "376:1221",
    title: "Withdraw this offer?",
    description:
      "The offer will be withdrawn and the candidate will no longer be able to accept it.",
    confirmLabel: "Withdraw offer",
    intent: "danger",
    requiresReason: false,
  },
  "376:1230": {
    id: "376:1230",
    title: "Extend offer deadline?",
    description:
      "{`The candidate's response deadline will be updated to the new date shown in the offer editor.`}",
    confirmLabel: "Extend deadline",
    intent: "warning",
    requiresReason: false,
  },
  "377:1082": {
    id: "377:1082",
    title: "Pass this candidate?",
    description:
      "The candidate will be marked as having passed the current pipeline stage. This decision will be recorded in the application history.",
    confirmLabel: "Pass candidate",
    intent: "warning",
    requiresReason: false,
  },
  "377:1091": {
    id: "377:1091",
    title: "Reject this candidate?",
    description:
      "The candidate will be marked as rejected at the current pipeline stage. Make sure the decision and feedback are ready before continuing.",
    confirmLabel: "Reject candidate",
    intent: "danger",
    requiresReason: false,
  },
  "376:1995": {
    id: "376:1995",
    title: "Publish pipeline version?",
    description:
      "The validated draft will become the published version for the selected hiring workflow. Check the stages and transitions before publishing.",
    confirmLabel: "Publish version",
    intent: "warning",
    requiresReason: false,
  },
  "376:2004": {
    id: "376:2004",
    title: "Confirm pipeline migration?",
    description:
      "Review the migration preview and stage mapping before applying the new pipeline version.",
    confirmLabel: "Confirm migration",
    intent: "warning",
    requiresReason: false,
  },
  "376:2013": {
    id: "376:2013",
    title: "Discard pipeline draft?",
    description: "Unsaved changes to this draft version will be discarded.",
    confirmLabel: "Discard draft",
    intent: "warning",
    requiresReason: false,
  },
  "377:1856": {
    id: "377:1856",
    title: "Submit interview scorecard?",
    description:
      "Your scorecard and private notes will be finalized. You will no longer be able to edit this submission.",
    confirmLabel: "Submit scorecard",
    intent: "warning",
    requiresReason: false,
  },
  "377:1865": {
    id: "377:1865",
    title: "Submit grading results?",
    description:
      "The points recorded for the assigned answers will be submitted to the hiring team.",
    confirmLabel: "Submit grading",
    intent: "warning",
    requiresReason: false,
  },
  "376:2109": {
    id: "376:2109",
    title: "Approve company registration?",
    description:
      "This company registration will be approved and the company can proceed to the active state.",
    confirmLabel: "Approve company",
    intent: "success",
    requiresReason: true,
  },
  "376:2121": {
    id: "376:2121",
    title: "Reject company registration?",
    description:
      "The company registration will be rejected. Provide a reason that can be recorded with the review.",
    confirmLabel: "Reject company",
    intent: "danger",
    requiresReason: true,
  },
  "376:2133": {
    id: "376:2133",
    title: "Ban this account?",
    description:
      "The account will be banned. This action requires an audit reason.",
    confirmLabel: "Ban account",
    intent: "danger",
    requiresReason: true,
  },
  "376:2145": {
    id: "376:2145",
    title: "Unban this account?",
    description:
      "The account ban will be lifted. Record the reason for this administrative action.",
    confirmLabel: "Unban account",
    intent: "warning",
    requiresReason: true,
  },
  "376:2157": {
    id: "376:2157",
    title: "Deactivate this account?",
    description:
      "The account will be deactivated. Record why this action is being taken.",
    confirmLabel: "Deactivate",
    intent: "warning",
    requiresReason: true,
  },
  "376:2169": {
    id: "376:2169",
    title: "Reactivate this account?",
    description:
      "The account will be reactivated. Record the reason for restoring access.",
    confirmLabel: "Reactivate",
    intent: "success",
    requiresReason: true,
  },
  "376:2181": {
    id: "376:2181",
    title: "Delete this account?",
    description:
      "This permanently deletes the account. The action requires an audit reason and should be used carefully.",
    confirmLabel: "Delete account",
    intent: "danger",
    requiresReason: true,
  },
} as const;
export type DialogSpecId = keyof typeof dialogSpecs;
