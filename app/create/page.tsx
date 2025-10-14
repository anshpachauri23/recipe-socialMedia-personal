'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Layout } from '@/components/Layout'
import { LoadingSpinner } from '@/components/LoadingSpinner'
import { FiPlus, FiX, FiUpload } from 'react-icons/fi'
import toast from 'react-hot-toast'
import axios from 'axios'

interface RecipeStep {
  id: string
  image: File | null
  imageUrl: string
  description: string
  stepNumber: number
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://760go4r862.execute-api.us-east-2.amazonaws.com/prod'

export default function CreatePage() {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [steps, setSteps] = useState<RecipeStep[]>([
    { id: '1', image: null, imageUrl: '', description: '', stepNumber: 1 }
  ])
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  const addStep = () => {
    const newStep: RecipeStep = {
      id: Date.now().toString(),
      image: null,
      imageUrl: '',
      description: '',
      stepNumber: steps.length + 1
    }
    setSteps([...steps, newStep])
  }

  const removeStep = (stepId: string) => {
    if (steps.length > 1) {
      const updatedSteps = steps
        .filter(step => step.id !== stepId)
        .map((step, index) => ({ ...step, stepNumber: index + 1 }))
      setSteps(updatedSteps)
    }
  }

  const updateStep = (stepId: string, field: keyof RecipeStep, value: any) => {
    setSteps(steps.map(step => 
      step.id === stepId ? { ...step, [field]: value } : step
    ))
  }

  const handleImageUpload = (stepId: string, file: File) => {
    // In a real app, you would upload to AWS S3 here
    // For now, we'll create a local URL
    const imageUrl = URL.createObjectURL(file)
    setSteps(steps.map(step => 
      step.id === stepId 
        ? { ...step, image: file, imageUrl: imageUrl }
        : step
    ))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!title.trim()) {
      toast.error('Please enter a recipe title')
      return
    }

    // Debug: Log the current steps state
    console.log('Current steps:', steps.map(step => ({
      stepNumber: step.stepNumber,
      hasImage: !!step.image,
      hasImageUrl: !!step.imageUrl,
      hasDescription: !!step.description.trim(),
      description: step.description
    })))

    // Check for missing images or descriptions
    const missingFields = steps.map((step, index) => {
      const issues = []
      if (!step.imageUrl) issues.push('image')
      if (!step.description.trim()) issues.push('description')
      return { stepNumber: index + 1, issues }
    }).filter(step => step.issues.length > 0)

    if (missingFields.length > 0) {
      const stepNumbers = missingFields.map(step => step.stepNumber).join(', ')
      toast.error(`Please add images and descriptions for steps: ${stepNumbers}`)
      return
    }

    setLoading(true)

    try {
      // Convert images to base64 and prepare for S3 upload
      const imageData = await Promise.all(steps.map(async (step) => {
        let imageData = ''
        if (step.image) {
          // Convert File to base64
          imageData = await new Promise((resolve, reject) => {
            const reader = new FileReader()
            reader.onload = () => {
              const result = reader.result as string
              // Remove data:image/jpeg;base64, prefix
              const base64 = result.split(',')[1]
              resolve(base64)
            }
            reader.onerror = reject
            reader.readAsDataURL(step.image!)
          })
        }
        
        return {
          image_data: imageData,
          image_url: step.imageUrl, // Fallback URL if no image data
          step_description: step.description,
          step_number: step.stepNumber
        }
      }))

      await axios.post(`${API_BASE_URL}/posts`, {
        title,
        description: description || null,
        images: imageData
      })

      toast.success('Recipe created successfully!')
      router.push('/feed')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to create recipe')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Layout>
      <div className="max-w-2xl mx-auto px-4 py-6">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-earth-800 mb-2">Create Recipe</h1>
          <p className="text-earth-600">Share your culinary masterpiece with the community</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Recipe Title */}
          <div className="form-group">
            <label htmlFor="title" className="form-label">
              Recipe Title *
            </label>
            <input
              id="title"
              type="text"
              required
              className="input-field"
              placeholder="What are you making?"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </div>

          {/* Recipe Description */}
          <div className="form-group">
            <label htmlFor="description" className="form-label">
              Description
            </label>
            <textarea
              id="description"
              className="input-field resize-none"
              rows={3}
              placeholder="Tell us about this recipe..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          {/* Recipe Steps */}
          <div className="form-group">
            <label className="form-label">Recipe Steps *</label>
            <div className="space-y-4">
              {steps.map((step, index) => (
                <div key={step.id} className="recipe-step">
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center space-x-2">
                      <span className="step-number">{step.stepNumber}</span>
                      <span className="font-medium text-earth-700">Step {step.stepNumber}</span>
                    </div>
                    {steps.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeStep(step.id)}
                        className="text-red-500 hover:text-red-700"
                      >
                        <FiX className="h-5 w-5" />
                      </button>
                    )}
                  </div>

                  {/* Image Upload */}
                  <div className="mb-3">
                    <label className="block text-sm font-medium text-earth-700 mb-2">
                      Step Photo *
                    </label>
                    {step.imageUrl ? (
                      <div className="relative">
                        <img
                          src={step.imageUrl}
                          alt={`Step ${step.stepNumber}`}
                          className="w-full h-48 object-cover rounded-lg"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            updateStep(step.id, 'image', null)
                            updateStep(step.id, 'imageUrl', '')
                          }}
                          className="absolute top-2 right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600"
                        >
                          <FiX className="h-4 w-4" />
                        </button>
                      </div>
                    ) : (
                      <div className="image-upload">
                        <input
                          type="file"
                          accept="image/*"
                          onChange={(e) => {
                            const file = e.target.files?.[0]
                            if (file) handleImageUpload(step.id, file)
                          }}
                          className="hidden"
                          id={`image-${step.id}`}
                        />
                        <label
                          htmlFor={`image-${step.id}`}
                          className="cursor-pointer flex flex-col items-center justify-center py-8"
                        >
                          <FiUpload className="h-8 w-8 text-earth-400 mb-2" />
                          <span className="text-earth-600">Upload step photo</span>
                        </label>
                      </div>
                    )}
                  </div>

                  {/* Step Description */}
                  <div>
                    <label className="block text-sm font-medium text-earth-700 mb-2">
                      Step Description *
                    </label>
                    <textarea
                      className="input-field resize-none"
                      rows={3}
                      placeholder="Describe what to do in this step..."
                      value={step.description}
                      onChange={(e) => updateStep(step.id, 'description', e.target.value)}
                    />
                  </div>
                </div>
              ))}

              {/* Add Step Button */}
              <button
                type="button"
                onClick={addStep}
                className="w-full border-2 border-dashed border-earth-300 rounded-lg p-4 text-earth-600 hover:border-earth-400 hover:text-earth-700 transition-colors duration-200"
              >
                <FiPlus className="h-5 w-5 mx-auto mb-2" />
                Add Another Step
              </button>
            </div>
          </div>

          {/* Submit Button */}
          <div className="flex space-x-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="btn-outline flex-1"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="btn-primary flex-1 flex items-center justify-center"
            >
              {loading ? (
                <>
                  <LoadingSpinner size="sm" />
                  <span className="ml-2">Creating...</span>
                </>
              ) : (
                'Create Recipe'
              )}
            </button>
          </div>
        </form>
      </div>
    </Layout>
  )
}
