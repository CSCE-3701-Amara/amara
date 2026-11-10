"use client";
import { useEffect, useState } from "react";
import styles from "./candidate.module.css";

export function AttemptClock({
  deadline,
  active,
}: {
  deadline?: string;
  active: boolean;
}) {
  const [now, setNow] = useState(0);
  useEffect(() => {
    if (!active) return;
    const firstTick = setTimeout(() => setNow(Date.now()), 0);
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => {
      clearTimeout(firstTick);
      clearInterval(timer);
    };
  }, [active]);
  const milliseconds =
    deadline && now ? Math.max(0, Date.parse(deadline) - now) : undefined;
  return (
    <section>
      <p>Remaining active time</p>
      <strong className={styles.clock}>
        {milliseconds !== undefined && Number.isFinite(milliseconds)
          ? `${Math.floor(milliseconds / 60000)}:${String(Math.floor(milliseconds / 1000) % 60).padStart(2, "0")}`
          : "Time unavailable"}
      </strong>
      {deadline && (
        <p>Completion deadline · {new Date(deadline).toLocaleString()}</p>
      )}
    </section>
  );
}
