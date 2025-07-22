#!/usr/bin/env python3
import re

# Read the file
with open('catalog/implementation_agent.go', 'r') as f:
    content = f.read()

# Pattern to match escaped backticks in struct tags
pattern = r'\\`json:"([^"]+)"\\`'
replacement = r'` + "`json:\"\1\"`" + `'

# Replace all occurrences
fixed_content = re.sub(pattern, replacement, content)

# Write back
with open('catalog/implementation_agent.go', 'w') as f:
    f.write(fixed_content)

print("Fixed all JSON struct tags")