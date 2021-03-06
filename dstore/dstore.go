package dstore

import (
	"log"
	"fmt"
	"context"

	"cloud.google.com/go/datastore"
)

type DStore struct {
	Ctx context.Context
	Client *datastore.Client
}

type PutStrct struct {
	Key *datastore.Key
	Src interface{}
}

func Create( project string ) ( *DStore, error ) {
	var err error
	ds := &DStore{
		Ctx: context.Background(),
	}
	if ds.Client, err = datastore.NewClient( ds.Ctx, project ); err != nil {
		return ds, fmt.Errorf( "dstore.Create:datastore.NewClient failure project:%v\n\t%v", project, err )
	}
	return ds, nil
}

func ( ds *DStore ) IDKey( kind string, id int64, parent *datastore.Key ) *datastore.Key {
	return datastore.IDKey( kind, id, parent )
}

func ( ds *DStore ) IncompleteKey( kind string, parent *datastore.Key ) *datastore.Key {
	return datastore.IncompleteKey( kind, parent )
}

func ( ds *DStore ) IDKeyGet( kind string, id int64, src interface{} ) ( *datastore.Key, error ) {
	key := datastore.IDKey( kind, id, nil )

	if err := ds.Client.Get( ds.Ctx, key, src ); err != nil {
    	return key, fmt.Errorf( "dstore.IDKeyGet:Failed to Get: %#v", err )
	}

	return key, nil
}

func ( ds *DStore ) IncompleteKeyPut( kind string, src interface{} ) ( int64, error ) {
	key := datastore.IncompleteKey( kind, nil )

	if _, err := ds.Client.Put( ds.Ctx, key, src ); err != nil {
		log.Fatalf( "IncompleteKeyPut:Failed to put: %#v", err )
		return key.ID, err
	}

	return key.ID, nil
}

/*
func ( ds *DStore ) Insert( kind, name string, src interface{} ) ( error, bool ) {
	exist := false
	key := datastore.NameKey( kind, name, nil )

	_, err := ds.Client.RunInTransaction( ds.Ctx, func(tx *datastore.Transaction ) error {
		// We first check that there is no entity stored with the given key.
		type Empty struct {
			Dummy            int
		}
		var empty Empty
		if err := tx.Get( key, &empty ); err == nil {
			log.Printf( "dstore.Insert:Already Exist key:%q.", key )
			exist = true
			return nil
		} else if err != datastore.ErrNoSuchEntity {
			log.Printf( "dstore.Insert:Get failure:%v.", err )
			return err
		}
		// If there was no matching entity, store it now.
		log.Printf( "before Put dstore.Insert kind:%q, name:%q, exist:%v.", kind, name, exist )
		log.Printf( "before Put src:%v", src) 
		if _, err := tx.Put( key, src ); err != nil {
			log.Printf( "after Put dstore.Insert kind:%q, name:%q, exist:%v.", kind, name, exist )
			log.Printf( "after Put src:%v.", src )
			return err
		}
		return nil
	} )

	log.Printf( "End dstore.Insert kind:%q, name:%q, exist:%v, err:%v.", kind, name, exist, err )
	return err, exist
}
*/

// property???value???????????????????????????Insert??????
func ( ds *DStore ) InsertUnique( kind, property, value string, src interface{} ) ( int64, error ) {
	log.Printf( "Start dstore.InsertUnique kind:%q, property:%q, value:%q, src:%v.", kind, property, value, src )

	var keys []*datastore.Key
	exist := false
	_, err := ds.Client.RunInTransaction( ds.Ctx, func(tx *datastore.Transaction ) error {
		// ????????????????????????????????????????????????????????????1?????????????????????????????????????????????
		q := datastore.NewQuery( kind) .Filter( property + " =", value ).KeysOnly().Limit( 1 )
		var terr error
		keys, terr = ds.Client.GetAll( ds.Ctx, q, nil )
		if terr != nil {
			return fmt.Errorf( "dstore.InsertUnique client.GetAll1 failure terr:%v.", terr )
		} else if len( keys ) == 0 {
			key := datastore.IncompleteKey( kind, nil )
			// If there was no matching entity, store it now.
			log.Printf( "before Put dstore.InsertUnique kind:%q property:%q value:%q terr:%v.", kind, property, value, terr )
			log.Printf( "before Put src:%v", src )
			_, terr = tx.Put( key, src )
			log.Printf( "after Put dstore.InsertUnique kind:%q property:%q value:%q terr:%v.", kind, property, value, terr )
			log.Printf( "after Put src:%v.", src )
			if terr != nil {
				return fmt.Errorf( "dstore.InsertUnique tx.Put failure terr:%v.", terr )
			}
			return nil
		} else {
			exist = true
			return fmt.Errorf( "dstore.InsertUnique:Already Exist Entity:%q.", value )
		}
	})

	log.Printf( "dstore.InsertUnique client.RunInTransaction kind:%q, property:%q, value:%q, src err:%v.", kind, property, value, src, err )
	if exist {
		return -1, nil	// ????????????????????????????????????????????????????????????-1?????????
	} else if err != nil {
		return -1, fmt.Errorf( "dstore.InsertUnique err:%vq.", err )
	}

	return keys[ 0 ].ID, nil
}

