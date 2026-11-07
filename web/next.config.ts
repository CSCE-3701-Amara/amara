import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  /* config options here */
  cacheComponents: true,
  partialPrefetching: true,
  async rewrites() {
    const apiOrigin = (
      process.env.API_URL ||
      (process.env.NODE_ENV === "development"
        ? "http://127.0.0.1:8080"
        : "http://api:8080")
    ).replace(/\/$/, "");
    return [
      {
        source: "/v1/:path*",
        destination: `${apiOrigin}/v1/:path*`,
      },
      {
        source: "/health/:path*",
        destination: `${apiOrigin}/health/:path*`,
      },
    ];
  },
  turbopack: {
    rules: {
      // Keep native CSS Module processing; only the global entry imports Tailwind.
      "globals.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
};

export default nextConfig;
