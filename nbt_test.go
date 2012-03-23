package nbt

import (
	"bytes"
	"io/ioutil"
	"testing"
)

/*
 TAG_Compound('Level'): 11 entries
    TAG_Compound('nested compound test'): 2 entries
		TAG_Compound('egg'): 2 entries
			TAG_String('name'): 'Eggbert'
			TAG_Float('value'): 0.5
		TAG_Compound('ham'): 2 entries
			TAG_String('name'): 'Hampus'
			TAG_Float('value'): 0.75
	TAG_Int('intTest'): 2147483647
	TAG_Byte('byteTest'): 127
	TAG_String('stringTest'): 'HELLO WORLD THIS IS A TEST STRING \xc3\x85\xc3\x84\xc3\x96!'
	TAG_List('listTest (long)'): 5 entires
		TAG_Long(None): 11
		TAG_Long(None): 12
		TAG_Long(None): 13
		TAG_Long(None): 14
		TAG_Long(None): 15
	TAG_Double('doubleTest'): 0.49312871321823148
	TAG_Float('floatTest'): 0.49823147058486938
	TAG_Long('longTest'): 9223372036854775807L
    TAG_List('listTest (compound)'): 2 entires
		TAG_Compound(None): 2 entries
			TAG_Long('created-on'): 1264099775885L
			TAG_String('name'): 'Compound tag #0'
		TAG_Compound(None): 2 entries
			TAG_Long('created-on'): 1264099775885L
			TAG_String('name'): 'Compound tag #1'
	TAG_Byte_Array('byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))'): [1000 bytes]
	TAG_Short('shortTest'): 32767
*/

func TestDecodeGzip(t *testing.T) {
	file, err := ioutil.ReadFile("bigtest.nbt")
	if err != nil {
		t.Fatal("Couldn't open bigtest.nbt:", err)
	}

	data, err := DecodeGzip(bytes.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	name := data.Name()
	if name != "Level" {
		t.Errorf("in (*Compound).Name(): expected 'Level', got %s", name)
		t.FailNow()
	}

	ham := data.Compound("nested compound test").Compound("ham").String("name")
	if ham != "Hampus" {
		t.Errorf("in /nested compound test/ham/name: expected 'Hampus', got %s", ham)
		t.FailNow()
	}

	list := data.List("listTest (long)")
	longs := list.Longs()
	if n := longs[3]; n != 14 {
		t.Errorf("in /listTest (long)[3]: expected 14, got %d", n)
		t.FailNow()
	}

	data.PrettyPrint()
}
