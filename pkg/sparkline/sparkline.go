package sparkline

// Sparkline processes humidity measurements into a sparkline-friendly series.
type Sparkline struct {
	Height int // pixel height of sparkline (e.g., 14)
}

// NewSparkline constructs a Sparkline with dynamic buffer size and height.
func NewSparkline(height int) *Sparkline {
	return &Sparkline{Height: height}
}

// percentile returns the p-th percentile (0–100) from raw data using index method without floats.
// raw values are expected in the range [0,100].
func (s *Sparkline) percentile(raw []int16, p int) int16 {
	// copy and sort
	sorted := make([]int16, len(raw))
	copy(sorted, raw)
	// insertion sort (O(N^2) but N is small)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j-1] > sorted[j]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	idx := p * (len(sorted) - 1) / 100
	return sorted[idx]
}

// median3 computes the median of three values without allocations.
func median3(a, b, c int16) int16 {
	if a > b {
		a, b = b, a
	}
	if b > c {
		b, c = c, b
	}
	if a > b {
		a, b = b, a
	}
	return b
}

// Process applies 1% trim, median smoothing, and resamples to length N.
// Input raw values should be in the range [0,100].
func (s *Sparkline) Process(raw []int16) []int16 {
	rawLen := len(raw)

	// 1. compute percentiles
	p1 := s.percentile(raw, 1)
	p99 := s.percentile(raw, 99)

	// 2. trim extremes
	filtered := make([]int16, 0, len(raw))
	for _, v := range raw {
		if v < p1 || v > p99 {
			continue
		}
		filtered = append(filtered, v)
	}
	if len(filtered) == 0 {
		filtered = append(filtered, raw...)
	}

	// 3. median smoothing (window 3)
	smoothed := make([]int16, len(filtered))
	for i := range filtered {
		var a, b, c int16
		b = filtered[i]
		if i == 0 {
			a = filtered[0]
			c = filtered[1]
		} else if i == len(filtered)-1 {
			a = filtered[i-1]
			c = filtered[i]
		} else {
			a = filtered[i-1]
			c = filtered[i+1]
		}
		smoothed[i] = median3(a, b, c)
	}

	// 4. linear resampling back to s.N points using integer interpolation
	final := make([]int16, rawLen)
	denom := len(smoothed) - 1
	if denom <= 0 {
		// constant series
		for i := range rawLen {
			final[i] = smoothed[0]
		}
		return final
	}
	for i := range rawLen {
		// interpolation position
		num := i * denom
		idx := num / (rawLen - 1)
		rem := num % (rawLen - 1)
		v0 := smoothed[idx]
		next := idx + 1
		if next >= len(smoothed) {
			next = len(smoothed) - 1
		}
		v1 := smoothed[next]
		// interpolate: (v0*(denom-rem) + v1*rem) / denom
		interp := (int32(v0)*(int32(denom)-int32(rem)) + int32(v1)*int32(rem)) / int32(denom)
		final[i] = int16(interp)
	}

	return final
}
