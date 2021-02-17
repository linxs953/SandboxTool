package utils

type TokenByUnionid struct {
	Data struct {
		JWT string `json:"jwt"`
	}
}

type RecyclePayload struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	RecycleStorage bool   `json:"recycle_storage"`
}

type Sandbox struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		Endpoint          string `json:"endpoint"`
		HeartbeatDuration int64  `json:"heartbeat_duration"`
		HeartbeatEndpoint string `json:"heartbeat_endpoint"`
		StorageID         string `json:"storage_id"`
		Token             string `json:"token"`
		WorkDir           string `json:"workdir"`
	}
}

type RecycleError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SandboxConfig struct {
	SandboxID   int    `json:"sandboxid"`
	SandboxType string `json:"sandboxtype"`
	Protect     int    `json:"protect"`
}

type SandboxGroupConfig struct {
	GroupConfig []SandboxConfig
}

type RecycleByPuidPayload struct {
	PUID      string `json:"puid"`
	ISRECYCLE bool   `json:"recycle_storage"`
}

type SBPools struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
	Data struct {
		Jupyter       string `json:"app-jupyter"`
		Pc            string `json:"pc-normal"`
		BoosttrapTest string `json:"sbp-boostrap-test"`
	}
}
