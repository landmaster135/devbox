# EXIF Viewer Library

This library provides functionality to extract and display EXIF data from image files including JPEG, PNG, and TIFF formats.

## Features

- Extract EXIF data from JPEG, PNG, and TIFF image files
- Support for recursive directory scanning
- Configurable property filtering
- Table-formatted output display
- PNG metadata extraction (chunks, color type, etc.)
- File information extraction (size, modification date, etc.)

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    exif_viewer "github.com/nov/devbox/internal/exif_viewer"
)

func main() {
    // Create a new configuration
    config := exif_viewer.NewConfig()
    config.Directory = "./photos"
    config.Recursive = true
    config.Verbose = true

    // Process images and extract EXIF data
    exifDataList, err := exif_viewer.ProcessImages(config)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    // Display results in table format
    exif_viewer.DisplayExifTable(exifDataList, config)
}
```

### Configuration Options

```go
config := exif_viewer.NewConfig()

// Set target directory
config.Directory = "./images"

// Set file extensions (default: jpg, jpeg, tiff, tif, png)
config.SetExtensions("jpg,png")

// Set specific properties to display
config.SetProperties("DateTime,Camera,ImageWidth,ImageHeight")

// Limit number of properties displayed
config.MaxProps = 10

// Enable verbose logging
config.Verbose = true

// Enable recursive directory scanning
config.Recursive = true
```

### Extract EXIF from a Single File

```go
config := exif_viewer.NewConfig()
exifData, err := exif_viewer.ExtractSingleFileExif("photo.jpg", config)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("File: %s\n", exifData.FilePath)
for key, value := range exifData.Properties {
    fmt.Printf("%s: %s\n", key, value)
}
```

### Using the Service Layer

```go
import (
    "github.com/nov/devbox/internal/exif_viewer/usecases"
)

func main() {
    service := usecases.NewExifViewerService()
    config := service.CreateDefaultConfig()
    config.Directory = "./photos"

    // Validate directory exists
    err := service.ValidateDirectory(config.Directory)
    if err != nil {
        fmt.Printf("Directory validation failed: %v\n", err)
        return
    }

    // Process images
    exifDataList, err := service.ProcessImages(config)
    if err != nil {
        fmt.Printf("Error processing images: %v\n", err)
        return
    }

    // Display results
    service.DisplayResults(exifDataList, config)
}
```

## Package Structure

- `types.go` - Core data structures (Config, ExifData)
- `api.go` - Main API functions and configuration helpers
- `extractor.go` - EXIF extraction logic for different file formats
- `file_finder.go` - File discovery and filtering
- `display.go` - Table formatting and output display
- `utils.go` - Utility functions and helpers
- `exifutil/` - Wrapper package for backward compatibility
- `usecases/` - Business logic and service layer

## Supported File Formats

- **JPEG/JPG** - Full EXIF data extraction
- **PNG** - Metadata chunks (IHDR, pHYs, gAMA, sRGB) and file information
- **TIFF/TIF** - EXIF data extraction

## Properties Extracted

### Common Properties (All Formats)
- File Name
- File Size
- File Modification Date/Time
- Directory
- File Type
- File Type Extension
- MIME Type

### JPEG-Specific Properties
- Camera make and model
- Exposure settings
- GPS coordinates
- Timestamp information
- Orientation
- And many more EXIF tags

### PNG-Specific Properties
- Image Width/Height
- Bit Depth
- Color Type
- Compression
- Filter
- Interlace
- Gamma
- Pixel density information

## Error Handling

The library is designed to be robust and will continue processing even if individual files cannot be read. Errors are logged when verbose mode is enabled, but they won't stop the overall processing.

## Dependencies

- `github.com/dsoprea/go-exif/v3` - EXIF data parsing
- `github.com/dsoprea/go-jpeg-image-structure/v2` - JPEG structure parsing

## Testing

Run tests with:

```bash
go test ./...
```

## CLI Tool

This library is also used by the CLI tool located at `~/devbox/cmd/cli/exif-viewer/`. The CLI provides command-line access to all the library functionality.
