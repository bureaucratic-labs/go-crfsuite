package crfsuite

/*
#cgo LDFLAGS: -lcrfsuite
#include <stdlib.h>
#include "crfsuite.h"
*/
import "C"

import (
	"unsafe"
)

type Feature struct {
	Key   string
	Value float32
}

type FeatureExtractor func(item string) []Feature

type Dictionary struct {
	Original *C.struct_tag_crfsuite_dictionary
}

// Obtain the number of strings in the dictionary.
func (d *Dictionary) Length() int {
	return int(C.DictionaryLength(d.Original))
}

// Assign and obtain the integer ID for the string.
func (d *Dictionary) Get(key string) int {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return int(C.DictionaryGet(d.Original, cKey))
}

// Obtain the integer ID for the string.
func (d *Dictionary) ToID(key string) int {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return int(C.DictionaryToID(d.Original, cKey))
}

type Model struct {
	Labels     Dictionary
	Attributes Dictionary
	Original   *C.struct_tag_crfsuite_model
}

func (m *Model) getAttributes() Dictionary {
	dictionary := C.GetModelAttributes(m.Original)
	return Dictionary{Original: dictionary}
}

func (m *Model) getLabels() Dictionary {
	dictionary := C.GetModelLabels(m.Original)
	return Dictionary{Original: dictionary}
}

func (m *Model) GetTagger() Tagger {
	tagger := C.GetModelTagger(m.Original)
	return Tagger{
		Labels:     &m.Labels,
		Attributes: &m.Attributes,
		Original:   tagger,
	}
}

func NewModelFromFile(path string) Model {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	model := C.NewModelFromFile(cPath)
	m := Model{Original: model}
	m.Labels = m.getLabels()
	m.Attributes = m.getAttributes()
	return m
}

type Tagger struct {
	Labels     *Dictionary
	Attributes *Dictionary
	Original   *C.struct_tag_crfsuite_tagger
}

func (t *Tagger) Set(inst Instance) {
	C.SetTaggerInstance(t.Original, inst.Original)
}

func (t *Tagger) Tag(items []string, extractor FeatureExtractor) {
	inst := NewInstance()
	for i := 0; i < len(items); i++ {
		item := NewItem()
		features := extractor(items[i])
		for i := 0; i < len(features); i++ {
			feature := features[i]
			id := t.Attributes.ToID(feature.Key) // TODO:
			attribute := NewAttribute(id, feature.Value)
			item.AddAttribute(attribute)
		}
		inst.AddItem(item, 0) // TODO:
	}
	if !inst.Empty() {
		t.Set(inst)
	}
}

type Attribute struct {
	Original unsafe.Pointer
}

func NewAttribute(id int, value float32) Attribute {
	attribute := C.NewAttribute(C.int(id), C.float(value))
	return Attribute{Original: unsafe.Pointer(&attribute)}
}

type Item struct {
	Original unsafe.Pointer
}

func (i *Item) AddAttribute(attr Attribute) {
	C.AppendAttributeToItem(i.Original, attr.Original)
}

func NewItem() Item {
	item := C.NewItem()
	return Item{Original: unsafe.Pointer(&item)}
}

type Instance struct {
	Original unsafe.Pointer
}

func (i *Instance) Empty() bool {
	status := int(C.EmptyInstance(i.Original))
	if status == 0 {
		return true
	} else {
		return false
	}
}

func (i *Instance) AddItem(item Item, label_id int) {
	C.AddItemToInstance(i.Original, item.Original, C.int(label_id))
}

func NewInstance() Instance {
	inst := C.NewInstance()
	return Instance{Original: unsafe.Pointer(&inst)}
}
