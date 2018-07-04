package fabric

import (
	"fmt"
	"time"
	"github.com/Sirupsen/logrus"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"

	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)


var (
	sdkConfig = config.FromFile("./config/config.yaml")
)


func createChannelClient(channelID, userName, orgName string) (*fabsdk.FabricSDK, *channel.Client, error) {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		return  nil, nil, errors.New(fmt.Sprintf("Failed to create new SDK: %s", err))
	}

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions
	client, err := channel.New(clientChannelContext)
	if err != nil {
		sdk.Close()
		return nil, nil, errors.New(fmt.Sprintf("Failed to create new channel client: %s", err))
	}

	return sdk, client, nil
}

func ChaincodeQuery(channelID, userName, orgName, chaincodeID, fcn string, queryArgs [][]byte) []byte {
	sdk, client, err := createChannelClient(channelID, userName, orgName)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil
	}
	defer sdk.Close()

	response, err := client.Query(channel.Request{ChaincodeID: chaincodeID, Fcn: fcn, Args: queryArgs}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to query funds: %s", err))
		return nil
	}

	return response.Payload
}

func ChaincodeExecute(channelID, userName, orgName, chaincodeID, fcn string, txArgs [][]byte) []byte {
	sdk, client, err := createChannelClient(channelID, userName, orgName)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil
	}
	defer sdk.Close()

	response, err := client.Execute(channel.Request{ChaincodeID: chaincodeID, Fcn: fcn, Args: txArgs}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to query funds: %s", err))
		return nil
	}

	return response.Payload
}

func QueryInstalledChaincode(orgName, orgAdmin, targetEndpoint string) string {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new SDK: %s", err))
		return ""
	}
	defer sdk.Close()

	orgAdminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(orgAdminContext)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new resource management client: %s", err))
		return ""
	}
	res, err := orgResMgmt.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(targetEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to query installed chaincodes: %s", err))
		return ""
	}
	return  res.String()
}

func CreateChannel(channelID, channelConfigPath, orgName, orgAdmin, ordererOrgName, ordererEndpoint, ordererAdmin string) bool {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new SDK: %s", err))
		return false
	}
	defer sdk.Close()

	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(ordererAdmin), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create channel management client: %s", err))
		return false
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		logrus.Errorln(err.Error())
		return false
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		logrus.Errorln(err.Error())
		return false
	}
	req := resmgmt.SaveChannelRequest{
		ChannelID: channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererEndpoint))

	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create channel: %s", err))
		return false
	}
	logrus.Info(fmt.Sprintf("Transaction ID: %s", txID))
	return true
}


func CreateChaincode(chaincodeID, version, chaincodePath, chaincodeGoPath, orgName, orgAdmin string) bool {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new SDK: %s", err))
		return false
	}
	defer sdk.Close()

	orgAdminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(orgAdminContext)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new resource management client: %s", err))
		return false
	}

	ccPkg, err := packager.NewCCPackage(chaincodePath, chaincodeGoPath)
	if err != nil {
		logrus.Errorln(err.Error())
		return false
	}
	// Install chaincode to org peers
	installCCReq := resmgmt.InstallCCRequest{
		Name: chaincodeID,
		Path: chaincodePath,
		Version: version,
		Package: ccPkg,
	}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logrus.Fatalln(fmt.Sprintf("Failed to install chaincode: %s", err))
		return false
	}
	return true
}

func InstantiateChaincode(channelID, chaincodeID, version, chaincodePath, orgName, orgAdmin string, ccInitArgs [][]byte) bool {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new SDK: %s", err))
		return false
	}
	defer sdk.Close()

	orgAdminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(orgAdminContext)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new resource management client: %s", err))
		return false
	}

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	_, err = orgResMgmt.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{
			Name: chaincodeID,
			Path: chaincodePath,
			Version: version,
			Args: ccInitArgs,
			Policy: ccPolicy,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to instantiate chaincode: %s", err))
		return false
	}
	return true
}

func UpgradeChaincode(channelID, chaincodeID, version, chaincodePath, orgName, orgAdmin string, ccUpgradeArgs [][]byte) bool {
	sdk, err := fabsdk.New(sdkConfig)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new SDK: %s", err))
		return false
	}
	defer sdk.Close()

	orgAdminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	// Org resource management client
	orgResMgmt, err := resmgmt.New(orgAdminContext)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to create new resource management client: %s", err))
		return false
	}

	// New chaincode policy (both orgs have to approve)
	ccPolicy, err := cauthdsl.FromString("AND ('Org1MSP.member','Org2MSP.member')")
	if err != nil {
		logrus.Errorln(err.Error())
		return false
	}

	_, err = orgResMgmt.UpgradeCC(channelID,
		resmgmt.UpgradeCCRequest{
			Name: chaincodeID,
			Path: chaincodePath,
			Version: version,
			Args: ccUpgradeArgs,
			Policy: ccPolicy,
		})
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to upgrade chaincode: %s", err))
		return false
	}
	return true
}

func RegisterChaincodeEvent(channelID, userName, orgName, chaincodeID, eventID string) {
	sdk, client, err := createChannelClient(channelID, userName, orgName)
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}
	defer sdk.Close()

	// Register chaincode event (pass in channel which receives event details when the event is complete)
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("Failed to register cc event: %s", err))
		return
	}
	defer client.UnregisterChaincodeEvent(reg)

	select {
	case ccEvent := <-notifier:
		logrus.Info(fmt.Sprintf("Received CC event: %#v\n", ccEvent))
	case <-time.After(time.Second * 20):
		logrus.Errorln(fmt.Sprintf("Did NOT receive CC event for eventId(%s)\n", eventID))
	}
}