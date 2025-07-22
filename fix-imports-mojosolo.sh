#!/bin/bash
# Fix all import statements to use mojosolo instead of yourusername

echo "Updating import statements from yourusername to mojosolo..."

# Find all Go files and update imports
find . -name "*.go" -type f -exec sed -i '' 's|github.com/mojosolo/mobot2025|github.com/mojosolo/mobot2025|g' {} \;

echo "✅ Import statements updated!"

# Also update documentation files
echo "Updating documentation..."
find . -name "*.md" -type f -exec sed -i '' 's|github.com/mojosolo/mobot2025|github.com/mojosolo/mobot2025|g' {} \;
find . -name "*.md" -type f -exec sed -i '' 's|mojosolo/mobot2025|mojosolo/mobot2025|g' {} \;

echo "✅ Documentation updated!"

# Update shell scripts
find . -name "*.sh" -type f -exec sed -i '' 's|mojosolo/mobot2025|mojosolo/mobot2025|g' {} \;

echo "✅ All references updated!"