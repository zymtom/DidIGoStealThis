package codegrammar

regex := map[string]map[string][]string{
  "php":map[string][]string{
    "spacing": []string{
      "\\$\\w*? ="
      "( ?\\()"
      "\\$\\w* ?= ?.*?;"
    },
    "quote-usage": []string{
      "(\".*?\")"
      "('.*?')"
    }
    "variable": []string{
      "\\$(\\w*) ?="
    },
    "comments": []string{
      "\\/\\/.*?\\n"
    }
  },
  "py":[]string{},
  "go":[]string{},
}
type codeReport struct {
  results report
  file  fileobj
}
type report struct {
  quoteUsage quoteUsageReport
  spacing spacingReport
}
type fileobj {
  filetype string
  length  int
  lines []string
}
type quoteUsageReport struct {
  doubleQuotesPercentage float
  singleQuotesPercentage float
  properUsagePercentage float
}
type spacingReport struct {
  beforeParenthesis float
  afterParenthesis float
  beforeCurlyBracket float
  afterCurlyBracket float
  beforeQuote float
  afterQuote float
}
type variableReport struct {
  firstLetterCapital float
  wordsSeparatedCapital float
  wordsSeparatedUnderscore float
  
}
func main() {

}
func (c *codeReport) quoteUsage() quoteUsageReport {

}
func (c *codeReport) variableStats() variableReport {

}
