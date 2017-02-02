package bot

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type botStore struct {
	bot *Bot
	Store
}

func (bs *botStore) Get(key string) (interface{}, error) {
	switch key {

	default:
		return bs.Store.Get(key)
	}
}
