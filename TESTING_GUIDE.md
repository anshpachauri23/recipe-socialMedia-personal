# 🧪 Recipe Social Media Testing Guide

## 🎯 **Testing Checklist**

### **✅ Step 1: Database Setup (COMPLETED)**
- ✅ RDS PostgreSQL database created
- ✅ All tables created successfully
- ✅ Indexes and triggers working
- ✅ Database connection working

### **🔧 Step 2: Lambda Functions Testing**

#### **2.1 Test Individual Lambda Functions**

**Login Function Test:**
1. Go to AWS Lambda Console
2. Click on `recipe-login` function
3. Go to "Test" tab
4. Create test event:
```json
{
  "body": "{\"email\":\"test@example.com\",\"password\":\"password123\"}",
  "headers": {
    "Content-Type": "application/json"
  }
}
```
5. Click "Test" - should return error (user doesn't exist yet)

**Register Function Test:**
1. Click on `recipe-register` function
2. Create test event:
```json
{
  "body": "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\",\"full_name\":\"Test User\"}",
  "headers": {
    "Content-Type": "application/json"
  }
}
```
3. Click "Test" - should create user successfully

### **🌐 Step 3: API Gateway Testing**

#### **3.1 Get Your API Gateway URL**
1. Go to AWS API Gateway Console
2. Click on your API (recipe-social-api)
3. Copy the "Invoke URL"
4. Format: `https://abc123def4.execute-api.us-east-2.amazonaws.com/Prod`

#### **3.2 Test API Endpoints**

**Test Registration:**
```bash
curl -X POST https://YOUR-API-GATEWAY-URL/Prod/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

**Test Login:**
```bash
curl -X POST https://YOUR-API-GATEWAY-URL/Prod/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Test Get Posts:**
```bash
curl -X GET https://YOUR-API-GATEWAY-URL/Prod/posts
```

### **💻 Step 4: Frontend Testing**

#### **4.1 Start Frontend (Already Running)**
```bash
cd recipe-socialMedia-personal
npm run dev
```
Frontend is running at: http://localhost:3000

#### **4.2 Update Frontend API URL**
1. Open `.env` file` in your project root
2. Update the API URL:
```bash
NEXT_PUBLIC_API_URL=https://YOUR-ACTUAL-API-GATEWAY-URL/Prod
```

#### **4.3 Test Frontend Features**

**Test User Registration:**
1. Go to http://localhost:3000/auth/register
2. Fill out the registration form
3. Click "Sign Up"
4. Should redirect to login page

**Test User Login:**
1. Go to http://localhost:3000/auth/login
2. Enter credentials
3. Click "Sign In"
4. Should redirect to feed page

**Test Feed Page:**
1. After login, should see feed page
2. Should show "No posts yet" message
3. Navigation should work

**Test Create Post:**
1. Click "Create" in navigation
2. Should see create post form
3. Fill out form and submit
4. Should create post successfully

### **🔍 Step 5: Database Verification**

#### **5.1 Check Users Table**
```bash
PGPASSWORD=Chahek.231104 /opt/homebrew/opt/postgresql@14/bin/psql -h recipe-db.cbcy4q2c4epi.us-east-2.rds.amazonaws.com -U postgres -d postgres -c "SELECT * FROM users;"
```

#### **5.2 Check Posts Table**
```bash
PGPASSWORD=Chahek.231104 /opt/homebrew/opt/postgresql@14/bin/psql -h recipe-db.cbcy4q2c4epi.us-east-2.rds.amazonaws.com -U postgres -d postgres -c "SELECT * FROM posts;"
```

### **🚨 Common Issues & Solutions**

#### **Issue 1: Lambda Function Timeout**
- **Solution**: Increase timeout in Lambda configuration
- **Check**: CloudWatch logs for errors

#### **Issue 2: Database Connection Failed**
- **Solution**: Check VPC configuration in Lambda
- **Check**: Security group allows Lambda access

#### **Issue 3: CORS Errors**
- **Solution**: Enable CORS in API Gateway
- **Check**: API Gateway CORS configuration

#### **Issue 4: Frontend API Calls Failing**
- **Solution**: Check API Gateway URL in .env file
- **Check**: Network tab in browser developer tools

### **✅ Success Criteria**

Your application is working correctly if:
- ✅ Users can register successfully
- ✅ Users can login and get JWT token
- ✅ Users can view the feed page
- ✅ Users can create recipe posts
- ✅ Database stores all data correctly
- ✅ No console errors in browser
- ✅ All API endpoints return proper responses

### **🎯 Next Steps After Testing**

1. **Fix any issues** found during testing
2. **Deploy frontend** to Vercel/Netlify
3. **Set up custom domain** (optional)
4. **Add more features** (image upload, search, etc.)
5. **Monitor performance** in AWS CloudWatch

## 🎉 **Congratulations!**

If all tests pass, your Recipe Social Media platform is working perfectly! 🚀
