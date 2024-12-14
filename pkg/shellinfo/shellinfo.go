package shellinfo

// ShellInfo holds the shell information
type ShellInfo struct {
	RootIndicator         string
	ColorTemplate         string
	EscapedDollar         string
	EscapedBacktick       string
	EscapedBackslash      string
	EvalPromptPrefix      string
	EvalPromptSuffix      string
	EvalPromptRightPrefix string
	EvalPromptRightSuffix string
}
