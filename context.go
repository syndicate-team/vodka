package vodka

import "net/http"

/*
Context - context that contains information about request/response and middlewares results
*/
type Context struct {
	Raw     RawContext
	Query   KeyStorage
	Params  KeyStorage
	Body    KeyStorage
	Options KeyStorage

	Handler     Handler
	HandlerFunc HandlerFunc
	iterator    int
	Next        func(*Context)
	Request     *http.Request
	Writer      http.ResponseWriter
	Validation  interface{}
}

// RawContext - raw context struct to save raw data
type RawContext struct {
	Query  KeyStorage
	Params KeyStorage
	Body   []byte
}

// Validation - validation
type Validation struct {
	Query  interface{}
	Params interface{}
	Body   interface{}
}

// KeyStorage - main key storage for Validated and raw data
type KeyStorage struct {
	keys map[string]interface{}
}

func (q *KeyStorage) Get(key string) interface{} {
	return q.keys[key]
}
func (q *KeyStorage) GetInt(key string) int {
	return q.keys[key].(int)
}
func (q *KeyStorage) GetInt64(key string) int64 {
	return q.keys[key].(int64)
}
func (q *KeyStorage) GetFloat32(key string) float32 {
	return q.keys[key].(float32)
}
func (q *KeyStorage) GetFloat64(key string) float64 {
	return q.keys[key].(float64)
}
func (q *KeyStorage) GetString(key string) string {
	return q.keys[key].(string)
}
func (q *KeyStorage) GetBool(key string) bool {
	return q.keys[key].(bool)
}
func (q *KeyStorage) GetArray(key string) []interface{} {
	return q.keys[key].([]interface{})
}
func (q *KeyStorage) GetArrayInt(key string) []int64 {
	return q.keys[key].([]int64)
}

// GetArrayFloat - get data in []float64
func (q *KeyStorage) GetArrayFloat(key string) []float64 {
	return q.keys[key].([]float64)
}

// GetArrayString - get []string
func (q *KeyStorage) GetArrayString(key string) []string {
	return q.keys[key].([]string)
}

func (q *KeyStorage) GetBytes(key string) []byte {
	return q.keys[key].([]byte)
}

func (q *KeyStorage) Delete(key string) {
	delete(q.keys, key)
}
func (q *KeyStorage) Set(key string, value interface{}) {
	if q.keys == nil {
		q.keys = make(map[string]interface{})
	}
	q.keys[key] = value
}

func (q *KeyStorage) Map() map[string]interface{} {
	return q.keys
}
