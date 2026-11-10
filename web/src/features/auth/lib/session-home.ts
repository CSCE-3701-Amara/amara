/** Fixed local destinations; roles never become arbitrary redirect URLs. */
export function sessionHome(role: unknown): string | undefined {
  switch (role) {
    case "candidate": return "/candidate/dashboard";
    case "company":
    case "company_representative": return "/employer/jobs";
    case "interviewer": return "/interviewer/interviews";
    case "admin": return "/admin/dashboard";
    default: return undefined;
  }
}
