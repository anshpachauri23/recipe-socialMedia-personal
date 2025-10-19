# AWS Credentials Setup Guide

## The Problem
You're getting this error when trying to upload images:
```
Failed to upload image: failed to upload to S3: InvalidAccessKeyId: The AWS Access Key Id you provided does not exist in our records.
```

This means your AWS credentials are not configured properly.

## Solution: Get Your AWS Credentials

### Option 1: Use AWS CLI (Recommended)
If you have AWS CLI installed, you can get your credentials:

```bash
# Check if you have AWS CLI
aws --version

# If you have it, get your credentials
aws configure list
```

### Option 2: Get from AWS Console
1. Go to [AWS Console](https://console.aws.amazon.com/)
2. Click on your username (top right)
3. Go to "Security credentials"
4. Scroll down to "Access keys"
5. Click "Create access key"
6. Choose "Application running outside AWS"
7. Copy the Access Key ID and Secret Access Key

### Option 3: Use Existing IAM User
1. Go to IAM in AWS Console
2. Find your user
3. Go to "Security credentials" tab
4. Create new access key if needed

## Configure Your Credentials

Once you have your credentials, update your `backend/local.env` file:

```env
# AWS S3 Configuration
S3_BUCKET=recipe-social-images
S3_REGION=us-east-2

# AWS Credentials (replace with your actual credentials)
AWS_ACCESS_KEY_ID=AKIA...your_access_key_here
AWS_SECRET_ACCESS_KEY=your_secret_key_here
```

## Alternative: Use AWS Profile
If you have AWS CLI configured, you can also set the profile:

```env
# AWS S3 Configuration
S3_BUCKET=recipe-social-images
S3_REGION=us-east-2

# Use AWS profile instead of individual credentials
AWS_PROFILE=default
```

## Test Your Setup

After adding your credentials, restart the backend:

```bash
# Kill the current backend
pkill -f "go run main.go"

# Start it again
cd backend && go run main.go
```

Then try creating a new post with an image.

## Security Note
- Never commit your AWS credentials to git
- The `local.env` file should be in your `.gitignore`
- For production, use IAM roles instead of access keys
