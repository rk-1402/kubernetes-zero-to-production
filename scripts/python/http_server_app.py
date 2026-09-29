#!/usr/bin/env python3
# Topic 2 — Containers vs VMs vs Kubernetes
# Minimal HTTP service used to demonstrate Docker -> Kubernetes flow.
# Build:  docker build -t hello:1 .
# Run:    docker run --rm -p 8080:8080 hello:1  ->  curl localhost:8080
from http.server import BaseHTTPRequestHandler, HTTPServer


class H(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"hello from container\n")


if __name__ == "__main__":
    HTTPServer(("", 8080), H).serve_forever()
