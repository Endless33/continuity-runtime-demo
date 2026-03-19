package runtime

func SelectBestTransport(current Transport, candidates []Transport) Transport {
	best := candidates[0]

	for _, t := range candidates {
		if t.Score > best.Score {
			best = t
		}
	}

	return best
}