package sup

import (
    "log"
    "regexp"
)

var regexUser = regexp.MustCompile("^(?i)[a-z0-9]+$")
var regexTags = regexp.MustCompile("^(?i)([a-z0-9-]+ ?)+$")

func validateUser(str string) string {
    if !regexUser.MatchString(str) {
        log.Panicf("User name should be letters and numbers only: %s", str)
    }
    return str
}
func validateTags(str string) string {
    if !regexTags.MatchString(str) {
        log.Panicf("Tags should be letters, numbers, and dashes only: %s", str)
    }
    return str
}
