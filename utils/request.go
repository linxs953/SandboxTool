package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func Delete(url string, code int, payload string) ([]byte, error) {
	var err error
	reader := strings.NewReader(payload)
	req, err := http.NewRequest("DELETE", url, reader)
	if err != nil {
		log.Error().Err(err).Msg("new request error")
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("send request error")
		return nil, err
	}
	if resp.StatusCode != code {
		log.Printf("Unexpected Status Code %d", resp.StatusCode)
		return nil, nil
	}
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("read resp error")
		return nil, err
	}
	recycleError := RecycleError{}
	if string(bodyContent) != "" {
		err = json.Unmarshal(bodyContent, &recycleError)
		if err != nil {
			log.Error().Err(err).Msg("Parse json error")
			return nil, err
		}
		log.Printf("%d --- %s", recycleError.Code, recycleError.Message)
		return nil, err
	}
	return bodyContent, nil

}

func Get(url string, code int, polling bool) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Get resp error %s", strings.Split(url, "?")[0])
		return []byte("")
	}
	if resp.StatusCode == 401 {
		log.Printf("Unexpected Status Code %d", resp.StatusCode)
		fmt.Println(resp.Request.URL.String())
		return []byte("")
	}
	if resp.StatusCode != code {
		if polling {
			schedulerWithoutForce := strings.Replace(url, "&force=1", "", -1)
			for resp.StatusCode != code {
				resp, err = http.Get(schedulerWithoutForce)
				if err != nil {
					log.Error().Err(err).Msg("Get resp error restart")
					break
				}
			}
		} else {
			log.Printf("Unexpected Status Code %d", resp.StatusCode)
			return []byte("")
		}
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Read resp error")
		return []byte("")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	return bodyByte
}
