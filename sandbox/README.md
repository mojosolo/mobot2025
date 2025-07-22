# Sandbox Directory

This directory contains all temporary files, test scripts, and one-time use files that were previously cluttering the root directory. These files are kept for historical reference but are not part of the main codebase.

## Directory Structure

### `/scripts/`
Contains all shell scripts and utility scripts that were used for testing, building, or running various components:
- Shell scripts (*.sh) for running tests and demos
- Python scripts for utilities and fixes
- One-time use automation scripts

### `/tests/`
Contains test files that were in the root directory:
- `real_data_test.go` - Real data testing implementation

### `/temp/`
Contains temporary and backup files:
- Backup files (*.bak)
- Temporary work files
- Files created during development but not needed in production

### `/reports/`
Contains generated HTML reports from various test runs:
- Parser reports
- Detailed analysis reports  
- Enhanced UX reports

## Note
These files are preserved for reference but should NOT be used in production. They represent:
- Development experiments
- One-time testing utilities
- Historical test runs
- Temporary solutions that have been superseded

For current development, use the proper project structure:
- `/cmd/` for command-line tools
- `/tests/` for organized test suite
- `/docs/` for documentation
- `/demo/` for demonstration scripts