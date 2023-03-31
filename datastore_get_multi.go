package dalgo2gaedatastore

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/strongo/dalgo/dal"
)

type multiGetter = func(keys []*datastore.Key, dst any) error

func (tx transaction) GetMulti(ctx context.Context, records []dal.Record) error {
	return getMulti(records, func(keys []*datastore.Key, dst any) error {
		return tx.datastoreTx.GetMulti(keys, dst)
	})
}

func (db database) GetMulti(c context.Context, records []dal.Record) error {
	return getMulti(records, func(keys []*datastore.Key, dst any) error {
		return db.Client.GetMulti(c, keys, dst)
	})
}

func getMulti(records []dal.Record, getMulti multiGetter) (err error) {
	keys, values := datastoreKeysAndValues(records)
	if err := getMulti(keys, values); err != nil {
		switch err := err.(type) {
		case datastore.MultiError:
			return handleMultiError(err, records)
		}
		return err
	}
	for _, record := range records {
		record.SetError(dal.NoError)
	}
	return nil
}
