/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: '/joe/:path*',
        destination: 'http://control-backend.cubeuniverse.svc.cluster.local:30401/:path*'
      },
      {
        source: "/joe2/:path*",
        destination: 'http://object-storage.cubeuniverse.svc.cluster.local:30402/:path*'
      }
    ]
  }
}

module.exports = nextConfig
