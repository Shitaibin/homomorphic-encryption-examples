package service

import (
	"log"
	"os"

	"github.com/astaxie/beego/logs"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 跟区块链交互的客户端
var CLI *Client

const (
	// org1CfgPath = "./../sdk-config/org1sdk-config.yaml"
	org1CfgPath = "/Users/shitaibin/Workspace/homomorphic-encryption-examples/lattigo-examples/transfer/chaincode-example/api-server/sdk-config/org1sdk-config.yaml"
	org2CfgPath = "/Users/shitaibin/Workspace/homomorphic-encryption-examples/lattigo-examples/transfer/chaincode-example/api-server/sdk-config/org2sdk-config.yaml"
)

var (
	peer0Org1 = "peer0.org1.example.com"
	peers     = []string{peer0Org1}
)

type Client struct {
	// Fabric network information
	ConfigPath string
	OrgName    string
	OrgAdmin   string
	OrgUser    string

	// sdk clients
	SDK *fabsdk.FabricSDK
	rc  *resmgmt.Client
	cc  *channel.Client

	// Same for each peer
	ChannelID string
	CCID      string // chaincode ID, eq name
	CCPath    string // chaincode source path, 是GOPATH下的某个目录
	CCGoPath  string // GOPATH used for chaincode
}

func NewOrg1Peer0Client() *Client {
	logs.Info("Create fabric client for org1 (bank001)")

	return NewClient(org1CfgPath, "Org1", "Admin", "User1")
}

func NewOrg2Peer0Client() *Client {
	logs.Info("Create fabric client for org2 (bank002)")

	return NewClient(org2CfgPath, "Org2", "Admin", "User1")
}

func NewClient(cfg, org, admin, user string) *Client {
	c := &Client{
		ConfigPath: cfg,
		OrgName:    org,
		OrgAdmin:   admin,
		OrgUser:    user,

		CCID:      "example4",
		CCPath:    "github.com/hyperledger/fabric/fabric-samples/chaincode/chaincode_example02/go/", // 相对路径是从GOPAHT/src开始的
		CCGoPath:  os.Getenv("GOPATH"),
		ChannelID: "mychannel",
	}

	// create sdk
	sdk, err := fabsdk.New(config.FromFile(c.ConfigPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	c.SDK = sdk
	log.Println("Initialized fabric sdk")

	c.rc, c.cc = NewSdkClient(sdk, c.ChannelID, c.OrgName, c.OrgAdmin, c.OrgUser)

	return c
}

// NewSdkClient create resource client and channel client
func NewSdkClient(sdk *fabsdk.FabricSDK, channelID, orgName, orgAdmin, OrgUser string) (rc *resmgmt.Client, cc *channel.Client) {
	var err error

	// create rc
	rcp := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	rc, err = resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}
	log.Println("Initialized resource client")

	// create cc
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(OrgUser))
	cc, err = channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	log.Println("Initialized channel client")

	return rc, cc
}

// RegisterChaincodeEvent more easy than event client to registering chaincode event.
func (c *Client) RegisterChaincodeEvent(ccid, eventName string) (fab.Registration, <-chan *fab.CCEvent, error) {
	return c.cc.RegisterChaincodeEvent(ccid, eventName)
}

func (c *Client) Close() {
	c.SDK.Close()
}
