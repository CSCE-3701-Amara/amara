"use client";
import Image from "next/image";
import { useEffect, useRef, useState } from "react";
import { workspaceStyles as s } from "@/shared/ui/workspace";
import styles from "./candidate.module.css";

export function VideoAnswer({ disabled }: { disabled: boolean }) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const recorderRef = useRef<MediaRecorder | null>(null);
  const recordingUrl = useRef<string | null>(null);
  const generation = useRef(0);
  const pending = useRef(false);
  const [requesting, setRequesting] = useState(false);
  const [recording, setRecording] = useState(false);
  const [preview, setPreview] = useState("");
  const [error, setError] = useState("");
  useEffect(
    () => () => {
      generation.current += 1;
      pending.current = false;
      const recorder = recorderRef.current;
      if (recorder) {
        recorder.onstop = null;
        recorder.ondataavailable = null;
        if (recorder.state !== "inactive") recorder.stop();
      }
      streamRef.current?.getTracks().forEach((track) => track.stop());
      if (recordingUrl.current) URL.revokeObjectURL(recordingUrl.current);
      recorderRef.current = null;
      streamRef.current = null;
      recordingUrl.current = null;
      // Activity can hide this route without destroying its retained state.
      setRequesting(false);
      setRecording(false);
      setPreview("");
    },
    [],
  );
  async function record() {
    if (pending.current || recording || disabled) return;
    const current = generation.current;
    pending.current = true;
    setRequesting(true);
    setError("");
    try {
      if (
        !navigator.mediaDevices?.getUserMedia ||
        typeof MediaRecorder === "undefined"
      )
        throw new Error(
          "Recording isn’t supported in this browser. Use a supported browser with camera access.",
        );
      const stream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true,
      });
      if (current !== generation.current) {
        stream.getTracks().forEach((track) => track.stop());
        return;
      }
      streamRef.current = stream;
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        videoRef.current.muted = true;
        await videoRef.current.play();
      }
      if (current !== generation.current) {
        stream.getTracks().forEach((track) => track.stop());
        return;
      }
      const recorder = new MediaRecorder(stream);
      const chunks: BlobPart[] = [];
      recorder.ondataavailable = (event) => {
        if (event.data.size) chunks.push(event.data);
      };
      recorder.onstop = () => {
        stream.getTracks().forEach((track) => track.stop());
        if (current !== generation.current) return;
        if (recordingUrl.current) URL.revokeObjectURL(recordingUrl.current);
        const url = URL.createObjectURL(
          new Blob(chunks, { type: recorder.mimeType || "video/webm" }),
        );
        recordingUrl.current = url;
        setPreview(url);
        setRecording(false);
        if (videoRef.current) {
          videoRef.current.srcObject = null;
          videoRef.current.muted = false;
        }
      };
      recorderRef.current = recorder;
      setPreview("");
      recorder.start();
      setRecording(true);
    } catch (cause) {
      if (current !== generation.current) return;
      streamRef.current?.getTracks().forEach((track) => track.stop());
      setError(
        cause instanceof Error
          ? cause.message
          : "Camera or microphone access was denied. Update your browser permissions and try again.",
      );
    } finally {
      if (current === generation.current) {
        pending.current = false;
        setRequesting(false);
      }
    }
  }
  return (
    <>
      <div className={styles.videoPreview}>
        {!recording && !preview && (
          <Image src="/figma/video.svg" width={40.32} height={40.32} alt="" />
        )}
        <video
          ref={videoRef}
          src={preview || undefined}
          controls={!!preview}
          playsInline
          hidden={!recording && !preview}
          aria-label="Your recorded video response"
        />
      </div>
      <div className={s.actions}>
        <button
          type="button"
          className={s.primary}
          disabled={disabled || recording || requesting}
          onClick={record}
        >
          {preview ? "Retake" : "Record"}
        </button>
        <button
          type="button"
          disabled={!recording}
          onClick={() => recorderRef.current?.stop()}
        >
          Stop
        </button>
        <button type="button" disabled aria-describedby="video-upload-note">
          Save recording
        </button>
      </div>
      {error && (
        <p role="alert" className={s.error}>
          {error}
        </p>
      )}
      <p>
        Changing the prompt or leaving this screen discards unsaved recordings.
      </p>
      <p id="video-upload-note">
        Recording previews stay on this screen. Uploading and saving recordings
        aren’t available yet.
      </p>
    </>
  );
}
