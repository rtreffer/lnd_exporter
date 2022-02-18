package lnc

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"

	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

type Client grpc.ClientConn

func NewClient(addr, crt, macroon string) (*Client, error) {
	usr, _ := user.Current()
	lndAddress, present := os.LookupEnv("LND_ADDRESS")
	if !present {
		lndAddress = "//127.0.0.1:10009"
	}
	if addr != "" {
		lndAddress = addr
	}
	certPath, present := os.LookupEnv("CERT_PATH")
	if !present {
		certPath = path.Join(usr.HomeDir, ".lnd/tls.cert")
	}
	if crt != "" {
		certPath = crt
	}
	macaroonPath, present := os.LookupEnv("MACAROON_PATH")
	if !present {
		macaroonPath = path.Join(usr.HomeDir, ".lnd/data/chain/bitcoin/mainnet/admin.macaroon")
	}
	if macroon != "" {
		macaroonPath = macroon
	}

	macaroonBytes, err := ioutil.ReadFile(macaroonPath)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(macaroonBytes)
	if err != nil {
		return nil, err
	}

	constrainedMac, err := macaroons.AddConstraints(mac, macaroons.TimeoutConstraint(60))
	if err != nil {
		return nil, err
	}

	cred, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(lndAddress)
	if err != nil {
		return nil, err
	}

	macaroon, err := macaroons.NewMacaroonCredential(constrainedMac)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithPerRPCCredentials(macaroon),
		grpc.WithContextDialer(lncfg.ClientAddressDialer(u.Port())),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800)),
	}

	conn, err := grpc.Dial(u.Hostname(), opts...)
	if err != nil {
		return nil, err
	}

	return (*Client)(conn), nil
}

func (c *Client) Close() error {
	return ((*grpc.ClientConn)(c)).Close()
}
