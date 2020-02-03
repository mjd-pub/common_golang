package logs

import "testing"

/**
 *
 *
 * @category   category
 * @package    logs
 * @subpackage Documentation\API
 * @author     zhangrubing  <zhangrubing@mioji.com>
 * @license    GPL https://mioji.com
 * @link       https://mioji.com
 */

func TestGenerateFileLogForTime(t *testing.T) {
	string := generateFileLogForTime("")

	t.Error(string)
}
