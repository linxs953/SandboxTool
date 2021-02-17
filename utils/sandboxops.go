package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/AlecAivazis/survey.v1"
	"net/http"
	"strings"
)


//func GetSandboxType(sandboxId int, typename string) string {
//	sandboxGroup := GetSandboxGroupEnv("SANDBOXGROUP", sandboxId, []string{})
//	if len(sandboxGroup) < 3 {
//		return ""
//	}
//	switch typename {
//	case "SANDBOXTYPE":
//		return sandboxGroup[1]
//	case "PROTECTTYPE":
//		return sandboxGroup[2]
//	}
//	return ""
//}
//
//func GetNamespace(sandboxID int, tokenString string) string {
//	var err error
//	var sandboxType string
//	sandboxType = GetSandboxType(sandboxID, "SANDBOXTYPE")
//	if sandboxType == "" {
//		return ""
//	}
//	entryResp := GetV2Sandbox(tokenString, false)
//	sandboxEndpoint := entryResp.Data.Endpoint
//	if sandboxEndpoint == "" {
//		return ""
//	}
//	expression, err := regexp.Compile("/(sandbox(.*?))/")
//	if err != nil {
//		log.Error().Err(err).Msg("compile regex error")
//		return ""
//	}
//	target := expression.FindString(sandboxEndpoint)
//	namespace := strings.Replace(target, "/", "", -1)
//	return namespace
//}

//func RecycleByPuid(puid string) bool {
//	var token string
//	rcpuidAPI := GetStringEnv("RECYCLEAPI", "")
//	if rcpuidAPI == "" {
//		log.Print("recycle by puid api get nil")
//		return false
//	}
//	token = GenerateToken(puid,"")
//	if token == "" {
//		return false
//	}
//	sandbox := GetV2Sandbox(token,false)
//	emp := Sandbox{}
//	if sandbox == emp {
//		return false
//	}
//	sbToken,sbEndpoint := sandbox.Data.Token,sandbox.Data.Endpoint
//	if sbToken == "" {
//		log.Print("Get sandbox token empty")
//		return false
//	}
//	if sbEndpoint == "" {
//		log.Print("Get sandbox endpoint empty")
//		return false
//	}
//	log.Print(sbEndpoint)
//	urlParsed  := strings.Split(sbEndpoint,"/")
//	if len(urlParsed) < 5 {
//		log.Printf("url invalid %s",sbEndpoint)
//		return false
//	}
//	poolName := urlParsed[3]
//	randomName := urlParsed[4]
//	recycleURL := fmt.Sprintf(rcpuidAPI,poolName,randomName,sbToken,RecycleActionTime)
// 	client := &http.Client{}
//	req, _ := http.NewRequest(http.MethodPost,recycleURL,nil)
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Error().Err(err).Msgf("POST Req Error %v",err)
//		return false
//	}
//	if resp.StatusCode  != 200 {
//		log.Printf("Unexpected Status Code %d",resp.StatusCode)
//		return false
//	}
//	return true
//}
//
//func RecycleByUnionID(args []string) {
//	var recycleStorage bool
//	unionid := args[0]
//	sandboxIDStr := args[1]
//	isRecycle := args[2]
//	if isRecycle == "" {
//		return
//	} else if isRecycle == "0" {
//		recycleStorage = false
//	} else if isRecycle == "1" {
//		recycleStorage = true
//	} else {
//		log.Print("Can not support value out of 0 and 1")
//		return
//	}
//	if unionid == "" || sandboxIDStr == "" {
//		log.Print("unionid or sandbox type is nil")
//		return
//	}
//	sandboxID, err := strconv.Atoi(sandboxIDStr)
//	if err != nil {
//		log.Error().Err(err).Msg("str to int error")
//		return
//	}
//	token := GetTokenByUnionid(unionid)
//	if token == "" {
//		log.Print("Token is nil")
//		return
//	}
//	puid := ParseToken(args[0])
//	if puid == nil {
//		log.Print("parse puid from token error")
//		return
//	}
//	podName := fmt.Sprintf("%s-1-%d", puid, sandboxID)
//	namespace := GetNamespace(sandboxID, token)
//	if namespace == "" {
//		log.Print("Namespace is nil")
//		return
//	}
//	recycleAPI := GetStringEnv("RECYCLEAPI", "")
//	if recycleAPI == "" {
//		log.Print("RECYCLEAPI is nil")
//		return
//	}
//	recyclePayload := RecyclePayload{}
//	recyclePayload.RecycleStorage = recycleStorage
//	recyclePayload.Name = podName
//	recyclePayload.Namespace = namespace
//	payloadByte, err := json.Marshal(recyclePayload)
//	if err != nil {
//		log.Error().Err(err).Msg("Byte to string error")
//		return
//	}
//	respContent, err := Delete(recycleAPI, 204, string(payloadByte))
//	if err != nil {
//		log.Error().Err(err).Msg("Recycle error")
//		return
//	}
//	if respContent != nil {
//		log.Print("Recycle Successfully")
//	}
//}

