# 🍳 Recipe Social Media Platform

A full-stack social media application for sharing and discovering recipes, built with Next.js, Go, and PostgreSQL. This platform allows users to create, share, and discover delicious recipes with a social media experience.

## ✨ Features

- **User Authentication**: Secure login/register with JWT tokens
- **Recipe Sharing**: Create posts with multiple images and step-by-step instructions
- **Social Features**: Like posts, comment, follow users, share recipes
- **Profile Management**: Upload profile photos, edit bio, view user posts
- **Smart Feed**: Personalized feed for followed users, all recipes for new users
- **Search Functionality**: Search for recipes, users, and ingredients
- **Image Storage**: AWS S3 integration for scalable image storage
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Step-by-Step Recipes**: Detailed cooking instructions with images
- **Share Functionality**: Copy post links or use native share on mobile
- **Avatar Generation**: Automatic SVG avatars with user initials

## 🛠️ Tech Stack

### Frontend
- **Next.js 14.2.33** - React framework with App Router
- **TypeScript 5.2.0** - Type-safe JavaScript
- **React 18.2.0** - UI library
- **Tailwind CSS 3.3.0** - Utility-first CSS framework
- **React Icons 4.12.0** - Beautiful icons
- **Axios 1.6.0** - HTTP client
- **js-cookie 3.0.5** - Cookie management
- **React Hook Form 7.47.0** - Form handling
- **React Hot Toast 2.4.1** - Notifications

### Backend
- **Go 1.21** - High-performance backend language
- **Gin 1.9.1** - HTTP web framework
- **PostgreSQL** - Relational database
- **AWS S3** - Object storage for images
- **JWT** - JSON Web Tokens for authentication
- **AWS SDK Go 1.50.0** - AWS services integration

