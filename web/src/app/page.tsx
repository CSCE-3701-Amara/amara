import { Suspense } from "react";
import { GuestEntry } from "@/features/auth/components/guest-entry";
import Image from "next/image";
import Link from "next/link";
import { Brand } from "@/shared/ui/brand";
import { ShaderFill } from "@/shared/ui/shader-fill";
import styles from "./welcome.module.css";

const steps = [
  {
    title: "Build a profile",
    description:
      "Bring your experience, projects, and skills together in one place.",
    action: "Create profile",
    image: "candidate-profile-illustration",
    href: "/register",
  },
  {
    title: "Explore roles",
    description:
      "Find openings that match your skills and what you want to do next.",
    action: "Find roles",
    image: "job-opportunities-illustration",
    href: "/login",
  },
  {
    title: "Show your skills",
    description:
      "Take job-related assessments and follow your application through each stage.",
    action: "Join amara",
    image: "branching-assessment-pipeline-illustration",
    href: "/register",
  },
];
const toolkit = [
  ["Job postings", "job-postings-icon"],
  ["Coding assessments", "coding-assessments-icon"],
  ["Video responses", "video-responses-icon"],
  ["Live interviews", "live-interviews-icon"],
];

function Welcome() {
  return (
    <main className={styles.page}>
      <section id="home" className={styles.hero}>
        <ShaderFill className={styles.shader} />
        <header className={styles.header}>
          <Brand variant="header" />
          <nav className={styles.nav} aria-label="Main navigation">
            <a href="#home">Home</a>
            <a href="#why-amara">Why amara</a>
            <a href="#how-it-works">How it works</a>
            <a href="#for-teams">For teams</a>
          </nav>
          <div className={styles.accountLinks}>
            <Link href="/login">Log in</Link>
            <Link href="/register">Sign up</Link>
          </div>
        </header>
        <div className={styles.heroCopy}>
          <h1>
            Beyond Resumes.
            <br />{" "}
            Towards Skills.
          </h1>
          <p>
            Find opportunities that value your skills.
            <br />{" "}
            Give hiring teams a clearer picture of what you can do.
          </p>
          <a className={styles.heroAction} href="#how-it-works">
            See how amara works
          </a>
        </div>
        <div className={styles.origin}>
          <Image
            src="/egypt-illustration.svg"
            width={69.6}
            height={69.6}
            alt=""
          />
          <div>
            <h2>Built in Egypt</h2>
            <p>
              An Egyptian startup helping people get hired for their skills.
            </p>
          </div>
        </div>
      </section>

      <section className={styles.about} aria-labelledby="show-skills">
        <div className={styles.sectionCopy}>
          <h2 id="show-skills">
            Show what
            <br />{" "}
            you can do
          </h2>
          <p>
            Build your profile, explore open roles, and show your skills through
            job-related assessments.
          </p>
          <Link className={styles.greenLink} href="/register">
            Join amara
          </Link>
        </div>
        <article className={styles.jobPreview} aria-label="Example job opening">
          <header>
            <span className={styles.iconTile}>
              <Image
                src="/briefcase.svg"
                width={39.6206}
                height={39.6206}
                alt=""
              />
            </span>
            <div>
              <p>Example job opening</p>
              <h3>Frontend Engineer</h3>
            </div>
          </header>
          <p>Product engineering · Full-time</p>
          <hr />
          <h4>Skills you’ll use</h4>
          <div className={styles.skills}>
            <span>React</span>
            <span>TypeScript</span>
            <span>Accessibility</span>
          </div>
        </article>
      </section>

      <section
        id="how-it-works"
        className={styles.how}
        aria-labelledby="how-heading"
      >
        <h2 id="how-heading">
          How amara
          <br />{" "}
          works
        </h2>
        <div className={styles.steps}>
          {steps.map((step, index) => (
            <article className={styles.step} key={step.title}>
              <div
                className={`${styles.stepImage} ${index === 2 ? styles.branchImage : ""}`}
              >
                <Image
                  src={`/${step.image}.svg`}
                  width={index === 2 ? 489.6 : 422.4}
                  height={index === 2 ? 319.2 : 252}
                  alt=""
                />
              </div>
              <h3>{step.title}</h3>
              <p>{step.description}</p>
              <Link href={step.href}>{step.action}</Link>
            </article>
          ))}
        </div>
        <Link className={styles.getStarted} href="/register">
          Get started
        </Link>
      </section>

      <section
        id="for-teams"
        className={styles.hiring}
        aria-labelledby="hiring-heading"
      >
        <div className={styles.hiringIntro}>
          <div className={styles.sectionCopy}>
            <span className={styles.teamEyebrow}>For hiring teams</span>
            <h2 id="hiring-heading">
              Build a better
              <br />{" "}
              hiring process
            </h2>
            <p>
              Create job postings, define the skills you want to see, and build
              assessment and interview stages around them.
            </p>
            <Link className={styles.greenLink} href="/register?account=company">
              Join as a team
            </Link>
          </div>
          <div
            className={styles.pipeline}
            aria-label="Example hiring pipeline: application branches into coding test and skills review, both precede interview"
          >
            <header>
              <Image
                src="/branching-pipeline-icon.svg"
                width={34.7806}
                height={34.7806}
                alt=""
              />
              <div>
                <h3>Hiring pipeline</h3>
                <p>Parallel assessments can feed one interview.</p>
              </div>
            </header>
            <div className={styles.pipelineDiagram}>
              <Image
                className={styles.arrows}
                src="/directed-prerequisite-arrows.svg"
                width={741.986}
                height={301.432}
                alt=""
              />
              <span className={styles.applicationStage}>Application</span>
              <span className={styles.codingStage}>Coding test</span>
              <span className={styles.reviewStage}>Skills review</span>
              <span className={styles.interviewStage}>Interview</span>
            </div>
            <p className={styles.pipelineNote}>
              Both assessment stages come before the interview.
            </p>
          </div>
        </div>
        <div className={styles.toolkit}>
          <h3>Your hiring toolkit</h3>
          <div>
            {toolkit.map(([label, icon]) => (
              <div key={label}>
                <Image src={`/${icon}.svg`} width={39.6} height={39.6} alt="" />
                <span>{label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="why-amara" className={styles.why} aria-label="Why amara">
        <div>
          <p>
            Candidates see where they stand. Hiring teams get evidence to guide
            their decisions. amara brings both sides together around the skills
            that matter.
          </p>
          <div className={styles.whyActions}>
            <Link className={styles.greenLink} href="/register">
              Create an account
            </Link>
            <Link href="/login">Log in to amara</Link>
          </div>
        </div>
        <article className={styles.progress}>
          <header>
            <span className={styles.iconTile}>
              <Image
                src="/application-checklist.svg"
                width={34.2178}
                height={34.2178}
                alt=""
              />
            </span>
            <h3>Application progress</h3>
          </header>
          {[
            ["Application", "Submitted"],
            ["Skills assessment", "In progress"],
            ["Interview", "Upcoming"],
          ].map(([label, status], index) => (
            <div className={styles.progressRow} key={label}>
              <i data-stage={index} />
              <strong>{label}</strong>
              <span>{status}</span>
            </div>
          ))}
        </article>
      </section>

      <footer className={styles.footer}>
        <Image
          className={styles.footerLines}
          src="/group.svg"
          width={1468.8}
          height={300.6}
          alt=""
        />
        <Brand variant="footer" />
        <div className={styles.footerGrid}>
          <div>
            <h2>
              Hiring built
              <br />{" "}
              around skills.
            </h2>
            <p>
              amara is an Egypt-based recruitment platform. Candidates can
              discover jobs and show their skills. Hiring teams can build
              assessments and run structured hiring pipelines.
            </p>
            <div className={styles.footerActions}>
              <Link href="/register">Join as a candidate</Link>
              <Link href="/register?account=company">
                Join as a hiring team
              </Link>
            </div>
          </div>
          <div className={styles.explore}>
            <h3>Explore amara</h3>
            <div>
              <section>
                <h4>For candidates</h4>
                <p>Find roles and track your applications.</p>
              </section>
              <Link href="/login">Candidate login</Link>
            </div>
            <div>
              <section>
                <h4>For hiring teams</h4>
                <p>Post openings and manage a branching pipeline.</p>
              </section>
              <Link href="/login?account=company">Hiring team login</Link>
            </div>
            <div>
              <section>
                <h4>How amara works</h4>
                <p>Learn how profiles, jobs, and assessments connect.</p>
              </section>
              <a href="#how-it-works">See the steps</a>
            </div>
          </div>
        </div>
        <div className={styles.footerBottom}>
          <span>Built in Egypt</span>
          <span>amara / SKILLS-FIRST RECRUITMENT</span>
        </div>
      </footer>
    </main>
  );
}

export default function Home() {
  return (
    <Suspense fallback={<main aria-busy="true"><p role="status">Checking your session…</p></main>}>
      <GuestEntry><Welcome /></GuestEntry>
    </Suspense>
  );
}