func RecycleSandbox(sb Sandbox,args ...string) bool{
	var (
		endpoint string
		sbToken string
	)
	if len(args) > 0 {
		if len(args) < 2 {
			return false
		}
		endpoint = args[0]
		sbToken = args[1]
	}else {
		endpoint = sb.Data.Endpoint
		sbToken = sb.Data.Token
	}
	log.Print(endpoint)
	urlParts := strings.Split(endpoint,"/")
	if len(urlParts) < 5 {
		log.Printf("invalid url %s",endpoint)
		return false
	}
	poolName := urlParts[3]
	randomName := urlParts[4]
	rcpuidAPI := GetStringEnv("RECYCLEAPI", "")
	if rcpuidAPI == "" {
		log.Print("recycle by puid api get nil")
		return false
	}
	recycleURL := fmt.Sprintf(rcpuidAPI,poolName,randomName,sbToken,RecycleActionTime)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost,recycleURL,nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("POST Req Error %v",err)
		return false
	}
	if resp.StatusCode  != 200 {
		log.Printf("Unexpected Status Code %d",resp.StatusCode)
		return false
	}
	return true
}



func StartSandboxByPuid(puid string, storageid string, version string, entrance bool) {
	tokenString := GenerateToken(puid, storageid)
	if tokenString == "" {
		log.Print("generate token nil")
		return
	}
	rs := Sandbox{}
	switch version {
	case "v1":
		rs = GetV1Sandbox(tokenString, false)
	case "v2":
		rs = GetV2Sandbox(tokenString, false)
	default:
		return
	}
	emp := Sandbox{}
	if rs != emp {
		log.Print("Start Sandbox Successfully")
		if entrance {
			wssURL := strings.Replace(fmt.Sprintf("%s/terminal/ws?token=%s", rs.Data.Endpoint, rs.Data.Token), "http", "ws", -1)
			DebugForURL(wssURL, DefaultOrigin, true)
		}
	}
}

func StartSandboxByUnionID(unionid string, version string, entrance bool) {
	var rs Sandbox
	token := GetTokenByUnionid(unionid)
	if token == "" {
		log.Print("Get token nil")
		return
	}
	switch version {
	case "v1":
		rs = GetV1Sandbox(token, false)
	case "v2":
		rs = GetV2Sandbox(token, false)
	default:
		return
	}
	emp := Sandbox{}
	if rs != emp {
		log.Print("Start Sandbox Successfully")
		if entrance {
			wssURL := strings.Replace(fmt.Sprintf("%s/terminal/ws?token=%s", rs.Data.Endpoint, rs.Data.Token), "http", "ws", -1)
			DebugForURL(wssURL, DefaultOrigin, true)
		}
	}
}

