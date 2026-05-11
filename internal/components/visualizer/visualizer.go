package visualizer

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

const (
	sampleSize   = 512    // Number of samples for FFT (must be power of 2)
	smoothFactor = 0.8    // Smoothing factor for visual transitions (higher = smoother)
	decayFactor  = 0.7    // How fast peaks decay (higher = slower decay)
	minMagnitude = 0.0001 // Minimum magnitude to avoid log(0)
)

// Model represents the visualizer state
type Model struct {
	width      int
	height     int
	numBars    int
	buffer     [][2]float64
	freqBins   []float64
	smoothBins []float64
	peakBins   []float64 // Peak values for each bar (decay over time)
	mu         sync.RWMutex
	enabled    bool
}

// New creates a new visualizer model
func New() Model {
	width := 100
	numBars := width / 2
	return Model{
		width:      width,
		height:     5,
		numBars:    numBars,
		buffer:     make([][2]float64, 0, sampleSize),
		freqBins:   make([]float64, numBars),
		smoothBins: make([]float64, numBars),
		peakBins:   make([]float64, numBars),
		enabled:    true,
	}
}

func (m *Model) SetSize(width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.width = width
	m.height = height
	// Calculate number of bars based on width (2 characters per bar)
	numBars := width / 2
	if numBars < 1 {
		numBars = 1
	}

	// Always update numBars to match width, resize arrays if needed
	if numBars != m.numBars {
		m.numBars = numBars
		// Create new arrays with the correct size
		newFreqBins := make([]float64, numBars)
		newSmoothBins := make([]float64, numBars)
		newPeakBins := make([]float64, numBars)

		// Copy existing data if we're resizing (interpolate if needed)
		oldBars := len(m.freqBins)
		if oldBars > 0 && oldBars == numBars {
			// Same size, just copy
			copy(newFreqBins, m.freqBins)
			copy(newSmoothBins, m.smoothBins)
			copy(newPeakBins, m.peakBins)
		} else if oldBars > 0 {
			// Different size, interpolate
			for i := 0; i < numBars; i++ {
				// Map new bar index to old bar index
				oldIdx := int(float64(i) * float64(oldBars) / float64(numBars))
				if oldIdx >= oldBars {
					oldIdx = oldBars - 1
				}
				if oldIdx < len(m.freqBins) {
					newFreqBins[i] = m.freqBins[oldIdx]
					newSmoothBins[i] = m.smoothBins[oldIdx]
					newPeakBins[i] = m.peakBins[oldIdx]
				}
			}
		}

		m.freqBins = newFreqBins
		m.smoothBins = newSmoothBins
		m.peakBins = newPeakBins
	} else {
		// Width/number of bars didn't change, but make sure numBars is set correctly
		m.numBars = numBars
	}
}

