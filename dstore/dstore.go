package dstore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
)

type DStore struct {
	Ctx    context.Context
	Client *datastore.Client
}

type PutStrct struct {
	Key *datastore.Key
	Src interface{}
}

func Create(project string) (*DStore, error) {
	var err error
	ds := &DStore{
		Ctx: context.Background(),
	}
	if ds.Client, err = datastore.NewClient(ds.Ctx, project); err != nil {
		return ds, fmt.Errorf("dstore.Create:datastore.NewClient failure project:%v\n\t%v", project, err)
	}
	return ds, nil
}

func (ds *DStore) IDKey(kind string, id int64, parent *datastore.Key) *datastore.Key {
	return datastore.IDKey(kind, id, parent)
}

func (ds *DStore) IncompleteKey(kind string, parent *datastore.Key) *datastore.Key {
	return datastore.IncompleteKey(kind, parent)
}

func (ds *DStore) IDKeyGet(kind string, id int64, src interface{}) (*datastore.Key, error) {
	key := datastore.IDKey(kind, id, nil)

	if err := ds.Client.Get(ds.Ctx, key, src); err != nil {
		return key, fmt.Errorf("dstore.IDKeyGet:Failed to Get: %v", err)
	}

	return key, nil
}

func (ds *DStore) IDKeyGetMulti(kind string, ids []int64, src interface{}) ([]*datastore.Key, error) {
	keys := make([]*datastore.Key, len(ids))
	for i, v := range ids {
		keys[i] = datastore.IDKey(kind, v, nil)
	}

	if err := ds.Client.GetMulti(ds.Ctx, keys, src); err != nil {
		return keys, fmt.Errorf("dstore.IDKeyGetMulti:Failed to GetMulti: %v", err)
	}

	return keys, nil
}

func (ds *DStore) IncompleteKeyPut(kind string, src interface{}) (int64, error) {
	newkey := datastore.IncompleteKey(kind, nil)
	if key, err := ds.Client.Put(ds.Ctx, newkey, src); err != nil {
		log.Fatalf("IncompleteKeyPut:Failed to put: %v", err)
		return key.ID, err
	} else {
		return key.ID, nil
	}
}

// どうみても返り値のint64の意味がないので暇を見つけて修正する
func (ds *DStore) IDKeyPut(kind string, id int64, src interface{}) (int64, error) {
	key := datastore.IDKey(kind, id, nil)

	if _, err := ds.Client.Put(ds.Ctx, key, src); err != nil {
		log.Fatalf("IDKeyPut:Failed to put: %#v", err)
		return key.ID, err
	}

	return key.ID, nil
}

// propertyにvalueがない場合にのみ、Insertする
func (ds *DStore) InsertUnique(kind, property, value string, src interface{}) (int64, error) {
	log.Printf("Start dstore.InsertUnique kind:%q, property:%q, value:%q, src:%v.", kind, property, value, src)

	var keys []*datastore.Key
	exist := false
	_, err := ds.Client.RunInTransaction(ds.Ctx, func(tx *datastore.Transaction) error {
		// 存在しているかどうかだけを知りたいので、1つより多く取得する必要はない。
		q := datastore.NewQuery(kind).FilterField(property, "=", value).KeysOnly().Limit(1)
		var terr error
		keys, terr = ds.Client.GetAll(ds.Ctx, q, nil)
		if terr != nil {
			return fmt.Errorf("dstore.InsertUnique client.GetAll1 failure terr:%v", terr)
		} else if len(keys) == 0 {
			key := datastore.IncompleteKey(kind, nil)
			// If there was no matching entity, store it now.
			log.Printf("before Put dstore.InsertUnique kind:%q property:%q value:%q terr:%v", kind, property, value, terr)
			log.Printf("before Put src:%v", src)
			_, terr = tx.Put(key, src)
			log.Printf("after Put dstore.InsertUnique kind:%q property:%q value:%q terr:%v", kind, property, value, terr)
			log.Printf("after Put src:%v.", src)
			if terr != nil {
				return fmt.Errorf("dstore.InsertUnique tx.Put failure terr:%v", terr)
			}
			return nil
		} else {
			exist = true
			return fmt.Errorf("dstore.InsertUnique:Already Exist Entity:%q", value)
		}
	})

	log.Printf("dstore.InsertUnique client.RunInTransaction kind:%q, property:%q, value:%q, src:%v err:%v", kind, property, value, src, err)
	if exist {
		return -1, nil // すでに同名のアカウントが存在した場合は、-1を返す
	} else if err != nil {
		return -1, fmt.Errorf("dstore.InsertUnique err:%vq", err)
	}

	return keys[0].ID, nil
}