func GetV2Sandbox(token string, restart bool) Sandbox {
	empEntry := Sandbox{}
	sbPoolAPI := GetStringEnv("SANDBOXPOOL", "")
	if sbPoolAPI == "" {
		log.Print("get [SANDBOXPOOL] env error")
		return empEntry
	}
	getPoolsURL := fmt.Sprintf(sbPoolAPI, token)
	respContent := Get(getPoolsURL, 200, false)
	if string(respContent) == "" {
		return empEntry
	}
	pools := SBPools{}
	if err := json.Unmarshal(respContent, &pools); err != nil {
		log.Printf("parse json error: ", err)
		return empEntry
	}
	tags := GetTagName(pools.Data)
	if len(tags) == 0 {
		return empEntry
	}
	selectedPool := ""
	prompt := &survey.Select{
		Message: "Choose Sandbox Pool:",
		Options: tags,
	}
	if err := survey.AskOne(prompt, &selectedPool, nil); err != nil {
		log.Printf("choose option error %v", err)
		return empEntry
	}
	fmt.Printf("switched %s\n", selectedPool)
	entryUserAPI := GetStringEnv("ENTRYV2SB", "")
	if entryUserAPI == "" {
		log.Print("get [entry-v2-api] env error")
		return empEntry
	}
	entryURL := fmt.Sprintf(entryUserAPI, selectedPool, token)
	if restart {
		entryURL = fmt.Sprintf("%s&force=1", entryURL)
	}
	respContent = Get(entryURL, 200, false)
	if string(respContent) == "" {
		return empEntry
	}
	entryResp := Sandbox{}
	if err := json.Unmarshal(respContent, &entryResp); err != nil {
		log.Printf("parse entry resp error,", err)
		return empEntry
	}
	return entryResp
}

func GetV1Sandbox(token string, restart bool) Sandbox {
	sb := Sandbox{}
	var err error
	v1SBAPI := GetStringEnv("EntryV1SB", "")
	if v1SBAPI == "" {
		log.Print("get [EntryV1SB] env error")
		return Sandbox{}
	}
	if restart {
		sb, err = RestartV1Sandbox(token, V1SandboxID)
		if err != nil {
			return Sandbox{}
		}
		return sb
	}
	v1SBURL := fmt.Sprintf(v1SBAPI, V1SandboxID, token)
	respBts := Get(v1SBURL, 201, false)
	if string(respBts) == "" {
		return Sandbox{}
	}
	if err := json.Unmarshal(respBts, &sb); err != nil {
		log.Print("parse resp error")
		return Sandbox{}
	}
	return sb
}

func RestartByUnionID(unionid string, version string, entrance bool) {
	var rs Sandbox
	if unionid == "" {
		log.Print("unionid is empty")
		return
	}
	token := GetTokenByUnionid(unionid)
	if token == "" {
		return
	}
	if strings.Contains(token, "user not found") {
		log.Print("user not found")
		return
	}
	switch version {
	case "v1":
		rs = GetV1Sandbox(token, true)
	case "v2":
		rs = GetV2Sandbox(token, true)
	default:
		return
	}
	emp := Sandbox{}
	if rs != emp {
		log.Print("Restart Successfully")
		if entrance {
			wssURL := strings.Replace(fmt.Sprintf("%s/terminal/ws?token=%s", rs.Data.Endpoint, rs.Data.Token), "http", "ws", -1)
			DebugForURL(wssURL, DefaultOrigin, true)
		}
	}
}

func RestartByPuid(puid string, version string, entrance bool) {
	if puid == "" {
		log.Print("puid is empty")
		return
	}
	token := GenerateToken(puid, "")
	if token == "" {
		return
	}
	var rs Sandbox
	switch version {
	case "v1":
		rs = GetV1Sandbox(token, true)
	case "v2":
		rs = GetV2Sandbox(token, true)
	default:
		return
	}
	emp := Sandbox{}
	if rs != emp {
		log.Print("Restart Sandbox Successfully")
		if entrance {
			wssURL := strings.Replace(fmt.Sprintf("%s/terminal/ws?token=%s", rs.Data.Endpoint, rs.Data.Token), "http", "ws", -1)
			DebugForURL(wssURL, DefaultOrigin, false)
		}
	}
}

