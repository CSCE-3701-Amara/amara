"use client";

import Image from "next/image";
import { useId, useState, type InputHTMLAttributes } from "react";
import styles from "./auth.module.css";

/** Native validation, explicit descriptions, and a keyboard-accessible password toggle. */
export function FormField({
  label,
  icon,
  hint,
  error,
  ...props
}: InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  icon?: "email" | "lock";
  hint?: string;
  error?: string;
}) {
  const uniqueId = useId();
  const inputId = `${props.id}-${uniqueId}`;
  const [visible, setVisible] = useState(false);
  const password = props.type === "password";
  const description =
    [
      props["aria-describedby"],
      hint && `${inputId}-hint`,
      error && `${inputId}-error`,
    ]
      .filter(Boolean)
      .join(" ") || undefined;
  return (
    <div className={styles.field}>
      <label htmlFor={inputId}>{label}</label>
      <div className={styles.inputWrap}>
        <input
          {...props}
          id={inputId}
          type={password && visible ? "text" : props.type}
          aria-invalid={error ? true : props["aria-invalid"]}
          aria-describedby={description}
          className={
            password
              ? styles.passwordInput
              : icon
                ? styles.iconInput
                : undefined
          }
        />
        {password ? (
          <button
            className={styles.passwordToggle}
            type="button"
            aria-label={visible ? "Hide password" : "Show password"}
            aria-pressed={visible}
            onClick={() => setVisible((value) => !value)}
          >
            {visible ? "Hide" : "Show"}
          </button>
        ) : (
          icon && (
            <Image src={`/${icon}.svg`} width={21.6} height={21.6} alt="" />
          )
        )}
      </div>
      {hint && (
        <p id={`${inputId}-hint`} className={styles.fieldHint}>
          {hint}
        </p>
      )}
      {error && (
        <p id={`${inputId}-error`} className={styles.fieldError} role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
