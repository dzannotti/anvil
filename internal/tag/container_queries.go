package tag

func (tc Container) MatchTag(target Tag) bool {
	for _, existing := range tc.tags {
		if existing.Match(target) {
			return true
		}
	}
	return false
}

func (tc Container) MatchAnyTag(other Container) bool {
	for _, tag := range other.tags {
		if tc.MatchTag(tag) {
			return true
		}
	}
	return false
}

func (tc Container) MatchAllTag(other Container) bool {
	for _, tag := range other.tags {
		if !tc.MatchTag(tag) {
			return false
		}
	}
	return true
}

func (container Container) HasTag(target Tag) bool {
	for _, t := range container.tags {
		if t.MatchExact(target) {
			return true
		}
	}
	return false
}

func (container Container) HasAnyTag(other Container) bool {
	for _, t := range other.tags {
		if container.HasTag(t) {
			return true
		}
	}
	return false
}

func (container Container) HasAllTag(other Container) bool {
	for _, t := range other.tags {
		if !container.HasTag(t) {
			return false
		}
	}
	return true
}
