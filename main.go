package fabric

import (

	"github.com/ChunmengYang/fabric-sdk-go/http"
)

func main()  {

	go http.Start()

	select {}
}
