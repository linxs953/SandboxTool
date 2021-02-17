package utils

//func SensorCmdProcess(args []string) []string{
//	if len(args) == 0 {
//		return nil
//	}
//	prodURL := GetStringEnv("PROD","")
//	if prodURL == "" {
//		return nil
//	}
//	getTokenAPI := GetStringEnv("TOKENBYOPENIDAPI","")
//	if getTokenAPI == "" {
//		return nil
//	}
//	var rs []string
//	openid := args[0]
//	path := args[1]
//	getTokenUrl := fmt.Sprintf(getTokenAPI,openid)
//	resp,err := http.Get(getTokenUrl)
//	if err != nil {
//		return nil
//	}
//	bytes,err  := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil
//	}
//	token := string(bytes)
//	sb := GetSandbox(token)
//	sbEP,sbToken,sbAccToken := sb.Data.Endpoint,sb.Data.Token,sb.Data.AccessToken
//	debugURL := fmt.Sprintf("%s/%s?token=%s&access_token=%s",sbEP,path,sbToken,sbAccToken)
//	parts := parse(path)
//	if len(parts) == 0 {
//		return nil
//	}
//	frontendURL := fmt.Sprintf("%sapps/%s?token=%s",prodURL,strings.Join(parts,"/"),token)
//	rs = append(rs,debugURL)
//	rs = append(rs,frontendURL)
//	return rs
//}

//func parse(path string) []string {
//	var rs []string
//	if path == "" {
//		return nil
//	}
//	parts := strings.Split(path,"/")
//	if len(parts) == 0 {
//		return nil
//	}
//	parts = parts[2:]
//
//	idG := parts[1]
//	idGParts := strings.Split(idG,"-")
//	if len(idGParts) == 0 {
//		return nil
//	}
//	rs = append(rs,idGParts[1])
//	tp := parts[0]
//	rs = append(rs,tp + "s")
//	rs = append(rs,idGParts[3])
//	return rs
//}
