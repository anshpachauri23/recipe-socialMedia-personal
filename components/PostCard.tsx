'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { FiHeart, FiMessageCircle, FiShare2, FiMoreHorizontal, FiSend, FiTrash2 } from 'react-icons/fi'
import { formatDistanceToNow } from 'date-fns'
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

interface Comment {
  id: number
  user_id: number
  post_id: number
  content: string
  created_at: string
  user: {
    id: number
    username: string
    full_name: string
    profile_photo_url?: string
  }
}

interface PostCardProps {
  post: Post
  onLike: () => void
  onDelete?: () => void
  showDelete?: boolean
}

export function PostCard({ post, onLike, onDelete, showDelete = false }: PostCardProps) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0)
  const [showComments, setShowComments] = useState(false)
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [loadingComments, setLoadingComments] = useState(false)
  const [submittingComment, setSubmittingComment] = useState(false)

  const nextImage = () => {
    setCurrentImageIndex((prev) => 
      prev === post.images.length - 1 ? 0 : prev + 1
    )
  }

  const prevImage = () => {
    setCurrentImageIndex((prev) => 
      prev === 0 ? post.images.length - 1 : prev - 1
    )
  }

  const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

  const fetchComments = async () => {
    if (loadingComments) return
    setLoadingComments(true)
    try {
      const token = Cookies.get('token')
      
      if (!token) {
        console.error('No authentication token found')
        return
      }
      
      const response = await axios.get(`${API_BASE_URL}/posts/${post.id}/comments`, {
        headers: { 
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
      setComments(response.data)
    } catch (error: any) {
      console.error('Error fetching comments:', error)
      if (error.response?.status === 401) {
        console.error('Authentication failed - token may be expired')
        // Optionally redirect to login or refresh token
      }
    } finally {
      setLoadingComments(false)
    }
  }

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newComment.trim() || submittingComment) return

    setSubmittingComment(true)
    try {
      const token = Cookies.get('token')
      
      if (!token) {
        console.error('No authentication token found')
        return
      }
      
      const response = await axios.post(`${API_BASE_URL}/posts/${post.id}/comments`, {
        content: newComment.trim()
      }, {
        headers: { 
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
      setComments([response.data, ...(comments || [])])
      setNewComment('')
    } catch (error: any) {
      console.error('Error submitting comment:', error)
      if (error.response?.status === 401) {
        console.error('Authentication failed - token may be expired')
      }
    } finally {
      setSubmittingComment(false)
    }
  }

  const toggleComments = () => {
    if (!showComments) {
      fetchComments()
    }
    setShowComments(!showComments)
  }

  return (
    <article className="post-card">
      {/* Header */}
      <div className="flex items-center justify-between p-4">
        <div className="flex items-center space-x-3">
          <Link href={`/profile/${post.user.id}`}>
            <img
              src={post.user.profile_photo_url || 'https://via.placeholder.com/40x40/8B7355/FFFFFF?text=' + post.user.username.charAt(0).toUpperCase()}
              alt={post.user.username}
              className="avatar"
            />
          </Link>
          <div>
            <Link 
              href={`/profile/${post.user.id}`}
              className="font-medium text-earth-800 hover:text-earth-700"
            >
              {post.user.username}
            </Link>
            <p className="text-sm text-earth-500">
              {formatDistanceToNow(new Date(post.created_at), { addSuffix: true })}
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-2">
          {showDelete && onDelete && (
            <button 
              onClick={onDelete}
              className="text-red-400 hover:text-red-600 transition-colors"
              title="Delete post"
            >
              <FiTrash2 className="h-5 w-5" />
            </button>
          )}
          <button className="text-earth-400 hover:text-earth-600">
            <FiMoreHorizontal className="h-5 w-5" />
          </button>
        </div>
      </div>

      {/* Images */}
      {post.images && post.images.length > 0 && (
        <div className="relative">
          <img
            src={post.images[currentImageIndex]?.image_url}
            alt={`Step ${post.images[currentImageIndex]?.step_number}`}
            className="post-image"
          />
          
          {/* Image navigation */}
          {post.images.length > 1 && (
            <>
              <button
                onClick={prevImage}
                className="absolute left-2 top-1/2 transform -translate-y-1/2 bg-black bg-opacity-50 text-white p-2 rounded-full hover:bg-opacity-75"
              >
                ‹
              </button>
              <button
                onClick={nextImage}
                className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-black bg-opacity-50 text-white p-2 rounded-full hover:bg-opacity-75"
              >
                ›
              </button>
              
              {/* Image indicators */}
              <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2 flex space-x-1">
                {post.images.map((_, index) => (
                  <button
                    key={index}
                    onClick={() => setCurrentImageIndex(index)}
                    className={`w-2 h-2 rounded-full ${
                      index === currentImageIndex ? 'bg-white' : 'bg-white bg-opacity-50'
                    }`}
                  />
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="post-actions">
        <div className="flex items-center space-x-4">
          <button
            onClick={onLike}
            className={`action-btn ${post.is_liked ? 'like-btn liked' : ''}`}
          >
            <FiHeart className={`h-6 w-6 ${post.is_liked ? 'fill-current' : ''}`} />
            <span>{post.total_likes}</span>
          </button>
          
          <button
            onClick={toggleComments}
            className="action-btn"
          >
            <FiMessageCircle className="h-6 w-6" />
            <span>{post.total_comments}</span>
          </button>
          
          <button className="action-btn">
            <FiShare2 className="h-6 w-6" />
            <span>Share</span>
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 pb-4">
        <Link href={`/posts/${post.id}`} className="block hover:bg-earth-50 -mx-4 px-4 py-2 rounded-lg transition-colors">
          <h3 className="font-semibold text-earth-800 mb-2">{post.title}</h3>
          {post.description && (
            <p className="text-earth-600 mb-3">{post.description}</p>
          )}
        </Link>
        
        {/* Step description for current image */}
        {post.images && post.images[currentImageIndex]?.step_description && (
          <div className="bg-earth-50 rounded-lg p-3">
            <div className="flex items-center mb-2">
              <span className="step-number">
                {post.images[currentImageIndex]?.step_number}
              </span>
              <span className="ml-2 text-sm font-medium text-earth-700">Step</span>
            </div>
            <p className="text-earth-600 text-sm">
              {post.images[currentImageIndex]?.step_description}
            </p>
          </div>
        )}
      </div>

      {/* Comments section */}
      {showComments && (
        <div className="border-t border-earth-100 px-4 py-3">
          <div className="space-y-3">
            
            {/* Comment form */}
            <form onSubmit={handleCommentSubmit} className="flex space-x-3">
              <img
                src="https://via.placeholder.com/32x32/8B7355/FFFFFF?text=U"
                alt="Your avatar"
                className="w-8 h-8 rounded-full"
              />
              <div className="flex-1 flex items-center space-x-2">
                <input
                  type="text"
                  placeholder="Add a comment..."
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  disabled={submittingComment}
                  className="flex-1 px-3 py-2 border border-earth-200 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-earth-500 disabled:opacity-50"
                />
                <button
                  type="submit"
                  disabled={!newComment.trim() || submittingComment}
                  className="p-2 rounded-full text-earth-500 hover:text-earth-700 hover:bg-earth-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                >
                  <FiSend className="h-4 w-4" />
                </button>
              </div>
            </form>
            
            {/* Comments list */}
            {loadingComments ? (
              <div className="text-center text-earth-500 text-sm py-2">
                Loading comments...
              </div>
            ) : comments && comments.length > 0 ? (
              <div className="space-y-3">
                {comments.map((comment) => (
                  <div key={comment.id} className="flex space-x-3">
                    <img
                      src={comment.user.profile_photo_url || `data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32"><rect width="32" height="32" fill="%238B7355"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="white" font-size="16" font-family="sans-serif">${comment.user.full_name?.charAt(0)?.toUpperCase() || 'U'}</text></svg>`}
                      alt={comment.user.full_name || 'User'}
                      className="w-8 h-8 rounded-full bg-earth-400"
                    />
                    <div className="flex-1">
                      <div className="flex items-center space-x-2">
                        <span className="font-medium text-earth-800 text-sm">
                          {comment.user.full_name || 'User'}
                        </span>
                        <span className="text-earth-500 text-xs">
                          {formatDistanceToNow(new Date(comment.created_at), { addSuffix: true })}
                        </span>
                      </div>
                      <p className="text-earth-700 text-sm mt-1">{comment.content}</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center text-earth-500 text-sm py-2">
                 No comments yet. Be the first to comment!
              </div>
            )}
          </div>
        </div>
      )}
    </article>
  )
}
