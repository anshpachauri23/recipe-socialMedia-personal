# 🍳 Recipe Social Media Platform

A full-stack social media application for sharing and discovering recipes, built with Next.js, Go, and PostgreSQL.

## ✨ Features

- **User Authentication**: Secure login/register with JWT tokens
- **Recipe Sharing**: Create posts with multiple images and step-by-step instructions
- **Social Features**: Like posts, comment, follow users
- **Profile Management**: Upload profile photos, edit bio, view user posts
- **Real-time Feed**: Personalized feed with your posts and followed users
- **Image Storage**: AWS S3 integration for scalable image storage
- **Responsive Design**: Mobile-first design with Tailwind CSS

## 🛠️ Tech Stack

### Frontend
- **Next.js 14** - React framework with App Router
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework
- **React Icons** - Beautiful icons
- **Axios** - HTTP client
- **js-cookie** - Cookie management

### Backend
- **Go** - High-performance backend language
- **Gin** - HTTP web framework
- **PostgreSQL** - Relational database
- **AWS S3** - Object storage for images
- **JWT** - JSON Web Tokens for authentication

### Database
- **PostgreSQL** - Primary database
- **AWS RDS** - Managed PostgreSQL service

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ and npm
- Go 1.21+
- PostgreSQL database
- AWS account (for S3)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd recipe-social-media
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

## 📁 Project Structure

```
recipe-social-media/
├── app/                    # Next.js app directory
│   ├── auth/              # Authentication pages
│   ├── create/            # Create post page
│   ├── feed/              # Main feed page
│   ├── profile/           # User profiles
│   ├── posts/             # Individual post pages
│   └── settings/          # User settings
├── backend/               # Go backend
│   ├── database/          # Database models and queries
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # Custom middleware
│   ├── models/            # Data models
│   └── services/          # Business logic (S3, etc.)
├── components/            # Reusable React components
├── contexts/              # React context providers
└── public/                # Static assets
```

## 🔧 Environment Variables

### Frontend (.env)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### Backend (backend/local.env)
```env
# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=your-database
DB_SSLMODE=require

# JWT
JWT_SECRET=your-jwt-secret

# AWS S3
S3_BUCKET=your-bucket-name
S3_REGION=your-region
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

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
- Comment system
- Follow/unfollow users
- Personalized feeds

### Database Design
- Users, Posts, Likes, Comments tables
- Triggers for automatic count updates
- Optimized queries with proper indexing

## 🚀 Deployment

### Local Development
- Backend runs on port 8080
- Frontend runs on port 3000
- Database: PostgreSQL (local or AWS RDS)

### Production Considerations
- Use environment variables for all secrets
- Set up proper CORS policies
- Configure SSL/TLS
- Use a reverse proxy (nginx)
- Set up monitoring and logging

## 📱 Screenshots

*Add screenshots of your application here*

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
- LinkedIn: [Ansh Pachauri](htps://www.linkedin.com/in/ansh-pachauri)

## 🙏 Acknowledgments

- Next.js team for the amazing framework
- Go team for the excellent language
- PostgreSQL community for the robust database
- AWS for cloud services