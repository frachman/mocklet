import type { NextConfig } from 'next';
const nextConfig: NextConfig = {
  reactStrictMode: true,
  turbopack: { root: process.cwd() },
  async redirects() {
    return [
      { source: "/docs", destination: "https://docs.mikrolyt.com/mocklet", permanent: true },
      { source: "/docs/:path*", destination: "https://docs.mikrolyt.com/mocklet/:path*", permanent: true },
    ];
  },
};
export default nextConfig;
