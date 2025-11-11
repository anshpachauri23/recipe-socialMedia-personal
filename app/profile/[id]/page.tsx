'use client'

import { useState, useEffect, useRef } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { Layout } from '@/components/Layout'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { PostCard } from '@/components/PostCard'
import { FiUserPlus, FiUserMinus, FiEdit, FiCamera } from 'react-icons/fi'
import axios from 'axios'
import Cookies from 'js-cookie'

interface User {
  id: number
  username: string
  full_name: string
  bio?: string
  profile_photo_url?: string
  follower_count: number
  following_count: number
  is_following: boolean
  is_own_profile: boolean
}

interface Post {
  id: number
  user_id: number
  title: string
  description?: string
  total_likes: number
  total_comments: number
  created_at: string
  user: {
    id: number
    username: string
    full_name: string
    profile_photo_url?: string
    follower_count: number
  }
  images: Array<{
    id: number
    image_url: string
    step_description?: string
    step_number: number
  }>
  is_liked: boolean
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export default function ProfilePage() {
  const params = useParams()
  const userId = params.id as string
  
  const [user, setUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [following, setFollowing] = useState(false)
  const [uploadingPhoto, setUploadingPhoto] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (userId && userId !== 'undefined') {
      fetchProfile()
    }
  }, [userId])

  const fetchProfile = async () => {
    if (!userId || userId === 'undefined') {
      return
    }

    try {
      setLoading(true)
      const token = Cookies.get('token')
      
      const headers = token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      } : {}

      const [profileResponse, postsResponse] = await Promise.all([
        axios.get(`${API_BASE_URL}/users/${userId}`, { headers }),
        axios.get(`${API_BASE_URL}/users/${userId}/posts`, { headers })
      ])

      setUser(profileResponse.data)
      setPosts(postsResponse.data)
    } catch (error: any) {
      setError('Failed to load profile')
      console.error('Profile error:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleFollow = async () => {
    if (!user) return

    try {
      if (user.is_following) {
        await axios.delete(`${API_BASE_URL}/users/${userId}/follow`)
        setUser({
          ...user,
          is_following: false,
          follower_count: user.follower_count - 1
        })
      } else {
        await axios.post(`${API_BASE_URL}/users/${userId}/follow`)
        setUser({
          ...user,
          is_following: true,
          follower_count: user.follower_count + 1
        })
      }
    } catch (error) {
      console.error('Follow error:', error)
    }
  }

  const handleLike = async (postId: number) => {
    try {
      const post = posts.find(p => p.id === postId)
      if (!post) return

      if (post.is_liked) {
        await axios.delete(`${API_BASE_URL}/posts/${postId}/like`)
        setPosts(posts.map(p => 
          p.id === postId 
            ? { ...p, is_liked: false, total_likes: p.total_likes - 1 }
            : p
        ))
      } else {
        await axios.post(`${API_BASE_URL}/posts/${postId}/like`)
        setPosts(posts.map(p => 
          p.id === postId 
            ? { ...p, is_liked: true, total_likes: p.total_likes + 1 }
            : p
        ))
      }
    } catch (error) {
      console.error('Error liking post:', error)
    }
  }

  const handlePhotoUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file || !user?.is_own_profile) return

