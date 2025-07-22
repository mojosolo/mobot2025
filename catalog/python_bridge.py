#!/usr/bin/env python3
"""
Python Bridge for Go AEP Parser
Provides compatibility layer between existing Python cataloging system and new Go parser
"""

import json
import subprocess
import os
import sys
from pathlib import Path
from typing import Dict, List, Optional, Any
from datetime import datetime
import tempfile

# Add parent directory to path to import mobot modules
sys.path.insert(0, str(Path(__file__).parent.parent.parent / "mobot"))

try:
    from mobot import NexRenderSession
    MOBOT_AVAILABLE = True
except ImportError:
    MOBOT_AVAILABLE = False
    print("Warning: mobot module not available. Some features will be limited.")


class AEPCatalogBridge:
    """Bridge between Go parser and Python cataloging system"""
    
    def __init__(self, go_parser_path: Optional[str] = None):
        """Initialize the bridge with optional custom parser path"""
        self.go_parser_path = go_parser_path or self._find_go_parser()
        self.cache = {}  # Cache parsed results
        
    def _find_go_parser(self) -> str:
        """Locate the Go parser executable"""
        # Look for the compiled parser in various locations
        possible_paths = [
            Path(__file__).parent / "aep_parser",
            Path(__file__).parent / "bin" / "aep_parser",
            Path(__file__).parent.parent / "aep_parser",
            Path(__file__).parent.parent / "bin" / "aep_parser",
        ]
        
        for path in possible_paths:
            if path.exists():
                return str(path)
                
        # If not found, try to build it
        return self._build_go_parser()
        
    def _build_go_parser(self) -> str:
        """Build the Go parser if not already built"""
        parser_dir = Path(__file__).parent
        parser_main = parser_dir / "cmd" / "parser" / "main.go"
        
        if not parser_main.exists():
            # Create the parser command
            self._create_parser_command()
            
        # Build the parser
        output_path = parser_dir / "bin" / "aep_parser"
        output_path.parent.mkdir(exist_ok=True)
        
        try:
            subprocess.run(
                ["go", "build", "-o", str(output_path), str(parser_main)],
                check=True,
                cwd=str(parser_dir.parent)
            )
            return str(output_path)
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"Failed to build Go parser: {e}")
            
    def _create_parser_command(self):
        """Create the Go command-line parser"""
        cmd_dir = Path(__file__).parent / "cmd" / "parser"
        cmd_dir.mkdir(parents=True, exist_ok=True)
        
        main_go = cmd_dir / "main.go"
        main_go.write_text('''package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    
    "github.com/boltframe/aftereffects-aep-parser/catalog"
)

func main() {
    var (
        aepPath      = flag.String("file", "", "Path to AEP file")
        outputPath   = flag.String("output", "", "Output JSON path (default: stdout)")
        extractText  = flag.Bool("text", true, "Extract text layers")
        extractMedia = flag.Bool("media", true, "Extract media assets")
        deepAnalysis = flag.Bool("deep", true, "Perform deep analysis")
    )
    
    flag.Parse()
    
    if *aepPath == "" {
        log.Fatal("Please provide an AEP file path with -file")
    }
    
    // Create parser with options
    parser := catalog.NewParser()
    parser.ExtractText = *extractText
    parser.ExtractMedia = *extractMedia
    parser.DeepAnalysis = *deepAnalysis
    
    // Parse the project
    metadata, err := parser.ParseProject(*aepPath)
    if err != nil {
        log.Fatalf("Failed to parse project: %v", err)
    }
    
    // Convert to JSON
    jsonData, err := metadata.ToJSON()
    if err != nil {
        log.Fatalf("Failed to convert to JSON: %v", err)
    }
    
    // Output results
    if *outputPath != "" {
        err = os.WriteFile(*outputPath, jsonData, 0644)
        if err != nil {
            log.Fatalf("Failed to write output: %v", err)
        }
    } else {
        fmt.Println(string(jsonData))
    }
}
''')
        
    def parse_aep(self, aep_path: str, use_cache: bool = True) -> Dict[str, Any]:
        """Parse an AEP file using the Go parser"""
        aep_path = Path(aep_path).resolve()
        
        if not aep_path.exists():
            raise FileNotFoundError(f"AEP file not found: {aep_path}")
            
        # Check cache
        cache_key = str(aep_path)
        if use_cache and cache_key in self.cache:
            return self.cache[cache_key]
            
        # Create temporary output file
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as tmp:
            tmp_path = tmp.name
            
        try:
            # Run the Go parser
            result = subprocess.run(
                [self.go_parser_path, "-file", str(aep_path), "-output", tmp_path],
                capture_output=True,
                text=True,
                check=True
            )
            
            # Read the JSON output
            with open(tmp_path, 'r') as f:
                metadata = json.load(f)
                
            # Cache the result
            if use_cache:
                self.cache[cache_key] = metadata
                
            return metadata
            
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"Go parser failed: {e.stderr}")
        finally:
            # Clean up temp file
            if os.path.exists(tmp_path):
                os.unlink(tmp_path)
                
    def catalog_directory(self, directory: str, pattern: str = "*.aep") -> List[Dict[str, Any]]:
        """Catalog all AEP files in a directory"""
        directory = Path(directory)
        results = []
        
        for aep_file in directory.rglob(pattern):
            try:
                metadata = self.parse_aep(str(aep_file))
                results.append(metadata)
            except Exception as e:
                print(f"Error parsing {aep_file}: {e}")
                continue
                
        return results
        
    def convert_to_mobot_format(self, metadata: Dict[str, Any]) -> Dict[str, Any]:
        """Convert Go parser output to mobot catalog format"""
        # Map Go parser format to mobot's expected format
        mobot_format = {
            "template_path": metadata["file_path"],
            "template_name": metadata["file_name"],
            "analyzed_at": metadata["parsed_at"],
            
            # Basic info
            "compositions": len(metadata["compositions"]),
            "total_layers": sum(c["layer_count"] for c in metadata["compositions"]),
            
            # Capabilities (matching mobot's format)
            "capabilities": {
                "text_replacement": metadata["capabilities"]["has_text_replacement"],
                "image_replacement": metadata["capabilities"]["has_image_replacement"],
                "color_control": metadata["capabilities"]["has_color_control"],
                "audio_replacement": metadata["capabilities"]["has_audio_replacement"],
                "data_driven": metadata["capabilities"]["has_data_driven"],
                "expressions": metadata["capabilities"]["has_expressions"],
                "modular": metadata["capabilities"]["is_modular"],
            },
            
            # Categories and tags
            "categories": metadata["categories"],
            "tags": metadata["tags"],
            
            # Customizable elements
            "customizable_elements": {
                "text_layers": [
                    {
                        "name": t["layer_name"],
                        "default_text": t["source_text"],
                        "comp": t["comp_id"],
                    }
                    for t in metadata["text_layers"]
                ],
                "media_placeholders": [
                    {
                        "name": m["name"],
                        "type": m["type"],
                        "placeholder": m["is_placeholder"],
                    }
                    for m in metadata["media_assets"]
                    if m["is_placeholder"]
                ],
            },
            
            # Usage scenarios (from opportunities)
            "usage_scenarios": [
                {
                    "type": o["type"],
                    "description": o["description"],
                    "difficulty": o["difficulty"],
                    "impact": o["impact"],
                }
                for o in metadata["opportunities"]
            ],
            
            # Technical details
            "technical_details": {
                "bit_depth": metadata["bit_depth"],
                "expression_engine": metadata["expression_engine"],
                "effects_used": [e["name"] for e in metadata["effects"]],
                "resolutions": list(set(
                    f"{c['width']}x{c['height']}" 
                    for c in metadata["compositions"]
                )),
            },
        }
        
        return mobot_format
        
    def generate_catalog_report(self, catalog_data: List[Dict[str, Any]], output_path: str):
        """Generate a catalog report compatible with mobot's format"""
        report = {
            "catalog_version": "2.0",
            "generated_at": datetime.now().isoformat(),
            "total_templates": len(catalog_data),
            "parser": "Go-Python Bridge",
            "templates": catalog_data,
            
            # Summary statistics
            "summary": {
                "total_templates": len(catalog_data),
                "with_text_replacement": sum(
                    1 for t in catalog_data 
                    if t.get("capabilities", {}).get("text_replacement", False)
                ),
                "with_image_replacement": sum(
                    1 for t in catalog_data 
                    if t.get("capabilities", {}).get("image_replacement", False)
                ),
                "modular_templates": sum(
                    1 for t in catalog_data 
                    if t.get("capabilities", {}).get("modular", False)
                ),
            },
            
            # Category breakdown
            "categories": self._aggregate_categories(catalog_data),
            "tags": self._aggregate_tags(catalog_data),
        }
        
        # Save as JSON
        output_path = Path(output_path)
        with open(output_path, 'w') as f:
            json.dump(report, f, indent=2)
            
        # Also generate markdown report
        md_path = output_path.with_suffix('.md')
        self._generate_markdown_report(report, md_path)
        
    def _aggregate_categories(self, catalog_data: List[Dict[str, Any]]) -> Dict[str, int]:
        """Aggregate category counts"""
        categories = {}
        for template in catalog_data:
            for cat in template.get("categories", []):
                categories[cat] = categories.get(cat, 0) + 1
        return categories
        
    def _aggregate_tags(self, catalog_data: List[Dict[str, Any]]) -> Dict[str, int]:
        """Aggregate tag counts"""
        tags = {}
        for template in catalog_data:
            for tag in template.get("tags", []):
                tags[tag] = tags.get(tag, 0) + 1
        return tags
        
    def _generate_markdown_report(self, report: Dict[str, Any], output_path: Path):
        """Generate a human-readable markdown report"""
        md_content = f"""# AEP Template Catalog Report

Generated: {report['generated_at']}
Total Templates: {report['total_templates']}
Parser: {report['parser']}

## Summary

- Templates with text replacement: {report['summary']['with_text_replacement']}
- Templates with image replacement: {report['summary']['with_image_replacement']}
- Modular templates: {report['summary']['modular_templates']}

## Categories

"""
        
        for cat, count in sorted(report['categories'].items(), key=lambda x: x[1], reverse=True):
            md_content += f"- **{cat}**: {count} templates\n"
            
        md_content += "\n## Tags\n\n"
        
        for tag, count in sorted(report['tags'].items(), key=lambda x: x[1], reverse=True):
            md_content += f"- `{tag}`: {count} templates\n"
            
        md_content += "\n## Template Details\n\n"
        
        for template in report['templates']:
            md_content += f"### {template['template_name']}\n\n"
            md_content += f"- **Path**: `{template['template_path']}`\n"
            md_content += f"- **Categories**: {', '.join(template['categories'])}\n"
            md_content += f"- **Tags**: {', '.join(template['tags'])}\n"
            
            caps = template['capabilities']
            md_content += f"- **Capabilities**:\n"
            for cap, enabled in caps.items():
                if enabled:
                    md_content += f"  - ✅ {cap.replace('_', ' ').title()}\n"
                    
            if template.get('usage_scenarios'):
                md_content += f"- **Usage Scenarios**:\n"
                for scenario in template['usage_scenarios']:
                    md_content += f"  - {scenario['description']} ({scenario['difficulty']} difficulty, {scenario['impact']} impact)\n"
                    
            md_content += "\n"
            
        with open(output_path, 'w') as f:
            f.write(md_content)
            
    def integrate_with_nexrender(self, metadata: Dict[str, Any]) -> Dict[str, Any]:
        """Generate nexrender configuration from metadata"""
        # Use the Go parser's nexrender config generation
        # This is already implemented in the Go code
        return metadata.get("nexrender_config", {})


