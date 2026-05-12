from mcp.server.fastmcp import FastMCP
import uvicorn
import math
import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(levelname)s | %(message)s"
)
logger = logging.getLogger("math-mcp")

mcp = FastMCP("Math MCP Server")


@mcp.tool()
def add(a: float, b: float) -> float:
    """Add two numbers"""
    result = a + b
    logger.info(f"add({a}, {b}) = {result}")
    return result


@mcp.tool()
def subtract(a: float, b: float) -> float:
    """Subtract b from a"""
    result = a - b
    logger.info(f"subtract({a}, {b}) = {result}")
    return result


@mcp.tool()
def multiply(a: float, b: float) -> float:
    """Multiply two numbers"""
    result = a * b
    logger.info(f"multiply({a}, {b}) = {result}")
    return result


@mcp.tool()
def divide(a: float, b: float) -> float:
    """Divide a by b"""
    if b == 0:
        raise ValueError("Cannot divide by zero")
    result = a / b
    logger.info(f"divide({a}, {b}) = {result}")
    return result


@mcp.tool()
def sqrt(x: float) -> float:
    """Calculate the square root of x"""
    if x < 0:
        raise ValueError("x must be >= 0")
    result = math.sqrt(x)
    logger.info(f"sqrt({x}) = {result}")
    return result


@mcp.tool()
def power(base: float, exponent: float) -> float:
    """Calculate base raised to the power of exponent"""
    result = math.pow(base, exponent)
    logger.info(f"power({base}, {exponent}) = {result}")
    return result


# Module-level app so uvicorn can reference it via "server:app".
# The Starlette app returned here carries the SessionManager lifespan,
# which uvicorn will start automatically — this is the critical fix.
app = mcp.streamable_http_app()


if __name__ == "__main__":
    uvicorn.run(
        "server:app",
        host="0.0.0.0",
        port=8000,
        reload=True,
    )
