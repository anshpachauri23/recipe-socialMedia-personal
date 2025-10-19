# Local Development Setup

This guide explains how to run the RecipeShare application locally with your AWS database.

## Architecture

- **Frontend**: Next.js running on `localhost:3000`
- **Backend**: Go server running on `localhost:8080`
- **Database**: AWS RDS PostgreSQL (remote)
- **Image Storage**: AWS S3 (remote)

## Prerequisites

- Node.js 18+ and npm
- Go 1.21+
- AWS RDS PostgreSQL instance (already configured)
- AWS S3 bucket for image storage (already configured)
- AWS credentials configured locally

## Quick Start

### Option 1: Start Everything at Once
```bash
./start-dev.sh
```

### Option 2: Start Backend and Frontend Separately

**Terminal 1 - Backend:**
```bash
./start-backend.sh
```

**Terminal 2 - Frontend:**
```bash
./start-frontend.sh
```

## Manual Setup

### 1. Backend Setup

```bash
cd backend
go mod tidy
go run main.go
```

The backend will:
- Load environment variables from `backend/local.env`
- Connect to your AWS RDS database
- Start server on `http://localhost:8080`

### 2. Frontend Setup

```bash
npm install
export NEXT_PUBLIC_API_URL=http://localhost:8080/api
npm run dev
```

The frontend will start on `http://localhost:3000`

## Environment Configuration

### Backend Environment (`backend/local.env`)
```env
PORT=8080
DB_HOST=recipe-db.cbcy4q2c4epi.us-east-2.rds.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Chahek.231104
DB_NAME=postgres
DB_SSLMODE=require
JWT_SECRET=85817ddb71e4d1f9dfc9b69ec7917e01b7776250a3b996194b874401659bdc97
S3_BUCKET=recipe-social-images
S3_REGION=us-east-2
```

### Frontend Environment
The frontend automatically uses `http://localhost:8080/api` as the API URL when running locally.

## API Endpoints

All API endpoints are available at `http://localhost:8080/api`:

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Users
- `GET /api/users/me` - Get current user
- `PUT /api/users/me` - Update profile
- `GET /api/users/:id` - Get user profile
- `POST /api/users/:id/follow` - Follow user
- `DELETE /api/users/:id/follow` - Unfollow user
- `GET /api/users/search` - Search users

### Posts
- `GET /api/posts/feed` - Get user feed
- `POST /api/posts` - Create post
- `GET /api/posts/:id` - Get post details
- `POST /api/posts/:id/like` - Like post
- `DELETE /api/posts/:id/like` - Unlike post
- `GET /api/posts/search` - Search posts

## S3 Image Storage Setup

Your local backend now supports S3 image uploads! See [S3_LOCAL_SETUP.md](./S3_LOCAL_SETUP.md) for detailed setup instructions.

### Quick S3 Setup:
1. **Configure AWS credentials** (one-time setup)
2. **Your S3 bucket is already configured** from Lambda setup
3. **Images upload automatically** when creating posts

## Troubleshooting

### Backend Issues

1. **Database Connection Failed**
   - Check your AWS RDS instance is running
   - Verify credentials in `backend/local.env`
   - Ensure your IP is whitelisted in RDS security groups

2. **S3 Upload Failed**
   - Check AWS credentials are configured
   - Verify S3 bucket permissions
   - See [S3_LOCAL_SETUP.md](./S3_LOCAL_SETUP.md) for details

3. **Port 8080 Already in Use**
   ```bash
   lsof -ti:8080 | xargs kill -9
   ```

4. **Go Dependencies**
   ```bash
   cd backend
   go mod tidy
   ```

### Frontend Issues

1. **API Connection Failed**
   - Ensure backend is running on port 8080
   - Check `NEXT_PUBLIC_API_URL` environment variable

2. **Node Dependencies**
   ```bash
   npm install
   ```

## Development Workflow

1. **Start Development Environment**
   ```bash
   ./start-dev.sh
   ```

2. **Access Application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api

3. **Make Changes**
   - Frontend changes auto-reload
   - Backend changes require restart

4. **Stop Development**
   - Press `Ctrl+C` in terminal running `start-dev.sh`

## Benefits of Local Development

✅ **Faster Development**: No Lambda cold starts  
✅ **Better Debugging**: Full access to logs and debugging tools  
✅ **Easier Testing**: Direct API testing with tools like Postman  
✅ **Cost Effective**: No Lambda execution costs during development  
✅ **Familiar Environment**: Standard web development workflow  

## Migration from AWS Lambda

This setup replaces the AWS Lambda functions with a local Go server while keeping:
- AWS RDS PostgreSQL database
- All existing data and schema
- Same API endpoints and functionality
- JWT authentication
- CORS configuration

The only change is the deployment method - from serverless Lambda to local server.
