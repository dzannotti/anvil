package tag

func (tc *Container) Add(other Container) {
	for _, t := range other.tags {
		tc.AddTag(t)
	}
}

func (tc *Container) AddTag(newTag Tag) {
	if tc.HasTag(newTag) {
		return
	}
	tc.tags = append(tc.tags, newTag)
}

func (tc *Container) RemoveTag(target Tag) {
	newTags := make([]Tag, 0, len(tc.tags))
	for _, existing := range tc.tags {
		if !existing.MatchExact(target) {
			newTags = append(newTags, existing)
		}
	}
	tc.tags = newTags
}
