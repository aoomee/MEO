package handler

import (
	"context"

	"github.com/xtls/xray-core/app/proxyman/command"
	cnet "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/blackhole"
	"github.com/xtls/xray-core/proxy/dns"
	"github.com/xtls/xray-core/proxy/freedom"
	"github.com/xtls/xray-core/proxy/http"
	"github.com/xtls/xray-core/proxy/socks"
	"github.com/xtls/xray-core/transport/internet"
)

func AddFreedomOutbound(ctx context.Context, client command.HandlerServiceClient, tag string) error {
	cfg := outboundConfig(
		tag,
		senderSettings(),
		serial.ToTypedMessage(&freedom.Config{
			DomainStrategy: internet.DomainStrategy_AS_IS,
			UserLevel:      0,
			Fragment: &freedom.Fragment{
				PacketsFrom: 5,
				PacketsTo:   10,
				LengthMin:   50,
				LengthMax:   150,
				IntervalMin: 10,
				IntervalMax: 20,
			},
		}),
	)
	_, err := client.AddOutbound(ctx, &command.AddOutboundRequest{Outbound: cfg})
	return err
}

func AddBlackholeOutbound(ctx context.Context, client command.HandlerServiceClient, tag string) error {
	cfg := outboundConfig(
		tag,
		senderSettings(),
		serial.ToTypedMessage(&blackhole.Config{
			Response: serial.ToTypedMessage(&blackhole.HTTPResponse{}),
		}),
	)
	_, err := client.AddOutbound(ctx, &command.AddOutboundRequest{Outbound: cfg})
	return err
}

func AddDNSOutbound(ctx context.Context, client command.HandlerServiceClient, tag string, upstream string) error {
	endpointCfg := &cnet.Endpoint{
		Network: cnet.Network_UDP,
		Address: cnet.NewIPOrDomain(cnet.ParseAddress(upstream)),
		Port:    53,
	}
	cfg := outboundConfig(
		tag,
		senderSettings(),
		serial.ToTypedMessage(&dns.Config{
			Server:     endpointCfg,
			UserLevel:  0,
			BlockTypes: []int32{1, 28},
		}),
	)
	_, err := client.AddOutbound(ctx, &command.AddOutboundRequest{Outbound: cfg})
	return err
}

func AddHTTPOutbound(ctx context.Context, client command.HandlerServiceClient, tag string) error {
	cfg := outboundConfig(
		tag,
		senderSettings(),
		serial.ToTypedMessage(&http.ClientConfig{
			Server: endpoint("example.com", 80, nil),
			Header: []*http.Header{
				{Key: "User-Agent", Value: "miaomiaowu"},
			},
		}),
	)
	_, err := client.AddOutbound(ctx, &command.AddOutboundRequest{Outbound: cfg})
	return err
}

func AddSocksOutbound(ctx context.Context, client command.HandlerServiceClient, tag string) error {
	cfg := outboundConfig(
		tag,
		senderSettings(),
		serial.ToTypedMessage(&socks.ClientConfig{
			Server: endpoint("127.0.0.1", 1080, nil),
		}),
	)
	_, err := client.AddOutbound(ctx, &command.AddOutboundRequest{Outbound: cfg})
	return err
}
