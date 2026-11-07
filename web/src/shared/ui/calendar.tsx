"use client";
import { useId, useState } from "react";
import styles from "./calendar.module.css";
/** Date-only keys avoid UTC/local-day drift in the interview picker. */
export function localDateKey(date: Date) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
}
export function Calendar({
  dates,
  selected,
  onSelect,
}: {
  dates: string[];
  selected: string;
  onSelect: (date: string) => void;
}) {
  const id = useId();
  const [offset, setOffset] = useState(0);
  const base = new Date((dates[0] || localDateKey(new Date())) + "T12:00:00");
  const month = new Date(base.getFullYear(), base.getMonth() + offset, 1);
  const start = month.getDay();
  const days = new Date(month.getFullYear(), month.getMonth() + 1, 0).getDate();
  const available = new Set(dates);
  return (
    <section className={styles.calendar} aria-labelledby={id}>
      <div className={styles.heading}>
        <button
          type="button"
          aria-label="Previous month"
          onClick={() => setOffset((x) => x - 1)}
        >
          ‹
        </button>
        <h3 id={id} aria-live="polite">
          {month.toLocaleDateString(undefined, {
            month: "long",
            year: "numeric",
          })}
        </h3>
        <button
          type="button"
          aria-label="Next month"
          onClick={() => setOffset((x) => x + 1)}
        >
          ›
        </button>
      </div>
      <div className={styles.grid}>
        {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map((day) => (
          <span className={styles.weekday} key={day} aria-hidden="true">
            {day}
          </span>
        ))}
        {Array.from({ length: start }, (_, i) => (
          <span key={"empty-" + i} />
        ))}
        {Array.from({ length: days }, (_, i) => {
          const date = new Date(month.getFullYear(), month.getMonth(), i + 1);
          const key = localDateKey(date);
          return (
            <button
              type="button"
              key={key}
              disabled={!available.has(key)}
              aria-pressed={selected === key}
              aria-label={date.toLocaleDateString(undefined, {
                weekday: "long",
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
              onClick={() => onSelect(key)}
            >
              {i + 1}
            </button>
          );
        })}
      </div>
      <p>Choose a highlighted date to see available times.</p>
    </section>
  );
}
