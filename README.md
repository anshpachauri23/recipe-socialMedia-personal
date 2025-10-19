# Recipe Social Media Platform

A full-stack social media application for sharing recipes, built with Next.js (React) for the frontend and Go (Gin framework) for the backend, connected to an AWS RDS PostgreSQL database and AWS S3 for image storage.

## ✨ Features

- **User Authentication**: Secure registration and login with JWT tokens
- **Recipe Management**: Create, view, update, and delete recipe posts
- **Image Uploads**: Upload images to AWS S3 with automatic optimization
- **Social Features**: Like and comment on posts, follow other users
- **User Profiles**: Customizable profiles with photo uploads
- **Search Functionality**: Search for recipes and users
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Real-time Updates**: Dynamic like counts and comment systems

## 🛠️ Tech Stack

### Frontend
- **Next.js 14** - React framework with App Router
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework
- **React Icons** - Icon library
- **Axios** - HTTP client
- **js-cookie** - Cookie management
- **React Hot Toast** - Notifications

### Backend
- **Go** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Database (AWS RDS)
- **AWS S3** - Image storage
- **JWT** - Authentication

### Infrastructure
- **AWS RDS** - Managed PostgreSQL database
- **AWS S3** - Object storage for images
- **AWS IAM** - Identity and access management

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ and npm
- Go 1.19+
- AWS Account with RDS and S3 access

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/recipe-social-media.git
cd recipe-social-media
```

### 2. Frontend Setup
```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Edit .env with your API URL
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Start development server
npm run dev
```

### 3. Backend Setup
```bash
# Navigate to backend directory
cd backend

# Copy environment file
cp local.env.example local.env

# Edit local.env with your AWS credentials
# Fill in your database and S3 configuration

# Install Go dependencies
go mod tidy

# Start the backend server
go run main.go
```

### 4. Database Setup
Run the SQL schema from `backend/database/schema.sql` on your PostgreSQL database.

## 📁 Project Structure

```
recipe-social-media/
├── app/                    # Next.js app directory
│   ├── auth/              # Authentication pages
│   ├── create/            # Create post page
│   ├── feed/              # Main feed
│   ├── profile/           # User profiles
│   └── search/            # Search functionality
├── components/            # Reusable React components
├── contexts/              # React context providers
├── backend/               # Go backend
│   ├── handlers/          # HTTP handlers
│   ├── database/          # Database operations
│   ├── models/            # Data models
│   ├── middleware/        # Custom middleware
│   └── services/          # Business logic
└── public/               # Static assets
```

## 🔐 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: Server-side validation for all inputs
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Proper cross-origin resource sharing
- **Environment Variables**: Sensitive data stored securely
- **User Authorization**: Users can only edit their own content

## 🌐 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update profile
- `GET /api/users/:id` - Get user profile
- `POST /api/users/me/profile-photo` - Upload profile photo

### Posts
- `GET /api/posts/feed` - Get user feed
- `POST /api/posts` - Create new post
- `GET /api/posts/:id` - Get specific post
- `DELETE /api/posts/:id` - Delete post
- `POST /api/posts/:id/like` - Like post
- `DELETE /api/posts/:id/like` - Unlike post

### Comments
- `POST /api/posts/:id/comments` - Add comment
- `GET /api/posts/:id/comments` - Get comments
- `DELETE /api/comments/:id` - Delete comment

## 🚀 Deployment

### Frontend (Vercel)
1. Connect your GitHub repository to Vercel
2. Set environment variables in Vercel dashboard
3. Deploy automatically on push to main branch

### Backend (AWS EC2)
1. Launch EC2 instance with Go runtime
2. Clone repository and build the application
3. Set up environment variables
4. Configure reverse proxy (nginx)
5. Set up SSL certificates

### Database (AWS RDS)
1. Create PostgreSQL RDS instance
2. Configure security groups
3. Run database migrations
4. Set up automated backups

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Ansh Pachauri**
- GitHub: [@anshpachauri](https://github.com/anshpachauri)
- LinkedIn: [Ansh Pachauri](https://linkedin.com/in/anshpachauri)

## 🙏 Acknowledgments

- Next.js team for the amazing framework
- Go team for the excellent language
- AWS for cloud infrastructure
- All open-source contributors