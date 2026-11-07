"use client";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { performAccountAction } from "@/features/auth/api/auth.api";
import {
  ConfirmDialog,
  StatusBadge,
  UnavailableNotice,
} from "@/shared/ui/product-ui";
import { workspaceStyles as s } from "@/shared/ui/workspace";
import styles from "./candidate.module.css";
interface CandidateProfile {
  name?: string;
  email?: string;
  address?: string;
  biography?: string;
  phoneNumber?: string;
  githubProfile?: string;
  linkedinProfile?: string;
  portfolioLink?: string;
  available?: boolean;
}
/** Profile presentation remains editable locally while the profile API is planned. */
export function CandidateProfileScreen({
  edit = false,
  profile,
}: {
  edit?: boolean;
  profile?: CandidateProfile;
}) {
  const router = useRouter();
  const [dialog, setDialog] = useState<string | null>(null);
  const [biography, setBiography] = useState(profile?.biography || "");
  const action = dialog === "delete" ? "delete" : "deactivate";
  return (
    <section data-figma-node={edit ? "38:22" : "99:513"}>
      <header className={styles.actions}>
        <Image
          src={
            edit
              ? "/figma/edit-profile-heading.png"
              : "/figma/profile-heading.png"
          }
          alt=""
          width={edit ? 36 : 50}
          height={edit ? 36 : 50}
        />
        <h1>{edit ? "Edit profile" : "Your profile"}</h1>
      </header>
      <UnavailableNotice>
        Your profile information can’t be loaded or saved yet. Account
        deactivation and deletion are available after signing in.
      </UnavailableNotice>
      {edit ? (
        <form
          onSubmit={(event) => event.preventDefault()}
          className={s.card}
          aria-label="Edit candidate profile"
        >
          <section className={styles.formSection}>
            <h2>Personal details</h2>
            <div className={styles.formGrid}>
              <label className={s.field}>
                Name
                <input
                  name="name"
                  autoComplete="name"
                  required
                  maxLength={100}
                  defaultValue={profile?.name}
                />
              </label>
              <label className={s.field}>
                Phone number
                <input
                  name="phoneNumber"
                  autoComplete="tel"
                  type="tel"
                  maxLength={50}
                  defaultValue={profile?.phoneNumber}
                />
              </label>
              <label className={s.field}>
                Address
                <input
                  name="address"
                  autoComplete="street-address"
                  maxLength={500}
                  defaultValue={profile?.address}
                />
              </label>
              <label className={s.field}>
                Nationality
                <input name="nationality" maxLength={100} />
              </label>
              <label className={s.field}>
                Birth date
                <input name="birthDate" type="date" autoComplete="bday" />
              </label>
              <label className={s.field}>
                Gender
                <input name="gender" autoComplete="sex" maxLength={50} />
              </label>
            </div>
          </section>
          <section className={styles.formSection}>
            <h2>Professional profile</h2>
            <div className={styles.formGrid}>
              <label className={s.field}>
                Years of experience
                <input
                  name="yearsOfExperience"
                  type="number"
                  min={0}
                  step={1}
                />
              </label>
              <label className={s.field}>
                Available for opportunities
                <select
                  name="available"
                  defaultValue={profile?.available ? "true" : "false"}
                >
                  <option value="true">Yes</option>
                  <option value="false">No</option>
                </select>
              </label>
              <label className={`${s.field} ${styles.fullWidth}`}>
                Biography
                <textarea
                  name="biography"
                  maxLength={5000}
                  value={biography}
                  onChange={(event) => setBiography(event.target.value)}
                  aria-describedby="biography-hint"
                />
              </label>
              <p id="biography-hint">
                Up to 5,000 characters. {biography.length.toLocaleString()} /
                5,000
              </p>
            </div>
          </section>
          <section className={styles.formSection}>
            <h2>Professional links</h2>
            <div className={styles.formGrid}>
              <label className={s.field}>
                GitHub profile
                <input
                  type="url"
                  name="githubProfile"
                  defaultValue={profile?.githubProfile}
                />
              </label>
              <label className={s.field}>
                LinkedIn profile
                <input
                  type="url"
                  name="linkedinProfile"
                  defaultValue={profile?.linkedinProfile}
                />
              </label>
              <label className={s.field}>
                Portfolio
                <input
                  type="url"
                  name="portfolioLink"
                  defaultValue={profile?.portfolioLink}
                />
              </label>
            </div>
          </section>
          <div className={s.actions}>
            <button
              className={s.primary}
              type="submit"
              disabled
              aria-describedby="profile-save-note"
            >
              Save changes
            </button>
            <Link href="/candidate/profile" className={s.button}>
              Cancel
            </Link>
          </div>
          <p id="profile-save-note" className={s.muted}>
            Profile saving isn’t available yet.
          </p>
          <div className={s.actions}>
            <button type="button" onClick={() => setDialog("deactivate")}>
              Deactivate account
            </button>
            <button
              type="button"
              className={s.danger}
              onClick={() => setDialog("delete")}
            >
              Delete account
            </button>
          </div>
        </form>
      ) : (
        <>
          <div className={styles.detailGrid}>
            <section className={styles.panel}>
              <div className={s.actions}>
                <h2>{profile?.name || "Your profile details"}</h2>
                {profile?.available !== undefined && (
                  <StatusBadge>
                    {profile.available
                      ? "Available for opportunities"
                      : "Not currently available"}
                  </StatusBadge>
                )}
                <Link
                  href="/candidate/profile/edit"
                  className={`${s.button} ${s.primary}`}
                >
                  Edit profile
                </Link>
              </div>
              <p>{profile?.address || "Location unavailable"}</p>
              <p>
                {profile?.biography ||
                  "Your biography will appear here when profile information is available."}
              </p>
            </section>
            <aside className={`${styles.panel} ${styles.contactPanel}`}>
              <h2>Contact and links</h2>
              {profile?.email && <p>{profile.email}</p>}
              {profile?.phoneNumber && <p>{profile.phoneNumber}</p>}
              {!profile && <p>Your contact details aren’t available yet.</p>}
              {[
                profile?.githubProfile,
                profile?.linkedinProfile,
                profile?.portfolioLink,
              ]
                .filter(
                  (url): url is string => !!url && /^https?:\/\//.test(url),
                )
                .map((url) => (
                  <p key={url}>
                    <a href={url} target="_blank" rel="noopener noreferrer">
                      {url}
                    </a>
                  </p>
                ))}
            </aside>
          </div>
          <section className={`${styles.panel} ${styles.earned}`}>
            <h2>Earned amarat</h2>
            <p>Your verified awards and credentials</p>
            <p>No awards can be loaded yet.</p>
          </section>
        </>
      )}
      {dialog && (
        <ConfirmDialog
          title={
            action === "delete"
              ? "Delete your account?"
              : "Deactivate your account?"
          }
          description={
            action === "delete"
              ? "Deleting your account permanently removes your candidate profile and ends your current session. This action cannot be undone."
              : "Deactivating your account ends your session and makes your candidate profile inactive."
          }
          confirmLabel={
            action === "delete" ? "Delete account" : "Deactivate account"
          }
          intent="danger"
          onClose={() => setDialog(null)}
          onConfirm={async () => {
            await performAccountAction(action);
            router.replace("/login");
            router.refresh();
          }}
        />
      )}
    </section>
  );
}
