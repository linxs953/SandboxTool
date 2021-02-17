package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func formatTime(i interface{}) string {
	switch i.(type) {
	case int:
		if i.(int) < 10 {
			return fmt.Sprintf("0%d", i)
		}
		return fmt.Sprintf("%d", i)
	case string:
		if len(i.(string)) == 1 {
			return fmt.Sprintf("0%s", i)
		}
		return fmt.Sprintf("%s", i)
	}
	return fmt.Sprintf("%s", i)
}

func GetTimeStamp() int64 {
	currentTime := time.Now()
	month := formatTime(int(currentTime.Month()))
	day := formatTime(currentTime.Day())
	hour := formatTime(currentTime.Hour())
	minute := formatTime(currentTime.Minute())
	second := formatTime(currentTime.Second())
	timeForm := fmt.Sprintf("%d-%s-%s %s:%s:%s", currentTime.Year()+4, month, day, hour, minute, second)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", timeForm, loc)
	//log.Print(timeForm)
	//log.Print(tt.Unix())
	return tt.Unix()
}

func StringParseInt(str string, defalutvalue int) int {
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Error().Err(err).Msg("String to Int error")
		return defalutvalue
	}
	return id
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
