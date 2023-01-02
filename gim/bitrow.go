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

func ( b *BitRow ) Set( index, shift, size int, value int64 ) error {
	if shift + size >= BIT_NUM_IN_VALUE {
		return fmt.Errorf( "gim.BitRow.Set shift,size(%d) is too large.", shift + size )
	}
	if index >= len( b.Body ) {
		empty := make( []int64, index - len( b.Body ) + 1 )
		b.Body = append( b.Body, empty... )
	}
	mask := int64( math.Pow( 2, float64(size) ) ) - 1
	b.Body[ index ] = ( b.Body[ index ] & ^( mask << shift ) ) | ( ( mask & value ) << shift )
	return nil
}

func ( b *BitRow ) Get( index, shift, size int ) ( int64, error ) {
	if shift + size >= BIT_NUM_IN_VALUE {
		return 0, fmt.Errorf( "gim.BitRow.Get shift,size(%d) is too large.", shift + size )
	}
	if index < len( b.Body ) {
		mask := (int64)( math.Pow( 2, float64( size ) ) ) - 1
		return ( b.Body[ index ] >> shift ) & mask, nil
	}
	return 0, nil
}
