package deprecated

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func jsonDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	var oldJSON, newJSON interface{}
	err := json.Unmarshal([]byte(old), &oldJSON)
	if err != nil {
		log.Printf("[ERROR] Version of %q in state is not valid JSON: %s", k, err)
		return false
	}
	err = json.Unmarshal([]byte(new), &newJSON)
	if err != nil {
		log.Printf("[ERROR] Version of %q in config is not valid JSON: %s", k, err)
		return true
	}
	return reflect.DeepEqual(oldJSON, newJSON)
}

func toStringArray(input []interface{}) []string {
	output := make([]string, len(input))

	for i, item := range input {
		output[i] = item.(string)
	}

	return output
}

func is404(err error) bool {
	return strings.Contains(err.Error(), "Code: 404")
}

func calculateConflictsWith(self string, group []string) []string {
	if len(group) < 2 {
		return []string{}
	}
	results := make([]string, 0, len(group)-2)
	for _, item := range group {
		if item == self {
			continue
		}
		results = append(results, item)
	}
	return results
}

func arrayToTerraformList(values []string) string {
	output := make([]string, len(values))
	for idx, value := range values {
		output[idx] = fmt.Sprintf(`"%s"`, value)
	}
	return fmt.Sprintf("[%s]", strings.Join(output, ", "))
}

func terraformSetToStringArray(set interface{}) []string {
	list := set.(*schema.Set).List()
	arr := make([]string, 0, len(list))
	for _, v := range list {
		arr = append(arr, v.(string))
	}
	return arr
}

func jsonStringArrayToStringArray(jsonList []interface{}) []string {
	strList := make([]string, 0, len(jsonList))
	for _, v := range jsonList {
		strList = append(strList, v.(string))
	}
	return strList
}

func isExpiredTokenErr(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "invalid accessor") {
		return true
	}
	if strings.Contains(err.Error(), "failed to find accessor entry") {
		return true
	}
	return false
}

func leaseExpiringSoon(d *schema.ResourceData) bool {
	startedStr := d.Get("lease_started").(string)
	duration := d.Get("lease_duration").(int)
	if startedStr == "" {
		return false
	}
	started, err := time.Parse(time.RFC3339, startedStr)
	if err != nil {
		log.Printf("[DEBUG] lease_started %q for %q is an invalid value, removing: %s", startedStr, d.Id(), err)
		d.Set("lease_started", "")
		return false
	}
	// whether the time the lease started plus the number of seconds specified in the duration
	// plus five minutes of buffer is before the current time or not. If it is, we don't need to
	// renew just yet.
	if started.Add(time.Second * time.Duration(duration)).Add(time.Minute * 5).Before(time.Now()) {
		return false
	}
	// if the lease duration expired more than five minutes ago, we can't renew anyways, so don't
	// bother even trying.
	if started.Add(time.Second * time.Duration(duration)).After(time.Now().Add(time.Minute * -5)) {
		return false
	}

	// the lease will expire in the next five minutes, or expired less than five minutes ago, in
	// which case renewing is worth a shot
	return true
}