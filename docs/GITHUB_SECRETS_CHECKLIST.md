# GitHub Secrets Setup Checklist

This checklist ensures all GitHub Secrets are properly configured for MoBot 2025.

## ✅ Required Secrets Checklist

### 🗄️ Database (Neon PostgreSQL)
- [ ] `NEON_DATABASE_URL` - Main database connection string
  - Format: `postgresql://[user]:[password]@[host].neon.tech/[database]?sslmode=require`
  - Get from: Neon Console → Connection Details → Connection string

### ☁️ AWS S3 Storage
- [ ] `AWS_ACCESS_KEY_ID` - AWS IAM user access key
  - Format: `AKIA[0-9A-Z]{16}`
  - Get from: AWS Console → IAM → Users → Security credentials
- [ ] `AWS_SECRET_ACCESS_KEY` - AWS IAM user secret key
  - Format: 40 character string
  - Get from: AWS Console → IAM → Users → Create access key
- [ ] `AWS_DEFAULT_REGION` - AWS region for S3 bucket
  - Example: `us-east-2`, `us-west-1`, `eu-west-1`
  - Get from: AWS Console → S3 → Your bucket → Properties
- [ ] `AWS_BUCKET` - S3 bucket name
  - Example: `mobot2025-storage`
  - Get from: AWS Console → S3 → Buckets

## 📋 Optional Secrets Checklist

### 🗄️ Environment-Specific Databases
- [ ] `NEON_DATABASE_URL_STAGING` - Staging database
- [ ] `NEON_DATABASE_URL_PRODUCTION` - Production database

### 🔍 Pinecone (Future Vector Search)
- [ ] `PINECONE_HOST` - Pinecone index URL
- [ ] `PINECONE_INDEX` - Index name
- [ ] `PINECONE_SECRET` - API key

### 🔧 Third-Party Services
- [ ] `ELEVENLABS_KEY` - Text-to-speech API
- [ ] `FIREFLIES_API_KEY` - Meeting transcription
- [ ] `FIREFLIES_WEBHOOK_SECRET` - Webhook verification
- [ ] `VIMEO_CLIENT` - Video hosting client ID
- [ ] `VIMEO_SECRET` - Video hosting secret
- [ ] `VIMEO_ACCESS` - Video hosting access token

## 🚀 Setup Instructions

### Step 1: Navigate to GitHub Secrets
1. Go to your repository: `https://github.com/[your-username]/mobot2025`
2. Click **Settings** tab
3. In the left sidebar, click **Secrets and variables** → **Actions**
4. You'll see the **Repository secrets** page

### Step 2: Add Each Secret
For each secret in the checklist:
1. Click **New repository secret**
2. Enter the **Name** exactly as shown (case-sensitive!)
3. Enter the **Value** (will be masked after saving)
4. Click **Add secret**

### Step 3: Verify Setup
Run the verification script locally:
```bash
# Clone and navigate to repo
git clone https://github.com/[your-username]/mobot2025.git
cd mobot2025

# Run verification
./scripts/verify-github-secrets.sh
```

### Step 4: Test in GitHub Actions
1. Push a commit or create a PR
2. Check the **Actions** tab
3. Look for the "Test Cloud Integration" workflow
4. Verify all steps pass

## 🔒 Security Best Practices

### DO:
- ✅ Use strong, unique credentials for each service
- ✅ Rotate credentials every 90 days
- ✅ Use IAM roles with minimal permissions
- ✅ Enable MFA on AWS and Neon accounts
- ✅ Monitor access logs regularly

### DON'T:
- ❌ Share credentials via email or chat
- ❌ Commit credentials to code
- ❌ Use personal AWS accounts for production
- ❌ Grant broad permissions (like `s3:*`)
- ❌ Reuse passwords across services

## 🛠️ Troubleshooting

### Secret Not Found in Actions
- Check exact spelling (case-sensitive)
- Ensure no leading/trailing spaces
- Verify secret is saved (shows last updated time)

### AWS Authentication Failures
- Verify IAM user is active
- Check IAM policy includes required permissions
- Ensure bucket exists and is in correct region

### Neon Connection Issues
- Check connection string format
- Verify database exists
- Ensure SSL mode is set to `require`
- Check IP allowlist if configured

### Verification Script Issues
```bash
# Make script executable
chmod +x scripts/verify-github-secrets.sh

# Run with bash explicitly
bash scripts/verify-github-secrets.sh

# Debug mode
bash -x scripts/verify-github-secrets.sh
```

## 📊 Quick Status Check

Run this command to see which secrets are set:
```bash
env | grep -E "NEON|AWS|PINECONE|ELEVENLABS|FIREFLIES|VIMEO" | awk -F= '{print $1}' | sort
```

## 🔗 Related Documentation
- [Deployment Guide](DEPLOYMENT.md)
- [GitHub Secrets Setup Guide](GITHUB_SECRETS_SETUP.md)
- [API S3 Integration](API_S3_INTEGRATION.md)

---

**Remember**: Secrets are encrypted and only exposed to GitHub Actions runners during workflow execution. They are never visible in logs (shown as `***`).