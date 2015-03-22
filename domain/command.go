package domain

type Command interface {
	Command() []string
}

type ShellCommand struct {
	shellCommand string
}

type ExecCommand struct {
	execCommand []string
}

func (s ShellCommand) Command() []string {
	return []string{"/bin/sh", "-c", s.shellCommand}
}

func (e ExecCommand) Command() []string {
	return e.execCommand
}
