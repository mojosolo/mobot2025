#!/usr/bin/env python3
"""
Test suite for Python Bridge integration with Go AEP parser
"""

import unittest
import tempfile
import json
import os
import sys
from pathlib import Path
from unittest.mock import Mock, patch, MagicMock

# Add the catalog directory to Python path
sys.path.insert(0, str(Path(__file__).parent.parent.parent / "catalog"))

from python_bridge import AEPCatalogBridge


class TestAEPCatalogBridge(unittest.TestCase):
    """Test cases for AEP Catalog Bridge"""
    
    def setUp(self):
        """Set up test environment"""
        self.temp_dir = tempfile.mkdtemp()
        self.bridge = AEPCatalogBridge()
        
    def tearDown(self):
        """Clean up after tests"""
        import shutil
        shutil.rmtree(self.temp_dir, ignore_errors=True)
        
    def test_bridge_initialization(self):
        """Test bridge initializes correctly"""
        self.assertIsNotNone(self.bridge)
        self.assertIsInstance(self.bridge.cache, dict)
        
    def test_find_go_parser(self):
        """Test Go parser discovery"""
        parser_path = self.bridge._find_go_parser()
        # Should either find existing parser or attempt to build
        self.assertIsNotNone(parser_path)
        
    @patch('subprocess.run')
    def test_parse_aep_file_success(self, mock_run):
        """Test successful AEP file parsing"""
        # Mock successful Go parser output
        mock_output = {
            "name": "Test Project",
            "version": "2021",
            "compositions": 3,
            "text_layers": [
                {"text": "Title", "layer_name": "title_layer"}
            ]
        }
        
        mock_run.return_value = MagicMock(
            returncode=0,
            stdout=json.dumps(mock_output)
        )
        
        result = self.bridge.parse_aep_file("test.aep")
        
        self.assertEqual(result["name"], "Test Project")
        self.assertEqual(len(result["text_layers"]), 1)
        
    @patch('subprocess.run')
    def test_parse_aep_file_error(self, mock_run):
        """Test AEP parsing error handling"""
        # Mock Go parser error
        mock_run.return_value = MagicMock(
            returncode=1,
            stderr="Error: Invalid AEP file"
        )
        
        with self.assertRaises(Exception) as context:
            self.bridge.parse_aep_file("invalid.aep")
            
        self.assertIn("Invalid AEP", str(context.exception))
        
    def test_cache_functionality(self):
        """Test caching of parsed results"""
        test_file = "test_cache.aep"
        test_result = {"name": "Cached Project"}
        
        # Store in cache
        self.bridge.cache[test_file] = test_result
        
        # Should return cached result without parsing
        with patch.object(self.bridge, '_run_go_parser') as mock_parser:
            result = self.bridge.get_cached_or_parse(test_file)
            
        # Parser should not be called
        mock_parser.assert_not_called()
        self.assertEqual(result["name"], "Cached Project")
        
    def test_batch_processing(self):
        """Test batch processing of multiple files"""
        test_files = ["file1.aep", "file2.aep", "file3.aep"]
        
        with patch.object(self.bridge, 'parse_aep_file') as mock_parse:
            mock_parse.return_value = {"name": "Test"}
            
            results = self.bridge.batch_process(test_files)
            
        self.assertEqual(len(results), 3)
        self.assertEqual(mock_parse.call_count, 3)
        
    def test_export_to_catalog_format(self):
        """Test export to legacy catalog format"""
        parsed_data = {
            "name": "Export Test",
            "version": "2021",
            "compositions": [
                {"name": "Comp1", "width": 1920, "height": 1080}
            ],
            "text_layers": [
                {"text": "Hello", "comp_name": "Comp1"}
            ]
        }
        
        catalog_format = self.bridge.export_to_catalog_format(parsed_data)
        
        # Verify format conversion
        self.assertIn("project_info", catalog_format)
        self.assertIn("compositions", catalog_format)
        self.assertIn("text_elements", catalog_format)
        
    @patch('subprocess.run')
    def test_build_go_parser(self, mock_run):
        """Test Go parser building"""
        mock_run.return_value = MagicMock(returncode=0)
        
        parser_path = self.bridge._build_go_parser()
        
        # Should call go build
        mock_run.assert_called()
        call_args = mock_run.call_args[0][0]
        self.assertIn("go", call_args[0])
        self.assertIn("build", call_args)
        
    def test_error_recovery(self):
        """Test error recovery mechanisms"""
        # Test with non-existent file
        result = self.bridge.safe_parse("non_existent.aep")
        self.assertIsNone(result)
        
        # Test with corrupted data
        with patch.object(self.bridge, '_run_go_parser') as mock_parser:
            mock_parser.side_effect = Exception("Parser crashed")
            
            result = self.bridge.safe_parse("corrupted.aep")
            self.assertIsNone(result)
            
    def test_mobot_integration(self):
        """Test integration with mobot module when available"""
        if not self.bridge.MOBOT_AVAILABLE:
            self.skipTest("mobot module not available")
            
        # Test mobot session creation
        with patch('mobot.NexRenderSession') as mock_session:
            session = self.bridge.create_mobot_session()
            mock_session.assert_called_once()
            
    def test_data_transformation(self):
        """Test data transformation utilities"""
        # Test color conversion
        rgb = [255, 128, 0]
        hex_color = self.bridge.rgb_to_hex(rgb)
        self.assertEqual(hex_color, "#FF8000")
        
        # Test timestamp conversion
        timestamp = "2024-01-15T10:30:00Z"
        parsed = self.bridge.parse_timestamp(timestamp)
        self.assertEqual(parsed.year, 2024)
        self.assertEqual(parsed.month, 1)
        
    def test_performance_metrics(self):
        """Test performance tracking"""
        with patch.object(self.bridge, 'parse_aep_file') as mock_parse:
            mock_parse.return_value = {"name": "Test"}
            
            # Enable performance tracking
            self.bridge.enable_performance_tracking()
            
            # Parse file
            self.bridge.parse_aep_file("perf_test.aep")
            
            # Check metrics
            metrics = self.bridge.get_performance_metrics()
            self.assertIn("parse_count", metrics)
            self.assertIn("avg_parse_time", metrics)


