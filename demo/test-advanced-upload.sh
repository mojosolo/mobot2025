#!/bin/bash

# Test script to check Advanced Mode upload response

echo "Testing Advanced Mode upload..."

# Upload the sample AEP file and save response
curl -X POST \
  -F "file=@sample-aep/Ai Text Intro.aep" \
  http://localhost:8081/upload/advanced \
  | python3 -m json.tool > advanced-response.json

echo "Response saved to advanced-response.json"
echo ""
echo "Text layers found:"
cat advanced-response.json | python3 -c "
import json, sys
data = json.load(sys.stdin)
if 'textLayers' in data:
    for i, layer in enumerate(data['textLayers']):
        print(f'{i+1}. Text: {layer.get(\"text\", \"N/A\")}')
        print(f'   Layer: {layer.get(\"layerName\", \"N/A\")}')
        print(f'   Comp: {layer.get(\"compName\", \"N/A\")}')
        print()
"