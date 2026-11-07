import Image from "next/image";
import Link from "next/link";
import styles from "./brand.module.css";

/** The shared Figma logo lockup, with its original mark dimensions. */
export function Brand({
  variant = "badge",
}: {
  variant?: "badge" | "header" | "footer" | "workspace";
}) {
  const mark =
    variant === "workspace"
      ? { src: "/figma/workspace-logo.svg", width: 30.68, height: 26.52 }
      : variant === "header"
        ? { src: "/frame.svg", width: 53.1277, height: 45.9239 }
        : variant === "footer"
          ? { src: "/amara-logo-mark.svg", width: 49.9765, height: 43.2 }
          : { src: "/amara-logo.svg", width: 47.2, height: 40.8 };
  return (
    <Link
      href="/"
      aria-label="amara home"
      className={`${styles.brand} ${styles[variant]}`}
    >
      <Image {...mark} alt="" priority />
      <span>amara</span>
    </Link>
  );
}
