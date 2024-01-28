package persistence

/*
Serializable is a data object that should be persisted in the table

The core interface is involved Marshal and Unmarshal. These two methods should work together to first
serialize an object to its byte representation and later deserialize the same byte representation to
an in-memory copy of the original object. The ByteLength method give information about how many bytes
the Marshal method will produce and is used to optimize storage characteristics for the objects.

The method key is required as the primary key for the table and needs to be separately accessible for
the table to handle serialization properly. You could think of Key returning the ID of the object while
Marshal returns the actual value of the object.

The information returned by Deleted gives information whether the object has been deleted. This information
is relevant due to the tree-like structure of the table, where objects may be deleted in newer tables but
the object may exist in older tables.
*/
type Serializable interface {
	/*
		Marshal transforms the object to its byte representation.

		A returned error indicates that serialisation failed. Implementors should keep the number of failure
		cases for serialisation to a minimum as there currently is no plan to remedy failures and objects
		will be omitted.

		The implementation of ByteLength should always return the exact length of byte returned by Marshal
	*/
	Marshal() ([]byte, error)
	/*
		Unmarshal populates the object's fields from its byte representation.

		A returned error indicates a failure during deserialization. As with Marshal the number failure cases
		should be kept to a minimum. Failure during deserialization usually indicates a failure during serialization.
	*/
	Unmarshal([]byte) error

	/*
		Key returns the identifier for the object
	*/
	Key() []byte

	/*
		Returns whether the object is marked as deleted
	*/
	Deleted() bool

	/*
		ByteLength indicates the number of byte returned by the Marshal method

		This method is used to work around storage limits and should always work correctly with Marshal.
	*/
	ByteLength() (uint64, error)
}

/*
Iterator iterates through a data structure returning one object at a time.

Users should always check that an object is available with HasNext. When HasNext returns true, Next should
return the next object in the data structure. The behaviour of calling Next after HasNext returned false is
undefined.
*/
type Iterator[Value any] interface {
	/*
		Next returns the next object in the data structure when the prior call to HasNext returned true.
	*/
	Next() Value

	/*
		HasNext indicates whether another object is available when calling Next.
	*/
	HasNext() bool
}

/*
Table is readonly key-value data structure.

Values can be retrieved by calling Get with the respective key. When an object can be found with the
respective key it returned and error is nil. When the object does not exist NotFoundError is returned as the
error; when the object is deleted MarkedDeletedError is returned as the error.

Call close when the Table is no longer needed.
*/
type Table[Value Serializable] interface {
	/*
		Get returns the value associated with key

		When the object is found, the value is returned and error is nil. Otherwise, the object will be
		zero value of the type and an error is returned. When the object does not exist NotFoundError is
		returned as the error; when the object is deleted MarkedDeletedError is returned as the error.
	*/
	Get(key []byte) (Value, error)

	/*
		See io.Closer
	*/
	Close() error
}
