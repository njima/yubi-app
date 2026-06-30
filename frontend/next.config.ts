import { withSentryConfig } from "@sentry/nextjs";

import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  basePath: "/web",
  async redirects() {
    return [
      {
        source: "/",
        destination: "/web",
        permanent: false,
        basePath: false,
      },
    ];
  },
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
};

export default withSentryConfig(nextConfig, {
  // Suppress source map upload logs in CI
  silent: !process.env.CI,

  // Upload larger set of source maps for readable error stack traces
  widenClientFileUpload: true,

  // Route browser requests to Sentry through a Next.js rewrite to circumvent ad-blockers
  tunnelRoute: "/monitoring",

  // Disable Sentry SDK logger to reduce bundle size
  disableLogger: true,

  // Automatically tree-shake Sentry logger statements to reduce bundle size
  automaticVercelMonitors: false,
});
