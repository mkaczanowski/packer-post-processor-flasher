package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Flasher struct {
	config FlasherConfig
	ui     packer.Ui
}

func NewFlasher() packer.PostProcessor {
	return &Flasher{}
}

func (f *Flasher) selectDevice() (string, string, error) {
	out, err := exec.Command("lsblk", "-do", "name,tran").CombinedOutput()
	if err != nil {
		return "", "", err
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)

		if filepath.Base(f.config.Device) == fields[0] {
			return f.config.Device, strings.Trim(fields[1], "\n"), nil
		}
	}

	return "", "", fmt.Errorf("lsblk didn't find the device: %s", f.config.Device)
}

func (f *Flasher) ask(question string) error {
	answer, err := f.ui.Ask(question + " [Y/n]")
	if err != nil {
		return err
	}

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" {
		return errors.New("canceled by user")
	}

	return nil
}

// ConfigSpec returns the config spec
func (f *Flasher) ConfigSpec() hcldec.ObjectSpec {
	return f.config.FlatMapstructure().HCL2Spec()
}

// Flash writes image to the selected device
func (f *Flasher) Flash(image string) error {
	device, deviceType, err := f.selectDevice()
	if err != nil {
		return err
	}

	f.ui.Say(fmt.Sprintf("Found device %s of type: %s", device, deviceType))

	if f.config.Interactive {
		if err := f.ask("Do you want to proceed?"); err != nil {
			return err
		}

		if deviceType != "usb" {
			err := f.ask("Device type is not 'usb', are you sure you want to proceed?")
			if err != nil {
				return err
			}
		}
	}

	f.ui.Say("Flashing...")
	_, err = exec.Command(
		"dd",
		fmt.Sprintf("if=%s", image),
		fmt.Sprintf("of=%s", device),
		fmt.Sprintf("bs=%d", f.config.BlockSize),
	).CombinedOutput()

	if err != nil {
		return err
	}

	f.ui.Say("Syncing ...")
	syscall.Sync()

	return nil
}

func (f *Flasher) Configure(cfgs ...interface{}) error {
	if err := config.Decode(&f.config, &config.DecodeOpts{
		Interpolate:       true,
		InterpolateFilter: &interpolate.RenderFilter{},
	}, cfgs...); err != nil {
		return err
	}

	return nil
}

func (f *Flasher) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	files := artifact.Files()
	if len(files) != 1 {
		return nil, true, true, errors.New("expected only one file")
	}

	f.ui = ui

	return nil, true, true, f.Flash(files[0])
}
