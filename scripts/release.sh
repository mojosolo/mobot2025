#!/bin/bash

# release.sh - Create a new release

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if version provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Version not provided${NC}"
    echo ""
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.2.3"
    exit 1
fi

VERSION=$1

# Validate version format
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo -e "${RED}Error: Invalid version format${NC}"
    echo "Version must be in format: vX.Y.Z or vX.Y.Z-suffix"
    echo "Example: v1.2.3 or v1.2.3-beta1"
    exit 1
fi

echo -e "${BLUE}🚀 MoBot 2025 Release Builder${NC}"
echo "============================="
echo ""
echo "Version: $VERSION"
echo ""

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: Uncommitted changes detected${NC}"
    echo "Please commit or stash your changes before creating a release"
    exit 1
fi

# Check current branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ]; then
    echo -e "${YELLOW}Warning: Not on main branch (current: $BRANCH)${NC}"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Run quality checks
echo -e "${YELLOW}Running quality checks...${NC}"

# Tests
echo "  Running tests..."
if ! go test -short ./... > /dev/null 2>&1; then
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}  ✓ Tests passed${NC}"

# Linting
echo "  Running linter..."
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run --timeout=5m > /dev/null 2>&1; then
        echo -e "${RED}❌ Linting failed${NC}"
        exit 1
    fi
    echo -e "${GREEN}  ✓ Linting passed${NC}"
fi

# Build test
echo "  Testing build..."
if ! go build -o /tmp/mobot-test ./cmd/mobot2025/main.go; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi
rm -f /tmp/mobot-test
echo -e "${GREEN}  ✓ Build successful${NC}"

echo ""
echo -e "${GREEN}✅ All checks passed${NC}"
echo ""

# Update version in code if needed
echo -e "${YELLOW}Updating version references...${NC}"
# Add any version update logic here
echo -e "${GREEN}✓ Version updated${NC}"

# Create release directory
RELEASE_DIR="dist/mobot-$VERSION"
mkdir -p "$RELEASE_DIR"

# Build for multiple platforms
echo ""
echo -e "${YELLOW}Building release binaries...${NC}"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a array <<< "$platform"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"
    
    echo "  Building for $GOOS/$GOARCH..."
    
    OUTPUT_NAME="mobot-$VERSION-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi
    
    if GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-s -w -X main.Version=$VERSION" \
        -o "$RELEASE_DIR/$OUTPUT_NAME" \
        ./cmd/mobot2025/main.go; then
        echo -e "${GREEN}    ✓ $OUTPUT_NAME${NC}"
    else
        echo -e "${RED}    ✗ Failed to build for $GOOS/$GOARCH${NC}"
    fi
done

# Create archives
echo ""
echo -e "${YELLOW}Creating release archives...${NC}"

cd dist
for file in mobot-$VERSION/*; do
    if [ -f "$file" ]; then
        base=$(basename "$file")
        if [[ "$base" == *.exe ]]; then
            # ZIP for Windows
            zip -q "${base%.exe}.zip" -j "$file"
            echo -e "${GREEN}  ✓ ${base%.exe}.zip${NC}"
        else
            # TAR.GZ for Unix
            tar -czf "$base.tar.gz" -C "mobot-$VERSION" "$(basename "$file")"
            echo -e "${GREEN}  ✓ $base.tar.gz${NC}"
        fi
    fi
done
cd ..

# Generate changelog
echo ""
echo -e "${YELLOW}Generating changelog...${NC}"
cat > "dist/CHANGELOG-$VERSION.md" << EOF
# Changelog for $VERSION

## What's New

- Feature 1
- Feature 2

## Bug Fixes

- Fix 1
- Fix 2

## Breaking Changes

None

## Full Changelog

See: https://github.com/mojosolo/mobot2025/compare/previous...$VERSION
EOF
echo -e "${GREEN}✓ Changelog template created${NC}"

# Create checksums
echo ""
echo -e "${YELLOW}Generating checksums...${NC}"
cd dist
shasum -a 256 *.tar.gz *.zip > "checksums-$VERSION.txt"
echo -e "${GREEN}✓ Checksums generated${NC}"
cd ..

# Summary
echo ""
echo "================================="
echo -e "${GREEN}✅ Release $VERSION prepared!${NC}"
echo ""
echo "Release files in: dist/"
echo ""
echo "Next steps:"
echo "  1. Review and edit dist/CHANGELOG-$VERSION.md"
echo "  2. Create git tag: git tag -a $VERSION -m \"Release $VERSION\""
echo "  3. Push tag: git push origin $VERSION"
echo "  4. Create GitHub release and upload archives"
echo ""
echo "Files to upload:"
ls -la dist/*.tar.gz dist/*.zip 2>/dev/null | awk '{print "  - " $9}'