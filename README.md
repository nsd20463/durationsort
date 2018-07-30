# durationsort

Sort text strings so that text strings containing time durations sort as
expected by humans.

For example, the values

  "1ns"
  "2us"
  "3ms"
  "2m"
  "2m1.5s"
  "1h"

are sorted smallest to largest.

The style of duration matches what the Go standard library's time packet's
Duration type generates, with additional support for d,w and y with
approximate[1] but still useful (to humans) values.


[1] a day is not 24 hours when a leap-second is added, thus a week is not always 7*24 hours. 
A month is so varied it's useless. A year varies slightly due to astronomy
reasons, as does the length of a day.
