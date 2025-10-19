# 🚀 Lambda Function Deployment Files

## ✅ **All Lambda Functions Built and Ready for Deployment**

### **📦 Zip Files Created (x86_64 Architecture)**

| Function | Zip File Location | Size | Purpose |
|----------|------------------|------|---------|
| **Login** | `lambda/functions/auth/login/function-x86.zip` | 6.8MB | User authentication |
| **Register** | `lambda/functions/auth/register/function-x86.zip` | 6.8MB | User registration |
| **Feed** | `lambda/functions/posts/feed/function-x86.zip` | 6.8MB | Get user feed |
| **Create Post** | `lambda/functions/create-post/function-x86.zip` | 6.7MB | Create recipe posts |
| **Get Posts** | `lambda/functions/get-posts/function-x86.zip` | 6.7MB | Get all posts |

## 🔧 **Deployment Instructions**

### **1. Create Lambda Functions in AWS Console**

Create these Lambda functions with these settings:

| Function Name | Runtime | Architecture | Handler |
|---------------|---------|--------------|---------|
| `recipe-login` | Go 1.x | x86_64 | main |
| `recipe-register` | Go 1.x | x86_64 | main |
| `recipe-feed` | Go 1.x | x86_64 | main |
| `recipe-create-post` | Go 1.x | x86_64 | main |
| `recipe-get-posts` | Go 1.x | x86_64 | main |

### **2. Upload Zip Files**

For each Lambda function:
1. Go to the function in AWS Console
2. Click "Code" tab
3. Click "Upload from" → ".zip file"
4. Upload the corresponding `function-x86.zip` file

### **3. Set Environment Variables**

Add these environment variables to **ALL** Lambda functions:

```bash
DB_HOST=recipe-db.cbcy4q2c4epi.us-east-2.rds.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Chahek.231104
DB_NAME=postgres
DB_SSLMODE=require
JWT_SECRET=85817ddb71e4d1f9dfc9b69ec7917e01b7776250a3b996194b874401659bdc97
```

### **4. Configure VPC Settings**

For **ALL** Lambda functions:
1. Go to "Configuration" → "VPC"
2. Select the same VPC as your RDS database
3. Select subnets (at least 2)
4. Select the security group that allows access to RDS

### **5. Set Up API Gateway**

Create these API Gateway endpoints:

| Method | Path | Lambda Function |
|--------|------|-----------------|
| POST | `/auth/login` | recipe-login |
| POST | `/auth/register` | recipe-register |
| GET | `/posts/feed` | recipe-feed |
| POST | `/posts` | recipe-create-post |
| GET | `/posts` | recipe-get-posts |

## 🧪 **Testing**

After deployment, test each endpoint:

### **Test Registration:**
```bash
curl -X POST https://760go4r862.execute-api.us-east-2.amazonaws.com/prod/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### **Test Login:**
```bash
curl -X POST https://760go4r862.execute-api.us-east-2.amazonaws.com/prod/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### **Test Get Posts:**
```bash
curl -X GET https://760go4r862.execute-api.us-east-2.amazonaws.com/prod/posts
```

## ✅ **All Files Ready!**

All Lambda functions are built with the correct x86_64 architecture and ready for deployment to AWS! 🚀
