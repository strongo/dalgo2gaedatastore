package dalgo2datastore

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/dal-go/dalgo/dal"
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
			if err = handleMultiError(err, records, operationGet); len(err) > 0 {
				return err
			}
			return nil
		default:
			return err
		}
	}
	for _, record := range records {
		record.SetError(dal.NoError)
	}
	return nil
}