func RestartV1Sandbox(token string, sandboxid string) (Sandbox, error) {
	restartAPI := GetStringEnv("EntryV1SB", "")
	if restartAPI == "" {
		log.Print("Get [EntryV1SB] env error")
		return Sandbox{}, errors.New("[EntryV1] env empty")
	}
	restartURL := fmt.Sprintf(restartAPI, sandboxid, token) + "&force=1"
	respBts := Get(restartURL, 201, true)
	if string(respBts) == "" {
		return Sandbox{}, errors.New("request occur error")
	}
	rs := Sandbox{}
	if err := json.Unmarshal(respBts, &rs); err != nil {
		log.Print("parse resp error")
		return Sandbox{}, err
	}
	return rs, nil
}

//func BlockByUnionID(args []string) {
//	if len(args) > 2 {
//		log.Print("too many args")
//		return
//	}
//	if args[0] == "" {
//		log.Print("unionid can not be none")
//		return
//	}
//	if args[1] == "" {
//		log.Print("resource group can not be none")
//		return
//	}
//	blockAPI := GetStringEnv("BLOCKUSERAPI", "")
//	if blockAPI == "" {
//		log.Print("get [block-api] env error\n")
//		return
//	}
//	puid := ParseToken(args[0])
//	if puid == nil {
//		log.Print("parse puid by token error")
//		return
//	}
//	blockUserURL := fmt.Sprintf(blockAPI, puid.(string), args[1])
//	resp, err := http.Post(blockUserURL, "application/json", strings.NewReader(""))
//	if err != nil {
//		log.Error().Err(err).Msgf("post occur %s", resp.Status)
//		return
//	}
//	if resp.StatusCode != 200 {
//		log.Printf("Unexpected Status Code %s", resp.StatusCode)
//		return
//	}
//	log.Printf("Block User %s Successfully", args[0])
//}
//
//func BlockByPuid(args []string) {
//	if len(args) > 2 {
//		log.Print("too many args")
//		return
//	}
//	if args[0] == "" {
//		log.Print("puid can not be none")
//		return
//	}
//	if args[1] == "" {
//		log.Print("resource group can not be none")
//		return
//	}
//	blockAPI := GetStringEnv("BLOCKAPI", "")
//	if blockAPI == "" {
//		log.Print("get [block-api] env error\n")
//		return
//	}
//	blockUserURL := fmt.Sprintf(blockAPI, args[0], args[1])
//	resp, err := http.Post(blockUserURL, "application/json", strings.NewReader(""))
//	if err != nil {
//		log.Error().Err(err).Msgf("post occur %s", resp.Status)
//		return
//	}
//	if resp.StatusCode != 200 {
//		log.Printf("Unexpected Status Code %s", resp.StatusCode)
//		return
//	}
//	log.Printf("Block User %s Successfully", args[0])
//}
//
//func UnbanByPuid(args []string) {
//	if len(args) > 2 {
//		log.Print("too many args\n")
//		return
//	}
//	if args[0] == "" {
//		log.Print("puid can not be none")
//		return
//	}
//	if args[1] == "" {
//		log.Print("resource group can not be none")
//		return
//	}
//	blockAPI := GetStringEnv("BLOCKAPI", "")
//	if blockAPI == "" {
//		log.Print(errors.New("get [block-api] env error"))
//		return
//	}
//	unBlockUserURL := fmt.Sprintf(blockAPI, args[0], args[1])
//	resp, _ := Delete(unBlockUserURL, 204, "")
//	if resp == nil {
//		log.Print("Delete user block qa error")
//		return
//	}
//	log.Printf("Unblock User %s Successfully", args[0])
//}
//
//func UnbanByUnionid(args []string) {
//	if args[0] == "" {
//		log.Print("unionid can not be none")
//		return
//	}
//	if args[1] == "" {
//		log.Print("resource group can not be none")
//		return
//	}
//	puid := ParseToken(args[0])
//	if puid == nil {
//		log.Print("parse puid from token error")
//		return
//	}
//	blockAPI := GetStringEnv("BLOCKAPI", "")
//	if blockAPI == "" {
//		log.Print("get block-api env error")
//		return
//	}
//	unBlockUserURL := fmt.Sprintf(blockAPI, puid.(string), args[1])
//	resp, _ := Delete(unBlockUserURL, 204, "")
//	if resp == nil {
//		log.Print("Delete user block qualification error")
//		return
//	}
//	log.Printf("Unblock User %s Successfully", args[0])
//}
