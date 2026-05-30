package analog

import (
	"testing"

	"github.com/michi-1221/tty-clock/internal/clock"
)

func TestEllipsePointDirections(t *testing.T) {
	const cx, cy = 50, 50
	const r = 20 // circle (rx==ry) for direction checks
	cases := []struct {
		deg    float64
		wx, wy int
	}{
		{0, 50, 30},   // 12 o'clock → straight up
		{90, 70, 50},  // 3 o'clock → right
		{180, 50, 70}, // 6 o'clock → down
		{270, 30, 50}, // 9 o'clock → left
	}
	for _, c := range cases {
		x, y := ellipsePoint(cx, cy, r, r, c.deg)
		if x != c.wx || y != c.wy {
			t.Errorf("ellipsePoint(deg=%.0f) = (%d,%d), want (%d,%d)", c.deg, x, y, c.wx, c.wy)
		}
	}
}

func TestEllipseWiderThanTall(t *testing.T) {
	// rx > ry should reach farther horizontally than vertically.
	rightX, _ := ellipsePoint(50, 50, 40, 20, 90)
	_, downY := ellipsePoint(50, 50, 40, 20, 180)
	if (rightX - 50) <= (downY - 50) {
		t.Errorf("rx=40,ry=20: horizontal reach %d should exceed vertical %d", rightX-50, downY-50)
	}
}

func TestHandAngles(t *testing.T) {
	if got := hourAngle(clock.TimeSnapshot{Hour12: 1, Minute: 30}); got != 45 {
		t.Errorf("hourAngle(1:30) = %v, want 45", got)
	}
	if got := hourAngle(clock.TimeSnapshot{Hour12: 12, Minute: 0}); got != 0 {
		t.Errorf("hourAngle(12:00) = %v, want 0", got)
	}
	if got := secAngle(clock.TimeSnapshot{Second: 15}); got != 90 {
		t.Errorf("secAngle(15) = %v, want 90", got)
	}
	if got := minAngle(clock.TimeSnapshot{Minute: 10, Second: 30}); got != 63 {
		t.Errorf("minAngle(10:30) = %v, want 63", got)
	}
}

func TestHandUpperRightAt130(t *testing.T) {
	const cx, cy = 50, 50
	x, y := ellipsePoint(cx, cy, 20, 20, hourAngle(clock.TimeSnapshot{Hour12: 1, Minute: 30}))
	if !(x > cx && y < cy) {
		t.Errorf("1:30 hour tip (%d,%d) not upper-right of (%d,%d)", x, y, cx, cy)
	}
}

func TestLineEndpoints(t *testing.T) {
	c := NewCanvas(40, 40)
	line(c, 5, 5, 30, 20)
	if !c.Test(5, 5) || !c.Test(30, 20) {
		t.Error("Bresenham line must set both endpoints")
	}
}
