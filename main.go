package main

import (
    "fmt"
    "image"
    "image/png"
    "os"
    "os/exec"
    "path/filepath"
    "github.com/atotto/clipboard"
    "github.com/vova616/screenshot"
)

const (
    version   = "1.0.0"
    maintainer = "Stevan Krstic <stevan@krstic.me>"
)

func showVersion() {
    fmt.Printf("screencolorpicker2 version %s\n", version)
    fmt.Printf("Maintainer: %s\n", maintainer)
    os.Exit(0)
}

func showHelp() {
    fmt.Println()
    fmt.Println("Usage: screencolorpicker2 [OPTION]")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("  --version   Show version information and exit")
    fmt.Println("  --help      Show this help message and exit")
    fmt.Println()
    os.Exit(0)
}

func getMousePosition() (int, int, error) {
    cmd := exec.Command("xdotool", "getmouselocation", "--shell")
    output, err := cmd.Output()
    if err != nil {
        return 0, 0, err
    }

    var x, y int
    _, err = fmt.Sscanf(string(output), "X=%d\nY=%d", &x, &y)
    if err != nil {
        return 0, 0, err
    }
    return x, y, nil
}

func main() {
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "--version":
            showVersion()
        case "--help":
            showHelp()
        default:
            fmt.Fprintf(os.Stderr, "Invalid option: %s\n", os.Args[1])
            showHelp()
        }
    }

    x, y, err := getMousePosition()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting mouse position: %v\n", err)
        os.Exit(1)
    }

    // Take screenshot
    img, err := screenshot.CaptureScreen()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error capturing screenshot: %v\n", err)
        os.Exit(1)
    }

    // Define the crop area
    rect := image.Rect(x, y, x+1, y+1)
    croppedImg := img.SubImage(rect)

    // Create a temporary screenshot
    imgFile, err := os.Create(filepath.Join(os.TempDir(), "screenshot.png"))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating temporary file: %v\n", err)
        os.Exit(1)
    }
    defer imgFile.Close()

    // Encode and save the cropped image
    err = png.Encode(imgFile, croppedImg)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
        os.Exit(1)
    }

    // Extract the color
    bounds := croppedImg.Bounds()
    color := croppedImg.At(bounds.Min.X, bounds.Min.Y)
    r, g, b, _ := color.RGBA()
    hexColor := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)

    fmt.Printf("Color at position (%d, %d): %s\n", x, y, hexColor)
    err = clipboard.WriteAll(hexColor)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error copying color to clipboard: %v\n", err)
        os.Exit(1)
    }

    // Delete temporary file
    os.Remove(filepath.Join(os.TempDir(), "screenshot.png"))
}
