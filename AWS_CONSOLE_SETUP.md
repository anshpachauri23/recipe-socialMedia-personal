# 🌐 AWS Console Setup Guide for Recipe Social Media

This guide will walk you through setting up your Recipe Social Media platform entirely using the AWS Management Console (web interface).

## 📋 Prerequisites

- AWS Account with appropriate permissions
- PostgreSQL database already created on AWS RDS
- Basic understanding of AWS services

## 🗄️ Step 1: Set Up RDS Data API

### 1.1 Enable RDS Data API

1. **Navigate to RDS Console**
   - Go to [AWS RDS Console](https://console.aws.amazon.com/rds/)
   - Click on "Databases" in the left sidebar

2. **Select Your Database**
   - Find your PostgreSQL database cluster
   - Click on the database identifier to open details

3. **Enable Data API**
   - Click on the "Connectivity & security" tab
   - Scroll down to "Data API" section
   - Click "Modify" button
   - Check the box for "Enable Data API"
   - Click "Continue" and then "Apply immediately"
   - Click "Modify cluster"

4. **Note Your Cluster ARN**
   - After enabling, note down the "Cluster ARN" from the "Connectivity & security" tab
   - It will look like: `arn:aws:rds:us-east-1:123456789012:cluster:your-cluster-name`

### 1.2 Create Secrets Manager Secret

1. **Navigate to Secrets Manager**
   - Go to [AWS Secrets Manager Console](https://console.aws.amazon.com/secretsmanager/)
   - Click "Store a new secret"

2. **Select Secret Type**
   - Choose "Credentials for RDS database"
   - Enter your database credentials:
     - Username: Your RDS master username
     - Password: Your RDS master password
   - Select your RDS database from the dropdown
   - Click "Next"

3. **Configure Secret**
   - Secret name: `recipe-db-secret`
   - Description: `Database credentials for Recipe Social Media`
   - Click "Next"

4. **Configure Rotation (Optional)**
   - Leave default settings for now
   - Click "Next"

5. **Review and Store**
   - Review your settings
   - Click "Store"
   - **Note the Secret ARN** - it will look like: `arn:aws:secretsmanager:us-east-1:123456789012:secret:recipe-db-secret-AbCdEf`

## 🗃️ Step 2: Set Up Database Schema

### 2.1 Using RDS Query Editor

1. **Open Query Editor**
   - Go back to your RDS database
   - Click on "Query Editor" tab
   - Click "Connect to database"

2. **Connect to Database**
   - Select your database
   - Choose "Connect with Secrets Manager"
   - Select the secret you just created
   - Click "Connect"

3. **Run Database Schema**
   - Copy the contents of `backend/database/schema.sql`
   - Paste into the Query Editor
   - Click "Run" to execute the schema

### 2.2 Alternative: Using pgAdmin or DBeaver

If you prefer a GUI tool:
1. Download pgAdmin or DBeaver
2. Connect using your RDS endpoint and credentials
3. Run the schema SQL file

## 🚀 Step 3: Deploy Lambda Functions

### 3.1 Prepare Lambda Function Code

1. **Create Lambda Functions in Console**
   - Go to [AWS Lambda Console](https://console.aws.amazon.com/lambda/)
   - Click "Create function"

2. **Create Login Function**
   - Function name: `recipe-login`
   - Runtime: Go 1.x
   - Click "Create function"

3. **Upload Code**
   - In the function, scroll down to "Code source"
   - Click "Upload from" → ".zip file"
   - Create a zip file with your Go code and upload

### 3.2 Configure Environment Variables

1. **Set Environment Variables**
   - In each Lambda function, go to "Configuration" tab
   - Click "Environment variables"
   - Add these variables:
     ```
     DB_CLUSTER_ARN = arn:aws:rds:us-east-1:123456789012:cluster:your-cluster-name
     DB_SECRET_ARN = arn:aws:secretsmanager:us-east-1:123456789012:secret:recipe-db-secret-AbCdEf
     DB_NAME = your_database_name
     JWT_SECRET = your-secure-jwt-secret-key
     ```

2. **Set IAM Permissions**
   - Go to "Configuration" → "Permissions"
   - Click on the execution role
   - Add these policies:
     - `AmazonRDSDataFullAccess`
     - `SecretsManagerReadWrite`

### 3.3 Create All Required Functions

Create these Lambda functions with the same configuration:

| Function Name | Description |
|---------------|-------------|
| `recipe-login` | User authentication |
| `recipe-register` | User registration |
| `recipe-get-posts` | Get all posts |
| `recipe-create-post` | Create new post |
| `recipe-feed` | Get user feed |

## 🌐 Step 4: Set Up API Gateway

### 4.1 Create API Gateway

1. **Navigate to API Gateway**
   - Go to [AWS API Gateway Console](https://console.aws.amazon.com/apigateway/)
   - Click "Create API"

2. **Choose API Type**
   - Select "REST API"
   - Click "Build"
   - Choose "New API"
   - API name: `recipe-social-api`
   - Description: `Recipe Social Media API`
   - Click "Create API"

### 4.2 Create Resources and Methods

1. **Create Auth Resources**
   - Click "Actions" → "Create Resource"
   - Resource name: `auth`
   - Click "Create Resource"
   - Create sub-resources: `login` and `register`

2. **Create Posts Resources**
   - Create resource: `posts`
   - Create sub-resource: `feed`

3. **Configure Methods**
   - For each resource, click "Actions" → "Create Method"
   - Choose HTTP method (GET, POST, etc.)
   - Integration type: "Lambda Function"
   - Select your Lambda function
   - Click "Save"

### 4.3 Enable CORS

1. **Enable CORS for Each Resource**
   - Select a resource
   - Click "Actions" → "Enable CORS"
   - Check "Enable CORS"
   - Access-Control-Allow-Origin: `*`
   - Click "Enable CORS and replace existing CORS headers"

### 4.4 Deploy API

1. **Deploy API**
   - Click "Actions" → "Deploy API"
   - Deployment stage: `[New Stage]`
   - Stage name: `prod`
   - Click "Deploy"

2. **Note Your API URL**
   - After deployment, you'll get a URL like:
   - `https://abc123def4.execute-api.us-east-1.amazonaws.com/prod`

## 🔧 Step 5: Configure Lambda Functions

### 5.1 Update Lambda Function Code

For each Lambda function, you'll need to:

1. **Create Go Binary**
   ```bash
   # For each function, create a zip file
   cd lambda/functions/auth/login
   go build -o main main.go
   zip function.zip main
   ```

2. **Upload to Lambda**
   - In Lambda console, go to "Code source"
   - Click "Upload from" → ".zip file"
   - Upload your zip file

### 5.2 Test Lambda Functions

1. **Test Each Function**
   - Go to "Test" tab in Lambda console
   - Create test events with sample JSON
   - Click "Test" to verify functions work

## 🔗 Step 6: Update Frontend Configuration

### 6.1 Update Environment Variables

1. **Create .env File**
   - In your project root, create `.env` file:
   ```
   NEXT_PUBLIC_API_URL=https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/prod
   ```

2. **Update API Calls**
   - The frontend is already configured to use the new endpoints
   - No code changes needed!

## 🧪 Step 7: Test Your Application

### 7.1 Test API Endpoints

1. **Test Login Endpoint**
   - Use Postman or curl to test:
   ```
   POST https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/prod/auth/login
   Content-Type: application/json
   
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```

2. **Test Register Endpoint**
   ```
   POST https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com/prod/auth/register
   Content-Type: application/json
   
   {
     "username": "testuser",
     "email": "test@example.com",
     "password": "password123",
     "full_name": "Test User"
   }
   ```

### 7.2 Test Frontend

1. **Start Frontend**
   ```bash
   npm run dev
   ```

2. **Test User Flow**
   - Register a new user
   - Login with credentials
   - Create a recipe post
   - View the feed

## 🔒 Step 8: Security Configuration

### 8.1 Configure IAM Roles

1. **Create Lambda Execution Role**
   - Go to [IAM Console](https://console.aws.amazon.com/iam/)
   - Click "Roles" → "Create role"
   - Select "AWS service" → "Lambda"
   - Attach policies:
     - `AWSLambdaBasicExecutionRole`
     - `AmazonRDSDataFullAccess`
     - `SecretsManagerReadWrite`

### 8.2 Configure RDS Security

1. **Update Security Groups**
   - Go to RDS Console
   - Click on your database
   - Go to "Connectivity & security" tab
   - Click on the security group
   - Add inbound rule for your IP address

## 📊 Step 9: Monitoring and Logs

### 9.1 CloudWatch Logs

1. **View Lambda Logs**
   - Go to [CloudWatch Console](https://console.aws.amazon.com/cloudwatch/)
   - Click "Logs" → "Log groups"
   - Find your Lambda function logs
   - Click to view recent logs

2. **Set Up Alarms**
   - Create CloudWatch alarms for errors
   - Set up notifications for failures

### 9.2 API Gateway Monitoring

1. **View API Metrics**
   - Go to API Gateway Console
   - Click on your API
   - Go to "Metrics" tab
   - Monitor request count, latency, errors

## 🎯 Step 10: Production Considerations

### 10.1 Custom Domain (Optional)

1. **Set Up Custom Domain**
   - In API Gateway, go to "Custom domain names"
   - Click "Create"
   - Enter your domain name
   - Configure SSL certificate

### 10.2 Environment-Specific Configuration

1. **Create Different Stages**
   - Create `dev`, `staging`, `prod` stages
   - Use different Lambda functions for each stage
   - Configure different environment variables

## 🛠️ Troubleshooting

### Common Issues:

1. **Lambda Timeout**
   - Increase timeout in Lambda configuration
   - Check database connection

2. **CORS Errors**
   - Verify CORS is enabled in API Gateway
   - Check preflight requests

3. **Database Connection Issues**
   - Verify RDS Data API is enabled
   - Check Secrets Manager permissions
   - Verify database credentials

4. **Authentication Errors**
   - Check JWT secret configuration
   - Verify token format

### Debug Steps:

1. **Check CloudWatch Logs**
   - Look for error messages in Lambda logs
   - Check API Gateway execution logs

2. **Test Individual Components**
   - Test Lambda functions individually
   - Test API Gateway endpoints
   - Test database connections

## 📚 Additional Resources

- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
- [API Gateway Documentation](https://docs.aws.amazon.com/apigateway/)
- [RDS Data API Documentation](https://docs.aws.amazon.com/rds/latest/AuroraUserGuide/data-api.html)
- [Secrets Manager Documentation](https://docs.aws.amazon.com/secretsmanager/)

## 🎉 Success!

Your Recipe Social Media platform is now running on AWS with:
- ✅ Serverless Lambda functions
- ✅ Scalable API Gateway
- ✅ Secure RDS database
- ✅ Modern React frontend
- ✅ Production-ready architecture

Happy cooking! 🍳

## 📞 Support

If you encounter any issues:
1. Check CloudWatch logs for error messages
2. Verify all environment variables are set correctly
3. Ensure IAM permissions are properly configured
4. Test each component individually

Your serverless Recipe Social Media platform is now ready to scale! 🚀
