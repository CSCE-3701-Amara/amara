"use client";

import Image from "next/image";
import Link from "next/link";
import type { ReactNode, RefObject } from "react";
import { Brand } from "@/shared/ui/brand";
import styles from "./auth.module.css";

/** Responsive shared layout for the login, registration, and auth state frames. */
export function AuthShell({
  company,
  registration,
  children,
  contentRef,
}: {
  company: boolean;
  registration: boolean;
  children: ReactNode;
  contentRef?: RefObject<HTMLElement | null>;
}) {
  return (
    <main className={styles.shell}>
      <aside
        className={styles.panel}
        aria-label={company ? "For hiring teams" : "For candidates"}
      >
        <div className={styles.gradient} aria-hidden="true" />
        <Brand />
        <div className={styles.pitch} key={company ? "company" : "candidate"}>
          <p className={styles.eyebrow}>
            {company ? "For hiring teams" : "For candidates"}
          </p>
          <h2>
            {company ? (
              <>
                Build teams
                <br />{" "}
                around skills.
              </>
            ) : (
              <>
                Show what
                <br />{" "}
                you can do.
              </>
            )}
          </h2>
          <p className={styles.description}>
            {company
              ? "Design assessments, manage your hiring pipeline, and make decisions with clearer evidence."
              : "Find roles, show your skills in assessments, and follow your progress in one place."}
          </p>
        </div>
        <div className={styles.artwork}>
          <div className={styles.artCopy}>
            <p className={styles.artLabel}>
              {company ? "Hiring workspace" : "Candidate workspace"}
            </p>
            <h3>
              {company ? (
                registration ? (
                  <>
                    Build your
                    <br />{" "}
                    pipeline
                  </>
                ) : (
                  <>
                    See every
                    <br />{" "}
                    hiring stage
                  </>
                )
              ) : registration ? (
                <>
                  Build your
                  <br />{" "}
                  skills profile
                </>
              ) : (
                <>
                  Your skills
                  <br />{" "}
                  in focus
                </>
              )}
            </h3>
            <p>
              {company
                ? "Run parallel assessments and bring the results together."
                : "Show your experience, projects, and skills in one place."}
            </p>
          </div>
          <Image
            className={styles.illustration}
            src={company ? "/hiring-pipeline.svg" : "/candidate-portfolio.svg"}
            width={333.6}
            height={231.6}
            alt=""
            priority
          />
        </div>
        <p className={styles.panelFooter}>
          amara&nbsp; / &nbsp;SKILLS-FIRST RECRUITMENT
        </p>
      </aside>
      <section
        ref={contentRef}
        className={styles.content}
        aria-label={registration ? "Create an account" : "Sign in"}
      >
        <Link href="/" className={styles.back}>
          <Image src="/back-arrow.svg" width={19.2} height={19.2} alt="" />
          Back to amara
        </Link>
        {children}
        <p className={styles.workspace}>
          amara&nbsp; / &nbsp;{company ? "HIRING" : "CANDIDATE"} WORKSPACE
        </p>
      </section>
    </main>
  );
}
