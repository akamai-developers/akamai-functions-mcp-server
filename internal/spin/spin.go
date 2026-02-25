package spin

import "os/exec"

const spinBinary = "spin"

func RunCommand(command ...string) ([]byte, error) {

	cmd := exec.Command(spinBinary, command...)
	spinOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return spinOut, nil

}
