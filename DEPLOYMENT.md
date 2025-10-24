# 🚀 Vercel Deployment Guide

This guide will help you deploy your full-stack Recipe Social Media application to Vercel.

## 📋 Prerequisites

- Vercel account (free tier available)
- GitHub repository with your code
- AWS account (for database and S3)
- Domain name (optional)

## 🏗️ Deployment Strategy

Since this is a full-stack application with:
- **Frontend**: Next.js (deploy to Vercel)
- **Backend**: Go (deploy to Vercel Functions or separate service)
- **Database**: PostgreSQL (AWS RDS)
- **Storage**: AWS S3

## 🎯 Step 1: Deploy Frontend to Vercel

### Option A: Deploy via Vercel CLI (Recommended)

```bash
# Install Vercel CLI
npm i -g vercel

# Login to Vercel
vercel login

# Deploy from project root
vercel

# Follow the prompts:
# - Set up and deploy? Yes
# - Which scope? (your account)
# - Link to existing project? No
# - Project name: recipe-social-media
# - Directory: ./
# - Override settings? No
```

### Option B: Deploy via Vercel Dashboard

1. Go to [vercel.com](https://vercel.com)
2. Click "New Project"
3. Import your GitHub repository
4. Configure:
   - **Framework Preset**: Next.js
   - **Root Directory**: `./`
   - **Build Command**: `npm run build`
   - **Output Directory**: `.next`

## 🔧 Step 2: Configure Environment Variables

In your Vercel dashboard, go to **Settings > Environment Variables** and add:

```env
# Frontend Environment Variables
NEXT_PUBLIC_API_URL=https://your-backend-url.vercel.app/api
```

## 🐹 Step 3: Deploy Go Backend

### Option A: Vercel Functions (Recommended)

Create `api/` directory in your project root:

```bash
mkdir api
```

Create `api/hello.go`:
```go
package main

import (
    "net/http"
    "github.com/vercel/vercel-go"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello from Vercel!"))
}

func main() {
    vercel.Start(Handler)
}
```

### Option B: Separate Vercel Project for Backend

1. Create a new Vercel project for your backend
2. Deploy the `backend/` directory as a separate project
3. Use the backend URL in your frontend environment variables

### Option C: Alternative Backend Hosting

Consider these alternatives for your Go backend:
- **Railway**: Great for Go applications
- **Render**: Free tier available
- **DigitalOcean App Platform**: Good for full-stack apps
- **AWS Lambda**: Serverless Go functions

## 🗄️ Step 4: Database Configuration

### Using AWS RDS (Current Setup)

Your current setup uses AWS RDS PostgreSQL. For production:

1. **Create Production Database**:
   ```bash
   # Use AWS Console or CLI to create a new RDS instance
   # Update your backend environment variables
   ```

2. **Run Database Migrations**:
   ```bash
   # Connect to your production database
   psql -h your-prod-db.amazonaws.com -U username -d database -f backend/database/schema.sql
   ```

### Alternative: Vercel Postgres

Consider using Vercel Postgres for simpler setup:

```bash
# Install Vercel Postgres
npm install @vercel/postgres

# Add to your Vercel project
vercel env add POSTGRES_URL
```

## ☁️ Step 5: AWS S3 Configuration

Update your S3 bucket for production:

1. **Create Production S3 Bucket**:
   ```bash
   # Use AWS Console to create a new bucket
   # Enable CORS for your domain
   ```

2. **Update CORS Policy**:
   ```json
   {
     "CORSRules": [
       {
         "AllowedOrigins": ["https://your-domain.vercel.app"],
         "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
         "AllowedHeaders": ["*"]
       }
     ]
   }
   ```

## 🔐 Step 6: Production Environment Variables

### Frontend (Vercel Dashboard)
```env
NEXT_PUBLIC_API_URL=https://your-backend.vercel.app/api
```

### Backend (Separate Vercel Project or Alternative Hosting)
```env
# Database
DB_HOST=your-production-db.amazonaws.com
DB_PORT=5432
DB_USER=your-username
DB_PASSWORD=your-secure-password
DB_NAME=your-database
DB_SSLMODE=require

# JWT
JWT_SECRET=your-super-secure-jwt-secret

# AWS S3
S3_BUCKET=your-production-bucket
S3_REGION=us-east-2
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

## 🚀 Step 7: Deploy Commands

### Deploy Frontend
```bash
# From project root
vercel --prod
```

### Deploy Backend (if using Vercel)
```bash
# From backend directory
cd backend
vercel --prod
```

## 🔍 Step 8: Testing Your Deployment

1. **Test Frontend**: Visit your Vercel URL
2. **Test API**: Check `https://your-domain.vercel.app/api/health`
3. **Test Database**: Create a user account
4. **Test S3**: Upload a profile photo

## 📊 Step 9: Monitoring and Analytics

### Vercel Analytics
- Enable Vercel Analytics in your dashboard
- Monitor performance and errors

### Database Monitoring
- Use AWS CloudWatch for RDS monitoring
- Set up alerts for database issues

## 🛠️ Troubleshooting

### Common Issues:

1. **CORS Errors**:
   ```javascript
   // Add to your Go backend
   r.Use(func(c *gin.Context) {
     c.Header("Access-Control-Allow-Origin", "https://your-domain.vercel.app")
     // ... other CORS headers
   })
   ```

2. **Environment Variables Not Loading**:
   - Check Vercel dashboard settings
   - Redeploy after adding variables

3. **Database Connection Issues**:
   - Verify database credentials
   - Check security groups and VPC settings

4. **S3 Upload Issues**:
   - Verify AWS credentials
   - Check bucket permissions

## 🎯 Production Checklist

- [ ] Frontend deployed to Vercel
- [ ] Backend deployed (Vercel/Railway/Render)
- [ ] Database configured and migrated
- [ ] S3 bucket configured with CORS
- [ ] Environment variables set
- [ ] Domain configured (optional)
- [ ] SSL certificates working
- [ ] Performance monitoring enabled

## 🔗 Useful Links

- [Vercel Documentation](https://vercel.com/docs)
- [Next.js Deployment](https://nextjs.org/docs/deployment)
- [Go on Vercel](https://vercel.com/docs/functions/serverless-functions/runtimes/go)
- [AWS RDS Setup](https://docs.aws.amazon.com/rds/)
- [S3 CORS Configuration](https://docs.aws.amazon.com/s3/latest/userguide/cors.html)

## 💡 Pro Tips

1. **Use Vercel's Preview Deployments** for testing
2. **Set up staging environment** before production
3. **Monitor your database connections** to avoid timeouts
4. **Use CDN** for static assets (Vercel handles this automatically)
5. **Enable Vercel Analytics** for performance insights

---

**Need Help?** Check the troubleshooting section or create an issue in your repository!
