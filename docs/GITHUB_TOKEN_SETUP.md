# GitHub Token Authentication Setup

This guide explains how to set up GitHub Personal Access Token (PAT) authentication for the GoHoliday sync functionality.

## Why Use a GitHub Token?

The GoHoliday syncer fetches holiday data from the [vacanza/holidays](https://github.com/vacanza/holidays) repository via GitHub's API. Without authentication:

- **Rate limit**: 60 requests per hour
- **IP-based limiting**: Shared across all unauthenticated requests from your IP

With a GitHub token:

- **Rate limit**: 5,000 requests per hour
- **User-based limiting**: Personal quota, not shared
- **Better reliability**: Less likely to hit rate limits during sync operations

## Required Permissions

Since we're only reading from **public repositories**, you need minimal permissions:

### Fine-Grained Personal Access Token (Recommended)

1. Go to [GitHub Settings > Developer settings > Personal access tokens > Fine-grained tokens](https://github.com/settings/personal-access-tokens/new)
2. Click **"Generate new token"**
3. Configure:
   - **Repository access**: "Public Repositories (read-only)" or specific repository `vacanza/holidays`
   - **Repository permissions**:
     - **Contents**: `Read` ✅
     - **Metadata**: `Read` ✅ (optional)
4. Set expiration date (recommended: 90 days or less)
5. Click **"Generate token"**

### Classic Personal Access Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens > Tokens (classic)](https://github.com/settings/tokens/new)
2. Click **"Generate new token (classic)"**
3. Select scopes:
   - **`public_repo`** ✅ (Access public repositories)
   - Leave all other scopes unchecked
4. Set expiration date
5. Click **"Generate token"**

## Setup Methods

### Method 1: Environment Variable (Recommended)

```bash
# Linux/macOS
export GITHUB_TOKEN="ghp_your_token_here"

# Windows PowerShell
$env:GITHUB_TOKEN="ghp_your_token_here"

# Windows Command Prompt
set GITHUB_TOKEN=ghp_your_token_here
```

### Method 2: Command Line Flag

```bash
# Using the sync tool
go run cmd/sync/main.go -token="ghp_your_token_here" -list

# Or build and run
go build -o sync cmd/sync/main.go
./sync -token="ghp_your_token_here" -country=US
```

### Method 3: Programmatic Usage

```go
package main

import (
    "context"
    "github.com/coredds/GoHoliday/updater"
)

func main() {
    // With token
    syncer := updater.NewGitHubSyncerWithToken("ghp_your_token_here")
    
    // Without token (rate limited)
    syncer := updater.NewGitHubSyncer()
    
    // Validate token
    ctx := context.Background()
    if err := syncer.ValidateToken(ctx); err != nil {
        panic(err)
    }
}
```

## Usage Examples

### List Available Countries

```bash
# With token from environment
export GITHUB_TOKEN="ghp_your_token_here"
go run cmd/sync/main.go -list -verbose

# With token from flag
go run cmd/sync/main.go -token="ghp_your_token_here" -list -verbose
```

### Sync Specific Country

```bash
# Sync US holidays with authentication
go run cmd/sync/main.go -country=US -verbose -output=./holiday_data

# Dry run to see what would be synced
go run cmd/sync/main.go -country=US -dry-run -verbose
```

### Sync All Countries

```bash
# This will make many API calls - token recommended!
go run cmd/sync/main.go -verbose -output=./all_holidays
```

## Token Security Best Practices

### ✅ Do:
- Store tokens in environment variables, not in code
- Set reasonable expiration dates (30-90 days)
- Use fine-grained tokens with minimal permissions
- Rotate tokens regularly
- Use different tokens for different projects/environments

### ❌ Don't:
- Commit tokens to version control
- Share tokens in chat/email
- Use tokens with excessive permissions
- Set tokens to never expire
- Reuse the same token across multiple applications

## Troubleshooting

### Token Validation Errors

```bash
# Test token validation
go run examples/github_token_usage.go
```

**Common issues:**

1. **"invalid GitHub token: unauthorized"**
   - Token is malformed or expired
   - Generate a new token

2. **"GitHub token lacks required permissions"**
   - Token doesn't have `public_repo` or `Contents: Read` permission
   - Update token permissions or generate new token

3. **"token validation failed with status 403"**
   - Rate limit exceeded (even for validation)
   - Wait and try again, or check token permissions

### Rate Limiting

- **Unauthenticated**: 60 requests/hour per IP
- **Authenticated**: 5,000 requests/hour per user
- **Current usage**: Check headers in API responses or GitHub settings

### Environment Variable Not Found

```bash
# Check if token is set
echo $GITHUB_TOKEN  # Linux/macOS
echo %GITHUB_TOKEN% # Windows CMD
echo $env:GITHUB_TOKEN # Windows PowerShell
```

## Example Output

### With Token:
```
GoHolidays Python Sync Tool
===========================
Using authenticated GitHub API access
✓ GitHub token validated successfully
Found 100+ countries: [US GB CA AU DE]
```

### Without Token:
```
GoHolidays Python Sync Tool
===========================
Using unauthenticated GitHub API access (rate limited)
Found 100+ countries: [US GB CA AU DE]
```

## Integration with CI/CD

### GitHub Actions

```yaml
name: Sync Holiday Data
on:
  schedule:
    - cron: '0 2 * * 0'  # Weekly on Sunday at 2 AM

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Sync Holiday Data
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}  # Built-in token
        run: |
          go run cmd/sync/main.go -verbose -output=./data
      
      - name: Commit Changes
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add data/
          git commit -m "Update holiday data" || exit 0
          git push
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o sync cmd/sync/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sync .

# Token provided via environment variable
CMD ["./sync", "-verbose", "-output=/data"]
```

```bash
# Run with token
docker run -e GITHUB_TOKEN="ghp_your_token_here" -v $(pwd)/data:/data your-image
```
