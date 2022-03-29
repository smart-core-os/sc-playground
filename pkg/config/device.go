package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Device struct {
	Name   string        `json:"name,omitempty"`
	Node   *NodeRef      `json:"node,omitempty"`
	Traits []TraitOrName `json:"traits,omitempty"`
}

func (d Device) IsLocal() bool {
	if d.Node == nil {
		return true
	}
	if d.Node.NameOnly {
		return false
	}
	return d.Node.Address == ""
}

type DeviceRef struct {
	Device
	NameOnly bool `json:"-"`
}

// Get resolves this reference using the given dict or it's own information if provided.
func (r DeviceRef) Get(dict map[string]Device) (Device, bool) {
	if r.NameOnly {
		found, ok := dict[r.Name]
		return found, ok
	}
	return r.Device, true
}

// Normalize converts r into a NameOnly ref, placing inline data into dict.
func (r *DeviceRef) Normalize(dict map[string]Device) (error, bool) {
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
	dict[name] = r.Device
	r.Device = Device{Name: name}
	return nil, true
}

func (r *DeviceRef) UnmarshalJSON(bytes []byte) error {
	var name string
	if err := json.Unmarshal(bytes, &name); err == nil {
		r.Name = name
		r.NameOnly = true
		return nil
	}
	return json.Unmarshal(bytes, &r.Device)
}

func (r *DeviceRef) MarshalJSON() ([]byte, error) {
	if r.NameOnly {
		return json.Marshal(r.Name)
	}
	return json.Marshal(r.Device)
}
