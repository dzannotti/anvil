package tag

func ContainerFromString(value string) Container {
	if len(value) == 0 {
		return NewContainer()
	}
	return Container{
		tags: []Tag{FromString(value)},
	}
}

func ContainerFromStrings(values []string) Container {
	tags := make([]Tag, len(values))
	for i, value := range values {
		tags[i] = FromString(value)
	}
	return Container{tags: tags}
}

func ContainerFromTag(value Tag) Container {
	return Container{
		tags: []Tag{value},
	}
}

func ContainerFromTags(values []Tag) Container {
	tags := make([]Tag, len(values))
	copy(tags, values)
	return Container{tags: tags}
}

func ContainerFromContainer(value Container) Container {
	return ContainerFromStrings(value.Strings())
}
