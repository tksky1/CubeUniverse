/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: '/joe/:path*',
        destination: 'http://192.168.177.201:30401/:path*'
      },
      {
        source: "/joe2/:path*",
        destination: 'http://192.168.177.201:30402/:path*'
      }
    ]
  }
}

module.exports = nextConfig
