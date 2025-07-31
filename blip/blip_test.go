package blip_test

import (
	"math"
	"testing"

	"github.com/jonathangjertsen/toyboy/blip"
)

func TestShiftingAndCasting(t *testing.T) {
	intmin := int(-0x8000_0000)
	if (intmin >> 1) != intmin/2 {
		t.Errorf("assumption failed: right shift of negative value preserves sign")
	}

	shortMax := int16(0x7fff)
	shortMin := int16(-0x8000)
	l := int((shortMax + 1) * 5)
	if int16(l) != shortMin {
		t.Errorf("assumption failed: casting to smaller signed type truncates bits and extends sign")
	}
}

func setUp() (*blip.Buffer, *blip.Synth) {
	bb := blip.NewBuffer(blip.DefaultBlipConfig)
	bs := blip.NewSynth(bb)
	return bb, bs
}

func TestExpectedWaveformShape(t *testing.T) {
	bb, bs := setUp()

	// Simplified version: just one step
	bs.Update(100, 1000)
	bb.EndFrame(200)

	out := make([]int16, 200)
	count := bb.Read(out)

	if count != 200 {
		t.Errorf("Expected 200 samples, got %d", count)
	}

	// In a working band-limited system, we should see:
	// - Samples before t=100 should be close to 0
	// - Samples after t=100 should be close to the target amplitude
	// - There should be a smooth transition around t=100

	beforeCount := 0
	afterCount := 0

	for i := range 90 { // Well before t=100
		if out[i] != 0 {
			beforeCount++
		}
	}

	for i := 110; i < 200; i++ { // Well after t=100
		if out[i] != 0 {
			afterCount++
		}
	}

	if beforeCount == 0 && afterCount == 0 {
		t.Error("No non-zero samples found - synthesis not working")
	}

	if afterCount == 0 {
		t.Error("No samples after transition - step response not working")
	}
}

