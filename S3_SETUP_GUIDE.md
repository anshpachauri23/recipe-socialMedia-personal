# 🚀 S3 Image Storage Setup Guide

## 📋 **Prerequisites:**
- AWS Account with S3 access
- Lambda functions with S3 permissions

## 🔧 **Step 1: Create S3 Bucket**

### **1.1 Create Bucket:**
1. Go to **AWS S3 Console**
2. Click **"Create bucket"**
3. Bucket name: `recipe-social-images` (or your preferred name)
4. Region: `us-east-2` (same as your Lambda functions)
5. **Uncheck "Block all public access"** (we need public read access for images)
6. Check **"I acknowledge that the current settings might result in this bucket and the objects within it becoming public"**
7. Click **"Create bucket"**

### **1.2 Configure Bucket Policy:**
1. Go to your bucket → **Permissions** tab
2. Scroll down to **"Bucket policy"**
3. Click **"Edit"** and add this policy:

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

## 🔧 **Step 2: Update Lambda Function Environment Variables**

### **2.1 Add S3 Environment Variables:**
For your `recipe-create-post` Lambda function, add these environment variables:

```
S3_BUCKET=recipe-social-images
AWS_REGION=us-east-2
```

### **2.2 Update IAM Role:**
Your Lambda function needs S3 permissions. Add this policy to your Lambda execution role:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:PutObjectAcl"
            ],
            "Resource": "arn:aws:s3:::recipe-social-images/*"
        }
    ]
}
```

## 🔧 **Step 3: Deploy Updated Lambda Function**

### **3.1 Upload New Function:**
1. Upload `create-post-s3.zip` to your `recipe-create-post` Lambda function
2. Set handler to: `bootstrap`
3. Add the environment variables from Step 2.1

### **3.2 Test the Function:**
The function will now:
- ✅ Accept base64 image data from frontend
- ✅ Upload images to S3
- ✅ Store S3 URLs in database
- ✅ Return proper image URLs

## 🎯 **How It Works:**

### **Frontend → Backend:**
1. **User selects image** → Converted to base64
2. **Frontend sends** `image_data` (base64) to Lambda
3. **Lambda uploads** to S3 and gets public URL
4. **Lambda stores** S3 URL in database

### **Image URLs Generated:**
```
https://recipe-social-images.s3.us-east-2.amazonaws.com/recipe-images/1705123456-step-1.jpg
```

## 📝 **Environment Variables Summary:**

### **Lambda Functions:**
```
# Database
DB_HOST=your-rds-endpoint.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=recipe_social
DB_SSLMODE=require

# JWT
JWT_SECRET=85817ddb71e4d1f9dfc9b69ec7917e01b7776250a3b996194b874401659bdc97

# S3 (NEW)
S3_BUCKET=recipe-social-images
AWS_REGION=us-east-2
```

### **Frontend:**
```
NEXT_PUBLIC_API_URL=https://760go4r862.execute-api.us-east-2.amazonaws.com/prod
```

## 🚀 **Testing:**

1. **Create a new recipe** with images
2. **Check S3 bucket** - images should appear
3. **Check database** - URLs should be S3 URLs
4. **View posts** - images should load from S3

## 💰 **Cost Considerations:**

- **S3 Storage**: ~$0.023 per GB per month
- **S3 Requests**: ~$0.0004 per 1,000 PUT requests
- **Data Transfer**: Free for first 1GB per month

**Estimated cost for 1,000 images (1GB): ~$0.05/month** 🎯
