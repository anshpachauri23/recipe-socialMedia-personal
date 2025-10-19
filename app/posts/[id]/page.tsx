'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { Layout } from '@/components/Layout'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { PostCard } from '@/components/PostCard'
import { FiArrowLeft } from 'react-icons/fi'
import axios from 'axios'
import Cookies from 'js-cookie'

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

export default function PostDetailPage() {
  const params = useParams()
  const postId = params.id as string
  
  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (postId) {
      fetchPost()
    }
  }, [postId])

  const fetchPost = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_BASE_URL}/posts/${postId}`)
      setPost(response.data)
    } catch (error: any) {
      setError('Failed to load post')
      console.error('Post error:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleLike = async () => {
    if (!post) return

    try {
      const token = Cookies.get('token')
      if (!token) {
        setError('Please login to like posts')
        return
      }

      if (post.is_liked) {
        await axios.delete(`${API_BASE_URL}/posts/${post.id}/like`, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        })
        setPost({
          ...post,
          is_liked: false,
          total_likes: post.total_likes - 1
        })
      } else {
        await axios.post(`${API_BASE_URL}/posts/${post.id}/like`, {}, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        })
        setPost({
          ...post,
          is_liked: true,
          total_likes: post.total_likes + 1
        })
      }
    } catch (error) {
      console.error('Error liking post:', error)
    }
  }

  const handleDeletePost = async () => {
    if (!post) return

    if (!confirm('Are you sure you want to delete this post?')) return

    try {
      const token = Cookies.get('token')
      await axios.delete(`${API_BASE_URL}/posts/${post.id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })

      // Redirect to feed after deletion
      window.location.href = '/feed'
    } catch (error) {
      console.error('Error deleting post:', error)
      setError('Failed to delete post')
    }
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

  if (error || !post) {
    return (
      <Layout>
        <div className="text-center py-12">
          <p className="text-red-600 mb-4">{error || 'Post not found'}</p>
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
        {/* Back Button */}
        <div className="mb-6">
          <Link 
            href="/feed" 
            className="inline-flex items-center space-x-2 text-earth-600 hover:text-earth-800 transition-colors"
          >
            <FiArrowLeft className="h-4 w-4" />
            <span>Back to Feed</span>
          </Link>
        </div>

        {/* Post */}
        <div className="space-y-6">
          <PostCard
            post={post}
            onLike={handleLike}
            onDelete={handleDeletePost}
            showDelete={true} // Show delete for now, can be made conditional based on ownership
          />
        </div>
      </div>
    </Layout>
  )
}