// valueを持つpropertyが一つだけの場合のみ取得する
func (ds *DStore) QueryUnique(kind, property, value string, src interface{}) (int64, error) {
	log.Printf("Start dstore.QueryUnique kind:%q, property:%q, value:%q, src:%+v", kind, property, value, src)

	// 重複しているかどうかだけを知りたいので、2つより多く取得する必要はない。
	q := datastore.NewQuery(kind).Filter(property+" =", value).Limit(2)

	if src == nil {
		q = q.KeysOnly()
	}
	if keys, err := ds.Client.GetAll(ds.Ctx, q, src); err != nil {
		return -1, fmt.Errorf("dstore.QueryUnique err:%v", err)
	} else {
		num := len(keys)
		switch num {
		case 1:
			return keys[0].ID, nil
		case -1:
			return -1, nil
		default:
			return -1, fmt.Errorf("dstore.QueryUnique:Invalid num:%d kind:%q, property:%q, value:%q, err:%v", num, kind, property, value, err)
		}
	}
}

// Keyを流用しながら、同時に2つのkindを追加
func (ds *DStore) IncompleteKeyDoublePut(kind1, kind2 string, src1, src2 interface{}) (int64, error) {
	key1 := datastore.IncompleteKey(kind1, nil)

	_, err := ds.Client.RunInTransaction(ds.Ctx, func(tx *datastore.Transaction) error {
		var err error
		if key1, err = ds.Client.Put(ds.Ctx, key1, src1); err != nil {
			return fmt.Errorf("dstore.IncompleteKeyDoublePut:Failed to put(1): %#v", err)
		}
		key2 := datastore.IDKey(kind2, key1.ID, nil) // 直前のPutでセットされたkey1を用いてkey2を設定する
		log.Printf("dstore.IncompleteKeyDoublePut: key1:%v key2:%v.", key1, key2)
		if _, err = ds.Client.Put(ds.Ctx, key2, src2); err != nil {
			return fmt.Errorf("dstore.IncompleteKeyDoublePut:Failed to put(2): %#v", err)
		}
		return nil
	})

	if err != nil {
		return -1, fmt.Errorf("dstore.IncompleteKeyDoublePut:client.RunInTransaction failed\n\t%v", err)
	}

	return key1.ID, nil
}

// transaction中に配列の要素をすべてPut
func (ds *DStore) PutKinds(ar []PutStrct) error {
	log.Printf("Start dstore.PutKinds ar:%v", ar)
	defer log.Printf("End dstore.PutKinds ar:%v", ar)
	_, err := ds.Client.RunInTransaction(ds.Ctx, func(tx *datastore.Transaction) error {
		for i, v := range ar {
			log.Printf("dstore.PutKinds i:%d v.Key:%v v.Src:%v", i, v.Key, v.Src)
			if _, err := ds.Client.Put(ds.Ctx, v.Key, v.Src); err != nil {
				log.Printf("dstore.PutKinds failure i:%d v.Key:%v v.Src:%v err:%v", i, v.Key, v.Src, err)
				return fmt.Errorf("dstore.IncompleteKeyDoublePut:Failed to datastore.client.Put(%d): %#v", i, err)
			}
			log.Printf("dstore.PutKinds success i:%d v.Key:%v v.Src:%v", i, v.Key, v.Src)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("dstore.PutKinds:client.RunInTransaction failed\n\t%v", err)
	}

	return nil
}

func (ds *DStore) GetAll(kind, sort string, parentKey *datastore.Key, dst interface{}) error {
	q := datastore.NewQuery(kind)
	if parentKey != nil {
		q = q.Ancestor(parentKey)
	}
	if len(sort) > 0 {
		q = q.Order(sort)
		//		q = q.Order( "-" + sort )
	}
	if _, err := ds.Client.GetAll(ds.Ctx, q, dst); err != nil {
		return fmt.Errorf("dstore.DStore.GetAll:datastore.Client.GetAll failure\n\t%v", err)
	}
	return nil
}

// Keyを流用しながら、同時に2つのkindを削除
func (ds *DStore) Delete(kind string, id int64) error {
	key := datastore.IDKey(kind, id, nil)

	if err := ds.Client.Delete(ds.Ctx, key); err != nil {
		return fmt.Errorf("dstore.Delete:Failed to Delete: %#v", err)
	}

	return nil
}