// value?????????property??????????????????????????????????????????
func ( ds *DStore ) QueryUnique( kind, property, value string, src interface{} ) ( int64, error ) {
	log.Printf( "Start dstore.QueryUnique kind:%q, property:%q, value:%q, src:%+v.", kind, property, value, src )

	// ????????????????????????????????????????????????????????????2?????????????????????????????????????????????
	q := datastore.NewQuery( kind ).Filter( property + " =", value ).Limit( 2 )

	if src == nil {
		q = q.KeysOnly()
	}
	if keys, err := ds.Client.GetAll( ds.Ctx, q, src );err != nil {
		return -1, fmt.Errorf( "dstore.QueryUnique err:%v.", err )
	} else {
		num := len( keys )
		if num == 1 {
			return keys[0].ID, nil
		} else if num == 0 {
			return -1, nil
		} else {
			return -1, fmt.Errorf( "dstore.QueryUnique:Invalid num:%d kind:%q, property:%q, value:%q, err:%v.", num, kind, property, value, err )
		}
	}
}

// Key?????????????????????????????????2??????kind?????????
func ( ds *DStore ) IncompleteKeyDoublePut( kind1, kind2 string, src1, src2 interface{} ) ( int64, error ) {
	key1 := datastore.IncompleteKey( kind1, nil )

	_, err := ds.Client.RunInTransaction( ds.Ctx, func(tx *datastore.Transaction ) error {
		var err error
		if key1, err = ds.Client.Put( ds.Ctx, key1, src1 ); err != nil {
			return fmt.Errorf( "dstore.IncompleteKeyDoublePut:Failed to put(1): %#v.", err )
		}
		key2 := datastore.IDKey( kind2, key1.ID, nil )	// ?????????Put?????????????????????key1????????????key2???????????????
		log.Printf( "dstore.IncompleteKeyDoublePut: key1:%v key2:%v.", key1, key2 )
		if _, err = ds.Client.Put( ds.Ctx, key2, src2 ); err != nil {
			return fmt.Errorf( "dstore.IncompleteKeyDoublePut:Failed to put(2): %#v.", err )
		}
		return nil
	})

	if err != nil {
		return -1, fmt.Errorf( "dstore.IncompleteKeyDoublePut:client.RunInTransaction failed.\n\t%v", err )
	}

	return key1.ID, nil
}

// transaction?????????????????????????????????Put
func ( ds *DStore ) PutKinds( ar []PutStrct ) error {
log.Printf( "Start dstore.PutKinds ar:%v", ar )
defer log.Printf( "End dstore.PutKinds ar:%v", ar )
	_, err := ds.Client.RunInTransaction( ds.Ctx, func(tx *datastore.Transaction ) error {
		for i, v := range ar {
log.Printf( "dstore.PutKinds i:%d v.Key:#v v.Src:%v", i, v.Key, v.Src )
			if _, err := ds.Client.Put( ds.Ctx, v.Key, v.Src ); err != nil {
log.Printf( "dstore.PutKinds failure i:%d v.Key:#v v.Src:%v err:%v", i, v.Key, v.Src, err )
				return fmt.Errorf( "dstore.IncompleteKeyDoublePut:Failed to datastore.client.Put(%d): %#v.", i, err )
			}
log.Printf( "dstore.PutKinds success i:%d v.Key:#v v.Src:%v", i, v.Key, v.Src )
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf( "dstore.PutKinds:client.RunInTransaction failed.\n\t%v", err )
	}

	return nil
}

func ( ds *DStore ) GetAll( kind, sort string, parentKey *datastore.Key, dst interface{} ) error {
	q := datastore.NewQuery( kind )
	if parentKey != nil {
		q = q.Ancestor( parentKey )
	}
	if len( sort ) > 0 {
		q = q.Order( sort )
//		q = q.Order( "-" + sort )
	}
	if _, err := ds.Client.GetAll( ds.Ctx, q, dst ); err != nil {
		return fmt.Errorf( "dstore.DStore.GetAll:datastore.Client.GetAll failure.\n\t%v", err )
	}
	return nil
}