### Database & Storage
- **PostgreSQL** - Primary database (managed by Northflank)
- **AWS S3** - Image storage and CDN

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL database
- AWS account (for S3)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/anshpachauri23/recipe-socialMedia-personal.git
   cd recipe-socialMedia-personal
   ```

2. **Install frontend dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   ```bash
   # Copy environment template
   cp env.example .env
   
   # Edit .env with your values
   NEXT_PUBLIC_API_URL=http://localhost:8080/api
   ```

4. **Set up backend environment**
   ```bash
   cd backend
   cp local.env.example local.env
   # Edit local.env with your database and AWS credentials
   ```

5. **Install Go dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

6. **Set up database**
   ```bash
   # Run the SQL schema
   psql -h your-db-host -U your-username -d your-database -f database/schema.sql
   ```

7. **Start the application**
   ```bash
   # Terminal 1: Start backend
   cd backend
   go run main.go
   
   # Terminal 2: Start frontend
   npm run dev
   ```

8. **Visit the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## 🌐 Live Deployment

The application is currently deployed and accessible at:

- **Frontend (Vercel)**: https://your-app.vercel.app
- **Backend (Northflank)**: https://your-backend.northflank.app
- **Database**: Northflank PostgreSQL addon (private network to backend)
- **Image Storage**: AWS S3 (us-east-1)

### Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Users → Vercel (Frontend)                                  │
│            ↓                                                │
│         Northflank (Go Backend) ◀──── Northflank Postgres   │
│            ↓                          (private network)     │
│         AWS S3 (Images, us-east-1)                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Deployment Services

- **Vercel**: Hosts the Next.js frontend with automatic deployments from GitHub
- **Northflank**: Runs the Go backend in a Docker container alongside a managed PostgreSQL addon (linked over Northflank's internal private network — no public DB exposure, no egress costs)
- **AWS S3**: Scalable object storage for recipe and profile images, with a least-privilege IAM user for the backend and a public-read policy scoped to the `recipe-images/*` prefix

## 📁 Project Structure

```
recipe-social-media/
├── app/                    # Next.js app directory
│   ├── auth/              # Authentication pages
│   │   ├── login/         # Login page
│   │   └── register/      # Registration page
│   ├── create/            # Create post page
│   ├── feed/              # Main feed page
│   ├── profile/           # User profiles
│   │   └── [id]/          # Dynamic user profile pages
│   ├── posts/             # Individual post pages
│   │   └── [id]/          # Dynamic post detail pages
│   ├── search/            # Search functionality
│   ├── settings/          # User settings
│   ├── globals.css        # Global styles
│   └── layout.tsx         # Root layout
├── backend/               # Go backend
│   ├── database/          # Database models and queries
│   │   ├── connection.go  # Database connection
│   │   ├── db.go         # Database operations
│   │   ├── posts.go      # Post-related queries
│   │   └── schema.sql    # Database schema
│   ├── handlers/          # HTTP request handlers
│   │   ├── auth.go       # Authentication handlers
│   │   ├── post.go       # Post handlers
│   │   └── user.go       # User handlers
│   ├── middleware/        # Custom middleware
│   │   └── auth.go       # Authentication middleware
│   ├── models/            # Data models
│   │   ├── post.go       # Post model
│   │   └── user.go       # User model
│   ├── services/          # Business logic
│   │   └── s3.go         # AWS S3 integration
│   ├── main.go           # Application entry point
│   ├── go.mod            # Go dependencies
│   └── local.env         # Environment variables
├── components/            # Reusable React components
│   ├── Layout.tsx        # Main layout component
│   ├── LoadingSpinner.tsx # Loading indicator
│   ├── PostCard.tsx      # Post display component
│   └── SearchBar.tsx     # Search functionality
├── contexts/              # React context providers
│   └── AuthContext.tsx   # Authentication context
├── public/                # Static assets
├── package.json          # Node.js dependencies
├── next.config.js        # Next.js configuration
├── tailwind.config.js    # Tailwind CSS configuration
├── tsconfig.json         # TypeScript configuration
├── Dockerfile            # Docker configuration for backend
├── .dockerignore         # Docker ignore file
└── northflank.yaml       # Northflank deployment config
```

## 🔧 Environment Variables

### Frontend (.env)
```env
# Local Development
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Production (Vercel)
NEXT_PUBLIC_API_URL=https://your-backend.northflank.app/api
```

### Backend (backend/local.env)
```env
# Server
PORT=8080

# Database (Northflank Postgres — use the addon's "public" connection string for local dev)
DB_HOST=your-addon-host.northflank.app
DB_PORT=5432
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=require

# JWT Authentication
JWT_SECRET=your-secure-jwt-secret-key

# AWS S3 Configuration
S3_BUCKET=your-s3-bucket-name
S3_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

**Note**: In production, the backend on Northflank reads `DB_*` from the Postgres addon's linked secrets (mapped from `POSTGRES_HOST`, `POSTGRES_USERNAME`, etc. → `DB_HOST`, `DB_USER`, etc.) and uses the addon's **internal** host. S3 credentials are set as plain env vars / secrets on the service.

## 🎯 Key Features Implementation

### Authentication
- JWT-based authentication
- Secure password hashing
- Protected routes with middleware

### Image Upload
- Base64 to S3 conversion
- Multiple images per post
- Profile photo uploads

### Social Features
- Like/unlike posts
- Real-time comment system with user avatars
- Follow/unfollow users
- Personalized feeds based on follows
- Share posts via native share API or clipboard
- User profile pages with post history

### Database Design
- Users, Posts, Likes, Comments tables
- Triggers for automatic count updates
- Optimized queries with proper indexing

## 🚀 Deployment

### Local Development
- Backend runs on port 8080
- Frontend runs on port 3000
- Database: PostgreSQL (local Docker container or the Northflank addon's public connection string)

### Available Scripts
```bash
# Frontend development
npm run dev          # Start Next.js development server
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint

# Backend development
npm run backend      # Start Go backend server
npm run backend:build # Build Go binary

# Setup
npm run setup        # Install all dependencies
```

### Production Deployment

#### Frontend (Vercel)
1. Connect GitHub repository to Vercel
2. Configure environment variables:
   - `NEXT_PUBLIC_API_URL`: Backend API URL
3. Deploy automatically on push to `main` branch
4. Vercel handles:
   - Automatic builds
   - CDN distribution
   - SSL certificates
   - Preview deployments for pull requests

#### Backend (Northflank)
1. Create a new service from GitHub repository
2. Use `Dockerfile` in backend directory
3. Configure environment variables (see above)
4. Deploy automatically on push to `main` branch
5. Northflank handles:
   - Docker container builds
   - Health checks
   - Automatic restarts
   - Zero downtime deployments

#### Database (Northflank Postgres addon)
- Created as a managed addon in the same Northflank project as the backend service
- Linked to the backend via "Linked secrets", remapping `POSTGRES_*` vars to the `DB_*` names the Go code expects
- Backend uses the addon's **internal** host (private network) — public host is only used for one-off `psql` access
- SSL required (`DB_SSLMODE=require`)
- Schema initialized once by running [backend/database/schema.sql](backend/database/schema.sql) in the Northflank Postgres console

#### Storage (AWS S3)
- Bucket in `us-east-1` with "Block all public access" disabled
- Bucket policy grants public `s3:GetObject` only on the `recipe-images/*` prefix
- IAM user `recipe-social-backend-s3` with an inline least-privilege policy (`s3:PutObject`, `s3:GetObject` on `recipe-images/*` only)
- CORS configured to allow GET from any origin

### Production Considerations
- ✅ All secrets managed via environment variables
- ✅ CORS policies configured for Vercel domain
- ✅ SSL/TLS enabled on all services
- ✅ Database connections use SSL
- ✅ Automatic deployments from GitHub
- ✅ Health check endpoints configured

## 📱 Screenshots

*Add screenshots of your application here*

## 🔄 Current Status

### ✅ Completed Features
- User authentication system (login/register)
- Recipe creation with multiple images
- Social features (likes, comments, follows, share)
- User profiles and settings
- Search functionality for recipes and users
- Responsive mobile-first design
- AWS S3 integration for image storage
- Smart feed (personalized or all recipes for new users)
- Avatar generation with user initials
- Share functionality (native share API + clipboard)
- Null-safe error handling throughout the app

### 🚧 Development Status
- **Frontend**: ✅ Deployed on Vercel (Production)
- **Backend**: ✅ Deployed on Northflank (Production)
- **Database**: ✅ Northflank Postgres addon (Production)
- **Storage**: ✅ AWS S3 — `us-east-1` (Production)
- **CI/CD**: ✅ Automatic deployments from GitHub

### 🎯 Recent Updates
- ✅ Migrated database from AWS RDS to the Northflank Postgres addon (linked over the internal private network for free, fast intra-cluster traffic)
- ✅ Migrated S3 to a new AWS account in `us-east-1` with a least-privilege IAM user and prefix-scoped public-read bucket policy
- ✅ Deployed to production (Vercel + Northflank)
- ✅ Fixed all null-safety issues in forms and comments
- ✅ Replaced external placeholder images with inline SVG avatars
- ✅ Added share functionality for posts
- ✅ Improved comment system with user avatars
- ✅ Enhanced new user experience (show all recipes)
- ✅ Fixed profile page rendering issues
- ✅ Updated error handling across all pages
- ✅ Migrated from Render to Northflank for better performance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Ansh Pachauri**
- GitHub: [@anshpachauri](https://github.com/anshpachauri)
- LinkedIn: [MyLinkedIn](https://linkedin.com/in/ansh-pachauri)

## 🙏 Acknowledgments

- Next.js team for the amazing framework
- Go team for the excellent language
- PostgreSQL community for the robust database
- AWS for cloud services