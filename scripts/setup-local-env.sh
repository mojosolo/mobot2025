#!/bin/bash
# Local Environment Setup Script for MoBot 2025
# This script helps set up your local development environment

set -e

echo "🚀 MoBot 2025 Local Environment Setup"
echo "====================================="
echo ""

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if .env exists
if [ -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file already exists${NC}"
    read -p "Do you want to backup and create a new one? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        backup_file=".env.backup.$(date +%Y%m%d_%H%M%S)"
        cp .env "$backup_file"
        echo -e "${GREEN}✓${NC} Backed up existing .env to $backup_file"
    else
        echo "Keeping existing .env file"
        exit 0
    fi
fi

# Copy from example
echo "Creating .env from .env.example..."
cp .env.example .env
echo -e "${GREEN}✓${NC} Created .env file"

echo ""
echo -e "${BLUE}📝 Configuration Options:${NC}"
echo ""

# Database configuration
echo "1. Database Configuration"
echo "   Choose your database type:"
echo "   a) SQLite (local development) - Default"
echo "   b) Neon PostgreSQL (cloud database)"
read -p "   Select (a/b): " db_choice

if [ "$db_choice" = "b" ] || [ "$db_choice" = "B" ]; then
    sed -i '' 's/MOBOT_DB_TYPE=sqlite/MOBOT_DB_TYPE=postgres/' .env
    echo ""
    echo "   Enter your Neon database URL:"
    echo "   Format: postgresql://user:pass@host.neon.tech/database?sslmode=require"
    read -p "   URL: " neon_url
    if [ -n "$neon_url" ]; then
        # Escape special characters in URL
        escaped_url=$(printf '%s\n' "$neon_url" | sed 's/[[\.*^$()+?{|]/\\&/g')
        sed -i '' "s|# NEON_DATABASE_URL=.*|NEON_DATABASE_URL=$escaped_url|" .env
    fi
fi

echo ""
echo "2. AWS S3 Storage Configuration"
read -p "   Enable S3 storage? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    sed -i '' 's/AWS_S3_ENABLED=false/AWS_S3_ENABLED=true/' .env
    
    echo "   Enter your AWS credentials:"
    read -p "   Access Key ID: " aws_key
    read -p "   Secret Access Key: " aws_secret
    read -p "   Region (default: us-east-2): " aws_region
    read -p "   Bucket name: " aws_bucket
    
    # Update .env with AWS credentials
    if [ -n "$aws_key" ]; then
        sed -i '' "s|# AWS_ACCESS_KEY_ID=.*|AWS_ACCESS_KEY_ID=$aws_key|" .env
    fi
    if [ -n "$aws_secret" ]; then
        sed -i '' "s|# AWS_SECRET_ACCESS_KEY=.*|AWS_SECRET_ACCESS_KEY=$aws_secret|" .env
    fi
    if [ -n "$aws_region" ]; then
        sed -i '' "s|# AWS_DEFAULT_REGION=.*|AWS_DEFAULT_REGION=${aws_region:-us-east-2}|" .env
    fi
    if [ -n "$aws_bucket" ]; then
        sed -i '' "s|# AWS_BUCKET=.*|AWS_BUCKET=$aws_bucket|" .env
    fi
fi

echo ""
echo "3. Server Configuration"
read -p "   Port (default: 8080): " port
if [ -n "$port" ]; then
    sed -i '' "s/PORT=8080/PORT=$port/" .env
fi

echo ""
echo -e "${GREEN}✓ Environment configuration complete!${NC}"
echo ""

# Create necessary directories
echo "Creating required directories..."
mkdir -p storage
mkdir -p logs
mkdir -p data
mkdir -p sample-aep
echo -e "${GREEN}✓${NC} Directories created"

echo ""
echo "📋 Next Steps:"
echo "1. Review your .env file and adjust any settings"
echo "2. For production, add these values to GitHub Secrets"
echo "3. Run './scripts/verify-github-secrets.sh' to verify setup"
echo "4. Build and run: go build -o mobot ./cmd/mobot2025/main.go && ./mobot serve"

echo ""
echo -e "${YELLOW}⚠️  Security Reminder:${NC}"
echo "- Never commit .env to version control"
echo "- Keep your credentials secure"
echo "- Rotate keys regularly"
echo "- Use GitHub Secrets for CI/CD"

echo ""
echo "✨ Setup complete! Happy coding!"