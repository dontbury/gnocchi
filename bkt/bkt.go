package bkt

import (
	"fmt"
	"log"
	"bytes"
	"strings"
	"context"

	"cloud.google.com/go/storage"
//	"golang.org/x/net/context"
)

type Bucket struct {
	ctx context.Context
	client *storage.Client
	bucket *storage.BucketHandle
	root string
}

func Create( project, root string ) ( *Bucket, error ) {
	bckt := &Bucket{
		ctx: context.Background(),
		client: nil,
		bucket: nil,
		root: root,
	}

	var err error
	bckt.client, err = storage.NewClient( bckt.ctx )
	if err != nil {
		return bckt, fmt.Errorf( "bsgp.CreateEnv:storage.NewClient failure err:%v", err )
	}

	bckt.bucket = bckt.client.Bucket( project )

	return bckt, nil
}


func ( bckt *Bucket ) Close() {
	if bckt.client != nil {
		bckt.client.Close()
	}
}

func ( bckt *Bucket ) Read( path string ) ( string, error ) {
	r, err := bckt.bucket.Object( bckt.root + "/" + path ).NewReader( bckt.ctx )
	if err != nil {
		return "", fmt.Errorf( "bkt.Read:ObjectHandle.NewReader failure bckt.root:%q path:%q.\n\t%v", bckt.root, path, err )
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom( r ); err != nil {
		log.Printf( "bkt.Read bytes.Buffer.Readfrom err:%v.", err )
	    return "", fmt.Errorf( "bkt.Read:bytes.Buffer.Readfrom failure bckt.root:%q path:%q.\n\t%v", bckt.root, path, err )
	}

	return string( buf.Bytes() ), nil
}

func ( bckt *Bucket ) ReadLine( path string ) ([]string, error) {
	buf, err := bckt.Read( path )
	if err != nil {
		return nil, fmt.Errorf( "bkt.ReadLine:bskt.Read failure path:%q.\n\t%v", path, err )
	}

	s := strings.Split( string( buf ), "\n" )
	t, n := make( []string, len( s ) ), 0
	for _, v := range s {
		r := strings.Trim( v, " \r\n" )
		if len( r ) > 0 { // 空白行（末尾など）を除く
			if r[ 0 ] == '#' {  // 行頭が#ならコメント行なので参照しない
				continue
			}
			t[ n ] = r
			n++
		}
	}

	return t[ :n ], nil
}
