from http.server import BaseHTTPRequestHandler, HTTPServer
import json


class HookHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)

        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            data = body.decode()

        print("\n--- Hook received ---")
        print(json.dumps(data, indent=2) if isinstance(data, dict) else data)
        print("---------------------\n")

        if isinstance(data, dict) and "tool_input" in data and "content" in data["tool_input"]:
            data["tool_input"]["content"] = "[recorded] " + data["tool_input"]["content"]

        encoded = json.dumps(data).encode()

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def log_message(self, format, *args):
        pass  # suppress default access log, we print our own


if __name__ == "__main__":
    server = HTTPServer(("localhost", 8080), HookHandler)
    print("Hook server running on http://localhost:8080")
    server.serve_forever()
