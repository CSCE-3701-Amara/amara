"use client";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useApiResource } from "@/shared/lib/use-api-resource";
import { useProfileUpdate } from "@/shared/lib/use-profile-update";
import { ApiFeedback } from "@/shared/ui/workspace";
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
  amarat?: string[];
}
export function CandidateProfileScreen({
  edit = false,
  profile,
}: {
  edit?: boolean;
  profile?: CandidateProfile;
}) {
  const router = useRouter();
  const update = useProfileUpdate<CandidateProfile>("/v1/candidates/me");
  const current = useApiResource<CandidateProfile>(
    edit ? null : "/v1/candidates/me",
  );
  profile = update.saved || current.data || profile;
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
      {edit && (
        <UnavailableNotice>
          Your current profile can’t be loaded yet. Save only the fields you
          want to update; other details will be kept.
        </UnavailableNotice>
      )}
      {edit ? (
        <form
          onSubmit={update.onSubmit}
          onChange={update.onChange}
          aria-busy={update.busy}
          className={s.card}
          aria-label="Edit candidate profile"
        >
          <fieldset disabled={update.busy} className={styles.profileFields}>
            <section className={styles.formSection}>
              <h2>Personal details</h2>
              <div className={styles.formGrid}>
                <label className={s.field}>
                  Name
                  <input
                    name="name"
                    autoComplete="name"
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
          </fieldset>
          <ApiFeedback error={update.error} />
          {update.message && <p role="status">{update.message}</p>}
          <div className={s.actions}>
            <button
              className={s.primary}
              type="submit"
              disabled={update.busy}
              aria-describedby="profile-save-note"
            >
              {update.busy ? "Saving…" : "Save changes"}
            </button>
            <Link
              href="/candidate/profile"
              onNavigate={(event) => {
                if (
                  update.dirty &&
                  !window.confirm("Discard your unsaved profile changes?")
                )
                  event.preventDefault();
              }}
              className={s.button}
            >
              Cancel
            </Link>
          </div>
          <p id="profile-save-note" className={s.muted}>
            Only edited fields are sent. Saved details cannot yet be reloaded
            from your profile.
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
            {current.loading ? (
              <p>Loading awards…</p>
            ) : current.error ? (
              <p>Your awards couldn’t be loaded.</p>
            ) : current.data?.amarat?.length ? (
              <ul>
                {current.data.amarat.map((title) => (
                  <li key={title}>{title}</li>
                ))}
              </ul>
            ) : (
              <p>No awards yet.</p>
            )}
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
