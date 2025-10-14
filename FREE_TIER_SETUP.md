# 💰 AWS Free Tier Setup for Recipe Social Media

This guide shows you how to run your Recipe Social Media platform **completely FREE** using AWS Free Tier.

## 🆓 **What's FREE in AWS Free Tier**

| Service | Free Tier Limit | Your Usage |
|---------|----------------|------------|
| **Lambda** | 1M requests/month | ~10K requests/month ✅ |
| **API Gateway** | 1M API calls/month | ~10K calls/month ✅ |
| **RDS PostgreSQL** | 750 hours + 20GB | 24/7 for 1 month ✅ |
| **Secrets Manager** | 10K API calls/month | ~1K calls/month ✅ |
| **CloudWatch** | 10GB logs/month | ~1GB/month ✅ |

## 🎯 **Free Tier Optimized Architecture**

```
Frontend (Vercel/Netlify FREE) → API Gateway (FREE) → Lambda (FREE) → RDS (FREE)
```

## 📋 **Step 1: Set Up RDS (FREE for 1 Month)**

### 1.1 Create RDS Instance (Free Tier)

1. **Go to RDS Console**
   - Navigate to [AWS RDS Console](https://console.aws.amazon.com/rds/)

2. **Create Database**
   - Click "Create database"
   - Choose "Standard create"
   - Engine type: **PostgreSQL**
   - Version: **PostgreSQL 15.x** (latest)

3. **Free Tier Configuration**
   - Templates: **Free tier**
   - DB instance identifier: `recipe-db`
   - Master username: `postgres`
   - Master password: `your-secure-password`

4. **Instance Configuration**
   - DB instance class: **db.t3.micro** (Free tier eligible)
   - Storage: **20 GB** (Free tier limit)
   - Storage type: **General Purpose SSD (gp2)**

5. **Connectivity**
   - VPC: Default VPC
   - Subnet group: Default
   - Public access: **Yes** (for easier setup)
   - VPC security group: Create new
   - Database port: **5432**

6. **Additional Configuration**
   - Initial database name: `recipe_social`
   - Backup retention: **7 days** (Free tier)
   - Monitoring: **Disable enhanced monitoring** (saves costs)

### 1.2 Enable RDS Data API

1. **Modify Database**
   - Go to your database
   - Click "Modify"
   - Scroll to "Data API"
   - Check "Enable Data API"
   - Click "Continue" → "Apply immediately"

2. **Note Your ARN**
   - Copy the "Cluster ARN" for later use

## 🔐 **Step 2: Set Up Secrets Manager (FREE)**

### 2.1 Create Secret

1. **Go to Secrets Manager**
   - Navigate to [AWS Secrets Manager](https://console.aws.amazon.com/secretsmanager/)

2. **Store New Secret**
   - Choose "Credentials for RDS database"
   - Username: `postgres`
   - Password: `your-secure-password`
   - Database: Select your RDS instance
   - Secret name: `recipe-db-secret`

3. **Note Secret ARN**
   - Copy the ARN for Lambda configuration

## 🚀 **Step 3: Create Lambda Functions (FREE)**

### 3.1 Create Login Function

1. **Go to Lambda Console**
   - Navigate to [AWS Lambda](https://console.aws.amazon.com/lambda/)

2. **Create Function**
   - Function name: `recipe-login`
   - Runtime: **Go 1.x**
   - Architecture: **x86_64**

3. **Upload Code**
   - Create a zip file with your Go code
   - Upload via "Code source" → "Upload from .zip file"

4. **Environment Variables**
   ```
   DB_CLUSTER_ARN = arn:aws:rds:us-east-1:123456789012:cluster:recipe-db
   DB_SECRET_ARN = arn:aws:secretsmanager:us-east-1:123456789012:secret:recipe-db-secret
   DB_NAME = recipe_social
   JWT_SECRET = your-jwt-secret-key
   ```

5. **IAM Permissions**
   - Go to "Configuration" → "Permissions"
   - Click execution role
   - Add policies:
     - `AmazonRDSDataFullAccess`
     - `SecretsManagerReadWrite`

### 3.2 Create Other Functions

Repeat the same process for:
- `recipe-register`
- `recipe-get-posts`
- `recipe-create-post`
- `recipe-feed`

## 🌐 **Step 4: Set Up API Gateway (FREE)**

### 4.1 Create REST API

1. **Go to API Gateway**
   - Navigate to [AWS API Gateway](https://console.aws.amazon.com/apigateway/)

2. **Create API**
   - Choose "REST API"
   - API name: `recipe-social-api`
   - Description: `Recipe Social Media API`

### 4.2 Create Resources

1. **Create Auth Resource**
   - Click "Actions" → "Create Resource"
   - Resource name: `auth`
   - Create sub-resources: `login`, `register`

2. **Create Posts Resource**
   - Resource name: `posts`
   - Create sub-resource: `feed`

### 4.3 Configure Methods

1. **Login Method**
   - Select `/auth/login`
   - Click "Actions" → "Create Method" → "POST"
   - Integration type: "Lambda Function"
   - Lambda function: `recipe-login`
   - Enable CORS

2. **Repeat for Other Endpoints**
   - Register: POST `/auth/register`
   - Get Posts: GET `/posts`
   - Create Post: POST `/posts`
   - Feed: GET `/posts/feed`

### 4.4 Deploy API

1. **Deploy**
   - Click "Actions" → "Deploy API"
   - Stage: `prod`
   - Note your API URL

## 🗄️ **Step 5: Set Up Database Schema**

### 5.1 Using RDS Query Editor

1. **Open Query Editor**
   - Go to your RDS database
   - Click "Query Editor"
   - Connect with your credentials

2. **Run Schema**
   - Copy contents from `backend/database/schema.sql`
   - Paste and execute in Query Editor

## 💻 **Step 6: Deploy Frontend (FREE)**

### 6.1 Deploy to Vercel (FREE)

1. **Install Vercel CLI**
   ```bash
   npm i -g vercel
   ```

2. **Deploy**
   ```bash
   cd recipe-socialMedia-personal
   vercel
   ```

3. **Set Environment Variables**
   - In Vercel dashboard, go to your project
   - Settings → Environment Variables
   - Add: `NEXT_PUBLIC_API_URL = https://your-api-gateway-url`

### 6.2 Alternative: Deploy to Netlify (FREE)

1. **Connect GitHub**
   - Push your code to GitHub
   - Connect to Netlify
   - Auto-deploy on push

## 💰 **Cost Breakdown (FREE TIER)**

| Service | Free Tier | Your Usage | Cost |
|---------|-----------|------------|------|
| **Lambda** | 1M requests | ~1K/month | **$0.00** |
| **API Gateway** | 1M calls | ~1K/month | **$0.00** |
| **RDS PostgreSQL** | 750 hours | 24/7 for 1 month | **$0.00** |
| **Secrets Manager** | 10K calls | ~100/month | **$0.00** |
| **CloudWatch** | 10GB logs | ~1GB/month | **$0.00** |
| **Vercel Hosting** | Unlimited | Personal use | **$0.00** |
| **Total Monthly Cost** | | | **$0.00** |

## ⚠️ **Important Free Tier Limits**

### **After 1 Month (RDS)**
- RDS will start costing ~$15/month
- **Solution**: Use Aurora Serverless (pay per use)
- **Alternative**: Use external PostgreSQL (ElephantSQL, Supabase)

### **If You Exceed Limits**
- Lambda: $0.20 per 1M requests
- API Gateway: $3.50 per 1M calls
- RDS: ~$15/month for db.t3.micro

## 🎯 **Optimization Tips**

### **Stay Within Free Tier**
1. **Monitor Usage**
   - Check AWS Billing dashboard
   - Set up billing alerts
   - Monitor CloudWatch metrics

2. **Optimize Lambda**
   - Use smaller memory allocation
   - Optimize code for faster execution
   - Use provisioned concurrency sparingly

3. **Optimize API Gateway**
   - Cache responses when possible
   - Use compression
   - Minimize payload sizes

## 🔄 **Migration Path (When Ready)**

### **Option 1: Keep Free Tier**
- Use Aurora Serverless (pay per use)
- Optimize for minimal usage
- Use external PostgreSQL

### **Option 2: Upgrade Gradually**
- Start with RDS t3.micro (~$15/month)
- Add more Lambda functions as needed
- Scale based on actual usage

### **Option 3: Hybrid Approach**
- Keep Lambda + API Gateway (serverless benefits)
- Use external PostgreSQL (Supabase, PlanetScale)
- Best of both worlds

## 🎉 **Success!**

Your Recipe Social Media platform is now running **completely FREE** for:
- ✅ 1 month (RDS free tier)
- ✅ Forever (Lambda, API Gateway, Vercel)
- ✅ Personal project scale
- ✅ Production-ready architecture

## 📊 **Monitoring Your Usage**

1. **AWS Billing Dashboard**
   - Monitor monthly costs
   - Set up billing alerts
   - Track free tier usage

2. **CloudWatch Metrics**
   - Lambda invocations
   - API Gateway requests
   - RDS connections

3. **Cost Optimization**
   - Use AWS Cost Explorer
   - Set up budgets
   - Monitor daily costs

Your personal project can run **completely free** for the first month, and then for just a few dollars per month after that! 🚀
