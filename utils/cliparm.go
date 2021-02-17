package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

func ParseCliParameter(args []string) ([]string, error) {
	var err error
	regex, err := regexp.Compile("(-|--)([a-z]+|-[A-Z]+)(=|\\s)")
	if err != nil {
		log.Error().Err(err).Msg("Compile regex error. Please check regex expression again")
		return nil, err
	}
	targetString := strings.Join(args[1:], " ")
	matchRS := regex.FindAllString(targetString, -1)
	for index, rs := range matchRS {
		if strings.Contains(rs, "=") {
			replacedStr := strings.Replace(rs, "=", "", -1)
			switch {
			case strings.Contains(replacedStr, "-"):
				replacedStr = strings.Replace(replacedStr, "-", "", -1)
			case strings.Contains(replacedStr, "--"):
				replacedStr = strings.Replace(replacedStr, "--", "", -1)
			default:
				log.Printf("%s can not contain `-` or `--`\n", replacedStr)
			}
			matchRS[index] = replacedStr
		}
	}

	return matchRS, err
}

func GetFlagValueStringP(cmd *cobra.Command, flagName string) (*string, error) {
	v, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Error().Err(err).Msgf("get flag name error: %v", err)
		return nil, err
	}
	return &v, nil
}



func GetFlagValueBool(cmd *cobra.Command,flagName string)(bool,error) {
	v,err := cmd.Flags().GetBool(flagName)
	if err != nil {
		log.Error().Err(err).Msgf("get flag name error: %v", err)
		return false, err
	}
	return v,nil
}