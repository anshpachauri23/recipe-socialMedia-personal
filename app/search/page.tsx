'use client'

import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import { Layout } from '@/components/Layout'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { PostCard } from '@/components/PostCard'
import { SearchBar } from '@/components/SearchBar'
import { FiUser, FiSearch } from 'react-icons/fi'
import axios from 'axios'

interface User {
  id: number
  username: string
  full_name: string
  profile_photo_url?: string
  follower_count: number
  following_count: number
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

export default function SearchPage() {
  const searchParams = useSearchParams()
  const query = searchParams.get('q') || ''
  
  const [activeTab, setActiveTab] = useState<'posts' | 'users'>('posts')
  const [posts, setPosts] = useState<Post[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (query) {
      performSearch()
    }
  }, [query])

  const performSearch = async () => {
    if (!query.trim()) return

    setLoading(true)
    setError(null)

    try {
      const [postsResponse, usersResponse] = await Promise.all([
        axios.get(`${API_BASE_URL}/posts/search?q=${encodeURIComponent(query)}`),
        axios.get(`${API_BASE_URL}/users/search?q=${encodeURIComponent(query)}`)
      ])

      setPosts(postsResponse.data)
      setUsers(usersResponse.data)
    } catch (error: any) {
      setError('Failed to search')
      console.error('Search error:', error)
    } finally {
      setLoading(false)
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

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 py-6">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-earth-800 mb-2">Search</h1>
          <p className="text-earth-600">Find recipes and people</p>
        </div>

        {/* Search Bar */}
        <div className="mb-6">
          <SearchBar />
        </div>

        {/* Results */}
        {query && (
          <>
            {/* Tabs */}
            <div className="flex space-x-1 mb-6 bg-earth-100 rounded-lg p-1">
              <button
                onClick={() => setActiveTab('posts')}
                className={`flex-1 flex items-center justify-center space-x-2 py-2 px-4 rounded-md transition-colors duration-200 ${
                  activeTab === 'posts'
                    ? 'bg-white text-earth-800 shadow-sm'
                    : 'text-earth-600 hover:text-earth-800'
                }`}
              >
                <FiSearch className="h-4 w-4" />
                <span>Recipes ({posts.length})</span>
              </button>
              <button
                onClick={() => setActiveTab('users')}
                className={`flex-1 flex items-center justify-center space-x-2 py-2 px-4 rounded-md transition-colors duration-200 ${
                  activeTab === 'users'
                    ? 'bg-white text-earth-800 shadow-sm'
                    : 'text-earth-600 hover:text-earth-800'
                }`}
              >
                <FiUser className="h-4 w-4" />
                <span>Users ({users.length})</span>
              </button>
            </div>

            {/* Loading */}
            {loading && (
              <div className="flex justify-center py-12">
                <LoadingSpinner size="lg" />
              </div>
            )}

            {/* Error */}
            {error && (
              <div className="text-center py-12">
                <p className="text-red-600 mb-4">{error}</p>
                <button 
                  onClick={performSearch}
                  className="btn-primary"
                >
                  Try Again
                </button>
              </div>
            )}

            {/* Results */}
            {!loading && !error && (
              <>
                {/* Posts Results */}
                {activeTab === 'posts' && (
                  <div>
                    {posts.length === 0 ? (
                      <div className="text-center py-12">
                        <div className="text-6xl text-earth-300 mb-4">🔍</div>
                        <h3 className="text-lg font-medium text-earth-800 mb-2">No recipes found</h3>
                        <p className="text-earth-600">
                          Try searching with different keywords
                        </p>
                      </div>
                    ) : (
                      <div className="space-y-6">
                        {posts.map((post) => (
                          <PostCard
                            key={post.id}
                            post={post}
                            onLike={() => handleLike(post.id)}
                          />
                        ))}
                      </div>
                    )}
                  </div>
                )}

                {/* Users Results */}
                {activeTab === 'users' && (
                  <div>
                    {users.length === 0 ? (
                      <div className="text-center py-12">
                        <div className="text-6xl text-earth-300 mb-4">👤</div>
                        <h3 className="text-lg font-medium text-earth-800 mb-2">No users found</h3>
                        <p className="text-earth-600">
                          Try searching with different keywords
                        </p>
                      </div>
                    ) : (
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {users.map((user) => (
                          <div key={user.id} className="card">
                            <div className="flex items-center space-x-4">
                              <img
                                src={user.profile_photo_url || '/default-avatar.png'}
                                alt={user.username}
                                className="avatar-lg"
                              />
                              <div className="flex-1">
                                <h3 className="font-semibold text-earth-800">
                                  {user.username}
                                </h3>
                                <p className="text-earth-600">{user.full_name}</p>
                                <div className="flex space-x-4 mt-2 text-sm text-earth-500">
                                  <span>{user.follower_count} followers</span>
                                  <span>{user.following_count} following</span>
                                </div>
                              </div>
                              <a
                                href={`/profile/${user.id}`}
                                className="btn-outline"
                              >
                                View Profile
                              </a>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </>
            )}
          </>
        )}

        {/* No search query */}
        {!query && (
          <div className="text-center py-12">
            <div className="text-6xl text-earth-300 mb-4">🔍</div>
            <h3 className="text-lg font-medium text-earth-800 mb-2">Start searching</h3>
            <p className="text-earth-600">
              Enter a search term above to find recipes and people
            </p>
          </div>
        )}
      </div>
    </Layout>
  )
}
