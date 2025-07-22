# MoBot 2025 Deployment Guide

This guide covers deploying MoBot 2025 with Neon (PostgreSQL) and AWS S3 storage using GitHub Secrets for secure credential management.

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         GitHub Actions CI/CD            │
│         (Uses GitHub Secrets)           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────┴───────────────────────┐
│           MoBot 2025 API                │
├─────────────────────────────────────────┤
│  Neon PostgreSQL │ AWS S3 Storage       │
│  (Serverless DB) │ (AEP Files)          │
└─────────────────────┴───────────────────┘
```

## Prerequisites

1. GitHub repository with admin access
2. Neon account and database
3. AWS account with S3 bucket
4. (Optional) Pinecone account for vector search

## Step 1: Set Up GitHub Secrets

Go to your repository → Settings → Secrets and variables → Actions

### Required Secrets

| Secret Name | Description | Example |
|------------|-------------|---------|
| `NEON_DATABASE_URL` | Neon PostgreSQL connection string | `postgresql://user:pass@host/db?sslmode=require` |
| `AWS_ACCESS_KEY_ID` | AWS access key | `AKIA...` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | `wJal...` |
| `AWS_DEFAULT_REGION` | AWS region | `us-east-2` |
| `AWS_BUCKET` | S3 bucket name | `mobot2025-storage` |

### Optional Secrets

| Secret Name | Description |
|------------|-------------|
| `PINECONE_HOST` | Pinecone host URL |
| `PINECONE_INDEX` | Pinecone index name |
| `PINECONE_SECRET` | Pinecone API key |
| `ELEVENLABS_KEY` | ElevenLabs API key |
| `FIREFLIES_API_KEY` | Fireflies.ai API key |
| `VIMEO_CLIENT` | Vimeo client ID |
| `VIMEO_SECRET` | Vimeo client secret |
| `VIMEO_ACCESS` | Vimeo access token |

## Step 2: Configure Neon Database

1. Create a Neon project at [console.neon.tech](https://console.neon.tech)
2. Create a database (e.g., `mobot2025`)
3. Copy the connection string
4. Add it as `NEON_DATABASE_URL` in GitHub Secrets

### Database Branching (Optional)

Create branches for different environments:
```bash
# Staging branch
neon branches create --name staging

# Development branch  
neon branches create --name development
```

## Step 3: Configure AWS S3

1. Create an S3 bucket (e.g., `mobot2025-storage`)
2. Create an IAM user with programmatic access
3. Attach the following policy:

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
        "arn:aws:s3:::mobot2025-storage/*",
        "arn:aws:s3:::mobot2025-storage"
      ]
    }
  ]
}
```

4. Save the access key and secret key
5. Add them to GitHub Secrets

## Step 4: Local Development

1. Clone the repository:
```bash
git clone https://github.com/mojosolo/mobot2025.git
cd mobot2025
```

2. Create `.env` file (DO NOT COMMIT):
```bash
cp .env.example .env
# Edit .env with your credentials
```

3. Run locally:
```bash
# With SQLite (default)
go run cmd/mobot2025/main.go

# With Neon
MOBOT_DB_TYPE=postgres go run cmd/mobot2025/main.go
```

## Step 5: Deploy with GitHub Actions

The repository includes automated deployment workflows:

### Automatic Deployment

- **Main branch** → Staging environment
- **Production branch** → Production environment

### Manual Deployment

1. Create a pull request
2. Tests run automatically
3. Merge to main (deploys to staging)
4. Create a release (deploys to production)

## Step 6: Verify Deployment

### Check Database Connection
```bash
# Test Neon connection
psql "$NEON_DATABASE_URL" -c "SELECT version();"
```

### Check S3 Access
```bash
# List S3 bucket
aws s3 ls s3://mobot2025-storage/
```

### Health Check
```bash
curl https://your-deployment-url/health
```

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `MOBOT_DB_TYPE` | Database type (`sqlite` or `postgres`) | `sqlite` |
| `NEON_DATABASE_URL` | Neon connection string | - |
| `AWS_S3_ENABLED` | Enable S3 storage | `false` |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - |
| `AWS_DEFAULT_REGION` | AWS region | `us-east-2` |
| `AWS_BUCKET` | S3 bucket name | - |
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment name | `development` |
| `LOG_LEVEL` | Log level | `info` |

## Monitoring and Logs

### GitHub Actions
- Check Actions tab for deployment logs
- Secrets are masked in logs

### Application Logs
```bash
# View logs (example with systemd)
journalctl -u mobot2025 -f
```

### Database Monitoring
- Use Neon console for query insights
- Monitor connection pool usage

### S3 Monitoring
- Enable S3 access logging
- Use CloudWatch for metrics

## Troubleshooting

### Database Connection Issues
1. Verify `NEON_DATABASE_URL` format
2. Check SSL mode is set to `require`
3. Ensure database exists
4. Check connection limits

### S3 Access Issues
1. Verify IAM permissions
2. Check bucket policy
3. Ensure region matches
4. Verify credentials are active

### GitHub Actions Failures
1. Check secret names (case-sensitive)
2. Verify no extra spaces in values
3. Check workflow syntax
4. Review masked error messages

## Security Best Practices

1. **Rotate credentials regularly**
   - Update GitHub Secrets quarterly
   - Use AWS IAM roles when possible

2. **Use least privilege**
   - Limit IAM permissions
   - Use database roles

3. **Enable audit logging**
   - S3 access logging
   - Database query logging

4. **Monitor access**
   - Review GitHub audit log
   - Check AWS CloudTrail

5. **Encrypt at rest**
   - S3 server-side encryption
   - Neon encrypts by default

## Scaling Considerations

### Database
- Neon auto-scales compute
- Consider read replicas for heavy loads
- Use connection pooling

### Storage
- S3 scales automatically
- Consider CloudFront for global distribution
- Use lifecycle policies for old files

### Application
- Deploy behind a load balancer
- Use horizontal scaling
- Implement caching

## Backup and Recovery

### Database Backups
- Neon provides automatic backups
- Create manual snapshots before major changes
- Test restore procedures

### S3 Backups
- Enable versioning
- Set up cross-region replication
- Implement lifecycle policies

## Cost Optimization

### Neon
- Use compute auto-suspend
- Right-size compute units
- Monitor storage usage

### S3
- Use appropriate storage classes
- Enable lifecycle transitions
- Monitor request patterns

## Support

- GitHub Issues: [github.com/mojosolo/mobot2025/issues](https://github.com/mojosolo/mobot2025/issues)
- Documentation: [docs/](./docs/)
- Neon Support: [neon.tech/support](https://neon.tech/support)
- AWS Support: [aws.amazon.com/support](https://aws.amazon.com/support)