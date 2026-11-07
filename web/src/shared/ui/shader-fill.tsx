"use client";

import { useEffect, useRef } from "react";
import { useReducedMotion } from "./use-reduced-motion";
import {
  setup,
  render,
} from "@/shared/lib/figma-shaders/CodeComponentId_8ce92017e53431a2f04b3574f4ba7c98f6f55f1e_8358";

// Parameters are copied from the welcome frame, rather than approximating its shader with CSS.
const params = {
  intensity: 3.5299999713897705,
  gradient: {
    stops: [
      {
        position: 0.5,
        color: { r: 1, g: 0.5568627715110779, b: 0.18039216101169586, a: 1 },
      },
      {
        position: 1,
        color: { r: 1, g: 0.3019607961177826, b: 0.1725490242242813, a: 1 },
      },
    ],
  },
  gradientBalance: 0,
  material: 0,
  morphSpeed: 3.740000009536743,
  detail: 3.879999876022339,
  twist: 4.940000057220459,
  zoom: 19,
  gradientMethod: 0,
  warp: 1.5,
  rotationSpeed: 24,
};

/** Hosts Figma's unmodified shader in a WebGPU canvas with resize and lifecycle cleanup. */
export function ShaderFill({ className }: { className: string }) {
  const paused = useReducedMotion();
  const ref = useRef<HTMLCanvasElement>(null);
  useEffect(() => {
    const canvas = ref.current;
    if (!canvas || !navigator.gpu) return;
    let disposed = false;
    let animation = 0;
    let device: GPUDevice | undefined;
    let observer: ResizeObserver | undefined;

    async function start() {
      const adapter = await navigator.gpu.requestAdapter();
      if (!adapter || disposed || !canvas) return;
      device = await adapter.requestDevice();
      device.addEventListener("uncapturederror", (event) => {
        if (canvas) canvas.dataset.error = event.error.message;
      });
      void device.lost.then((info) => {
        if (!disposed) canvas.dataset.error = info.message || "GPU device lost";
      });
      if (disposed) {
        device.destroy();
        return;
      }
      const context = canvas.getContext("webgpu");
      if (!context) return;
      context.configure({
        device,
        format: navigator.gpu.getPreferredCanvasFormat(),
        alphaMode: "premultiplied",
      });
      const frame = {
        state: {},
        params,
        time: 0,
        output: context.getCurrentTexture(),
      };
      setup(device, frame);
      const begin = performance.now();
      const paint = (time: number) => {
        if (disposed || !device) return;
        frame.output = context.getCurrentTexture();
        frame.time = paused ? 0 : time - begin;
        render(device, frame);
        void device.queue.onSubmittedWorkDone().then(
          () => {
            if (!disposed) canvas.dataset.ready = "true";
          },
          () => {
            if (!disposed) canvas.dataset.ready = "unavailable";
          },
        );
      };
      const resize = () => {
        const ratio = Math.min(window.devicePixelRatio || 1, 1);
        const width = Math.max(1, Math.round(canvas.clientWidth * ratio));
        const height = Math.max(1, Math.round(canvas.clientHeight * ratio));
        // Resetting even unchanged canvas dimensions clears a static GPU frame.
        if (canvas.width === width && canvas.height === height) return;
        canvas.width = width;
        canvas.height = height;
        paint(performance.now());
      };
      resize();
      // Reinitializing after a pause keeps the same canvas dimensions.
      paint(performance.now());
      observer = new ResizeObserver(resize);
      observer.observe(canvas);
      const draw = (time: number) => {
        if (disposed) return;
        paint(time);
        if (!paused) animation = requestAnimationFrame(draw);
      };
      if (!paused) animation = requestAnimationFrame(draw);
    }
    void start().catch(() => {
      if (canvas) canvas.dataset.ready = "unavailable";
    });
    return () => {
      disposed = true;
      cancelAnimationFrame(animation);
      observer?.disconnect();
      device?.destroy();
    };
  }, [paused]);
  return <canvas ref={ref} className={className} aria-hidden="true" />;
}
