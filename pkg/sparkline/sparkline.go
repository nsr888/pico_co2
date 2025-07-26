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
func (s *Sparkline) Process(raw []int16) []int16 {
	rawLen := len(raw)
	if rawLen == 0 {
		return []int16{}
	}

	// 1. compute percentiles
	p1 := s.percentile(raw, 1)
	p99 := s.percentile(raw, 99)

	// 2. trim extremes
	filtered := make([]int16, 0, rawLen)
	for _, v := range raw {
		if v < p1 || v > p99 {
			continue
		}
		filtered = append(filtered, v)
	}
	if len(filtered) == 0 {
		filtered = raw
	}

	// 3. median smoothing (window 3)
	smoothed := make([]int16, len(filtered))
	if len(filtered) < 2 {
		copy(smoothed, filtered)
	} else {
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
	}

	// max-min normalization to fit into sparkline height
	if len(smoothed) == 0 {
		return make([]int16, rawLen)
	}
	min := smoothed[0]
	max := smoothed[0]
	for _, v := range smoothed {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// if all values are the same, return a constant series
	if min == max {
		constant := int16(s.Height / 2)
		final := make([]int16, rawLen)
		for i := range final {
			final[i] = constant
		}
		return final
	}

	// normalize to height
	scale := max - min
	for i := range smoothed {
		// scale to height-1, to fit in 0..height-1 range
		smoothed[i] = int16((int32(smoothed[i]-min) * int32(s.Height-1)) / int32(scale))
	}

	// 4. linear resampling back to rawLen points
	final := make([]int16, rawLen)
	if len(smoothed) <= 1 {
		val := int16(s.Height / 2)
		if len(smoothed) == 1 {
			val = smoothed[0]
		}
		for i := range final {
			final[i] = val
		}
		return final
	}

	if rawLen <= 1 {
		if rawLen == 1 {
			final[0] = smoothed[0]
		}
		return final
	}

	inDenom := len(smoothed) - 1
	outDenom := rawLen - 1
	for i := range final {
		// interpolation position
		num := i * inDenom
		idx := num / outDenom
		rem := num % outDenom

		v0 := smoothed[idx]
		next := idx + 1
		if next >= len(smoothed) {
			next = len(smoothed) - 1
		}
		v1 := smoothed[next]

		// interpolate: (v0*(out_denom-rem) + v1*rem) / out_denom
		interp := (int32(v0)*(int32(outDenom)-int32(rem)) + int32(v1)*int32(rem)) / int32(outDenom)
		final[i] = int16(interp)
	}

	return final
}
