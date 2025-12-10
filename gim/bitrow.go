package gim

import (
	"fmt"
	"math"
)

const BIT_NUM_IN_VALUE = 64

type BitRow struct {
	Body []int64
}

func ( b *BitRow ) GetBody() []int64 {
	return b.Body
}

func ( b *BitRow ) Size() int {
	return len( b.Body )
}

func ( b *BitRow ) Set( shift, size int, value int64 ) error {
	if shift < 0 {
		return fmt.Errorf( "gim.BitRow.Set shift(%d) is too small", shift )
	} else if size <= 0 {
		return fmt.Errorf( "gim.BitRow.Set size(%d) is too small", size )
	} else if size >= BIT_NUM_IN_VALUE {
		return fmt.Errorf( "gim.BitRow.Set size(%d) is too large", size )
	}
	hindex := ( shift + size ) / BIT_NUM_IN_VALUE
	lindex := shift / BIT_NUM_IN_VALUE
	if hindex > lindex {	// ビットごとの更新が配列の区切りにまたがってしまう場合
		lshift := shift % BIT_NUM_IN_VALUE
		lsize := BIT_NUM_IN_VALUE - lshift
		hsize := size - lsize
		lmask := int64( math.Pow( 2, float64( lsize ) ) ) - 1
		if err := b.set( hindex, 0, hsize, value >> lsize ); err != nil {	// sliceが拡張される場合に備えて、上位インデックスからセットしていく
			return fmt.Errorf( "gim.BitRow.Set:gim.BitRow.set failed b:%+v, shift:%d size:%d value:%d\t\n%v", b, shift, size, value, err )
		} else if err = b.set( lindex, lshift, lsize, value & lmask ); err != nil {
			return fmt.Errorf( "gim.BitRow.Set:gim.BitRow.set failed b:%+v, shift:%d size:%d value:%d\t\n%v", b, shift, size, value, err )
		}
	} else if err := b.set( hindex, shift % BIT_NUM_IN_VALUE, size, value ); err != nil {
		return fmt.Errorf( "gim.BitRow.Set:gim.BitRow.set failed b:%+v, shift:%d size:%d value:%d\t\n%v", b, shift, size, value, err )
	}
	return nil
}

func ( b *BitRow ) set( index, shift, size int, value int64 ) error {
	if index < 0 {
		return fmt.Errorf( "gim.BitRow.set index(%d) is too small", index )
	} else if shift < 0 {
		return fmt.Errorf( "gim.BitRow.set shift(%d) is too small", shift )
	} else if size <= 0 {
		return fmt.Errorf( "gim.BitRow.set size(%d) is too small", size )
	} else if shift + size >= BIT_NUM_IN_VALUE {
		return fmt.Errorf( "gim.BitRow.set shift(%d) size(%d) is too large", shift, size )
	} else if index > len( b.Body ) {
		add := make( []int64, index - len( b.Body ) + 1 )
		b.Body = append( b.Body, add... )
	} else if mask := int64( math.Pow( 2, float64( size ) ) ) - 1; value > mask {
		return fmt.Errorf( "gim.BitRow.set value(%d) mask(%d) shift(%d) size(%d) is too large", value, mask, shift, size )
	} else {
		b.Body[ index ] = ( b.Body[ index ] & ^( mask << shift ) ) | ( ( mask & value ) << shift )
	}
	return nil
}

func ( b *BitRow ) Get( shift, size int ) ( int64, error ) {
	if shift < 0 {
		return 0, fmt.Errorf( "gim.BitRow.Get shift(%d) is too small", shift )
	} else if size <= 0 {
		return 0, fmt.Errorf( "gim.BitRow.Get size(%d) is too small", size )
	} else if size >= BIT_NUM_IN_VALUE {
		return 0, fmt.Errorf( "gim.BitRow.Get size(%d) is too large", size )
	}
	hindex := ( shift + size ) / BIT_NUM_IN_VALUE
	lindex := shift / BIT_NUM_IN_VALUE
	if hindex > lindex {	// ビットごとの更新が配列の区切りにまたがってしまう場合
		lshift := shift % BIT_NUM_IN_VALUE
		lsize := BIT_NUM_IN_VALUE - lshift
		hsize := size - lsize
		var hval, lval int64
		var err error
		if hval, err = b.get( hindex, 0, hsize ); err != nil {
			return 0, fmt.Errorf( "gim.BitRow.Get:gim.BitRow.get failed b:%+v, shift:%d size:%d\t\n%v", b, shift, size, err )
		} else if lval, err = b.get( lindex, lshift, lsize ); err != nil {
			return 0, fmt.Errorf( "gim.BitRow.Get:gim.BitRow.get failed b:%+v, shift:%d size:%d\t\n%v", b, shift, size, err )
		} else {
			return ( hval << lshift ) + lval, nil
		}
	} else if val, err := b.get( hindex, shift % BIT_NUM_IN_VALUE, size ); err != nil {
		return 0, fmt.Errorf( "gim.BitRow.Get:gim.BitRow.get failed b:%+v, shift:%d size:%d\t\n%v", b, shift, size, err )
	} else {
		return val, nil
	}
}

func ( b *BitRow ) get( index, shift, size int ) ( int64, error ) {
	if index < 0 {
		return 0, fmt.Errorf( "gim.BitRow.get index(%d) is too small", index )
	} else if shift < 0 {
		return 0, fmt.Errorf( "gim.BitRow.get shift(%d) is too small", shift )
	} else if size <= 0 {
		return 0, fmt.Errorf( "gim.BitRow.get size(%d) is too small", size )
	} else if shift + size >= BIT_NUM_IN_VALUE {
		return 0, fmt.Errorf( "gim.BitRow.Get shift,size(%d) is too large", shift + size )
	}
	if index < len( b.Body ) {
		mask := (int64)( math.Pow( 2, float64( size ) ) ) - 1
		return ( b.Body[ index ] >> shift ) & mask, nil
	}
	return 0, nil
}
