package pubsub

func GetFile(topic string) (*File, error) {
	file := filesStore.Get(topic)
	if file == nil {
		return nil, ErrFileNotFound
	}

	return file, nil
}
