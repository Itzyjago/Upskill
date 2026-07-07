"""Summarizes wordcount's structured JSON access log (middleware.go's one
`slog.Info("request", ...)` line per request) by route: request count, error
count, and p50/p95 latency. Stdlib only -- json, dataclasses, pathlib,
argparse, statistics, the actual "Python stdlib fluency" roadmap item.

    python scripts/parse_access_log.py scripts/testdata/wordcount_access_sample.log
"""

from __future__ import annotations

import argparse
import json
import statistics
import sys
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Iterator, TextIO


@dataclass(frozen=True)
class RequestLogEntry:
    method: str
    path: str
    status: int
    dur_ms: int
    trace_id: str


def parse_line(line: str) -> RequestLogEntry | None:
    """Parses one log line into a RequestLogEntry, or None if the line isn't
    a request log line -- wordcount's stdout can carry other slog lines
    (startup messages, otlp export failures) that aren't request records."""
    line = line.strip()
    if not line:
        return None
    try:
        record = json.loads(line)
    except json.JSONDecodeError:
        return None
    if record.get("msg") != "request":
        return None
    return RequestLogEntry(
        method=record["method"],
        path=record["path"],
        status=record["status"],
        dur_ms=record["dur_ms"],
        trace_id=record["trace_id"],
    )


def parse_lines(lines: Iterable[str]) -> Iterator[RequestLogEntry]:
    for line in lines:
        entry = parse_line(line)
        if entry is not None:
            yield entry


@dataclass
class RouteSummary:
    count: int = 0
    error_count: int = 0
    durations_ms: list[int] | None = None

    def __post_init__(self) -> None:
        if self.durations_ms is None:
            self.durations_ms = []

    @property
    def error_rate(self) -> float:
        return self.error_count / self.count if self.count else 0.0

    def percentile(self, pct: float) -> float:
        """pct in (0, 100). statistics.quantiles needs >=2 points; a
        single-sample route just reports that sample for every percentile."""
        if not self.durations_ms:
            return 0.0
        if len(self.durations_ms) == 1:
            return float(self.durations_ms[0])
        # n=100 buckets -> quantiles[k] is roughly the (k+1)th percentile.
        cuts = statistics.quantiles(self.durations_ms, n=100, method="inclusive")
        index = max(0, min(len(cuts) - 1, round(pct) - 1))
        return cuts[index]


def summarize(entries: Iterable[RequestLogEntry]) -> dict[str, RouteSummary]:
    summaries: dict[str, RouteSummary] = defaultdict(RouteSummary)
    for e in entries:
        s = summaries[e.path]
        s.count += 1
        if e.status >= 500:
            s.error_count += 1
        s.durations_ms.append(e.dur_ms)
    return dict(summaries)


def format_report(summaries: dict[str, RouteSummary]) -> str:
    lines = []
    for path in sorted(summaries):
        s = summaries[path]
        lines.append(
            f"{path:<12} count={s.count:<4} errors={s.error_count:<3} "
            f"error_rate={s.error_rate:.1%} p50={s.percentile(50):.0f}ms "
            f"p95={s.percentile(95):.0f}ms"
        )
    return "\n".join(lines)


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "logfile",
        nargs="?",
        type=Path,
        help="path to a wordcount access log (default: read stdin)",
    )
    args = parser.parse_args(argv)

    handle: TextIO
    if args.logfile is None:
        handle = sys.stdin
    else:
        handle = args.logfile.open(encoding="utf-8")

    try:
        entries = list(parse_lines(handle))
    finally:
        if handle is not sys.stdin:
            handle.close()

    if not entries:
        print("no request log lines found", file=sys.stderr)
        return 1

    print(format_report(summarize(entries)))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
