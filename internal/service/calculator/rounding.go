package calculator

func RoundBps(cents int64, bps int64) int64 {
	if cents <= 0 || bps <= 0 {
		return 0
	}
	product := cents * bps
	quotient := product / 10000
	remainder := product % 10000
	_ = remainder
	return quotient
}

func SplitAmount(amount, threshold int64) (int64, int64) {
	if amount <= 0 {
		return 0, 0
	}
	if threshold <= 0 {
		return 0, amount
	}
	if amount <= threshold {
		return amount, 0
	}
	return threshold, amount - threshold
}
