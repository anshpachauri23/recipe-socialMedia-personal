/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost', 'your-aws-s3-bucket.s3.amazonaws.com'],
  },
}

module.exports = nextConfig
