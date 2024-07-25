/**
 * 数据校验异常封装
 * @author duhaifeng
 * @date   2021/04/15
 */
package busierr

import "fmt"

type ValidationError struct {
	FieldName    string
	ViolatedRule string
}

func (this *ValidationError) Error() string {
	return fmt.Sprintf("field %s data validate failed. should: %s", this.FieldName, this.ViolatedRule)
}
