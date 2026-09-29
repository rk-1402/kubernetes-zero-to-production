#!/usr/bin/env python3
# Topic 49 — Liveness, Readiness, Startup Probes
# Separate /healthz (am I alive? no external deps) from /ready (can I serve? checks deps).
# Liveness should hit /healthz; readiness should hit /ready.
from fastapi import FastAPI, HTTPException

app = FastAPI()


class _DB:
    def is_connected(self) -> bool:
        # replace with a real connectivity check (pool ping, etc.)
        return True


db = _DB()


# /healthz — simple (no external deps). Used by livenessProbe + startupProbe.
@app.get("/healthz")
def healthz():
    return {"status": "alive"}


# /ready — checks dependencies. Used by readinessProbe so traffic stops if DB is down.
@app.get("/ready")
def ready():
    if not db.is_connected():
        raise HTTPException(503, "DB not ready")
    return {"status": "ready"}
