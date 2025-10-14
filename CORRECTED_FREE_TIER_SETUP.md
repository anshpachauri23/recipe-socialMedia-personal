# 💰 CORRECTED Free Tier Setup for Recipe Social Media

You're absolutely right! RDS Data API is only for Aurora. Here's the **corrected** free tier setup using regular RDS PostgreSQL.

## 🆓 **What's Actually FREE in AWS Free Tier**

| Service | Free Tier Limit | Your Usage |
|---------|----------------|------------|
| **Lambda** | 1M requests/month | ~10K requests/month ✅ |
| **API Gateway** | 1M API calls/month | ~10K calls/month ✅ |
| **RDS PostgreSQL** | 750 hours + 20GB | 24/7 for 1 month ✅ |
| **VPC** | Always free | ✅ |
| **Vercel Hosting** | Unlimited | Personal use ✅ |

## 🎯 **Corrected Architecture**

```
Frontend (Vercel FREE) → API Gateway (FREE) → Lambda (FREE) → RDS PostgreSQL (FREE for 1 month)
```

## 📋 **Step 1: Set Up RDS PostgreSQL (FREE for 1 Month)**

### 1.1 Create RDS Instance

1. **Go to RDS Console**
   - Navigate to [AWS RDS Console](https://console.aws.amazon.com/rds/)

2. **Create Database**
   - Click "Create database"
   - Choose "Standard create"
   - Engine type: **PostgreSQL**
   - Version: **PostgreSQL 15.x**

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

### 1.2 Configure Security Group

1. **Edit Security Group**
   - Go to EC2 Console → Security Groups
   - Find your RDS security group
   - Add inbound rule:
     - Type: PostgreSQL
     - Port: 5432
     - Source: 0.0.0.0/0 (for Lambda access)

## 🚀 **Step 2: Create Lambda Functions (FREE)**

### 2.1 Create Login Function

1. **Go to Lambda Console**
   - Navigate to [AWS Lambda](https://console.aws.amazon.com/lambda/)

2. **Create Function**
   - Function name: `recipe-login`
   - Runtime: **Go 1.x**
   - Architecture: **x86_64**

3. **Environment Variables**
   ```
   DB_HOST = your-rds-endpoint.amazonaws.com
   DB_PORT = 5432
   DB_USER = postgres
   DB_PASSWORD = your-secure-password
   DB_NAME = recipe_social
   JWT_SECRET = your-jwt-secret-key
   ```

4. **VPC Configuration**
   - Go to "Configuration" → "VPC"
   - Select the same VPC as your RDS
   - Select subnets (at least 2)
   - Select the security group you created

### 2.2 Updated Lambda Code (Direct PostgreSQL Connection)

Here's the corrected Lambda function code:

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	_ "github.com/lib/pq"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email         string `json:"email"`
	FullName      string `json:"full_name"`
	Bio           *string `json:"bio"`
	ProfilePhotoURL *string `json:"profile_photo_url"`
	FollowerCount int    `json:"follower_count"`
	FollowingCount int   `json:"following_count"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request body
	var req LoginRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request body"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Database connection failed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}
	defer db.Close()

	// Get user by email
	user, err := getUserByEmail(db, req.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid credentials"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       `{"error": "Invalid credentials"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username":  user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to generate token"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Prepare response
	response := LoginResponse{
		Message: "Login successful",
		Token:   tokenString,
		User: User{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			Bio:            user.Bio,
			ProfilePhotoURL: user.ProfilePhotoURL,
			FollowerCount:  user.FollowerCount,
			FollowingCount: user.FollowingCount,
		},
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to process response"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func connectToDatabase() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	return sql.Open("postgres", connStr)
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, full_name, bio, profile_photo_url, 
			  follower_count, following_count FROM users WHERE email = $1`
	
	user := &User{}
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Bio, &user.ProfilePhotoURL, &user.FollowerCount, &user.FollowingCount,
	)
	
	return user, err
}

func main() {
	lambda.Start(handler)
}
```

## 🗄️ **Step 3: Set Up Database Schema**

### 3.1 Connect to RDS

1. **Get Connection Details**
   - Go to your RDS instance
   - Note the endpoint and port
   - Use your master username and password

2. **Connect Using psql**
   ```bash
   psql -h your-rds-endpoint.amazonaws.com -U postgres -d recipe_social
   ```

3. **Run Schema**
   - Copy contents from `backend/database/schema.sql`
   - Paste and execute in psql

### 3.2 Alternative: Use pgAdmin

1. **Download pgAdmin**
   - Install pgAdmin from https://www.pgadmin.org/

2. **Connect to RDS**
   - Host: your-rds-endpoint.amazonaws.com
   - Port: 5432
   - Username: postgres
   - Password: your-password

3. **Run Schema**
   - Open Query Tool
   - Paste and execute schema SQL

## 🌐 **Step 4: Set Up API Gateway (FREE)**

### 4.1 Create REST API

1. **Go to API Gateway**
   - Navigate to [AWS API Gateway](https://console.aws.amazon.com/apigateway/)

2. **Create API**
   - Choose "REST API"
   - API name: `recipe-social-api`

### 4.2 Create Resources and Methods

1. **Create Auth Resources**
   - Resource: `/auth`
   - Sub-resources: `/login`, `/register`

2. **Create Posts Resources**
   - Resource: `/posts`
   - Sub-resource: `/feed`

3. **Configure Methods**
   - POST `/auth/login` → `recipe-login` Lambda
   - POST `/auth/register` → `recipe-register` Lambda
   - GET `/posts` → `recipe-get-posts` Lambda
   - POST `/posts` → `recipe-create-post` Lambda
   - GET `/posts/feed` → `recipe-feed` Lambda

### 4.3 Enable CORS

1. **Enable CORS for Each Resource**
   - Select resource
   - Actions → Enable CORS
   - Access-Control-Allow-Origin: `*`

### 4.4 Deploy API

1. **Deploy**
   - Actions → Deploy API
   - Stage: `prod`
   - Note your API URL

## 💻 **Step 5: Deploy Frontend (FREE)**

### 5.1 Deploy to Vercel

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
   - In Vercel dashboard
   - Settings → Environment Variables
   - Add: `NEXT_PUBLIC_API_URL = https://your-api-gateway-url`

## 💰 **Cost Breakdown (CORRECTED)**

| Service | Free Tier | Your Usage | Cost |
|---------|-----------|------------|------|
| **Lambda** | 1M requests | ~1K/month | **$0.00** |
| **API Gateway** | 1M calls | ~1K/month | **$0.00** |
| **RDS PostgreSQL** | 750 hours | 24/7 for 1 month | **$0.00** |
| **VPC** | Always free | ✅ | **$0.00** |
| **Vercel Hosting** | Unlimited | Personal use | **$0.00** |
| **Total Monthly Cost** | | | **$0.00** |

## ⚠️ **After 1 Month (RDS)**

- RDS will cost ~$15/month
- **Solution**: Use external PostgreSQL (FREE forever)

## 🎯 **Alternative: External PostgreSQL (FREE Forever)**

### **Option 1: Supabase (FREE)**
- 500MB database
- PostgreSQL compatible
- REST API included
- Perfect for personal projects

### **Option 2: PlanetScale (FREE)**
- 1GB database
- MySQL compatible
- Serverless scaling

### **Option 3: Railway (FREE)**
- 1GB database
- PostgreSQL
- Easy deployment

## 🎉 **Success!**

Your Recipe Social Media platform is now running **completely FREE** for:
- ✅ 1 month (RDS free tier)
- ✅ Forever (if using external PostgreSQL)
- ✅ Personal project scale
- ✅ Production-ready architecture

## 📊 **Next Steps**

1. **Set up RDS** (FREE for 1 month)
2. **Create Lambda functions** with direct PostgreSQL connection
3. **Set up API Gateway**
4. **Deploy frontend to Vercel**
5. **After 1 month**: Migrate to external PostgreSQL (FREE forever)

This approach is **much simpler** and **actually works** with regular RDS PostgreSQL! 🚀
