package logging

import (
	"encoding/json"
	"go.uber.org/zap"
)

var zapConfigJson = []byte(`{
	  "level": "info",
	  "encoding": "console",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

var zapConfig zap.Config
var Logger *zap.Logger

func init() {
	if err := json.Unmarshal(zapConfigJson, &zapConfig); err != nil {
		panic(err)
	}
	Logger = zap.Must(zapConfig.Build())
}
