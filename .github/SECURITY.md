# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in MoBot 2025, please follow these steps:

1. **DO NOT** create a public GitHub issue
2. Email security details to: security@mobot2025.ai
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will acknowledge receipt within 48 hours and provide updates on the fix.

## Security Best Practices

### Credentials Management
- **Never** commit credentials to the repository
- Use environment variables for local development
- Use GitHub Secrets for CI/CD
- Rotate credentials every 90 days
- Use least-privilege IAM policies

### GitHub Secrets Configuration
All sensitive credentials MUST be stored as GitHub Secrets:
- Database URLs (Neon PostgreSQL)
- AWS credentials (S3 access)
- API keys (Pinecone, ElevenLabs, etc.)
- Webhook secrets

See [GitHub Secrets Setup Guide](../docs/GITHUB_SECRETS_SETUP.md) for details.

### AWS S3 Security
- Enable bucket encryption
- Use IAM policies, not bucket policies
- Enable access logging
- Use versioning for data recovery
- Implement lifecycle policies

### Database Security
- Always use SSL connections (`sslmode=require`)
- Use connection pooling
- Implement query timeouts
- Regular backups
- Monitor for unusual activity

### API Security
- Implement rate limiting
- Use HTTPS only
- Validate all inputs
- Sanitize file uploads
- Implement proper CORS policies

## Security Checklist

- [ ] All credentials in GitHub Secrets
- [ ] .env file in .gitignore
- [ ] No hardcoded credentials in code
- [ ] SSL enabled for database connections
- [ ] S3 bucket encryption enabled
- [ ] IAM policies follow least privilege
- [ ] API rate limiting implemented
- [ ] Input validation on all endpoints
- [ ] Regular security updates applied
- [ ] Monitoring and alerting configured

## Compliance

This project follows security best practices for:
- Credential management (GitHub Secrets)
- Data encryption (S3, Neon)
- Access control (IAM)
- Secure communications (HTTPS/SSL)

For questions about security, contact: security@mobot2025.ai