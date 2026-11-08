import { ApiFeedback, BusyState } from "@/shared/ui/workspace";

export function PlannedRead({
  loading,
  error,
  retry,
}: {
  loading: boolean;
  error?: Error & { status?: number };
  retry: () => void;
}) {
  return (
    <>
      <BusyState loading={loading} />
      {error?.status !== 404 && <ApiFeedback error={error} retry={retry} />}
    </>
  );
}
