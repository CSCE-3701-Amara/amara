"use client";

import { useRef, useState, type FormEvent } from "react";
import { apiRequest } from "./api-client";
import { useUnsavedChanges } from "./use-unsaved-changes";

/** Patch only edited fields: missing profile reads must never erase existing data. */
export function useProfileUpdate<T>(path: string, nullable = false) {
  const changed = useRef(new Set<string>());
  const pending = useRef(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error>();
  const [saved, setSaved] = useState<T>();
  const [message, setMessage] = useState("");
  const [dirty, setDirty] = useState(false);
  useUnsavedChanges(dirty);

  function onChange(event: FormEvent<HTMLFormElement>) {
    const field = event.target;
    if (
      field instanceof HTMLInputElement ||
      field instanceof HTMLTextAreaElement ||
      field instanceof HTMLSelectElement
    ) {
      changed.current.add(field.name);
      setDirty(true);
      setMessage("");
    }
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending.current) return;
    const form = event.currentTarget;
    const data = new FormData(form);
    const body: Record<string, unknown> = {};
    setError(undefined);
    setMessage("");
    for (const name of changed.current) {
      const value = String(data.get(name) ?? "").trim();
      if (["name", "legalName"].includes(name) && !value) {
        setError(new Error("Enter a name; spaces alone aren’t accepted."));
        form.querySelector<HTMLInputElement>(`[name="${name}"]`)?.focus();
        return;
      }
      if (name === "available") body[name] = value === "true";
      else if (name === "yearsOfExperience") {
        if (value !== "") body[name] = Number(value);
      } else if (nullable) body[name] = value || null;
      // The candidate handler cannot clear dates or URLs; empty optional fields are omitted.
      else if (
        value ||
        ![
          "birthDate",
          "githubProfile",
          "linkedinProfile",
          "portfolioLink",
        ].includes(name)
      )
        body[name] = value;
    }
    if (!Object.keys(body).length) {
      setMessage("No changes to save.");
      return;
    }
    pending.current = true;
    setBusy(true);
    try {
      const profile = await apiRequest<T>(path, { method: "PATCH", body });
      setSaved(profile);
      changed.current.clear();
      setDirty(false);
      setMessage("Profile saved.");
    } catch (cause) {
      setError(
        cause instanceof Error
          ? cause
          : new Error("Your profile couldn’t be saved. Try again."),
      );
    } finally {
      pending.current = false;
      setBusy(false);
    }
  }
  return { busy, error, saved, message, dirty, onChange, onSubmit };
}
