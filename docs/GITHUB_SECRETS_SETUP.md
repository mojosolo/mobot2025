# GitHub Secrets Setup Guide

This guide explains how to configure GitHub Secrets for MoBot 2025 to securely manage all credentials.

## Required Secrets

Navigate to your repository settings → Secrets and variables → Actions, then add the following secrets:

### Database (Neon)
- `NEON_DATABASE_URL` - Your Neon PostgreSQL connection string
- `NEON_DATABASE_URL_STAGING` - Staging database URL (optional)
- `NEON_DATABASE_URL_PRODUCTION` - Production database URL (optional)

### AWS S3 Storage
- `AWS_ACCESS_KEY_ID` - AWS access key for S3
- `AWS_SECRET_ACCESS_KEY` - AWS secret access key
- `AWS_DEFAULT_REGION` - AWS region (e.g., us-east-2)
- `AWS_BUCKET` - S3 bucket name for storing AEP files

### Pinecone (Future Vector Search)
- `PINECONE_HOST` - Pinecone host URL
- `PINECONE_INDEX` - Pinecone index name
- `PINECONE_SECRET` - Pinecone API key

### Optional Services
- `ELEVENLABS_KEY` - ElevenLabs API key
- `FIREFLIES_API_KEY` - Fireflies.ai API key
- `FIREFLIES_WEBHOOK_SECRET` - Fireflies webhook secret
- `VIMEO_CLIENT` - Vimeo client ID
- `VIMEO_SECRET` - Vimeo client secret
- `VIMEO_ACCESS` - Vimeo access token

## Setting Up Secrets

1. Go to your GitHub repository
2. Click on "Settings" → "Secrets and variables" → "Actions"
3. Click "New repository secret"
4. Add each secret with the exact name listed above
5. Paste the credential value (GitHub will mask it automatically)

## Environment-Specific Secrets

For different environments (staging/production), you can use GitHub Environments:

1. Go to Settings → Environments
2. Create environments: `staging`, `production`
3. Add environment-specific secrets
4. Configure protection rules (e.g., require approval for production)

## Local Development

For local development, create a `.env` file (never commit this!):

```bash
# Copy the example
cp .env.example .env

# Edit with your credentials
nano .env
```

## Security Best Practices

1. **Never commit credentials** - Always use environment variables
2. **Rotate keys regularly** - Update secrets periodically
3. **Use least privilege** - Create IAM users/roles with minimal permissions
4. **Audit access** - Review who has access to secrets
5. **Use environments** - Separate staging/production credentials

## Verifying Setup

After adding all secrets, you can verify they're working by:

1. Creating a pull request
2. Checking the Actions tab for test results
3. Reviewing logs (secrets will be masked as ***)

## Troubleshooting

If you encounter issues:

1. Check secret names match exactly (case-sensitive)
2. Ensure no extra spaces in secret values
3. Verify IAM permissions for AWS credentials
4. Check Neon connection string format
5. Review GitHub Actions logs for masked error messages

## Example Connection Strings

### Neon Database URL Format:
```
postgresql://[user]:[password]@[host]/[database]?sslmode=require
```

### AWS S3 Bucket Policy Example:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::your-bucket-name/*",
        "arn:aws:s3:::your-bucket-name"
      ]
    }
  ]
}
```