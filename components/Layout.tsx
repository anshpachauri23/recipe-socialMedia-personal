'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter, usePathname } from 'next/navigation'
import { useAuth } from '@/contexts/AuthContext'
import { 
  FiHome, 
  FiSearch, 
  FiPlus, 
  FiUser, 
  FiSettings, 
  FiLogOut,
  FiMenu,
  FiX
} from 'react-icons/fi'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { user, logout } = useAuth()
  const router = useRouter()
  const pathname = usePathname()

  const handleLogout = () => {
    logout()
  }

  const navigation = [
    { name: 'Home', href: '/feed', icon: FiHome },
    { name: 'Search', href: '/search', icon: FiSearch },
    { name: 'Create', href: '/create', icon: FiPlus },
    { name: 'Profile', href: `/profile/${user?.id}`, icon: FiUser },
    { name: 'Settings', href: '/settings', icon: FiSettings },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-earth-50 to-sage-50">
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-black bg-opacity-25" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 w-64 bg-white shadow-xl">
          <div className="flex items-center justify-between p-4 border-b border-earth-200">
            <h1 className="text-xl font-bold text-earth-800">RecipeShare</h1>
            <button onClick={() => setSidebarOpen(false)}>
              <FiX className="h-6 w-6 text-earth-600" />
            </button>
          </div>
          <nav className="mt-4">
            {navigation.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`nav-link ${isActive ? 'active' : ''}`}
                  onClick={() => setSidebarOpen(false)}
                >
                  <Icon className="h-5 w-5" />
                  <span>{item.name}</span>
                </Link>
              )
            })}
            <button
              onClick={handleLogout}
              className="nav-link text-red-600 hover:text-red-700 hover:bg-red-50"
            >
              <FiLogOut className="h-5 w-5" />
              <span>Logout</span>
            </button>
          </nav>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
        <div className="flex flex-col flex-grow bg-white shadow-sm border-r border-earth-200">
          <div className="flex items-center px-4 py-6 border-b border-earth-200">
            <h1 className="text-xl font-bold text-earth-800">RecipeShare</h1>
          </div>
          <nav className="mt-4 flex-1 px-2">
            {navigation.map((item) => {
              const Icon = item.icon
              const isActive = pathname === item.href
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`nav-link ${isActive ? 'active' : ''}`}
                >
                  <Icon className="h-5 w-5" />
                  <span>{item.name}</span>
                </Link>
              )
            })}
            <button
              onClick={handleLogout}
              className="nav-link text-red-600 hover:text-red-700 hover:bg-red-50"
            >
              <FiLogOut className="h-5 w-5" />
              <span>Logout</span>
            </button>
          </nav>
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Mobile header */}
        <div className="lg:hidden bg-white shadow-sm border-b border-earth-200">
          <div className="flex items-center justify-between px-4 py-3">
            <button onClick={() => setSidebarOpen(true)}>
              <FiMenu className="h-6 w-6 text-earth-600" />
            </button>
            <h1 className="text-lg font-semibold text-earth-800">RecipeShare</h1>
            <div className="w-6" />
          </div>
        </div>

        {/* Page content */}
        <main className="flex-1">
          {children}
        </main>
      </div>
    </div>
  )
}
