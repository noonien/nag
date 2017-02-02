package main

import "errors"

type MapStore map[string]interface{}

func (ms MapStore) Set(key string, value interface{}) error {
	ms[key] = value
	return nil
}

func (ms MapStore) Get(key string) (interface{}, error) {
	val, ok := ms[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	return val, nil
}

func (ms MapStore) Delete(key string) error {
	delete(ms, key)
	return nil
}
