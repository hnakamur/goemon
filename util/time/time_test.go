package time

import "testing"
import Ti "time"

func checkParseDuration(t *testing.T, input string, expected Ti.Duration) {
	d, err := ParseDuration(input)
	if err != nil {
		t.Fatalf("Parse error: %s", input)
	}
	if d != expected {
		t.Fatalf("Result error: input=%s, parsed=%s, expected=%s", input, d, expected)
	}
}

func TestParseDuration(t *testing.T) {
	checkParseDuration(t, "1s", Ti.Second);
	checkParseDuration(t, "1m", Ti.Minute);
	checkParseDuration(t, "1h", Ti.Hour);
	checkParseDuration(t, "1d", 24 * Ti.Hour);
	checkParseDuration(t, "1M", 31 * 24 * Ti.Hour);
	checkParseDuration(t, "1Y", 366 * 24 * Ti.Hour);
	checkParseDuration(t, "1M2h1m", (31 * 24 + 2) * Ti.Hour + Ti.Minute);
}
