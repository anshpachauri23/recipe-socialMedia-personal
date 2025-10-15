# Lambda Function Debugging Guide

## Current Issue: 500 Internal Server Error on Registration

### What We Fixed

1. **Improved Error Handling**: Added detailed error messages that will help identify the exact issue
2. **Database Connection Validation**: Added proper connection testing with `db.Ping()`
3. **Environment Variable Validation**: Added checks for required environment variables
4. **Better Logging**: Added more detailed logging for debugging

### Updated Lambda Functions

The following functions have been improved with better error handling:

- **Register Function**: `/lambda/functions/auth/register/main.go`
- **Login Function**: `/lambda/functions/auth/login/main.go`

### Key Improvements Made

1. **Database Connection**:
   ```go
   // Now validates environment variables and tests connection
   func connectToDatabase() (*sql.DB, error) {
       // Validate required environment variables
       requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
       for _, envVar := range requiredEnvVars {
           if os.Getenv(envVar) == "" {
               return nil, fmt.Errorf("missing required environment variable: %s", envVar)
           }
       }
       
       // Test the connection with db.Ping()
       if err := db.Ping(); err != nil {
           return nil, fmt.Errorf("failed to ping database: %v", err)
       }
   }
   ```

2. **Better Error Messages**:
   - Database connection errors now show the actual error
   - JWT secret validation
   - Duplicate username/email detection

3. **JWT Secret Validation**:
   ```go
   jwtSecret := os.Getenv("JWT_SECRET")
   if jwtSecret == "" {
       return events.APIGatewayProxyResponse{
           StatusCode: 500,
           Body: `{"error": "JWT secret not configured"}`,
       }
   }
   ```

### How to Deploy the Fixed Functions

1. **Build the functions**:
   ```bash
   cd lambda
   ./deploy-lambda.sh
   ```

2. **Upload to AWS Lambda**:
   - Go to AWS Lambda Console
   - For each function, upload the new zip file
   - Make sure environment variables are set correctly

### Environment Variables Required

Make sure your Lambda functions have these environment variables:

```
DB_HOST=recipe-db.cbcy4q2c4epi.us-east-2.rds.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Chahek.231104
DB_NAME=postgres
JWT_SECRET=85817ddb71e4d1f9dfc9b69ec7917e01b7776250a3b996194b874401659bdc97
```

### Debugging Steps

1. **Check CloudWatch Logs**:
   - Go to AWS Lambda Console
   - Click on your register function
   - Go to "Monitor" tab → "View CloudWatch logs"
   - Look for the specific error message

2. **Common Issues and Solutions**:

   **Issue**: "missing required environment variable"
   **Solution**: Set the missing environment variable in Lambda console

   **Issue**: "failed to ping database"
   **Solution**: 
   - Check if RDS instance is running
   - Verify security groups allow Lambda to access RDS
   - Check if Lambda is in the same VPC as RDS

   **Issue**: "JWT secret not configured"
   **Solution**: Set the JWT_SECRET environment variable

   **Issue**: "Username already exists" or "Email already exists"
   **Solution**: This is expected behavior - user needs to use different credentials

3. **VPC Configuration**:
   - Make sure your Lambda function is in the same VPC as your RDS instance
   - Check security groups allow inbound connections on port 5432

4. **Database Schema**:
   - Ensure the `users` table exists with the correct schema
   - Run the SQL from `setup_database.sql` if needed

### Testing the Fixed Functions

After deploying, test with:

```bash
# Test registration
curl -X POST https://760go4r862.execute-api.us-east-2.amazonaws.com/prod/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "password123",
    "full_name": "Test User"
  }'
```

### Expected Responses

**Success (201)**:
```json
{
  "message": "User created successfully",
  "token": "jwt-token-here",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "full_name": "Test User",
    "bio": null,
    "profile_photo_url": null,
    "follower_count": 0,
    "following_count": 0
  }
}
```

**Error (500)** - Now with detailed error message:
```json
{
  "error": "Database connection failed: failed to ping database: connection refused"
}
```

### Next Steps

1. Deploy the updated functions using the deployment script
2. Check CloudWatch logs for the specific error
3. Fix any environment variable or VPC issues
4. Test the registration endpoint again

The improved error handling should now give you the exact cause of the 500 error!
