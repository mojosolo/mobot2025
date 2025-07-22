#!/bin/bash
# GitHub Secrets Verification Script for MoBot 2025
# This script helps verify that all required GitHub Secrets are configured

set -e

echo "🔐 GitHub Secrets Verification for MoBot 2025"
echo "============================================"
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if environment variable exists
check_env() {
    local var_name=$1
    local description=$2
    local required=$3
    
    if [ -n "${!var_name}" ]; then
        echo -e "${GREEN}✓${NC} $var_name is set ($description)"
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "${RED}✗${NC} $var_name is NOT set ($description) - REQUIRED"
            return 1
        else
            echo -e "${YELLOW}⚠${NC} $var_name is NOT set ($description) - Optional"
            return 0
        fi
    fi
}

# Track if any required secrets are missing
MISSING_REQUIRED=0

echo "📋 Checking Required Secrets:"
echo "-----------------------------"

# Database Secrets
echo ""
echo "🗄️  Database Configuration:"
check_env "NEON_DATABASE_URL" "Neon PostgreSQL connection string" "true" || MISSING_REQUIRED=1

# AWS S3 Secrets
echo ""
echo "☁️  AWS S3 Configuration:"
check_env "AWS_ACCESS_KEY_ID" "AWS access key for S3" "true" || MISSING_REQUIRED=1
check_env "AWS_SECRET_ACCESS_KEY" "AWS secret access key" "true" || MISSING_REQUIRED=1
check_env "AWS_DEFAULT_REGION" "AWS region (e.g., us-east-2)" "true" || MISSING_REQUIRED=1
check_env "AWS_BUCKET" "S3 bucket name for storing AEP files" "true" || MISSING_REQUIRED=1

echo ""
echo "📋 Checking Optional Secrets:"
echo "-----------------------------"

# Optional Database Secrets
echo ""
echo "🗄️  Optional Database Configurations:"
check_env "NEON_DATABASE_URL_STAGING" "Staging database URL" "false"
check_env "NEON_DATABASE_URL_PRODUCTION" "Production database URL" "false"

# Pinecone Secrets (Future Vector Search)
echo ""
echo "🔍 Pinecone Configuration (Future):"
check_env "PINECONE_HOST" "Pinecone host URL" "false"
check_env "PINECONE_INDEX" "Pinecone index name" "false"
check_env "PINECONE_SECRET" "Pinecone API key" "false"

# Optional Service Secrets
echo ""
echo "🔧 Optional Services:"
check_env "ELEVENLABS_KEY" "ElevenLabs API key" "false"
check_env "FIREFLIES_API_KEY" "Fireflies.ai API key" "false"
check_env "FIREFLIES_WEBHOOK_SECRET" "Fireflies webhook secret" "false"
check_env "VIMEO_CLIENT" "Vimeo client ID" "false"
check_env "VIMEO_SECRET" "Vimeo client secret" "false"
check_env "VIMEO_ACCESS" "Vimeo access token" "false"

# Additional Configuration
echo ""
echo "⚙️  Additional Configuration:"
check_env "MOBOT_DB_TYPE" "Database type (sqlite or postgres)" "false"
check_env "AWS_S3_ENABLED" "Enable S3 storage (true/false)" "false"
check_env "PORT" "Server port (default: 8080)" "false"
check_env "ENVIRONMENT" "Environment name (development/staging/production)" "false"
check_env "LOG_LEVEL" "Log level (debug/info/warn/error)" "false"

echo ""
echo "============================================"

if [ $MISSING_REQUIRED -eq 0 ]; then
    echo -e "${GREEN}✅ All required secrets are configured!${NC}"
    echo ""
    echo "Next steps:"
    echo "1. If running locally, create a .env file with these values"
    echo "2. For GitHub Actions, add these as repository secrets"
    echo "3. For production, use environment-specific secrets"
else
    echo -e "${RED}❌ Some required secrets are missing!${NC}"
    echo ""
    echo "To fix this:"
    echo "1. Go to your GitHub repository"
    echo "2. Navigate to Settings → Secrets and variables → Actions"
    echo "3. Click 'New repository secret' for each missing secret"
    echo "4. Add the secret with the exact name shown above"
    exit 1
fi

echo ""
echo "📚 Documentation:"
echo "- Setup Guide: docs/GITHUB_SECRETS_SETUP.md"
echo "- Deployment Guide: docs/DEPLOYMENT.md"
echo "- API Integration: docs/API_S3_INTEGRATION.md"

# Optional: Test connections if secrets are set
echo ""
read -p "Would you like to test the connections? (y/N) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🧪 Testing Connections..."
    echo "------------------------"
    
    # Test Neon connection
    if [ -n "$NEON_DATABASE_URL" ]; then
        echo -n "Testing Neon database connection... "
        if psql "$NEON_DATABASE_URL" -c "SELECT 1;" >/dev/null 2>&1; then
            echo -e "${GREEN}✓ Connected${NC}"
        else
            echo -e "${RED}✗ Failed${NC}"
        fi
    fi
    
    # Test AWS S3 access
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ] && [ -n "$AWS_BUCKET" ]; then
        echo -n "Testing AWS S3 access... "
        if aws s3 ls "s3://$AWS_BUCKET" --max-items 1 >/dev/null 2>&1; then
            echo -e "${GREEN}✓ Accessible${NC}"
        else
            echo -e "${RED}✗ Failed${NC}"
        fi
    fi
    
    # Test Pinecone connection (if configured)
    if [ -n "$PINECONE_HOST" ] && [ -n "$PINECONE_SECRET" ]; then
        echo -n "Testing Pinecone connection... "
        if curl -s -H "Api-Key: $PINECONE_SECRET" "$PINECONE_HOST/databases" >/dev/null 2>&1; then
            echo -e "${GREEN}✓ Connected${NC}"
        else
            echo -e "${RED}✗ Failed${NC}"
        fi
    fi
fi

echo ""
echo "✨ Verification complete!"