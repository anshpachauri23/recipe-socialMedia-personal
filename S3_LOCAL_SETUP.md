# S3 Image Storage Setup for Local Development

This guide explains how to set up AWS S3 for image storage in your local development environment.

## 🎯 Overview

Your local backend now supports S3 image uploads just like the Lambda functions, but running locally for easier development.

## 📋 Prerequisites

- AWS Account with S3 access
- AWS CLI configured (optional, for testing)
- Your existing S3 bucket: `recipe-social-images`

## 🔧 Step 1: Verify S3 Bucket Configuration

### 1.1 Check Your Existing Bucket
Your S3 bucket should already be configured from the Lambda setup:
- **Bucket Name**: `recipe-social-images`
- **Region**: `us-east-2`
- **Public Read Access**: Enabled

### 1.2 Verify Bucket Policy
Ensure your bucket has this policy for public read access:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::recipe-social-images/*"
        }
    ]
}
```

## 🔧 Step 2: Configure AWS Credentials

### Option 1: AWS CLI (Recommended)
```bash
aws configure
```
Enter your:
- AWS Access Key ID
- AWS Secret Access Key
- Default region: `us-east-2`

### Option 2: Environment Variables
Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-2
```

### Option 3: AWS Credentials File
Create `~/.aws/credentials`:
```ini
[default]
aws_access_key_id = your_access_key
aws_secret_access_key = your_secret_key
```

## 🔧 Step 3: Test S3 Integration

### 3.1 Start Your Backend
```bash
cd backend
go run main.go
```

### 3.2 Test Image Upload
Create a test post with an image through your frontend, or use this curl command:

```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Test Recipe",
    "description": "Testing S3 upload",
    "images": [
      {
        "image_data": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD...",
        "step_description": "Test step",
        "step_number": 1
      }
    ]
  }'
```

## 🎯 How It Works

### Frontend → Backend → S3 Flow:
1. **User selects image** → Frontend converts to base64
2. **Frontend sends** `image_data` (base64) to local backend
3. **Backend uploads** to S3 and gets public URL
4. **Backend stores** S3 URL in database
5. **Frontend displays** image from S3 URL

### Image URLs Generated:
```
https://recipe-social-images.s3.us-east-2.amazonaws.com/recipe-images/1705123456-step-1.jpg
```

## 🔧 Step 4: Environment Configuration

Your `backend/local.env` should include:

```env
# AWS S3 Configuration
S3_BUCKET=recipe-social-images
S3_REGION=us-east-2
```

## 🚀 Benefits of Local S3 Development

✅ **Same S3 Integration**: Identical to Lambda functions  
✅ **Real Image Storage**: Images stored in your AWS S3 bucket  
✅ **Easy Debugging**: Full control over upload process  
✅ **Cost Effective**: No Lambda execution costs  
✅ **Faster Development**: No cold starts or deployment delays  

## 🐛 Troubleshooting

### Common Issues:

1. **AWS Credentials Not Found**
   ```
   Error: NoCredentialProviders: no valid providers in chain
   ```
   **Solution**: Configure AWS credentials (see Step 2)

2. **S3 Bucket Access Denied**
   ```
   Error: AccessDenied: Access Denied
   ```
   **Solution**: Check your AWS credentials have S3 permissions

3. **Images Not Loading**
   - Check bucket policy allows public read
   - Verify S3 URLs are correct
   - Check CORS settings if needed

### Debug Commands:

```bash
# Test AWS credentials
aws sts get-caller-identity

# List S3 buckets
aws s3 ls

# Check bucket contents
aws s3 ls s3://recipe-social-images/recipe-images/
```

## 💰 Cost Considerations

- **S3 Storage**: ~$0.023 per GB per month
- **S3 Requests**: ~$0.0004 per 1,000 PUT requests
- **Data Transfer**: Free for first 1GB per month

**Estimated cost for 1,000 images (1GB): ~$0.05/month** 🎯

## 🔄 Migration from Lambda

Your S3 setup is identical to the Lambda functions:
- Same bucket: `recipe-social-images`
- Same region: `us-east-2`
- Same image URLs and structure
- Same database schema

The only difference is the deployment method - from serverless Lambda to local server.

## 🎉 Success!

Once configured, your local development environment will:
- ✅ Upload images to S3
- ✅ Store S3 URLs in database
- ✅ Display images from S3
- ✅ Work identically to Lambda functions

You can now develop with full S3 integration without any AWS Lambda complexity!
