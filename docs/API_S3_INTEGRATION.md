# API S3 Integration Guide

This guide documents the S3 storage integration for the MoBot 2025 API, which enables cloud storage of AEP files with secure access controls.

## Configuration

The API automatically uses S3 storage when the following environment variables are set:

```bash
# Enable S3 storage
AWS_S3_ENABLED=true

# AWS credentials (use GitHub Secrets in production)
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_DEFAULT_REGION=us-east-2
AWS_BUCKET=mobot2025-storage
```

## New API Endpoints

### 1. Upload AEP File

Upload an AEP file to the system. The file will be parsed, stored in S3 (if enabled), and metadata saved to the database.

**Endpoint:** `POST /api/v1/upload`

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Form field: `file` (the AEP file to upload)

**Example:**
```bash
curl -X POST \
  http://localhost:8080/api/v1/upload \
  -F "file=@template.aep"
```

**Response:**
```json
{
  "success": true,
  "project_id": 123,
  "metadata": {
    "file_name": "template.aep",
    "parsed_at": "2025-01-22T10:30:00Z",
    "compositions": [...],
    "text_layers": [...],
    "s3_bucket": "mobot2025-storage",
    "s3_key": "projects/2025/01/22/template.aep"
  },
  "s3_location": {
    "bucket": "mobot2025-storage",
    "key": "projects/2025/01/22/template.aep"
  }
}
```

### 2. Download AEP File

Download an AEP file from S3 storage.

**Endpoint:** `GET /api/v1/download/{project_id}`

**Parameters:**
- `project_id`: The ID of the project to download
- `presigned` (optional): Set to `true` to get a presigned S3 URL instead of streaming the file

**Examples:**

Direct download:
```bash
curl -O \
  http://localhost:8080/api/v1/download/123
```

Get presigned URL:
```bash
curl http://localhost:8080/api/v1/download/123?presigned=true
```

**Presigned URL Response:**
```json
{
  "download_url": "https://mobot2025-storage.s3.amazonaws.com/...",
  "expires_in": "15m"
}
```

### 3. List/Get Projects

List all projects or get details for a specific project.

**Endpoint:** `GET /api/v1/projects` or `GET /api/v1/projects/{project_id}`

**Query Parameters (for listing):**
- `q`: Search query
- `limit`: Maximum number of results (default: 50)

**Examples:**

List all projects:
```bash
curl http://localhost:8080/api/v1/projects
```

Search projects:
```bash
curl "http://localhost:8080/api/v1/projects?q=motion+graphics&limit=10"
```

Get specific project:
```bash
curl http://localhost:8080/api/v1/projects/123
```

**Project Detail Response:**
```json
{
  "project": {
    "file_name": "template.aep",
    "parsed_at": "2025-01-22T10:30:00Z",
    "s3_bucket": "mobot2025-storage",
    "s3_key": "projects/2025/01/22/template.aep",
    "compositions": [...],
    "text_layers": [...]
  },
  "download_url": "/api/v1/download/123",
  "presigned_url": "/api/v1/download/123?presigned=true"
}
```

### 4. Enhanced Parse Endpoint

The parse endpoint now supports parsing files from S3.

**Endpoint:** `POST /api/v1/parse`

**Request Body Options:**

Parse from local file:
```json
{
  "file_path": "/path/to/template.aep",
  "options": {
    "extract_text": true,
    "extract_media": true,
    "deep_analysis": false
  }
}
```

Parse from S3 key:
```json
{
  "s3_key": "projects/2025/01/22/template.aep",
  "options": {
    "extract_text": true,
    "extract_media": true,
    "deep_analysis": false
  }
}
```

Parse from project ID:
```json
{
  "project_id": 123,
  "options": {
    "extract_text": true,
    "extract_media": true,
    "deep_analysis": false
  }
}
```

## S3 Storage Structure

Files are organized in S3 with the following structure:
```
mobot2025-storage/
├── projects/
│   ├── 2025/
│   │   ├── 01/
│   │   │   ├── 22/
│   │   │   │   ├── template1.aep
│   │   │   │   └── template2.aep
│   │   │   └── 23/
│   │   │       └── template3.aep
```

## Security Considerations

1. **Presigned URLs**: When using presigned URLs, they expire after 15 minutes for security
2. **Access Control**: S3 bucket should have proper IAM policies (see DEPLOYMENT.md)
3. **HTTPS**: Always use HTTPS in production for secure file transfers
4. **File Validation**: Only `.aep` files are accepted for upload

## Local Development

For local development without S3:
```bash
# Disable S3 (uses local storage)
AWS_S3_ENABLED=false
STORAGE_TYPE=local
LOCAL_STORAGE_PATH=./storage
```

## Error Handling

Common error responses:

- `400 Bad Request`: Invalid file type or missing parameters
- `404 Not Found`: Project or file not found
- `500 Internal Server Error`: S3 access error or parsing failure

## Example Integration

```go
// Upload a file
func uploadAEP(filePath string) (*ProjectMetadata, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    part, err := writer.CreateFormFile("file", filepath.Base(filePath))
    if err != nil {
        return nil, err
    }
    io.Copy(part, file)
    writer.Close()

    resp, err := http.Post(
        "http://localhost:8080/api/v1/upload",
        writer.FormDataContentType(),
        &buf,
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Success  bool             `json:"success"`
        Metadata *ProjectMetadata `json:"metadata"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result.Metadata, nil
}
```

## Monitoring

Monitor S3 usage through:
- AWS CloudWatch metrics for bucket operations
- Application logs for upload/download activities
- Database queries for metadata access patterns