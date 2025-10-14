'use client'

import { useEffect, useState } from 'react'
import { useAuth } from '@/contexts/AuthContext'
import { Layout } from '@/components/Layout'
import { PostCard } from '@/components/PostCard'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { SearchBar } from '@/components/SearchBar'
import axios from 'axios'

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

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://760go4r862.execute-api.us-east-2.amazonaws.com/prod'

export default function FeedPage() {
  const { user } = useAuth()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (user) {
      fetchPosts()
    }
  }, [user])

  const fetchPosts = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_BASE_URL}/posts/feed`)
      setPosts(response.data)
    } catch (error: any) {
      setError('Failed to load posts')
      console.error('Error fetching posts:', error)
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

  if (!user) {
    return null
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 py-6">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-earth-800 mb-2">Your Feed</h1>
          <p className="text-earth-600">Discover amazing recipes from your network</p>
        </div>

        {/* Search Bar */}
        <div className="mb-6">
          <SearchBar />
        </div>

        {/* Posts */}
        {loading ? (
          <div className="flex justify-center py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : error ? (
          <div className="text-center py-12">
            <p className="text-red-600 mb-4">{error}</p>
            <button 
              onClick={fetchPosts}
              className="btn-primary"
            >
              Try Again
            </button>
          </div>
        ) : posts.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-6xl text-earth-300 mb-4">🍳</div>
            <h3 className="text-lg font-medium text-earth-800 mb-2">No posts yet</h3>
            <p className="text-earth-600 mb-6">
              Follow some users or create your first recipe to see posts here!
            </p>
            <a href="/create" className="btn-primary">
              Create Your First Recipe
            </a>
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
    </Layout>
  )
}
