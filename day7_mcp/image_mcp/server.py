"""
Image Processing MCP Server
Provides tools for Claude Code to process images via plugin system.
"""
import base64
import io
from pathlib import Path

from mcp.server.fastmcp import FastMCP
from PIL import Image, ImageFilter, ImageEnhance

mcp = FastMCP("image-processor")


@mcp.tool()
def resize_image(image_path: str, width: int, height: int, output_path: str) -> str:
    """Resize an image to specified dimensions."""
    img = Image.open(image_path)
    resized = img.resize((width, height), Image.LANCZOS)
    resized.save(output_path)
    return f"Resized {image_path} to {width}x{height} -> {output_path}"


@mcp.tool()
def convert_to_grayscale(image_path: str, output_path: str) -> str:
    """Convert an image to grayscale."""
    img = Image.open(image_path).convert("L")
    img.save(output_path)
    return f"Converted {image_path} to grayscale -> {output_path}"


@mcp.tool()
def apply_blur(image_path: str, radius: float, output_path: str) -> str:
    """Apply Gaussian blur to an image."""
    img = Image.open(image_path)
    blurred = img.filter(ImageFilter.GaussianBlur(radius=radius))
    blurred.save(output_path)
    return f"Applied blur (radius={radius}) to {image_path} -> {output_path}"


@mcp.tool()
def adjust_brightness(image_path: str, factor: float, output_path: str) -> str:
    """
    Adjust image brightness.
    factor: 0.0=black, 1.0=original, 2.0=double brightness
    """
    img = Image.open(image_path)
    enhanced = ImageEnhance.Brightness(img).enhance(factor)
    enhanced.save(output_path)
    return f"Adjusted brightness (factor={factor}) -> {output_path}"


@mcp.tool()
def get_image_info(image_path: str) -> dict:
    """Get metadata about an image (size, mode, format)."""
    img = Image.open(image_path)
    return {
        "path": image_path,
        "width": img.width,
        "height": img.height,
        "mode": img.mode,
        "format": img.format,
    }


@mcp.tool()
def image_to_base64(image_path: str) -> str:
    """Encode an image file as a base64 string (for inline preview)."""
    with open(image_path, "rb") as f:
        return base64.b64encode(f.read()).decode()


if __name__ == "__main__":
    mcp.run()
