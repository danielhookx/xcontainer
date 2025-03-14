package xmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"reflect"
	"strconv"

	"github.com/danielhookx/xcontainer/list"
)

// MapNode represents a key-value pair in the OrderedMap.
type MapNode[K comparable, V any] struct {
	K K
	V V
}

// OrderedMap is a map that preserves the order of key-value pairs.
// It combines a map for fast lookups and a linked list for order preservation.
type OrderedMap[K comparable, V any] struct {
	kv map[K]*list.Element[*MapNode[K, V]]
	l  *list.List[*MapNode[K, V]]
}

// NewOrderedMap creates and initializes a new OrderedMap.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		kv: make(map[K]*list.Element[*MapNode[K, V]]),
		l:  list.New[*MapNode[K, V]](),
	}
}

// Get retrieves a value from the map by key.
// Returns the value and a boolean indicating whether the key was found.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	e, ok := m.kv[key]
	if ok {
		return e.Value.V, true
	}
	var v V
	return v, false
}

// Set adds or updates a key-value pair in the map.
// Returns true if the key is new, false if it already existed.
func (m *OrderedMap[K, V]) Set(key K, value V) bool {
	_, ok := m.kv[key]
	if ok {
		m.kv[key].Value.V = value
		return false
	}

	e := m.l.PushBack(&MapNode[K, V]{K: key, V: value})
	m.kv[key] = e
	return true
}

// Delete removes a key-value pair from the map.
// Returns true if the key was found and removed, false otherwise.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	e, ok := m.kv[key]
	if ok {
		m.l.Remove(e)
		delete(m.kv, key)
	}
	return ok
}

// Len returns the number of key-value pairs in the map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.kv)
}

// Copy creates a deep copy of the OrderedMap.
func (m *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	m2 := NewOrderedMap[K, V]()

	for e := m.l.Front(); e != nil; e = e.Next() {
		m2.Set(e.Value.K, e.Value.V)
	}

	return m2
}

// ToArray returns all values in the map as a slice, preserving the order.
func (m *OrderedMap[K, V]) ToArray() []V {
	array := make([]V, 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		array = append(array, e.Value.V)
	}
	return array
}

// Iter returns an iterator that yields key-value pairs in insertion order.
// The iterator creates a snapshot of the map to ensure consistent iteration.
func (m *OrderedMap[K, V]) Iter() iter.Seq2[K, V] {
	// Create a snapshot at the beginning of iteration
	snapshot := make([]*MapNode[K, V], 0, m.Len())
	for e := m.l.Front(); e != nil; e = e.Next() {
		snapshot = append(snapshot, e.Value)
	}

	return func(yield func(K, V) bool) {
		// Use the snapshot for iteration
		for _, node := range snapshot {
			if !yield(node.K, node.V) {
				return
			}
		}
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes a JSON object into an OrderedMap, preserving the order of keys.
func (m *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	// Clear the current map
	m.l = list.New[*MapNode[K, V]]()
	m.kv = make(map[K]*list.Element[*MapNode[K, V]])

	// Return early if input is empty
	if len(data) == 0 {
		return nil
	}

	// Create decoder and enable UseNumber
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	// Read the opening brace
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected {, got %v", t)
	}

	// Read key-value pairs one by one, preserving order
	for dec.More() {
		// Read the key
		t, err := dec.Token()
		if err != nil {
			return err
		}
		keyStr, ok := t.(string)
		if !ok {
			return fmt.Errorf("expected string key, got %v", t)
		}

		// Read the value
		var value any
		if err := dec.Decode(&value); err != nil {
			return err
		}

		// Process the value
		processedValue := processValue(value)

		// Convert string key to K type and set the key-value pair
		if err := m.setKeyFromString(keyStr, processedValue); err != nil {
			return fmt.Errorf("failed to convert key '%s': %v", keyStr, err)
		}
	}

	// Read the closing brace
	t, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected }, got %v", t)
	}

	return nil
}

// setKeyFromString converts a string key to K type and sets the key-value pair.
func (m *OrderedMap[K, V]) setKeyFromString(keyStr string, value any) error {
	var key K
	if err := convertKey(keyStr, &key); err != nil {
		return err
	}

	var v V
	if value != nil {
		v = any(value).(V)
	}
	m.Set(key, v)
	return nil
}

// convertKey converts a string key to the specified type K.
// Supports basic types like string, numeric types, bool, and attempts to use
// json.Unmarshal for other types.
func convertKey[K any](keyStr string, key *K) error {
	v := reflect.ValueOf(key).Elem()

	switch v.Kind() {
	case reflect.String:
		v.SetString(keyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(keyStr, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(keyStr, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(keyStr, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(keyStr)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		// For other types, try using json.Unmarshal
		return json.Unmarshal([]byte(`"`+keyStr+`"`), key)
	}
	return nil
}

// processValue recursively processes JSON values.
// Handles json.Number, arrays, and maps to ensure proper type conversion.
func processValue(v any) any {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case json.Number:
		if intVal, err := val.Int64(); err == nil {
			return float64(intVal)
		}
		if floatVal, err := val.Float64(); err == nil {
			return floatVal
		}
		return val.String()
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			result[i] = processValue(item)
		}
		return result
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[k] = processValue(v)
		}
		return result
	default:
		return v
	}
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes an OrderedMap into a JSON object, preserving the order of keys.
func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	// Create a buffer to build the JSON
	var buf bytes.Buffer
	buf.WriteByte('{')

	// Iterate through the OrderedMap, serializing key-value pairs in insertion order
	first := true
	for e := m.l.Front(); e != nil; e = e.Next() {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// Serialize the key - all key types are converted to strings
		keyStr := fmt.Sprintf("%v", e.Value.K)
		keyBytes, err := json.Marshal(keyStr)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key: %v", err)
		}
		buf.Write(keyBytes)

		buf.WriteByte(':')

		// Serialize the value
		valueBytes, err := json.Marshal(e.Value.V)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value: %v", err)
		}
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}
