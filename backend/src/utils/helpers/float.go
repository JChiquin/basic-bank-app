package helpers

func FloatsInDelta(actual, expected, delta float64) bool {
	dt := actual - expected

	if dt < -delta || dt > delta {
		return false
	}
	return true
}