class TestBridgeIntegration(unittest.TestCase):
    """Integration tests for Python Bridge with Go parser"""
    
    @classmethod
    def setUpClass(cls):
        """Set up integration test environment"""
        cls.test_aep_path = Path(__file__).parent.parent.parent / "data" / "ExEn-js.aep"
        if not cls.test_aep_path.exists():
            raise unittest.SkipTest("Test AEP file not found")
            
    def test_real_file_parsing(self):
        """Test parsing real AEP file"""
        bridge = AEPCatalogBridge()
        
        try:
            result = bridge.parse_aep_file(str(self.test_aep_path))
            
            # Verify result structure
            self.assertIsInstance(result, dict)
            self.assertIn("name", result)
            self.assertIn("version", result)
            
        except Exception as e:
            # If Go parser not available, skip
            if "parser not found" in str(e).lower():
                self.skipTest("Go parser not available")
            else:
                raise
                
    def test_concurrent_parsing(self):
        """Test concurrent file parsing"""
        import concurrent.futures
        
        bridge = AEPCatalogBridge()
        test_files = [self.test_aep_path] * 5  # Parse same file 5 times
        
        with concurrent.futures.ThreadPoolExecutor(max_workers=3) as executor:
            futures = [
                executor.submit(bridge.safe_parse, str(f))
                for f in test_files
            ]
            
            results = [f.result() for f in futures]
            
        # All should succeed or all should fail
        success_count = sum(1 for r in results if r is not None)
        self.assertIn(success_count, [0, 5])


def run_tests():
    """Run all tests"""
    unittest.main()


if __name__ == "__main__":
    run_tests()