package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smart-core-os/sc-golang/pkg/trait"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

type Trait struct {
	Name   trait.Name   `json:"name,omitempty"`
	Config *TraitConfig `json:"config,omitempty"`
}

// TraitConfig is an anypb.Any that uses protojson to marshal/unmarshal from json.
type TraitConfig anypb.Any

func (t *TraitConfig) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*anypb.Any)(t))
}

func (t *TraitConfig) UnmarshalJSON(bytes []byte) error {
	return protojson.Unmarshal(bytes, (*anypb.Any)(t))
}

func (t *TraitConfig) UnmarshalProto(opts proto.UnmarshalOptions) (proto.Message, error) {
	if t == nil {
		return nil, nil
	}
	return anypb.UnmarshalNew((*anypb.Any)(t), opts)
}

// ProtoReflect implements proto.Message in terms of the underlying anypb.Any.
// It's useful to refer to TraitConfig as a proto.Message.
func (t *TraitConfig) ProtoReflect() protoreflect.Message {
	return (*anypb.Any)(t).ProtoReflect()
}

type TraitOrName struct {
	Trait
	NameOnly bool `json:"-"`
}

// Get resolves this reference using the given dict or it's own information if provided.
func (r TraitOrName) Get(dict map[trait.Name]Trait) (Trait, bool) {
	if r.NameOnly {
		found, ok := dict[r.Name]
		return found, ok
	}
	return r.Trait, true
}

// Normalize converts r into a NameOnly ref, placing inline data into dict.
func (r *TraitOrName) Normalize(dict map[trait.Name]Trait) (error, bool) {
	if r.NameOnly {
		return nil, false
	}
	name := r.Name
	if name == "" {
		return errors.New("no name"), false
	}
	_, ok := dict[name]
	if ok {
		return fmt.Errorf("duplicate %v", name), false
	}
	dict[name] = r.Trait
	r.Trait = Trait{Name: name}
	return nil, true
}

func (r *TraitOrName) UnmarshalJSON(bytes []byte) error {
	var name trait.Name
	if err := json.Unmarshal(bytes, &name); err == nil {
		r.Name = name
		r.NameOnly = true
		return nil
	}
	return json.Unmarshal(bytes, &r.Trait)
}

func (r *TraitOrName) MarshalJSON() ([]byte, error) {
	if r.NameOnly {
		return json.Marshal(r.Name)
	}
	return json.Marshal(r.Trait)
}
