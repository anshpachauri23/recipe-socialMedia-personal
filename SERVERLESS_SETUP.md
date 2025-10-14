# 🚀 Serverless Recipe Social Media Setup Guide

This guide will help you deploy your Recipe Social Media platform using AWS API Gateway + Lambda functions instead of a traditional server.

## 🏗️ Architecture

```
Frontend (Next.js) → API Gateway → Lambda Functions → RDS Data API → PostgreSQL
```

## 📋 Prerequisites

1. **AWS Account** with appropriate permissions
2. **AWS CLI** installed and configured
3. **SAM CLI** installed (`brew install aws-sam-cli`)
4. **Go 1.21+** installed
5. **PostgreSQL database** on AWS RDS

## 🔧 Step 1: Set Up AWS RDS Data API

### 1.1 Enable RDS Data API
```bash
# Enable Data API for your RDS cluster
aws rds modify-db-cluster \
    --db-cluster-identifier your-cluster-name \
    --enable-http-endpoint \
    --apply-immediately
```

### 1.2 Create Secrets Manager Secret
```bash
# Create a secret for your database credentials
aws secretsmanager create-secret \
    --name recipe-db-secret \
    --description "Database credentials for Recipe Social Media" \
    --secret-string '{"username":"your_username","password":"your_password","engine":"postgres","host":"your-rds-endpoint","port":5432,"dbname":"your_database"}'
```

## 🗄️ Step 2: Set Up Database Schema

### 2.1 Connect to your RDS instance
```bash
# Using psql
psql -h your-rds-endpoint.amazonaws.com -U your_username -d your_database

# Or using AWS RDS Query Editor in the console
```

### 2.2 Run the schema
```sql
-- Copy and paste the contents of backend/database/schema.sql
-- This will create all the necessary tables, indexes, and triggers
```

## 🚀 Step 3: Deploy Lambda Functions

### 3.1 Navigate to Lambda directory
```bash
cd lambda
```

### 3.2 Configure deployment parameters
Edit `template.yaml` and update these parameters:
- `DatabaseClusterArn`: Your RDS cluster ARN
- `DatabaseSecretArn`: Your Secrets Manager secret ARN
- `DatabaseName`: Your database name
- `JwtSecret`: A secure JWT secret key

### 3.3 Deploy using SAM
```bash
# Build the functions
sam build

# Deploy (first time - guided)
sam deploy --guided

# Subsequent deployments
sam deploy
```

### 3.4 Alternative: Use the deployment script
```bash
chmod +x deploy.sh
./deploy.sh
```

## 🔗 Step 4: Update Frontend Configuration

### 4.1 Get your API Gateway URL
After deployment, you'll get an API Gateway URL like:
```
https://abc123def4.execute-api.us-east-1.amazonaws.com/Prod
```

### 4.2 Update environment variables
Create/update your `.env` file:
```bash
NEXT_PUBLIC_API_URL=https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/Prod
```

### 4.3 Update frontend code
The frontend is already configured to use the Lambda endpoints. The API calls will automatically use the new serverless endpoints.

## 📊 Step 5: Test Your Deployment

### 5.1 Test API endpoints
```bash
# Test login endpoint
curl -X POST https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/Prod/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Test register endpoint
curl -X POST https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/Prod/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","full_name":"Test User"}'
```

### 5.2 Test frontend
```bash
# Start your Next.js frontend
npm run dev
```

## 🎯 Available Lambda Functions

| Function | Endpoint | Method | Description |
|----------|----------|--------|-------------|
| LoginFunction | `/auth/login` | POST | User authentication |
| RegisterFunction | `/auth/register` | POST | User registration |
| GetPostsFunction | `/posts` | GET | Get all posts |
| CreatePostFunction | `/posts` | POST | Create new post |
| FeedFunction | `/posts/feed` | GET | Get user feed |

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **CORS Enabled**: Cross-origin requests allowed
- **RDS Data API**: Secure database access without exposing credentials
- **Secrets Manager**: Encrypted credential storage
- **IAM Roles**: Least privilege access

## 📈 Benefits of Serverless Architecture

1. **Cost Effective**: Pay only for what you use
2. **Auto Scaling**: Handles traffic spikes automatically
3. **No Server Management**: AWS handles infrastructure
4. **High Availability**: Built-in redundancy
5. **Security**: AWS handles security patches and updates

## 🛠️ Troubleshooting

### Common Issues:

1. **Lambda timeout**: Increase timeout in `template.yaml`
2. **CORS errors**: Check CORS configuration in API Gateway
3. **Database connection**: Verify RDS Data API is enabled
4. **Secrets access**: Check IAM permissions for Secrets Manager

### Debug Commands:
```bash
# View Lambda logs
sam logs -n LoginFunction --stack-name recipe-social-media

# Test locally
sam local start-api

# Check deployment status
aws cloudformation describe-stacks --stack-name recipe-social-media
```

## 🔄 Updating the Application

1. **Update Lambda functions**: Modify code and run `sam deploy`
2. **Update frontend**: Modify React components and redeploy
3. **Database changes**: Run new SQL scripts on RDS

## 📚 Additional Resources

- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
- [API Gateway Documentation](https://docs.aws.amazon.com/apigateway/)
- [RDS Data API Documentation](https://docs.aws.amazon.com/rds/latest/AuroraUserGuide/data-api.html)
- [SAM CLI Documentation](https://docs.aws.amazon.com/serverless-application-model/)

## 🎉 Success!

Your Recipe Social Media platform is now running on a serverless architecture with:
- ✅ Scalable Lambda functions
- ✅ Secure API Gateway
- ✅ Managed PostgreSQL database
- ✅ Modern React frontend
- ✅ Production-ready security

Happy cooking! 🍳
