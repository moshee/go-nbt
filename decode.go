package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"fmt"

//	"io/ioutil"
)

var (
	ErrInvalidTag   = errors.New("Invalid tag")
	ErrNotCompound  = errors.New("Invalid NBT file: root node is not a TAG_Compound")
	ErrStoppedShort = errors.New("Unexpected TAG_End")
	ErrTruncated    = errors.New("Unexpected EOF")
)

// Decodes a gzipped NBT file into a native Go structure.
func DecodeGzip(src io.Reader) (*Compound, error) {
	buf := new(bytes.Buffer)
	r, err := gzip.NewReader(src)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	io.Copy(buf, r)
	return Decode(buf)
}


// Decodes an NBT file into a native Go structure.
func Decode(src io.Reader) (*Compound, error) {
	var tag byte
	read(&tag, src)
	if tag != TagCompound {
		return nil, ErrNotCompound
	}

	name := read_string(src)
	return read_compound(src, name, nil)
}

func read(dest interface{}, src io.Reader) error {
	return binary.Read(src, binary.BigEndian, dest)
}

func read_string(src io.Reader) string {
	var strlen int16
	read(&strlen, src)
	str := make([]byte, strlen)
	read(str, src)
	return string(str)
}

func read_compound(src io.Reader, name string, parent *Compound) (*Compound, error) {
	current := &Compound{
		parent: parent,
		name:   name,
		data:   make(map[string]interface{}),
	}
	root := current

	var tag byte
	for {
		read(&tag, src)
		println("reading tag", tag)
		switch tag {
		case TagEnd:
			if current.parent == nil {
				return root, nil
			} else {
				current = current.parent
			}

		case TagByte:
			var value int8
			current.store(&value, src)

		case TagShort:
			var value int16
			current.store(&value, src)

		case TagInt:
			var value int32
			current.store(&value, src)

		case TagLong:
			var value int64
			current.store(&value, src)

		case TagFloat:
			var value float32
			current.store(&value, src)

		case TagDouble:
			var value float64
			current.store(&value, src)

		case TagByteArray:
			name := read_string(src)
			var length int32
			read(&length, src)
			bytea := make([]int8, length)
			read(bytea, src)
			current.data[name] = bytea

		case TagString:
			name := read_string(src)
			data := read_string(src)
			current.data[name] = data

		case TagList:
			list, err := read_list(src)
			if err != nil {
				return root, err
			}
			current.data[list.name] = list

		case TagCompound:
			// we need to go deeper
			// Create a NEW Compound pointer which will be the recipient of any
			// further calls to (*Compound).store. Once a TAG_End is reached,
			// appropriate action will be taken to move the target back to this
			// *Compound's parent.
			name := read_string(src)
			c := &Compound{
				parent: current,
				name:   name,
				data:   make(map[string]interface{}),
			}
			current.data[name] = c
			current = c

		case TagIntArray:
			// I'll assume for now that the length is also a signed int, like
			// TAG_ByteArray
			name := read_string(src)
			var length int32
			read(&length, src)
			inta := make([]int32, length)
			read(inta, src)
			current.data[name] = inta

		default:
			return root, errors.New(fmt.Sprintf("Unknown type: %v", tag))
		}
	}

	// not enough TAG_Ends, reached EOF already
	return root, ErrTruncated
}

func read_list(src io.Reader) (*List, error) {
	name := read_string(src)
	var list_type byte
	read(&list_type, src)
	var length int32
	read(&length, src)
	list := &List{
		name:      name,
		list_type: list_type,
		length:    length,
	}

	switch list_type {
	case TagCompound:
		data := make([]*Compound, length)
		for k, _ := range data {
			c, err := read_compound(src, "", nil)
			if err != nil {
				return nil, err
			}
			data[k] = c
		}
		list.data = data

	case TagByte:
		data := make([]int8, length)
		read(data, src)
		list.data = data

	case TagShort:
		data := make([]int16, length)
		read(data, src)
		list.data = data

	case TagInt:
		data := make([]int32, length)
		read(data, src)
		list.data = data

	case TagLong:
		data := make([]int64, length)
		read(data, src)
		list.data = data

	case TagFloat:
		data := make([]float32, length)
		read(data, src)
		list.data = data

	case TagDouble:
		data := make([]float64, length)
		read(data, src)
		list.data = data

	default:
		panic(fmt.Sprintf("%#v", list_type))
	}
	return list, nil
}
