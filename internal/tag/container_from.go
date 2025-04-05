package tag

func ContainerFromString(value string) Container {
	if len(value) == 0 {
		return Container{}
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

func ContainerFromTag(values ...Tag) Container {
	tc := Container{}
	for _, value := range values {
		tc.AddTag(value)
	}
	return tc
}

func ContainerFromContainer(value Container) Container {
	return ContainerFromStrings(value.Strings())
}
