# RecipeShare - Social Media for Home Cooks

A full-stack social media platform for sharing homemade recipes, built with Next.js, Go, and PostgreSQL, hosted on AWS.

## Features

- 🔐 **Authentication**: Secure user registration and login
- 📱 **Social Feed**: View recipes from users you follow and popular accounts
- 🍳 **Recipe Creation**: Step-by-step recipe posts with multiple photos
- 👥 **User Profiles**: Public profiles with follower/following system
- 🔍 **Search**: Find recipes and users by keywords
- ❤️ **Interactions**: Like, comment, and share recipes
- ⚙️ **Settings**: Update profile, change password, manage account

## Tech Stack

### Frontend
- **Next.js 14** with App Router
- **TypeScript** for type safety
- **Tailwind CSS** with custom earthy color palette
- **React Icons** for beautiful icons
- **Axios** for API calls
- **React Hook Form** for form handling

### Backend
- **Go** with Gin framework
- **PostgreSQL** database
- **JWT** authentication
- **AWS Lambda** functions for serverless operations
- **AWS RDS** for database hosting

### Database
- **PostgreSQL** with optimized schema
- **Triggers** for automatic count updates
- **Indexes** for performance optimization

## Project Structure

```
recipe-socialMedia-personal/
├── app/                    # Next.js app directory
│   ├── auth/              # Authentication pages
│   ├── feed/              # Main feed page
│   ├── create/            # Recipe creation page
│   ├── search/            # Search functionality
│   ├── profile/           # User profiles
│   └── settings/          # Account settings
├── components/            # Reusable React components
├── contexts/              # React contexts (Auth)
├── backend/               # Go backend
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Data models
│   ├── database/          # Database operations
│   └── middleware/        # Middleware functions
├── lambda/                # AWS Lambda functions
│   └── functions/         # Individual Lambda functions
└── database/             # Database schema and migrations
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL 13+
- AWS Account (for production deployment)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd recipe-socialMedia-personal
   ```

2. **Install frontend dependencies**
   ```bash
   npm install
   ```

3. **Install backend dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

4. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb recipe_social
   
   # Run schema
   psql recipe_social < backend/database/schema.sql
   ```

5. **Configure environment variables**
   ```bash
   # Copy example file
   cp env.example .env
   
   # Edit with your database credentials
   ```

6. **Start the development servers**
   ```bash
   # Terminal 1: Start backend
   cd backend
   go run main.go
   
   # Terminal 2: Start frontend
   npm run dev
   ```

7. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update profile
- `DELETE /api/users/me` - Delete account
- `GET /api/users/:id` - Get user profile
- `POST /api/users/:id/follow` - Follow user
- `DELETE /api/users/:id/follow` - Unfollow user
- `GET /api/users/search` - Search users

### Posts
- `GET /api/posts/feed` - Get user feed
- `POST /api/posts` - Create post
- `GET /api/posts/:id` - Get post details
- `PUT /api/posts/:id` - Update post
- `DELETE /api/posts/:id` - Delete post
- `POST /api/posts/:id/like` - Like post
- `DELETE /api/posts/:id/like` - Unlike post
- `GET /api/posts/search` - Search posts

### Comments
- `POST /api/posts/:id/comments` - Create comment
- `GET /api/posts/:id/comments` - Get comments
- `DELETE /api/comments/:id` - Delete comment

## AWS Deployment

### Database Setup
1. Create RDS PostgreSQL cluster
2. Configure security groups
3. Set up connection pooling

### Lambda Functions
1. Deploy Lambda functions for data operations
2. Configure API Gateway
3. Set up IAM roles and permissions

### S3 Storage
1. Create S3 bucket for image storage
2. Configure CORS and permissions
3. Set up CloudFront for CDN

## Design System

The application uses a custom earthy color palette inspired by natural ingredients:

- **Earth Tones**: Warm browns and tans
- **Sage Greens**: Fresh herb colors
- **Terracotta**: Warm orange-reds
- **Instagram-like Layout**: Clean, modern interface

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, email support@recipeshare.com or create an issue in the repository.