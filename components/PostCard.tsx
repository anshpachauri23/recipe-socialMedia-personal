'use client'

import { useState } from 'react'
import Link from 'next/link'
import { FiHeart, FiMessageCircle, FiShare2, FiMoreHorizontal } from 'react-icons/fi'
import { formatDistanceToNow } from 'date-fns'

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

interface PostCardProps {
  post: Post
  onLike: () => void
}

export function PostCard({ post, onLike }: PostCardProps) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0)
  const [showComments, setShowComments] = useState(false)

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

  return (
    <article className="post-card">
      {/* Header */}
      <div className="flex items-center justify-between p-4">
        <div className="flex items-center space-x-3">
          <Link href={`/profile/${post.user.id}`}>
            <img
              src={post.user.profile_photo_url || '/default-avatar.png'}
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
        <button className="text-earth-400 hover:text-earth-600">
          <FiMoreHorizontal className="h-5 w-5" />
        </button>
      </div>

      {/* Images */}
      {post.images.length > 0 && (
        <div className="relative">
          <img
            src={post.images[currentImageIndex].image_url}
            alt={`Step ${post.images[currentImageIndex].step_number}`}
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
            onClick={() => setShowComments(!showComments)}
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
        <h3 className="font-semibold text-earth-800 mb-2">{post.title}</h3>
        {post.description && (
          <p className="text-earth-600 mb-3">{post.description}</p>
        )}
        
        {/* Step description for current image */}
        {post.images[currentImageIndex]?.step_description && (
          <div className="bg-earth-50 rounded-lg p-3">
            <div className="flex items-center mb-2">
              <span className="step-number">
                {post.images[currentImageIndex].step_number}
              </span>
              <span className="ml-2 text-sm font-medium text-earth-700">Step</span>
            </div>
            <p className="text-earth-600 text-sm">
              {post.images[currentImageIndex].step_description}
            </p>
          </div>
        )}
      </div>

      {/* Comments section */}
      {showComments && (
        <div className="border-t border-earth-100 px-4 py-3">
          <div className="space-y-3">
            {/* Comment form */}
            <div className="flex space-x-3">
              <img
                src="/default-avatar.png"
                alt="Your avatar"
                className="w-8 h-8 rounded-full"
              />
              <div className="flex-1">
                <input
                  type="text"
                  placeholder="Add a comment..."
                  className="w-full px-3 py-2 border border-earth-200 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-earth-500"
                />
              </div>
            </div>
            
            {/* Comments list would go here */}
            <div className="text-center text-earth-500 text-sm">
              Comments feature coming soon...
            </div>
          </div>
        </div>
      )}
    </article>
  )
}
