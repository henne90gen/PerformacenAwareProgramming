package simulator8086

type Label struct {
	PositionInBytes int
}

func insert(a []Label, index int, value Label) []Label {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func insertLabel(labels []Label, position int) []Label {
	if len(labels) == 0 {
		return append(labels, Label{PositionInBytes: position})
	}

	for i, label := range labels {
		if label.PositionInBytes == position {
			return labels
		}

		if label.PositionInBytes < position {
			continue
		}

		return insert(labels, i, Label{PositionInBytes: position})
	}

	return append(labels, Label{PositionInBytes: position})
}
