# Changelog

## [v1.0.1] - 2025-06-17

### Added
- **New architecture**: Complete refactoring of the codebase with improved structure
- **Common module** (`common.go`): Centralized shared functionality including:
  - Configuration struct for sensor instances
  - Reading struct for sensor data representation
  - Command definitions for all sensor operations
  - Utility functions (`calculateChecksum`, `byteToInt`, `toHex`, `writeAndRead`, `write`)
  - Comprehensive package-level documentation with usage examples
  - Custom error types (`ErrChecksumMismatch`, `ErrInvalidFrame`, `ErrSensorCommunication`)
- **Interface definition** (`interface.go`): Clean `SensorInterface` with well-defined methods
- **Comprehensive test coverage**: 
  - `common_test.go` with utility function tests
  - Enhanced `zh07i_test.go` and `zh07q_test.go` with better test scenarios
  - Improved test patterns and helper functions
- **GitHub Actions CI/CD**: Automated testing, formatting checks, and linting
- **Enhanced documentation**: Godoc comments for all exported types and functions

### Changed
- **BREAKING**: Complete API redesign from single `ZH07` struct to separate `ZH07i` and `ZH07q` implementations
- **BREAKING**: Constructor pattern changed from `NewZH07(mode, rw)` to `NewZH07i(config)` and `NewZH07q(config)`
- **BREAKING**: Initialization now requires explicit `Init()` call after construction
- **BREAKING**: Reading struct field units corrected from `[g/m3]` to `[μg/m³]`
- **Improved error handling**: Structured error types with better context and wrapping
- **Enhanced sensor implementations**:
  - `zh07i.go`: Refactored initiative upload mode with dependency injection
  - `zh07q.go`: Refactored question & answer mode with improved error handling
- **Updated dependencies**: Go 1.23 support with testify for testing
- **README improvements**: Updated examples and documentation to reflect new API

### Removed
- **Legacy `zh07.go`**: Removed monolithic implementation in favor of modular approach
- **Example cleanup**: Removed temporary example files and build artifacts
- **Simplified dependencies**: Cleaned up unused dependencies in go.sum

### Technical Improvements
- **Dependency injection**: Functions are now injectable for better testability
- **Better separation of concerns**: Clear distinction between initiative and Q&A modes
- **Consistent error handling**: Structured error types across all implementations
- **Enhanced maintainability**: Modular architecture with clear interfaces
- **Improved test coverage**: More comprehensive test scenarios and edge cases
- **Documentation standards**: Full godoc compliance for all public APIs

### Migration Guide
To upgrade from v1.0.0 to v1.0.1:

**Before (v1.0.0)**:
```go
z, err := zh07.NewZH07(zh07.ModeQA, rw)
reading, err := z.Read()
```

**After (v1.0.1)**:
```go
z := zh07.NewZH07q(&zh07.Config{RW: rw})
if err := z.Init(); err != nil {
    // handle error
}
reading, err := z.Read()
```

---

## [v1.0.0] - 2025-06-17

### Added
- Initial implementation of ZH07 sensor library
- Support for both initiative upload and question & answer communication modes
- Basic ZH07, ZH07I, and ZH07Q sensor support
- Initial documentation and README
- MIT License

### Changed
- Corrected readme file name

---

**Note**: v1.0.1 represents a major architectural improvement over v1.0.0, providing better testability, maintainability, and adherence to Go best practices. While it introduces breaking changes, the new API is more intuitive and follows established Go patterns.