func TestWaveformValues(t *testing.T) {
	bb, bs := setUp()

	// Simple step from 0 to 1000 at t=100
	bs.Update(100, 1000)
	bb.EndFrame(200)

	out := make([]int16, 200)
	_ = bb.Read(out)

	// Check key sample ranges
	var beforeAvg, afterAvg float64

	// Average samples well before the transition (t=50-90)
	for i := 50; i < 90; i++ {
		beforeAvg += float64(out[i])
	}
	beforeAvg /= 40

	// Average samples well after the transition (t=120-160)
	for i := 120; i < 160; i++ {
		afterAvg += float64(out[i])
	}
	afterAvg /= 40

	// The step should create a clear difference
	if abs(afterAvg-beforeAvg) < 100 {
		t.Error("Step response too small - synthesis may not be working correctly")
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestZeroAmplitudeTransitions(t *testing.T) {
	bb, bs := setUp()

	// Test transitions with zero amplitude
	bs.Update(100, 0)
	bs.Update(150, 0)
	bs.Update(200, 0)
	bb.EndFrame(300)

	out := make([]int16, 300)
	count := bb.Read(out)

	// All samples should be zero or very close to zero
	maxAbs := int16(0)
	for i := range count {
		if abs16(out[i]) > maxAbs {
			maxAbs = abs16(out[i])
		}
	}

	if maxAbs > 10 { // Allow for small numerical errors
		t.Errorf("Expected near-zero output for zero amplitudes, got max absolute value %d", maxAbs)
	}
}

func TestNegativeAmplitudes(t *testing.T) {
	bb, bs := setUp()

	// Test negative amplitude step
	bs.Update(100, -1000)
	bb.EndFrame(200)

	out := make([]int16, 200)
	count := bb.Read(out)

	// Should have negative values after the transition
	hasNegative := false
	for i := 110; i < count; i++ {
		if out[i] < 0 {
			hasNegative = true
			break
		}
	}

	if !hasNegative {
		t.Error("Expected negative samples after negative amplitude transition")
	}
}

func TestMultipleConsecutiveTransitions(t *testing.T) {
	bb, bs := setUp()

	// Test rapid consecutive transitions
	bs.Update(100, 500)
	bs.Update(101, 1000)
	bs.Update(102, 750)
	bs.Update(103, 250)
	bs.Update(104, 0)
	bb.EndFrame(200)

	out := make([]int16, 200)
	count := bb.Read(out)

	if count == 0 {
		t.Error("No samples generated from consecutive transitions")
	}

	// Check that the output shows variation around the transition area
	hasVariation := false
	for i := 95; i < 110 && i < count-1; i++ {
		if out[i] != out[i+1] {
			hasVariation = true
			break
		}
	}

	if !hasVariation {
		t.Error("Expected variation in output during rapid transitions")
	}
}

func TestLargeAmplitudeValues(t *testing.T) {
	bb, bs := setUp()

	// Test with very large amplitude that should cause saturation/clipping
	largeAmp := int(100000) // Much larger than int16 range
	bs.Update(100, largeAmp)
	bb.EndFrame(200)

	out := make([]int16, 200)
	count := bb.Read(out)

	// Find the peak absolute value in the output
	maxAbs := int16(0)
	peakIdx := -1
	for i := range count {
		abs := out[i]
		if abs < 0 {
			abs = -abs
		}
		if abs > maxAbs {
			maxAbs = abs
			peakIdx = i
		}
	}

	// With such a large amplitude, we should see saturation at the limits
	// The blip buffer should clamp rather than wrap around
	if maxAbs < 20000 {
		t.Errorf("Expected large amplitude output near saturation, got max absolute value %d", maxAbs)
	}

	// Test that we don't see sudden sign flips that would indicate wrapping
	// (legitimate signal should have smooth transitions)
	signFlips := 0
	for i := 1; i < count; i++ {
		// Look for unrealistic jumps that suggest wrapping
		diff := int32(out[i]) - int32(out[i-1])
		if diff > 40000 || diff < -40000 {
			signFlips++
		}
	}

	if signFlips > 0 {
		t.Errorf("Found %d suspicious large jumps that may indicate wrapping at sample %d", signFlips, peakIdx)
	}
}

func TestPartialBufferReads(t *testing.T) {
	bb, bs := setUp()

	bs.Update(100, 1000)
	bb.EndFrame(200)

	totalAvailable := bb.SamplesAvailable()

	// Read in chunks
	chunkSize := (50)
	totalRead := (0)

	for bb.SamplesAvailable() > 0 {
		availableBefore := bb.SamplesAvailable()

		out := make([]int16, chunkSize)
		count := bb.Read(out)

		if count == 0 {
			break
		}

		totalRead += count

		// Verify we're actually consuming samples
		availableAfter := bb.SamplesAvailable()
		if availableAfter >= availableBefore {
			t.Errorf("ReadSamples not consuming samples: before=%d, after=%d, read=%d",
				availableBefore, availableAfter, count)
			break
		}

		// Verify the number consumed matches what was read
		consumed := availableBefore - availableAfter
		if consumed != count {
			t.Errorf("Samples consumed (%d) doesn't match samples read (%d)", consumed, count)
		}
	}

	if totalRead != totalAvailable {
		t.Errorf("Total read samples %d != initially available %d", totalRead, totalAvailable)
	}
}

// Helper function for absolute value of int16
func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
func TestSquareWaveWithSineMix(t *testing.T) {
	bb, bs := setUp()

	// Generate 1000 clocks of square wave (similar to C++ demo)
	length := 1000
	amplitude := 1
	for time := 0; time < length; time += 10 {
		bs.Update(time, amplitude)
		amplitude = -amplitude
	}

	// Find out how many samples of sine wave to generate
	count := bb.CountSamples(length)
	temp := make([]int16, count)

	// Generate sine wave samples
	for i := range count {
		y := math.Sin(float64(i) * (math.Pi / 100))
		temp[i] = int16(y * 0.30 * 32767) // convert to Sample range
	}

	// Mix sine wave's samples into Blip_Buffer
	bb.MixSamples(temp)

	// End frame
	bb.EndFrame(length)

	// Verify we have samples
	if bb.SamplesAvailable() == 0 {
		t.Error("No samples available after square wave + sine mix")
	}

	// Read samples and verify they're not all zero
	out := make([]int16, bb.SamplesAvailable())
	readCount := bb.Read(out)

	nonZeroCount := 0
	for i := range readCount {
		if out[i] != 0 {
			nonZeroCount++
		}
	}

	if nonZeroCount == 0 {
		t.Error("All output samples are zero after mixing square wave and sine")
	}

	// Verify the mixed signal has reasonable amplitude
	maxAbs := int16(0)
	for i := range readCount {
		abs := out[i]
		if abs < 0 {
			abs = -abs
		}
		if abs > maxAbs {
			maxAbs = abs
		}
	}

	if maxAbs < 1000 {
		t.Error("Mixed signal amplitude seems too low")
	}
}

func TestSawWaveGeneration(t *testing.T) {
	bb, bs := setUp()

	// Generate a simple saw wave
	length := 1000
	amp := 0

	// Generate saw wave: 0, 1, 2, 3, 4, 0, 1, 2, 3, 4, ...
	for time := 0; time < length; time += 50 {
		bs.Update(time, amp)
		amp = (amp + 1) % 5
	}

	bb.EndFrame(length)

	// Read samples
	out := make([]int16, bb.SamplesAvailable())
	count := bb.Read(out)

	if count == 0 {
		t.Error("No samples generated for saw wave")
	}

	// Verify the saw wave has the expected ramping behavior
	// Look for generally increasing trends followed by resets
	hasIncreasingTrend := false

	for i := 1; i < count-1; i++ {
		if out[i] > out[i-1] && out[i+1] > out[i] {
			hasIncreasingTrend = true
		}
	}

	if !hasIncreasingTrend {
		t.Error("Saw wave should have increasing trends")
	}
}

func TestMixSamplesFunction(t *testing.T) {
	bb := blip.NewBuffer(blip.DefaultBlipConfig)

	// Create some test samples to mix in
	testSamples := []int16{1000, 2000, -1500, 500, -500}

	// Mix the samples
	bb.MixSamples(testSamples)

	// End frame to make samples available
	bb.EndFrame(100)

	if bb.SamplesAvailable() == 0 {
		t.Error("No samples available after MixSamples")
	}

	// Read back and verify we get non-zero output
	out := make([]int16, bb.SamplesAvailable())
	count := bb.Read(out)

	hasNonZero := false
	for i := range count {
		if out[i] != 0 {
			hasNonZero = true
			break
		}
	}

	if !hasNonZero {
		t.Error("MixSamples produced all zero output")
	}
}
