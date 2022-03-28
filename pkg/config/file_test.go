package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/smart-core-os/sc-api/go/traits"
	"github.com/smart-core-os/sc-golang/pkg/trait"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestRead(t *testing.T) {
	content := `{
	"nodes": {
		"node1": {"address": "1.1.1.1:1234"}
	},
	"devices": {
		"device1": {"node": "node1"},
		"device2": {"node": {"address": "1.1.1.2"}},
		"device3": {"traits": [
			"smartcore.traits.Electric", 
			{"name": "smartcore.traits.OnOff", "config": {"@type": "type.googleapis.com/smartcore.traits.OnOff", "state": "ON"}}
		]}
	}
}`
	want := File{
		Nodes: map[string]Node{
			"node1": {Address: "1.1.1.1:1234"},
		},
		Devices: map[string]Device{
			"device1": {Node: &NodeRef{NameOnly: true, Node: Node{Name: "node1"}}},
			"device2": {Node: &NodeRef{Node: Node{Address: "1.1.1.2"}}},
			"device3": {Traits: []TraitOrName{
				{Trait: Trait{Name: trait.Electric}, NameOnly: true},
				{Trait: Trait{Name: trait.OnOff, Config: traitConfig(&traits.OnOff{State: traits.OnOff_ON})}},
			}},
		},
	}

	var file File
	err := Read(strings.NewReader(content), &file)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, file, protocmp.Transform()); diff != "" {
		t.Fatalf("Read (-want, +got)\n%v", diff)
	}
}

func TestWrite(t *testing.T) {
	content := File{
		Nodes: map[string]Node{
			"node1": {Address: "1.1.1.1:1234"},
		},
		Devices: map[string]Device{
			"device1": {Node: &NodeRef{NameOnly: true, Node: Node{Name: "node1"}}},
			"device2": {Node: &NodeRef{Node: Node{Address: "1.1.1.2"}}},
			"device3": {Traits: []TraitOrName{
				{Trait: Trait{Name: trait.Electric}, NameOnly: true},
				{Trait: Trait{Name: trait.OnOff, Config: traitConfig(&traits.OnOff{State: traits.OnOff_ON})}},
			}},
		},
	}
	want := `{
	"nodes": {
		"node1": {
			"address": "1.1.1.1:1234"
		}
	},
	"devices": {
		"device1": {
			"node": "node1"
		},
		"device2": {
			"node": {
				"address": "1.1.1.2"
			}
		},
		"device3": {
			"traits": [
				"smartcore.traits.Electric",
				{
					"name": "smartcore.traits.OnOff",
					"config": {
						"@type": "type.googleapis.com/smartcore.traits.OnOff",
						"state": "ON"
					}
				}
			]
		}
	}
}
`

	out := new(bytes.Buffer)
	err := content.Write(out, WithIndent("", "\t"))
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, out.String()); diff != "" {
		t.Fatalf("Write (-want, +got)\n%v", diff)
	}
}

func traitConfig(src proto.Message) *TraitConfig {
	return (*TraitConfig)(mustAny(src))
}

func mustAny(msg proto.Message) *anypb.Any {
	a, err := anypb.New(msg)
	if err != nil {
		panic(err)
	}
	return a
}
