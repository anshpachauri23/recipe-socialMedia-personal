import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { AuthProvider } from '@/contexts/AuthContext'
import { Toaster } from 'react-hot-toast'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'RecipeShare - Share Your Culinary Creations',
  description: 'A social media platform for sharing homemade recipes and culinary experiences',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AuthProvider>
          {children}
          <Toaster 
            position="top-right"
            toastOptions={{
              duration: 4000,
              style: {
                background: '#f7f5f0',
                color: '#2d2d2d',
                border: '1px solid #c4b693',
              },
            }}
          />
        </AuthProvider>
      </body>
    </html>
  )
}
