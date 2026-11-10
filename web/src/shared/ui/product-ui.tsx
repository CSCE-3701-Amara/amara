"use client";
import Image from "next/image";
import { useEffect, useId, useRef, useState, type ReactNode } from "react";
import styles from "./product-ui.module.css";
export { styles as productStyles };

/** A compact page heading with optional supporting brand artwork. */
export function ProductHeading({
  title,
  description,
  artwork,
  gradient = false,
  compact = false,
  className = "",
  children,
}: {
  title: string;
  description?: string;
  artwork?: { src: string; width: number; height: number };
  gradient?: boolean;
  compact?: boolean;
  className?: string;
  children?: ReactNode;
}) {
  return (
    <header
      className={`${styles.heading} ${gradient ? styles.gradientHeading : ""} ${artwork ? styles.illustratedHeading : ""} ${compact ? styles.compactHeading : ""} ${className}`}
    >
      {artwork && (
        <Image
          src={artwork.src}
          alt=""
          width={Math.round(artwork.width)}
          height={Math.round(artwork.height)}
          className={styles.headingArtwork}
          style={{ width: "auto", height: "auto" }}
        />
      )}
      <div>
        <h1>{title}</h1>
        {description && <p>{description}</p>}
        {children}
      </div>
    </header>
  );
}
/** Quiet product signature shared by hiring and administration screens. */
export function WorkspaceFooter() {
  return <footer className={styles.footer}>amara · Built in Egypt</footer>;
}
/** Clear integration boundary; never renders a successful persistence message. */
export function UnavailableNotice({
  children = "This feature isn’t available yet. You can review the details, but changes can’t be saved.",
}: {
  children?: ReactNode;
}) {
  return (
    <p className={styles.unavailable} role="note">
      {children}
    </p>
  );
}
export function StatusBadge({
  children,
  tone = "neutral",
}: {
  children: ReactNode;
  tone?: "neutral" | "success" | "warning" | "danger";
}) {
  return <span className={`${styles.badge} ${styles[tone]}`}>{children}</span>;
}
export function Stats({
  items,
}: {
  items: { label: string; value: string | number }[];
}) {
  return (
    <dl className={styles.stats}>
      {items.map((item) => (
        <div key={item.label}>
          <dt>{item.label}</dt>
          <dd
            className={
              typeof item.value === "string" &&
              item.value !== "—" &&
              !Number.isFinite(Number(item.value))
                ? styles.statNote
                : undefined
            }
          >
            {item.value}
          </dd>
        </div>
      ))}
    </dl>
  );
}
export function EmptyState({
  title,
  description,
  children,
  artwork,
  level = 2,
}: {
  title: string;
  description: string;
  children?: ReactNode;
  artwork?: { src: string; width: number; height: number };
  level?: 1 | 2;
}) {
  const Heading = level === 1 ? "h1" : "h2";
  return (
    <section className={styles.empty}>
      {artwork && (
        <Image
          src={artwork.src}
          width={artwork.width}
          height={artwork.height}
          alt=""
          style={{ width: "auto", height: "auto" }}
        />
      )}
      <Heading>{title}</Heading>
      <p>{description}</p>
      {children}
    </section>
  );
}

/** Native modal focus containment, Escape dismissal, and return to the trigger. */
export function ConfirmDialog({
  title,
  description,
  confirmLabel,
  intent = "warning",
  reason = false,
  unavailable,
  onConfirm,
  onClose,
}: {
  title: string;
  description: string;
  confirmLabel: string;
  intent?: "success" | "warning" | "danger";
  reason?: boolean;
  unavailable?: string;
  onConfirm?: (reason: string) => Promise<void>;
  onClose: () => void;
}) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const cancelRef = useRef<HTMLButtonElement>(null);
  const busyRef = useRef(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const titleId = useId();
  const descriptionId = useId();
  const reasonId = useId();
  useEffect(() => {
    const previous =
      document.activeElement instanceof HTMLElement
        ? document.activeElement
        : null;
    const dialog = dialogRef.current;
    dialog?.showModal();
    cancelRef.current?.focus();
    return () => {
      dialog?.close();
      if (previous?.isConnected && previous !== document.body) previous.focus();
      else document.getElementById("workspace-content")?.focus();
    };
  }, []);
  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busyRef.current || unavailable || !onConfirm) return;
    const value = String(
      new FormData(event.currentTarget).get("reason") || "",
    ).trim();
    if (reason && !value) {
      setError("Enter a reason to continue.");
      return;
    }
    busyRef.current = true;
    setBusy(true);
    setError("");
    try {
      await onConfirm(value);
      onClose();
    } catch (cause) {
      setError(
        cause instanceof Error
          ? cause.message
          : "The action couldn’t be completed. Try again.",
      );
    } finally {
      busyRef.current = false;
      setBusy(false);
    }
  }
  return (
    <dialog
      ref={dialogRef}
      className={styles.dialog}
      aria-labelledby={titleId}
      aria-describedby={descriptionId}
      onCancel={(event) => {
        event.preventDefault();
        if (!busyRef.current) onClose();
      }}
    >
      <form onSubmit={submit} aria-busy={busy}>
        <header className={styles.dialogHeader}>
          <span
            className={`${styles.intentIcon} ${styles[intent]}`}
            aria-hidden="true"
          >
            {intent === "success" ? "✓" : intent === "danger" ? "!" : "?"}
          </span>
          <h2 id={titleId}>{title}</h2>
        </header>
        <p id={descriptionId} className={styles.dialogDescription}>
          {description}
        </p>
        {reason && (
          <label htmlFor={reasonId} className={styles.reason}>
            Decision reason
            <textarea id={reasonId} name="reason" required maxLength={5000} />
          </label>
        )}
        {unavailable && <UnavailableNotice>{unavailable}</UnavailableNotice>}
        {error && (
          <p role="alert" className={styles.dialogError}>
            {error}
          </p>
        )}
        <div className={styles.dialogActions}>
          <button
            ref={cancelRef}
            type="button"
            onClick={onClose}
            disabled={busy}
          >
            Cancel
          </button>
          <button
            type="submit"
            className={styles[intent]}
            disabled={busy || !!unavailable || !onConfirm}
          >
            {busy ? "Please wait…" : confirmLabel}
          </button>
        </div>
      </form>
    </dialog>
  );
}
