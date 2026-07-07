"""Runs notes/regex.md's lookaround examples for real against Python's `re`
(PCRE-style, unlike Go's RE2 — see regexcheck) instead of trusting them as
written-from-memory prose. `python -m unittest scripts/regex_lookaround.py`.
"""

import re
import unittest


class LookaroundTests(unittest.TestCase):
    def test_lookahead_matches_digits_before_usd(self):
        pattern = re.compile(r"\d+(?= USD)")
        m = pattern.search("150 USD")
        self.assertIsNotNone(m)
        self.assertEqual(m.group(), "150")
        self.assertIsNone(pattern.search("150 EUR"))

    def test_negative_lookahead_blocks_only_the_specific_continuation(self):
        pattern = re.compile(r"foo(?!bar)")
        self.assertIsNotNone(pattern.search("foobaz"))
        self.assertIsNone(pattern.search("foobar"))
        # Narrower than "reject the whole string" -- "foobar" also contains
        # a non-"bar"-followed "foo" nowhere else, so this isn't a case where
        # the assertion happens to coincide with a whole-string rejection.
        self.assertIsNone(re.search(r"foo(?!bar)", "xfoobar"))

    def test_lookbehind_matches_digits_after_dollar(self):
        pattern = re.compile(r"(?<=\$)\d+")
        m = pattern.search("$150")
        self.assertIsNotNone(m)
        self.assertEqual(m.group(), "150")
        self.assertIsNone(pattern.search("USD150"))

    def test_composed_lookaround_pulls_price_only(self):
        pattern = re.compile(r"(?<=\$)\d+\.\d+(?= each)")
        m = pattern.search("Total: $42.50 each")
        self.assertIsNotNone(m)
        self.assertEqual(m.group(), "42.50")
        # Missing either marker and the whole composed assertion fails, not
        # just the half it's missing.
        self.assertIsNone(pattern.search("Total: $42.50"))
        self.assertIsNone(pattern.search("Total: 42.50 each"))


if __name__ == "__main__":
    unittest.main()
