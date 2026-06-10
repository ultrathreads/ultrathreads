// next.config.ts
import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.GO_API_BASE}/:path*`,
      },
    ];
  },
};

export default nextConfig;