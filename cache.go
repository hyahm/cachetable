package cachetable


type Table interface {
	Set(key string, value interface{},setkey string,setvalue interface{}) error
 	Add( interface{}) error
	Key(key string) error
	Get(key string, field string, value interface{}) (interface{}, error)
}