import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  /* config options here */
  cacheComponents: true,
  partialPrefetching: true,
  async rewrites() {
    const apiOrigin = process.env.API_URL || "http://api:8080";
    return [
      {
        source: "/v1/:path*",
        destination: `${apiOrigin}/v1/:path*`,
      },
      {
        source: "/health",
        destination: `${apiOrigin}/health`,
      },
    ];
  },
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
};

export default nextConfig;