// Enable enables or disables the visualizer
func (m *Model) Enable(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

// IsEnabled returns whether the visualizer is enabled
func (m *Model) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// Update processes audio samples and updates frequency bins
func (m *Model) Update(samples [][2]float64) {
	if !m.enabled {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Add samples to buffer
	m.buffer = append(m.buffer, samples...)

	// Keep only the last sampleSize samples
	if len(m.buffer) > sampleSize {
		m.buffer = m.buffer[len(m.buffer)-sampleSize:]
	}

	// If we have enough samples, perform FFT
	if len(m.buffer) >= sampleSize {
		m.performFFT()
	}
}

// performFFT performs Fast Fourier Transform on the audio buffer
func (m *Model) performFFT() {
	// Convert stereo to mono by averaging channels
	mono := make([]float64, sampleSize)
	for i := 0; i < sampleSize && i < len(m.buffer); i++ {
		mono[i] = (m.buffer[i][0] + m.buffer[i][1]) / 2.0
	}

	// Perform FFT
	fftResult := fft(mono)

	// Use logarithmic frequency binning for better musical representation
	// This gives more resolution to lower frequencies (bass/mids) and less to highs
	maxFreqBin := len(fftResult) / 2 // Nyquist frequency

	numBars := m.numBars
	rawValues := make([]float64, numBars)

	// First pass: calculate raw magnitude values for each bar
	for i := 0; i < numBars; i++ {
		// Use logarithmic frequency mapping (more resolution for bass)
		// Map bar index to frequency range using logarithmic scale
		// For log scale: map 0-1 range to log-space frequencies

		// Convert linear bar position (0 to 1) to logarithmic frequency position
		// Use a logarithmic scale that gives better distribution
		linearPos := float64(i) / float64(numBars)
		nextLinearPos := float64(i+1) / float64(numBars)

		// Map to logarithmic scale: log10 range from -2 (very low) to 0 (high)
		// This gives more resolution in lower frequencies
		minLogFreq := -1.0 // Start at 0.1 * maxFreq (avoids collapsing many bars onto same FFT bin)
		maxLogFreq := 0.0  // End at 1.0 * maxFreq

		logStart := minLogFreq + linearPos*(maxLogFreq-minLogFreq)
		logEnd := minLogFreq + nextLinearPos*(maxLogFreq-minLogFreq)

		// Convert back from log space to linear frequency bin indices
		start := int(math.Pow(10, logStart) * float64(maxFreqBin))
		end := int(math.Pow(10, logEnd) * float64(maxFreqBin))

		if end <= start {
			end = start + 1
		}
		if end > maxFreqBin {
			end = maxFreqBin
		}
		if start > maxFreqBin {
			start = maxFreqBin
		}
		if start < 0 {
			start = 0
		}

		// Calculate RMS (Root Mean Square) for better representation
		sumSquares := 0.0
		count := 0
		for j := start; j < end; j++ {
			magnitude := cmplx.Abs(fftResult[j])
			sumSquares += magnitude * magnitude
			count++
		}

		if count > 0 {
			// RMS = sqrt(mean of squares)
			rms := math.Sqrt(sumSquares / float64(count))
			rawValues[i] = rms
		}
	}

	// Find maximum value in this frame for adaptive normalization
	maxValue := 0.0
	for i := 0; i < numBars; i++ {
		if rawValues[i] > maxValue {
			maxValue = rawValues[i]
		}
	}

	// Normalize and update bins with adaptive scaling
	for i := 0; i < numBars; i++ {
		rawValue := rawValues[i]

		// Apply logarithmic scaling (dB scale) for better visual representation
		// This makes quieter sounds more visible
		var normalized float64
		if rawValue > minMagnitude && maxValue > minMagnitude {
			// Convert to dB scale: 20 * log10(magnitude)
			db := 20 * math.Log10(rawValue)
			maxDb := 20 * math.Log10(maxValue)

			// Normalize relative to current frame maximum
			// Scale from (maxDb - 60) to maxDb range
			rangeMin := maxDb - 60.0 // 60dB dynamic range
			if db < rangeMin {
				normalized = 0.0
			} else {
				normalized = (db - rangeMin) / 60.0
			}

			// Clamp to 0-1 range
			if normalized > 1.0 {
				normalized = 1.0
			}
			if normalized < 0.0 {
				normalized = 0.0
			}
		} else {
			normalized = 0.0
		}

		// Apply smoothing to the normalized value
		m.smoothBins[i] = m.smoothBins[i]*smoothFactor + normalized*(1-smoothFactor)

		// Track peaks with decay (bars stay at peak briefly, then decay)
		if m.smoothBins[i] > m.peakBins[i] {
			m.peakBins[i] = m.smoothBins[i]
		} else {
			// Decay peaks slowly
			m.peakBins[i] *= decayFactor
			// If smoothed value is higher than decayed peak, use smoothed value
			if m.smoothBins[i] > m.peakBins[i] {
				m.peakBins[i] = m.smoothBins[i]
			}
		}

		// Use peak value for display (this gives the characteristic visualizer look)
		m.freqBins[i] = m.peakBins[i]
	}
}

// View renders the visualizer
func (m *Model) View() string {
	if !m.enabled {
		return ""
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.height < 1 || m.width < 1 {
		return ""
	}

	var sb strings.Builder

	numBars := m.numBars

	// Define block characters for bar rendering (from empty to full)
	blocks := []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	// Pre-calculate colors for each bar (same color for all rows in a column)
	barColors := make([]lipgloss.Color, numBars)
	for barIdx := 0; barIdx < numBars; barIdx++ {
		// Calculate rainbow color based on frequency (barIdx) only
		// Each column (bar) gets one solid color, next column changes slightly
		// Frequency position: 0.0 (low) to 1.0 (high) - controls hue across columns
		freqHue := float64(barIdx) / float64(numBars)
		// Start at 270 degrees (magenta) and cycle through full rainbow spectrum
		hue := math.Mod(freqHue*360.0+270.0, 360.0)
		// Convert HSV to RGB for rainbow colors (full saturation and brightness)
		barColors[barIdx] = hsvToRgb(hue, 1.0, 1.0)
	}

	// Render each row from top to bottom
	for row := 0; row < m.height; row++ {
		for barIdx := 0; barIdx < numBars; barIdx++ {
			// Use the normalized frequency value directly (already 0-1 from FFT)
			normalized := m.freqBins[barIdx]
			if normalized > 1.0 {
				normalized = 1.0
			}
			if normalized < 0.0 {
				normalized = 0.0
			}

			// Calculate bar height in rows
			barHeight := normalized * float64(m.height)

			// Calculate which character to use for this position
			currentRow := float64(m.height - row - 1)

			var char rune
			if barHeight > currentRow+1 {
				// Full block
				char = blocks[len(blocks)-1]
			} else if barHeight > currentRow {
				// Partial block
				fraction := barHeight - currentRow
				blockIdx := int(fraction * float64(len(blocks)-1))
				if blockIdx >= len(blocks) {
					blockIdx = len(blocks) - 1
				}
				char = blocks[blockIdx]
			} else {
				// Empty
				char = blocks[0]
			}

			// Use pre-calculated color for this bar
			style := lipgloss.NewStyle().Foreground(barColors[barIdx])

			// Each bar is 2 characters wide to fill the screen
			sb.WriteString(style.Render(string(char) + string(char)))
		}
		if row < m.height-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

// hsvToRgb converts HSV color to lipgloss.Color (RGB hex string)
// Hue: 0-360, Saturation: 0-1, Value: 0-1
func hsvToRgb(h, s, v float64) lipgloss.Color {
	// Clamp values
	if h < 0 {
		h = 0
	}
	if h >= 360 {
		h = 359.999
	}
	if s < 0 {
		s = 0
	}
	if s > 1 {
		s = 1
	}
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}

	// HSV to RGB conversion
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2.0)-1))
	m := v - c

	var r, g, b float64

	if h < 60 {
		r, g, b = c, x, 0
	} else if h < 120 {
		r, g, b = x, c, 0
	} else if h < 180 {
		r, g, b = 0, c, x
	} else if h < 240 {
		r, g, b = 0, x, c
	} else if h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	// Convert to 0-255 range and format as hex
	R := int((r + m) * 255)
	G := int((g + m) * 255)
	B := int((b + m) * 255)

	// Ensure values are in valid range
	if R < 0 {
		R = 0
	}
	if R > 255 {
		R = 255
	}
	if G < 0 {
		G = 0
	}
	if G > 255 {
		G = 255
	}
	if B < 0 {
		B = 0
	}
	if B > 255 {
		B = 255
	}

	// Format as hex color
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", R, G, B))
}

// fft performs a simple FFT on the input data
// This is a straightforward implementation of the Cooley-Tukey algorithm
func fft(input []float64) []complex128 {
	n := len(input)

	// Base case
	if n <= 1 {
		result := make([]complex128, n)
		if n == 1 {
			result[0] = complex(input[0], 0)
		}
		return result
	}

	// Divide
	even := make([]float64, n/2)
	odd := make([]float64, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = input[i*2]
		odd[i] = input[i*2+1]
	}

	// Conquer
	fftEven := fft(even)
	fftOdd := fft(odd)

	// Combine
	result := make([]complex128, n)
	for k := 0; k < n/2; k++ {
		t := cmplx.Exp(complex(0, -2*math.Pi*float64(k)/float64(n))) * fftOdd[k]
		result[k] = fftEven[k] + t
		result[k+n/2] = fftEven[k] - t
	}

	return result
}
