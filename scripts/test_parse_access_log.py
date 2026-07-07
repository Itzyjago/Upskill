"""python -m unittest scripts.test_parse_access_log"""

import unittest
from pathlib import Path

from parse_access_log import RouteSummary, parse_line, parse_lines, summarize

FIXTURE = Path(__file__).parent / "testdata" / "wordcount_access_sample.log"


class ParseLineTests(unittest.TestCase):
    def test_parses_a_real_request_line(self):
        line = (
            '{"time":"2026-07-08T01:02:25.88Z","level":"INFO","msg":"request",'
            '"method":"POST","path":"/count","status":200,"dur_ms":3,'
            '"trace_id":"abc","span_id":"def"}'
        )
        entry = parse_line(line)
        self.assertIsNotNone(entry)
        self.assertEqual(entry.method, "POST")
        self.assertEqual(entry.path, "/count")
        self.assertEqual(entry.status, 200)
        self.assertEqual(entry.dur_ms, 3)

    def test_skips_non_request_log_lines(self):
        line = '{"time":"2026-07-08T01:02:25Z","level":"INFO","msg":"otlp export enabled","endpoint":"x"}'
        self.assertIsNone(parse_line(line))

    def test_skips_malformed_json(self):
        self.assertIsNone(parse_line("not json at all"))

    def test_skips_blank_lines(self):
        self.assertIsNone(parse_line("   "))
        self.assertIsNone(parse_line(""))


class SummarizeTests(unittest.TestCase):
    def test_against_the_real_captured_fixture(self):
        """testdata/wordcount_access_sample.log is real stdout captured from
        a running `wc -serve`, hit with real curl requests including a live
        Idempotency-Key replay -- not synthesized."""
        with FIXTURE.open(encoding="utf-8") as f:
            entries = list(parse_lines(f))

        # 4 successful POST /count (including the idempotent replay, which
        # the server still logs as a request), 1 rejected GET /count, 1
        # healthz -- 6 requests total, matching what the capture session ran.
        self.assertEqual(len(entries), 6)

        summaries = summarize(entries)
        self.assertEqual(set(summaries), {"/count", "/healthz"})
        self.assertEqual(summaries["/count"].count, 5)
        self.assertEqual(summaries["/healthz"].count, 1)
        # The 405 from the GET isn't a 5xx -- error_count counts server
        # errors, not any non-200, so it must stay 0 here.
        self.assertEqual(summaries["/count"].error_count, 0)

    def test_error_rate_only_counts_5xx(self):
        entries = [
            _entry(status=200),
            _entry(status=404),
            _entry(status=500),
            _entry(status=503),
        ]
        s = summarize(entries)["/count"]
        self.assertEqual(s.count, 4)
        self.assertEqual(s.error_count, 2)
        self.assertAlmostEqual(s.error_rate, 0.5)

    def test_percentile_with_a_single_sample_returns_that_sample(self):
        s = RouteSummary()
        s.count = 1
        s.durations_ms = [42]
        self.assertEqual(s.percentile(50), 42.0)
        self.assertEqual(s.percentile(95), 42.0)

    def test_percentile_with_no_samples_is_zero(self):
        s = RouteSummary()
        self.assertEqual(s.percentile(95), 0.0)


def _entry(status: int, path: str = "/count", dur_ms: int = 0):
    from parse_access_log import RequestLogEntry

    return RequestLogEntry(method="POST", path=path, status=status, dur_ms=dur_ms, trace_id="t")


if __name__ == "__main__":
    unittest.main()
