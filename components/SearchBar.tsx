'use client'

import { useState } from 'react'
import { FiSearch } from 'react-icons/fi'
import { useRouter } from 'next/navigation'

export function SearchBar() {
  const [query, setQuery] = useState('')
  const router = useRouter()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.trim()) {
      router.push(`/search?q=${encodeURIComponent(query.trim())}`)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="search-bar">
      <div className="relative">
        <FiSearch className="search-icon" />
        <input
          type="text"
          placeholder="Search recipes, users, or ingredients..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="search-input"
        />
      </div>
    </form>
  )
}
