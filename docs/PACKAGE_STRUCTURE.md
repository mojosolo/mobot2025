# Package Structure and Boundaries

This document defines the package organization and boundaries for MoBot 2025.

## Package Hierarchy

```
mobot2025/
├── cmd/                    # Executable commands
├── pkg/                    # Public packages (importable by external projects)
│   └── types/             # Shared type definitions
├── internal/              # Private packages (not importable externally)
│   ├── errors/           # Error handling utilities
│   ├── validation/       # Input validation
│   └── utils/            # General utilities
├── catalog/               # AEP cataloging and automation
├── test/                  # Test utilities and fixtures
│   └── fixtures/         # Test data generators
└── tests/                 # Integration and E2E tests
```

## Package Responsibilities

### Root Package (`github.com/mojosolo/mobot2025`)
**Purpose**: Core AEP parsing functionality
**Exports**: 
- `Open()` - Parse AEP files
- `Project` - AEP project structure
- Core types (Item, Layer, Property)

**Dependencies**: None (self-contained)

### `/cmd` - Command Line Tools
**Purpose**: Executable binaries
**Structure**: One subdirectory per command
- `cmd/mobot2025/` - Main CLI
- `cmd/parser/` - Standalone parser
- `cmd/scanner/` - Directory scanner

**Rules**:
- Each subdirectory must have single `main.go`
- Import from root and catalog packages only
- No business logic, only CLI handling

### `/catalog` - Cataloging System
**Purpose**: Advanced features built on core parser
**Exports**:
- Database operations
- Multi-agent orchestration
- Automation scoring
- Search functionality

**Dependencies**:
- Root package for parsing
- Internal packages for utilities
- External: database/sql, encoding/json

### `/pkg` - Public Packages
**Purpose**: Reusable code for external projects
**Current**: Empty (future use)

**Future packages**:
- `pkg/types/` - Shared type definitions
- `pkg/client/` - API client library

### `/internal` - Private Packages
**Purpose**: Shared internal utilities

#### `/internal/errors`
- Common error types
- Error wrapping utilities
- Error categorization

#### `/internal/validation`
- Input validation functions
- File path validation
- Data constraint checking

#### `/internal/utils`
- File operations
- String utilities
- Time formatting
- Generic helpers

**Rules**:
- Not importable outside mobot2025
- No circular dependencies
- Keep focused on specific utilities

### `/test/fixtures`
**Purpose**: Test data generation
**Exports**:
- Mock AEP generators
- Test data factories
- Test environment helpers

**Rules**:
- Used by tests only
- No production dependencies

### `/tests`
**Purpose**: Integration and E2E tests
**Structure**:
- `tests/integration/` - Cross-package tests
- `tests/e2e/` - Full system tests
- `tests/helpers/` - Test utilities

## Import Rules

### Allowed Dependencies

```
┌─────────────┐
│     cmd     │ ──imports──> catalog, root
└─────────────┘
      │
      ▼
┌─────────────┐
│   catalog   │ ──imports──> root, internal/*
└─────────────┘
      │
      ▼
┌─────────────┐
│    root     │ ──imports──> (none)
└─────────────┘
      │
      ▼
┌─────────────┐
│  internal   │ ──imports──> (minimal std lib only)
└─────────────┘
```

### Forbidden Dependencies

- ❌ Root package → catalog (would create cycle)
- ❌ Internal → root or catalog (keep utilities generic)
- ❌ Test fixtures → production code
- ❌ Circular imports of any kind

## Guidelines

### 1. Package Cohesion
- Each package should have single, clear purpose
- Related functionality stays together
- Unrelated functionality gets separated

### 2. Dependency Direction
- Dependencies flow downward only
- Higher-level packages import lower-level
- No circular dependencies

### 3. Interface Boundaries
- Define interfaces in consumer package
- Implement in provider package
- Minimize exported surface area

### 4. Testing
- Unit tests stay with code (`*_test.go`)
- Integration tests in `/tests/integration/`
- Test utilities in `/test/fixtures/`

### 5. Naming Conventions
- Package names are singular nouns
- Short, descriptive names
- No redundant prefixes (not `utilsPackage`)

## Migration Plan

To achieve this structure:

1. **Phase 1**: Create `/pkg/types/` (if needed)
   - Move shared types from root/catalog
   - Update imports

2. **Phase 2**: Refactor utilities
   - Use new internal packages
   - Remove duplicate code

3. **Phase 3**: Clarify boundaries
   - Move misplaced code
   - Update documentation

## Examples

### Good Package Design

```go
// internal/validation/validation.go
package validation

// Focused on validation only
func ValidateFilePath(path string) error { }
func ValidateAEPFile(path string) error { }
```

### Bad Package Design

```go
// utils/everything.go
package utils

// Too many responsibilities
func ValidateFile() { }
func ParseJSON() { }
func SendEmail() { }
func CalculateTax() { }
```

## Package Documentation

Each package must have:
1. Package comment describing purpose
2. Clear examples in `example_test.go`
3. README.md for complex packages
4. Godoc for all exported items

## Review Checklist

When adding code, verify:
- [ ] Code is in correct package
- [ ] Dependencies flow correctly
- [ ] No circular imports
- [ ] Package has single purpose
- [ ] Exports are minimized
- [ ] Tests are properly located