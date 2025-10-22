'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Layout } from '@/components/Layout'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { FiUser, FiLock, FiTrash2, FiSave } from 'react-icons/fi'
import toast from 'react-hot-toast'
import axios from 'axios'

interface User {
  id: number
  username: string
  email: string
  full_name: string
  bio?: string
  profile_photo_url?: string
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<'profile' | 'password' | 'danger'>('profile')
  const [loading, setLoading] = useState(false)
  const [user, setUser] = useState<User | null>(null)
  const router = useRouter()

  // Profile form
  const [profileData, setProfileData] = useState({
    full_name: '',
    bio: '',
    profile_photo_url: ''
  })

  // Password form
  const [passwordData, setPasswordData] = useState({
    current_password: '',
    new_password: '',
    confirm_password: ''
  })

  useEffect(() => {
    fetchUser()
  }, [])

  const fetchUser = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/users/me`)
      const userData = response.data
      setUser(userData)
      setProfileData({
        full_name: userData.full_name,
        bio: userData.bio || '',
        profile_photo_url: userData.profile_photo_url || ''
      })
    } catch (error) {
      console.error('Failed to fetch user:', error)
    }
  }

  const handleProfileUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      await axios.put(`${API_BASE_URL}/users/me`, profileData)
      toast.success('Profile updated successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to update profile')
    } finally {
      setLoading(false)
    }
  }

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (passwordData.new_password !== passwordData.confirm_password) {
      toast.error('New passwords do not match')
      return
    }

    if (passwordData.new_password.length < 6) {
      toast.error('New password must be at least 6 characters')
      return
    }

    setLoading(true)

    try {
      await axios.put(`${API_BASE_URL}/users/me/password`, {
        current_password: passwordData.current_password,
        new_password: passwordData.new_password
      })
      toast.success('Password updated successfully!')
      setPasswordData({
        current_password: '',
        new_password: '',
        confirm_password: ''
      })
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to update password')
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteAccount = async () => {
    if (!confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
      return
    }

    if (!confirm('This will permanently delete your account and all your data. Are you absolutely sure?')) {
      return
    }

    setLoading(true)

    try {
      await axios.delete(`${API_BASE_URL}/users/me`)
      toast.success('Account deleted successfully')
      router.push('/auth/login')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to delete account')
    } finally {
      setLoading(false)
    }
  }

  if (!user) {
    return (
      <Layout>
        <div className="flex justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-2xl mx-auto px-4 py-6">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-earth-800 mb-2">Account Settings</h1>
          <p className="text-earth-600">Manage your account preferences and security</p>
        </div>

        {/* Tabs */}
        <div className="flex space-x-1 mb-6 bg-earth-100 rounded-lg p-1">
          <button
            onClick={() => setActiveTab('profile')}
            className={`flex-1 flex items-center justify-center space-x-2 py-2 px-4 rounded-md transition-colors duration-200 ${
              activeTab === 'profile'
                ? 'bg-white text-earth-800 shadow-sm'
                : 'text-earth-600 hover:text-earth-800'
            }`}
          >
            <FiUser className="h-4 w-4" />
            <span>Profile</span>
          </button>
          <button
            onClick={() => setActiveTab('password')}
            className={`flex-1 flex items-center justify-center space-x-2 py-2 px-4 rounded-md transition-colors duration-200 ${
              activeTab === 'password'
                ? 'bg-white text-earth-800 shadow-sm'
                : 'text-earth-600 hover:text-earth-800'
            }`}
          >
            <FiLock className="h-4 w-4" />
            <span>Password</span>
          </button>
          <button
            onClick={() => setActiveTab('danger')}
            className={`flex-1 flex items-center justify-center space-x-2 py-2 px-4 rounded-md transition-colors duration-200 ${
              activeTab === 'danger'
                ? 'bg-white text-earth-800 shadow-sm'
                : 'text-earth-600 hover:text-earth-800'
            }`}
          >
            <FiTrash2 className="h-4 w-4" />
            <span>Danger</span>
          </button>
        </div>

        {/* Profile Tab */}
        {activeTab === 'profile' && (
          <form onSubmit={handleProfileUpdate} className="space-y-6">
            <div className="card">
              <h3 className="text-lg font-semibold text-earth-800 mb-4">Profile Information</h3>
              
              <div className="form-group">
                <label className="form-label">Full Name</label>
                <input
                  type="text"
                  className="input-field"
                  value={profileData.full_name}
                  onChange={(e) => setProfileData({ ...profileData, full_name: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label className="form-label">Bio</label>
                <textarea
                  className="input-field resize-none"
                  rows={3}
                  placeholder="Tell us about yourself..."
                  value={profileData.bio}
                  onChange={(e) => setProfileData({ ...profileData, bio: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label className="form-label">Profile Photo URL</label>
                <input
                  type="url"
                  className="input-field"
                  placeholder="https://example.com/photo.jpg"
                  value={profileData.profile_photo_url}
                  onChange={(e) => setProfileData({ ...profileData, profile_photo_url: e.target.value })}
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="btn-primary flex items-center space-x-2"
              >
                {loading ? (
                  <LoadingSpinner size="sm" />
                ) : (
                  <FiSave className="h-4 w-4" />
                )}
                <span>Save Changes</span>
              </button>
            </div>
          </form>
        )}

        {/* Password Tab */}
        {activeTab === 'password' && (
          <form onSubmit={handlePasswordChange} className="space-y-6">
            <div className="card">
              <h3 className="text-lg font-semibold text-earth-800 mb-4">Change Password</h3>
              
              <div className="form-group">
                <label className="form-label">Current Password</label>
                <input
                  type="password"
                  className="input-field"
                  value={passwordData.current_password}
                  onChange={(e) => setPasswordData({ ...passwordData, current_password: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label className="form-label">New Password</label>
                <input
                  type="password"
                  className="input-field"
                  value={passwordData.new_password}
                  onChange={(e) => setPasswordData({ ...passwordData, new_password: e.target.value })}
                />
              </div>

              <div className="form-group">
                <label className="form-label">Confirm New Password</label>
                <input
                  type="password"
                  className="input-field"
                  value={passwordData.confirm_password}
                  onChange={(e) => setPasswordData({ ...passwordData, confirm_password: e.target.value })}
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="btn-primary flex items-center space-x-2"
              >
                {loading ? (
                  <LoadingSpinner size="sm" />
                ) : (
                  <FiLock className="h-4 w-4" />
                )}
                <span>Update Password</span>
              </button>
            </div>
          </form>
        )}

        {/* Danger Tab */}
        {activeTab === 'danger' && (
          <div className="card">
            <h3 className="text-lg font-semibold text-earth-800 mb-4">Danger Zone</h3>
            
            <div className="border border-red-200 rounded-lg p-4 bg-red-50">
              <h4 className="font-medium text-red-800 mb-2">Delete Account</h4>
              <p className="text-red-600 text-sm mb-4">
                Once you delete your account, there is no going back. Please be certain.
              </p>
              <button
                onClick={handleDeleteAccount}
                disabled={loading}
                className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200 disabled:opacity-50"
              >
                {loading ? (
                  <div className="flex items-center space-x-2">
                    <LoadingSpinner size="sm" />
                    <span>Deleting...</span>
                  </div>
                ) : (
                  'Delete Account'
                )}
              </button>
            </div>
          </div>
        )}
      </div>
    </Layout>
  )
}
