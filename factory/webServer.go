package factory

import (
	"io"

	"github.com/dragos-rebegea/evaluare-tool/api/gin"
	"github.com/dragos-rebegea/evaluare-tool/config"
	"github.com/dragos-rebegea/evaluare-tool/facade"
)

// StartWebServer creates and starts a web server able to respond with the metrics holder information
func StartWebServer(configs config.Configs) (io.Closer, error) {
	argsFacade := facade.ArgsEvaluationFacade{
		ApiInterface: configs.FlagsConfig.RestApiInterface,
		PprofEnabled: configs.FlagsConfig.EnablePprof,
	}

	authFacade, err := facade.NewEvaluationFacade(argsFacade)
	if err != nil {
		return nil, err
	}

	httpServerArgs := gin.ArgsNewWebServer{
		Facade:          authFacade,
		ApiConfig:       configs.ApiRoutesConfig,
		AntiFloodConfig: configs.GeneralConfig.Antiflood.WebServer,
	}

	httpServerWrapper, err := gin.NewWebServerHandler(httpServerArgs)
	if err != nil {
		return nil, err
	}

	err = httpServerWrapper.StartHttpServer()
	if err != nil {
		return nil, err
	}

	return httpServerWrapper, nil
}
