/*
Package nbt provides facilities to encode and decode NBT (Named Binary Tag) data structures. From the Minecraft Coalition Wiki (http://wiki.vg):
    ID   Name           Size    Description
    0    TAG_End        0       This tag serves no purpose but to signify the
                                end of an open TAG_Compound. In most libraries,
                                this type is abstracted away and never seen.
    1    TAG_Byte       1       A single signed byte
    2    TAG_Short      2       A single signed short
    3    TAG_Int        4       A single signed integer
    4    TAG_Long       8       A single signed long (typically long long in
                                C/C++)
    5    TAG_Float      4       A single IEEE-754 single-precision floating
                                point number
    6    TAG_Double     8       A single IEEE-754 double-precision floating
                                point number
    7    TAG_Byte_Array ...     A length-prefixed array of signed bytes. The
                                prefix is a signed integer (thus 4 bytes)
    8    TAG_String     ...     A length-prefixed UTF-8 string. The prefix is an
                                unsigned short (thus 2 bytes)
    9    TAG_List       ...     A list of nameless tags, all of the same type.
                                The list is prefixed with the Type ID of the
                                items it contains (thus 1 byte), and the
                                length of the list as a signed integer (a
                                further 4 bytes).
    10   TAG_Compound   ...     Effectively a list of a named tags
    11   TAG_Int_Array  ...     A length-prefixed array of signed integers. The
                                prefix is presumably a signed integer.
*/
package nbt

import (
	"fmt"
	"io"
	"strings"
)

const (
	TagEnd byte = iota
	TagByte
	TagShort
	TagInt
	TagLong
	TagFloat
	TagDouble
	TagByteArray
	TagString
	TagList
	TagCompound
	TagIntArray
)

// Compound represents an NBT TAG_Compound structure.
type Compound struct {
	name   string
	data   map[string]interface{}
	parent *Compound
}

func (c *Compound) store(data interface{}, src io.Reader) {
	name := read_string(src)
	read(data, src)
	c.data[name] = data
}

func (self *Compound) Byte(name string) int8          { return self.data[name].(int8) }
func (self *Compound) Short(name string) int16        { return self.data[name].(int16) }
func (self *Compound) Int(name string) int32          { return self.data[name].(int32) }
func (self *Compound) Long(name string) int64         { return self.data[name].(int64) }
func (self *Compound) Float(name string) float32      { return self.data[name].(float32) }
func (self *Compound) Double(name string) float64     { return self.data[name].(float64) }
func (self *Compound) Compound(name string) *Compound { return self.data[name].(*Compound) }
func (self *Compound) List(name string) *List         { return self.data[name].(*List) }
func (self *Compound) String(name string) string      { return self.data[name].(string) }
func (self *Compound) Name() string                   { return self.name }
func (self *Compound) Len() int                       { return len(self.data) }

// Recursively print the compound's contents
func (self *Compound) PrettyPrint() {
	self.pretty_print(0)
}
func (self *Compound) pretty_print(indent_level int) {
	fmt.Printf("%sCompound \"%s\" (%d entries):\n", strings.Repeat("    ", indent_level), self.name, len(self.data))
	indent_level++
	for k, v := range self.data {
		spaces := strings.Repeat("    ", indent_level)

		switch v.(type) {
		case *Compound:
			v.(*Compound).pretty_print(indent_level)

		case *List:
			l := v.(*List)
			fmt.Printf("%sList \"%s\" (%d entries):\n", spaces, k, l.Len())
			spaces += "    "

			switch l.list_type {
			case TagCompound:
				for _, c := range l.Compounds() {
					c.pretty_print(indent_level + 1)
				}

			case TagByte:
				for _, v := range l.Bytes() {
					print_item(v, spaces, "Byte")
				}

			case TagShort:
				for _, v := range l.Shorts() {
					print_item(v, spaces, "Short")
				}

			case TagInt:
				for _, v := range l.Ints() {
					print_item(v, spaces, "Int")
				}

			case TagLong:
				for _, v := range l.Longs() {
					print_item(v, spaces, "Long")
				}

			case TagFloat:
				for _, v := range l.Floats() {
					print_item(v, spaces, "Float")
				}

			case TagDouble:
				for _, v := range l.Doubles() {
					print_item(v, spaces, "Double")
				}

			case TagString:
				for _, v := range l.Strings() {
					print_item(v, spaces, "String")
				}
			}
		default:
			switch v.(type) {
			case *int8:
				fmt.Printf("%sByte \"%s\": %v\n", spaces, k, *v.(*int8))
			case *int16:
				fmt.Printf("%sShort \"%s\": %v\n", spaces, k, *v.(*int16))
			case *int32:
				fmt.Printf("%sInt \"%s\": %v\n", spaces, k, *v.(*int32))
			case *int64:
				fmt.Printf("%sLong \"%s\": %v\n", spaces, k, *v.(*int64))
			case *float32:
				fmt.Printf("%sFloat \"%s\": %v\n", spaces, k, *v.(*float32))
			case *float64:
				fmt.Printf("%sDouble \"%s\": %v\n", spaces, k, *v.(*float64))
			case *string:
				fmt.Printf("%sString \"%s\": %v\n", spaces, k, *v.(*string))
			case []int8:
				fmt.Printf("%sByte Array \"%s\": [%d]\n", spaces, k, len(v.([]int8)))
			case []int32:
				fmt.Printf("%sInt Array \"%s\": [%d]\n", spaces, k, len(v.([]int32)))
			}
		}
	}
}

func print_item(thing interface{}, spaces, kind string) {
	fmt.Printf("%s%s: %v\n", spaces, kind, thing)
}

// List represents an NBT TAG_List structure. 
type List struct {
	name      string
	list_type byte
	data      interface{}
	length    int32
}

func (self *List) ListType() byte         { return self.list_type }
func (self *List) Len() int               { return int(self.length) }
func (self *List) Bytes() []int8          { return self.data.([]int8) }
func (self *List) Shorts() []int16        { return self.data.([]int16) }
func (self *List) Ints() []int32          { return self.data.([]int32) }
func (self *List) Longs() []int64         { return self.data.([]int64) }
func (self *List) Floats() []float32      { return self.data.([]float32) }
func (self *List) Doubles() []float64     { return self.data.([]float64) }
func (self *List) Strings() []string      { return self.data.([]string) }
func (self *List) Compounds() []*Compound { return self.data.([]*Compound) }
