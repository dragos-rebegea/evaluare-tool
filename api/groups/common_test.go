package groups

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dragos-rebegea/evaluare-tool/api/shared"
	"github.com/dragos-rebegea/evaluare-tool/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type generalResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func init() {
	gin.SetMode(gin.TestMode)
}

func startWebServer(group shared.GroupHandler, path string, apiConfig config.ApiRoutesConfig) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	routes := ws.Group(path)
	group.RegisterRoutes(routes, apiConfig)
	return ws
}

func getServiceRoutesConfig() config.ApiRoutesConfig {
	return config.ApiRoutesConfig{
		APIPackages: map[string]config.APIPackageConfig{
			"auth": {
				Routes: []config.RouteConfig{
					{Name: "/register", Open: true},
					{Name: "/sendTransaction", Open: true},
					{Name: "/debug", Open: true},
					{Name: "/peerinfo", Open: true},
				},
			},
		},
	}
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	logError(err)
}

func requestToReader(request interface{}) io.Reader {
	data, _ := json.Marshal(request)
	return bytes.NewReader(data)
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
