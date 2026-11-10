"use client";

import { useRef, useState } from "react";
import { ApiFeedback, BusyState } from "@/shared/ui/workspace";
import styles from "@/shared/ui/editor.module.css";
import { copyModule } from "../api/invw.api";
import type { ModuleSummary } from "../types/invw.types";

const moduleTypes = {
  coding: "Coding",
  programming_language: "Programming language",
  behavioral: "Behavioral",
  recorded_video: "Recorded video",
};

export function PipelineModulePicker({
  modules,
  library,
  selected,
  disabled,
  onToggle,
  onCopy,
  onCopyingChange,
}: {
  modules: ModuleSummary[];
  library: ModuleSummary[];
  selected: string[];
  disabled?: boolean;
  onToggle: (id: string, checked: boolean) => void;
  onCopy: (module: ModuleSummary) => void;
  onCopyingChange?: (copying: boolean) => void;
}) {
  const [source, setSource] = useState("company");
  const [type, setType] = useState("");
  const [copying, setCopying] = useState(false);
  const [error, setError] = useState<Error>();
  const pending = useRef(false);
  const visible = (source === "library" ? library : modules).filter(
    (module) => !type || module.type === type,
  );

  async function attach(id: string) {
    if (pending.current) return;
    pending.current = true;
    setCopying(true);
    onCopyingChange?.(true);
    setError(undefined);
    try {
      onCopy(await copyModule(id));
      setSource("company");
    } catch (cause) {
      setError(cause as Error);
    } finally {
      pending.current = false;
      setCopying(false);
      onCopyingChange?.(false);
    }
  }

  return (
    <fieldset>
      <legend>Assessment modules</legend>
      <label className={styles.field}>
        Source
        <select value={source} onChange={(event) => setSource(event.target.value)}>
          <option value="company">Company modules</option>
          <option value="library">Published module library</option>
        </select>
      </label>
      <label className={styles.field}>
        Module type
        <select value={type} onChange={(event) => setType(event.target.value)}>
          <option value="">All types</option>
          {Object.entries(moduleTypes).map(([value, label]) => (
            <option key={value} value={value}>{label}</option>
          ))}
        </select>
      </label>
      <BusyState loading={copying} />
      <ApiFeedback error={error} />
      {!visible.length && <p>No matching modules.</p>}
      {visible.map((module) => source === "library" && !modules.some((item) => item.id === module.id) ? (
        <button key={module.id} type="button" disabled={disabled || copying} onClick={() => void attach(module.id)}>
          Copy and attach {module.name} · {moduleTypes[module.type]}
        </button>
      ) : (
        <label className={styles.check} key={module.id}>
          <input type="checkbox" disabled={disabled || copying} checked={selected.includes(module.id)}
            onChange={(event) => onToggle(module.id, event.target.checked)} />
          {module.name} · {moduleTypes[module.type]}
        </label>
      ))}
    </fieldset>
  );
}
