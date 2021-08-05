package err2

type _Int struct{}
type _Int64 struct{}
type _UInt struct{}
type _UInt64 struct{}

// Int is a helper variable to generated
// 'type wrappers' to make Try function as fast as Check.
var Int _Int
var Int64 _Int64
var UInt _UInt
var UInt64 _UInt64

// Try is a helper method to call func() (int, error) functions
// with it and be as fast as Check(err).
func (o _Int) Try(v int, err error) int {
	Check(err)
	return v
}

func (o _Int64) Try(v int64, err error) int64 {
	Check(err)
	return v
}

func (o _UInt) Try(v uint, err error) uint {
	Check(err)
	return v
}

func (o _UInt64) Try(v uint64, err error) uint64 {
	Check(err)
	return v
}
