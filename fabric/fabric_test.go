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
	ccQueryArgs := [][]byte{[]byte("b")}
	res := ChaincodeQuery("mychannel", "User1", "org1", "mycc", ccQueryArgs)
	fmt.Println(string(res))
}

func TestChaincodeExecute(t *testing.T) {
	ccTxArgs :=[][]byte{[]byte("a"), []byte("b"), []byte("10")}
	res := ChaincodeExecute("mychannel", "User1", "org1", "mycc", ccTxArgs)
	fmt.Println(string(res))
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

func TestCreateChaincode(t *testing.T) {
	res := CreateChaincode(
		"example_cc",
		"1.0",
		"github.com/example_cc",
		"./chaincode",
		"org1",
		"Admin")

	fmt.Println(res)
}

func TestInstantiateChaincode(t *testing.T) {
	res := InstantiateChaincode(
		"mychannel",
		"example_cc",
		"1.0",
		"github.com/example_cc",
		"org1",
		"Admin",
		[][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})

	fmt.Println(res)
}

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