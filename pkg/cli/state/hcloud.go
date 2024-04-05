package state

import (
	"bytes"
)

func (s *HCloudServerCreateOpts) Args() ([]string, *bytes.Buffer, error) {
	args := []string{
		"server", "create",
		// required
		"--type", s.Type,
		"--image", s.Image,
		"--name", s.Name,
	}
	if s.Location != nil {
		args = append(args, "--location", *s.Location)
	}
	if s.Datacenter != nil {
		args = append(args, "--datacenter", *s.Datacenter)
	}
	if s.StartAfterCreate {
		args = append(args, "--start-after-create")
	}
	if s.Automount {
		args = append(args, "--automount")
	}
	for _, value := range s.Labels {
		args = append(args, "--label", value)
	}
	for _, value := range s.Volumes {
		args = append(args, "--volume", value)
	}
	for _, value := range s.Networks {
		args = append(args, "--network", value)
	}
	for _, value := range s.Firewalls {
		args = append(args, "--firewall", value)
	}
	if s.PlacementGroup != nil {
		args = append(args, "--placement-group", *s.PlacementGroup)
	}
	var data bytes.Buffer
	if s.UserData != nil {
		content, err := s.UserData.Content()
		if err != nil {
			return nil, nil, err
		}
		_, err = data.Write([]byte(content))
		if err != nil {
			return nil, nil, err
		}
		args = append(args, "--user-data-from-file", "-")
	}

	return args, &data, nil
}
