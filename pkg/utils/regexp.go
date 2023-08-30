package utils

import "regexp"

func FindGroups(re *regexp.Regexp, s string) map[string]string {
	matches := re.FindStringSubmatch(s)
	subNames := re.SubexpNames()
	if matches == nil || len(matches) != len(subNames) {
		return nil
	}

	matchMap := map[string]string{}
	for i := 1; i < len(matches); i++ {
		matchMap[subNames[i]] = matches[i]
	}
	return matchMap
}
