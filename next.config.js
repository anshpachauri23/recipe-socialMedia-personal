/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: [
      'localhost', 
      'recipe-social-images.s3.us-east-2.amazonaws.com',
      'your-aws-s3-bucket.s3.amazonaws.com'
    ],
  },
  // Enable static exports for better Vercel performance
  output: 'standalone',
  // Optimize for production
  compress: true,
  poweredByHeader: false,
}

module.exports = nextConfig
