package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func confirmPrompt(in io.Reader, out io.Writer, message string) (bool, error) {
	fmt.Fprintf(out, "%s Type 'yes' to confirm: ", message)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(line), "yes"), nil
}