    // Convert to base64
    const reader = new FileReader()
    reader.onload = async (e) => {
      const base64 = e.target?.result as string
      const base64Data = base64.split(',')[1] // Remove data:image/jpeg;base64, prefix

      try {
        setUploadingPhoto(true)
        const token = Cookies.get('token')
        
        const response = await axios.post(`${API_BASE_URL}/users/me/profile-photo`, {
          image_data: base64Data
        }, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        })

        // Update user profile photo
        if (user) {
          setUser({
            ...user,
            profile_photo_url: response.data.profile_photo_url
          })
        }
      } catch (error) {
        console.error('Error uploading photo:', error)
        setError('Failed to upload photo')
      } finally {
        setUploadingPhoto(false)
      }
    }
    reader.readAsDataURL(file)
  }

  const triggerPhotoUpload = () => {
    fileInputRef.current?.click()
  }

  const handleDeletePost = async (postId: number) => {
    if (!user?.is_own_profile) return

    if (!confirm('Are you sure you want to delete this post?')) return

    try {
      const token = Cookies.get('token')
      await axios.delete(`${API_BASE_URL}/posts/${postId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })

      // Remove post from local state
      setPosts(posts.filter(p => p.id !== postId))
    } catch (error) {
      console.error('Error deleting post:', error)
      setError('Failed to delete post')
    }
  }

  if (!userId || userId === 'undefined') {
    return (
      <Layout>
        <div className="text-center py-12">
          <p className="text-red-600 mb-4">Invalid user ID</p>
          <Link href="/feed" className="btn-primary">
            Back to Feed
          </Link>
        </div>
      </Layout>
    )
  }

  if (loading) {
    return (
      <Layout>
        <div className="flex justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      </Layout>
    )
  }

  if (error || !user) {
    return (
      <Layout>
        <div className="text-center py-12">
          <p className="text-red-600 mb-4">{error || 'User not found'}</p>
          <Link href="/feed" className="btn-primary">
            Back to Feed
          </Link>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 py-6">
        {/* Profile Header */}
        <div className="card mb-6">
          <div className="flex flex-col md:flex-row items-center md:items-start space-y-4 md:space-y-0 md:space-x-6">
            <div className="relative">
              <img
                src={user.profile_photo_url || `data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="150" height="150"><rect width="150" height="150" fill="%238B7355"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="white" font-size="60" font-family="sans-serif">${user.username?.charAt(0)?.toUpperCase() || 'U'}</text></svg>`}
                alt={user.username}
                className="avatar-xl bg-earth-400"
              />
              {user.is_own_profile && (
                <button
                  onClick={triggerPhotoUpload}
                  disabled={uploadingPhoto}
                  className="absolute bottom-0 right-0 bg-earth-500 text-white rounded-full p-2 hover:bg-earth-600 disabled:opacity-50 transition-colors"
                >
                  {uploadingPhoto ? (
                    <LoadingSpinner size="sm" />
                  ) : (
                    <FiCamera className="h-4 w-4" />
                  )}
                </button>
              )}
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                onChange={handlePhotoUpload}
                className="hidden"
              />
            </div>
            
            <div className="flex-1 text-center md:text-left">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-4">
                <div>
                  <h1 className="text-2xl font-bold text-earth-800">{user.username}</h1>
                  <p className="text-earth-600">{user.full_name}</p>
                </div>
                
                <div className="flex space-x-3 mt-4 md:mt-0">
                  {user.is_own_profile ? (
                    <Link href="/settings" className="btn-outline flex items-center space-x-2">
                      <FiEdit className="h-4 w-4" />
                      <span>Edit Profile</span>
                    </Link>
                  ) : (
                    <button
                      onClick={handleFollow}
                      className={`flex items-center space-x-2 ${
                        user.is_following ? 'btn-outline' : 'btn-primary'
                      }`}
                    >
                      {user.is_following ? (
                        <>
                          <FiUserMinus className="h-4 w-4" />
                          <span>Unfollow</span>
                        </>
                      ) : (
                        <>
                          <FiUserPlus className="h-4 w-4" />
                          <span>Follow</span>
                        </>
                      )}
                    </button>
                  )}
                </div>
              </div>

              {user.bio && (
                <p className="text-earth-600 mb-4">{user.bio}</p>
              )}

              <div className="flex justify-center md:justify-start space-x-6 text-sm">
                <div className="text-center">
                  <div className="font-semibold text-earth-800">{posts?.length || 0}</div>
                  <div className="text-earth-600">Recipes</div>
                </div>
                <div className="text-center">
                  <div className="font-semibold text-earth-800">{user.follower_count}</div>
                  <div className="text-earth-600">Followers</div>
                </div>
                <div className="text-center">
                  <div className="font-semibold text-earth-800">{user.following_count}</div>
                  <div className="text-earth-600">Following</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Posts */}
        <div>
          <h2 className="text-xl font-semibold text-earth-800 mb-4">
            {user.is_own_profile ? 'Your Recipes' : 'Recipes'}
          </h2>
          
          {!posts || posts.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-6xl text-earth-300 mb-4">🍳</div>
              <h3 className="text-lg font-medium text-earth-800 mb-2">
                {user.is_own_profile ? 'No recipes yet' : 'No recipes posted'}
              </h3>
              <p className="text-earth-600 mb-6">
                {user.is_own_profile 
                  ? 'Start sharing your culinary creations!'
                  : 'This user hasn\'t shared any recipes yet.'
                }
              </p>
              {user.is_own_profile && (
                <Link href="/create" className="btn-primary">
                  Create Your First Recipe
                </Link>
              )}
            </div>
          ) : (
            <div className="space-y-6">
              {posts.map((post) => (
                <PostCard
                  key={post.id}
                  post={post}
                  onLike={() => handleLike(post.id)}
                  onDelete={() => handleDeletePost(post.id)}
                  showDelete={user.is_own_profile}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  )
}