# Command-line interface
if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="AEP Catalog Bridge - Python interface to Go parser")
    parser.add_argument("command", choices=["parse", "catalog", "report"], 
                        help="Command to execute")
    parser.add_argument("path", help="Path to AEP file or directory")
    parser.add_argument("-o", "--output", help="Output file path")
    parser.add_argument("--mobot-format", action="store_true", 
                        help="Convert to mobot catalog format")
    
    args = parser.parse_args()
    
    bridge = AEPCatalogBridge()
    
    if args.command == "parse":
        # Parse single file
        metadata = bridge.parse_aep(args.path)
        
        if args.mobot_format:
            metadata = bridge.convert_to_mobot_format(metadata)
            
        if args.output:
            with open(args.output, 'w') as f:
                json.dump(metadata, f, indent=2)
        else:
            print(json.dumps(metadata, indent=2))
            
    elif args.command == "catalog":
        # Catalog directory
        catalog_data = bridge.catalog_directory(args.path)
        
        if args.mobot_format:
            catalog_data = [bridge.convert_to_mobot_format(m) for m in catalog_data]
            
        if args.output:
            with open(args.output, 'w') as f:
                json.dump(catalog_data, f, indent=2)
        else:
            print(f"Cataloged {len(catalog_data)} templates")
            
    elif args.command == "report":
        # Generate report
        catalog_data = bridge.catalog_directory(args.path)
        catalog_data = [bridge.convert_to_mobot_format(m) for m in catalog_data]
        
        output_path = args.output or "catalog_report.json"
        bridge.generate_catalog_report(catalog_data, output_path)
        print(f"Report generated: {output_path}")