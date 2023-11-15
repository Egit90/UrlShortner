package database

import (
	"egit90/urlShortner/model"

	"github.com/dgraph-io/badger/v4"
)

type DBManager struct {
	db *badger.DB
}

func NewDBManager(path string) (*DBManager, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &DBManager{db: db}, err
}

func (m *DBManager) CloseDB() error {

	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *DBManager) InsertKeyValue(d model.Data) error {
	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Set(d.Key, d.Value)
	})
}

func (m *DBManager) GetValue(s string) (string, error) {
	var outVal []byte
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(s))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			outVal = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		return "", err
	}
	return string(outVal), nil
}

func (c *DBManager) GetAll() ([]model.Data, error) {
	var cmds []model.Data
	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				cmds = append(cmds, model.Data{Key: k, Value: v})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmds, nil
}
