package fabric


import (
	"testing"
	"fmt"
	"path"
)

//var (
//
//	ccInitArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
//	ccUpgradeArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("400")}
//	ccQueryArgs := [][]byte{[]byte("query"), []byte("b")}
//	ccTxArgs =[][]byte{[]byte("invoke"), []byte("a"), []byte("b"), []byte("1")}
//)

func TestChaincodeQuery(t *testing.T) {
	ccQueryArgs := [][]byte{[]byte("a")}
	res := ChaincodeQuery("mychannel", "User1", "org1", "mycc", "query",  ccQueryArgs)
	fmt.Println(string(res))
}

func TestChaincodeExecute(t *testing.T) {
	ccTxArgs :=[][]byte{[]byte("b"), []byte("a"), []byte("10")}
	res := ChaincodeExecute("mychannel", "User1", "org1", "mycc", "invoke", ccTxArgs)

	//ccQueryArgs := [][]byte{[]byte("b")}
	//res := ChaincodeExecute("mychannel", "User1", "org1", "mycc", "query", ccQueryArgs)
	fmt.Println(string(res))
}

func TestQueryInstalledChaincode(t *testing.T) {
	res := QueryInstalledChaincode(
		"org2",
		"Admin",
		"peer0.org2.example.com")

	fmt.Println(res)
}

func TestCreateChannel(t *testing.T) {
	res := CreateChannel(
		"orgchannel",
		path.Join("./channel", "orgchannel.tx"),
		"org1",
		"Admin",
		"ordererorg",
		"orderer.example.com",
		"Admin")

	fmt.Println(res)
}

//在peer容器中可查找到
//find -name "example*"
//./var/hyperledger/production/chaincodes/example_cc.2.0
func TestCreateChaincode(t *testing.T) {
	//res := CreateChaincode(
	//	"example_cc",
	//	"2.0",
	//	"github.com/example_cc",
	//	"./chaincode",
	//	"org1",
	//	"Admin")
	//
	//fmt.Println(res)

	res := CreateChaincode(
		"example_cc",
		"2.1",
		"github.com/example_cc",
		"./chaincode",
		"org2",
		"Admin")

	fmt.Println(res)
}

//未测试通
func TestInstantiateChaincode(t *testing.T) {
	res := InstantiateChaincode(
		"mychannel",
		"example_cc",
		"2.1",
		"github.com/example_cc",
		"org1",
		"Admin",
		[][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})

	fmt.Println(res)
}
//未测试通
func TestUpgradeChaincode(t *testing.T) {
	res := UpgradeChaincode(
		"mychannel",
		"example_cc",
		"1.0",
		"github.com/example_cc",
		"org1",
		"Admin",
		[][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("400")})

	fmt.Println(res)
}