package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// File defines the structure of a config file.
type File struct {
	Nodes   map[string]Node   `json:"nodes,omitempty"`
	Devices map[string]Device `json:"devices,omitempty"`
}

// Normalize moves any inline references into top level maps and replaces them with NameOnly references.
func (f *File) Normalize() error {
	if f == nil {
		return nil
	}
	if f.Devices == nil {
		f.Devices = make(map[string]Device)
	}
	if f.Nodes == nil {
		f.Nodes = make(map[string]Node)
	}

	for name, device := range f.Devices {
		if err, _ := device.Node.Normalize(f.Nodes); err != nil {
			return fmt.Errorf("device %v node %w", name, err)
		}
	}

	return nil
}

// GetNode calls ref.Get using f.Nodes.
func (f File) GetNode(ref NodeRef) (Node, bool) {
	return ref.Get(f.Nodes)
}

// GetDevice calls ref.Get using f.Devices.
func (f File) GetDevice(ref DeviceRef) (Device, bool) {
	return ref.Get(f.Devices)
}

// Read reads a File from the json content exposed via r into dst.
func Read(r io.Reader, dst *File) error {
	dec := json.NewDecoder(r)
	return dec.Decode(dst)
}

// ReadFile calls os.Open then Read on the opened file.
func ReadFile(path string, dst *File) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return Read(file, dst)
}

// Write write f to w as formatted JSON.
func (f File) Write(w io.Writer, opts ...WriteOption) error {
	enc := json.NewEncoder(w)
	for _, opt := range opts {
		opt(enc)
	}
	return enc.Encode(f)
}

// WriteFile opens a file at the given path and calls Write using it.
func (f File) WriteFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return f.Write(file)
}

type WriteOption func(enc *json.Encoder)

func WithIndent(prefix, indent string) WriteOption {
	return func(enc *json.Encoder) {
		enc.SetIndent(prefix, indent)
	}
}
