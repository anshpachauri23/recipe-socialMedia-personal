# 🚀 Lambda Deployment Guide

## 📦 **Files to Upload to AWS Lambda:**

### **1. Authentication Functions:**
- **Login**: `auth/login/function.zip` → Upload to `recipe-login` Lambda
- **Register**: `auth/register/function.zip` → Upload to `recipe-register` Lambda

### **2. Post Functions:**
- **Create Post**: `create-post/` → Build and upload to `recipe-create-post` Lambda
- **Get Posts**: `get-posts/get-posts-fixed.zip` → Upload to `recipe-get-posts` Lambda
- **Feed**: `posts/feed/feed-fixed.zip` → Upload to `recipe-feed` Lambda

### **3. Diagnostic Function:**
- **Diagnose**: `diagnose-db/diagnose-fixed.zip` → Upload to `diagnose-db` Lambda

## 🔧 **Build Commands (if needed):**

```bash
# For create-post function
cd create-post
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip create-post.zip bootstrap
```

## ⚙️ **Lambda Configuration:**

### **Handler Settings:**
- Set handler to: `bootstrap`
- Runtime: Go 1.x

### **Environment Variables:**
```
DB_HOST=your-rds-endpoint.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=recipe_social
DB_SSLMODE=require
JWT_SECRET=85817ddb71e4d1f9dfc9b69ec7917e01b7776250a3b996194b874401659bdc97
```

## 🎯 **Current Status:**
- ✅ **Get Posts**: Working perfectly
- ✅ **Authentication**: Working perfectly  
- ✅ **Create Post**: Working perfectly
- ❌ **Feed**: Needs JWT secret verification
- ❌ **Diagnostic**: Needs upload of fixed version

## 📝 **Notes:**
- All functions use direct PostgreSQL connections
- NULL values are properly handled in fixed versions
- CORS is configured for all endpoints
- JWT tokens are validated correctly
