// SPDX-FileCopyrightText: 2023 Free Mobile
// SPDX-License-Identifier: AGPL-3.0-only

package flow

import (
	"net"
	"path/filepath"
	"testing"

	"github.com/gogo/protobuf/proto"

	"akvorado/common/helpers"
	"akvorado/common/reporter"
	"akvorado/inlet/flow/decoder"
	"akvorado/inlet/flow/decoder/netflow"
	"akvorado/inlet/flow/decoder/sflow"
)

// The goal is to benchmark flow decoding + encoding to protobuf

func BenchmarkDecodeEncodeNetflow(b *testing.B) {
	r := reporter.NewMock(b)
	nfdecoder := netflow.New(r)

	template := helpers.ReadPcapPayload(b, filepath.Join("decoder", "netflow", "testdata", "options-template-257.pcap"))
	got := nfdecoder.Decode(decoder.RawFlow{Payload: template, Source: net.ParseIP("127.0.0.1")})
	if got == nil || len(got) != 0 {
		b.Fatalf("Decode() error on options template")
	}
	data := helpers.ReadPcapPayload(b, filepath.Join("decoder", "netflow", "testdata", "options-data-257.pcap"))
	got = nfdecoder.Decode(decoder.RawFlow{Payload: data, Source: net.ParseIP("127.0.0.1")})
	if got == nil || len(got) != 0 {
		b.Fatalf("Decode() error on options data")
	}
	template = helpers.ReadPcapPayload(b, filepath.Join("decoder", "netflow", "testdata", "template-260.pcap"))
	got = nfdecoder.Decode(decoder.RawFlow{Payload: template, Source: net.ParseIP("127.0.0.1")})
	if got == nil || len(got) != 0 {
		b.Fatalf("Decode() error on template")
	}
	for i := 0; i < b.N; i++ {
		data = helpers.ReadPcapPayload(b, filepath.Join("decoder", "netflow", "testdata", "data-260.pcap"))
		got = nfdecoder.Decode(decoder.RawFlow{Payload: data, Source: net.ParseIP("127.0.0.1")})
		for _, flow := range got {
			buf := proto.NewBuffer([]byte{})
			if err := buf.EncodeMessage(flow); err != nil {
				b.Fatalf("EncodeMessage() error:\n%+v", err)
			}
		}
	}
}

func BenchmarkDecodeEncodeSflow(b *testing.B) {
	r := reporter.NewMock(b)
	sdecoder := sflow.New(r)

	for i := 0; i < b.N; i++ {
		data := helpers.ReadPcapPayload(b, filepath.Join("decoder", "sflow", "testdata", "data-1140.pcap"))
		got := sdecoder.Decode(decoder.RawFlow{Payload: data, Source: net.ParseIP("127.0.0.1")})
		for _, flow := range got {
			buf := proto.NewBuffer([]byte{})
			if err := buf.EncodeMessage(flow); err != nil {
				b.Fatalf("EncodeMessage() error:\n%+v", err)
			}
		}
	}
}