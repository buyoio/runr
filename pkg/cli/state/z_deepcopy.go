package state

// DeepCopy creates a deep copy of File.
func (f File) DeepCopy() File {
	// Strings are immutable in Go, so this is effectively a deep copy.
	return f
}

// DeepCopy creates a deep copy of Data.
func (d *Data) DeepCopy() *Data {
	if d == nil {
		return nil
	}

	cpy := *d // Start with a shallow copy.

	// Deep copy the File field.
	cpy.File = d.File.DeepCopy()

	// Deep copy the Args slice.
	cpy.Args = make(Args, len(d.Args))
	for i, arg := range d.Args {
		cpy.Args[i] = arg.DeepCopy() // Assuming Arg has a DeepCopy method.
	}

	// tmpl is unexported and thus not copied.

	return &cpy
}

// DeepCopy creates a deep copy of Arg.
func (a *Arg) DeepCopy() *Arg {
	if a == nil {
		return nil
	}

	cpy := *a // Start with a shallow copy.

	// Deep copy for *File.
	if a.File != nil {
		fileCopy := a.File.DeepCopy() // Using DeepCopy method of File.
		cpy.File = &fileCopy
	}

	// Deep copy for *Exec.
	if a.Exec != nil {
		execCopy := a.Exec.DeepCopy() // Using DeepCopy method of File.
		cpy.Exec = &execCopy
	}

	// Deep copy for *Value.
	if a.Value != nil {
		valueCopy := *a.Value
		cpy.Value = &valueCopy
	}

	// Deep copy for Envs slice.
	if a.Envs != nil {
		cpy.Envs = make([]string, len(a.Envs))
		copy(cpy.Envs, a.Envs)
	}

	return &cpy
}

// DeepCopy creates a deep copy of DockerRunner.
func (dr *DockerRunner) DeepCopy() *DockerRunner {
	if dr == nil {
		return nil
	}

	cpy := *dr // Start with a shallow copy.

	// Deep copy for Dockerfile, which is of type *Data.
	if dr.Dockerfile != nil {
		cpy.Dockerfile = dr.Dockerfile.DeepCopy() // Assuming Data has a DeepCopy method.
	}

	// Deep copy for Labels slice.
	if dr.Labels != nil {
		cpy.Labels = make([]string, len(dr.Labels))
		copy(cpy.Labels, dr.Labels)
	}

	// Quantity is a basic type, so it's already copied by the shallow copy.

	return &cpy
}

// DeepCopy creates a deep copy of RunnerSetup.
func (rs *RunnerSetup) DeepCopy() *RunnerSetup {
	if rs == nil {
		return nil
	}

	cpy := *rs // Start with a shallow copy for basic and embedded fields.

	// Deep copy for the Packages slice.
	if rs.Packages != nil {
		cpy.Packages = make([]string, len(rs.Packages))
		copy(cpy.Packages, rs.Packages)
	}

	// Deep copy for Installimage, which is of type *Data.
	if rs.Installimage != nil {
		cpy.Installimage = rs.Installimage.DeepCopy() // Assuming Data has a DeepCopy method.
	}

	// User, RunnerWorkDir are basic types and thus already handled by the shallow copy.
	// SCMPlatform is an embedded field, and its fields are basic types, so they're also handled by the shallow copy.

	return &cpy
}

// DeepCopy creates a deep copy of Runner.
func (r *Runner) DeepCopy() *Runner {
	if r == nil {
		return nil
	}

	cpy := *r // Start with a shallow copy.

	// Deep copy for pointer fields.
	if r.Name != nil {
		nameCopy := *r.Name
		cpy.Name = &nameCopy
	}
	if r.Host != nil {
		hostCopy := *r.Host
		cpy.Host = &hostCopy
	}
	if r.HCloud != nil {
		hCloudCopy := *r.HCloud
		cpy.HCloud = &hCloudCopy
	}
	if r.Password != nil {
		passwordCopy := *r.Password
		cpy.Password = &passwordCopy
	}
	if r.Pre != nil {
		preCopy := *r.Pre
		cpy.Pre = &preCopy
	}
	if r.Post != nil {
		postCopy := *r.Post
		cpy.Post = &postCopy
	}

	// Deep copy for Setup, which is of type *RunnerSetup.
	if r.Setup != nil {
		cpy.Setup = r.Setup.DeepCopy() // Assuming RunnerSetup has a DeepCopy method.
	}

	// Deep copy for the Docker map.
	if r.Docker != nil {
		cpy.Docker = make(map[string]*DockerRunner, len(r.Docker))
		for key, dockerRunner := range r.Docker {
			cpy.Docker[key] = dockerRunner.DeepCopy() // Assuming DockerRunner has a DeepCopy method.
		}
	}

	// Basic fields like User, Auth are copied by the shallow copy operation.

	return &cpy
}

// DeepCopy creates a deep copy of HCloudServerCreateOpts.
func (hc *HCloudServerCreateOpts) DeepCopy() *HCloudServerCreateOpts {
	if hc == nil {
		return nil
	}

	cpy := *hc // Start with a shallow copy for basic fields and the Name field which is explicitly excluded from deep copy.

	// Deep copy for the Setup field, which is of type *SCMPlatform.
	if hc.Setup != nil {
		cpy.Setup = hc.Setup.DeepCopy() // Assuming SCMPlatform has a DeepCopy method.
	}

	// Deep copy for SSHKeys slice.
	if hc.SSHKeys != nil {
		cpy.SSHKeys = make([]string, len(hc.SSHKeys))
		copy(cpy.SSHKeys, hc.SSHKeys)
	}

	// Deep copy for pointers to string fields.
	if hc.Location != nil {
		locationCopy := *hc.Location
		cpy.Location = &locationCopy
	}
	if hc.Datacenter != nil {
		datacenterCopy := *hc.Datacenter
		cpy.Datacenter = &datacenterCopy
	}

	// Deep copy for UserData, which is of type *Data.
	if hc.UserData != nil {
		cpy.UserData = hc.UserData.DeepCopy() // Assuming Data has a DeepCopy method.
	}

	// Fields like Type, Image, StartAfterCreate, Labels, Automount, Volumes, Networks, Firewalls, PlacementGroup are basic types or slices of basic types, so they are effectively handled by the shallow copy.

	return &cpy
}

// DeepCopy creates a deep copy of Environments.
func (e *Environments) DeepCopy() *Environments {
	if e == nil {
		return nil
	}

	cpy := *e // Start with a shallow copy.

	// Deep copy for the HCloud map.
	if e.HCloud != nil {
		cpy.HCloud = make(map[string]*HCloudServerCreateOpts, len(e.HCloud))
		for key, hCloudOpts := range e.HCloud {
			cpy.HCloud[key] = hCloudOpts.DeepCopy() // Assuming HCloudServerCreateOpts has a DeepCopy method.
		}
	}

	return &cpy
}

// DeepCopy creates a deep copy of HCloud.
func (h *HCloud) DeepCopy() *HCloud {
	if h == nil {
		return nil
	}
	cpy := *h // Basic types are effectively copied here.
	return &cpy
}

// DeepCopy creates a deep copy of Hetzner.
func (h *Hetzner) DeepCopy() *Hetzner {
	if h == nil {
		return nil
	}
	cpy := *h // Basic types are effectively copied here.
	return &cpy
}

// DeepCopy creates a deep copy of Github.
func (g *Github) DeepCopy() *Github {
	if g == nil {
		return nil
	}
	cpy := *g // Start with a shallow copy.

	// Deep copy for pointer to string fields.
	if g.Repository != nil {
		repoCopy := *g.Repository
		cpy.Repository = &repoCopy
	}

	return &cpy
}

// DeepCopy creates a deep copy of Gitlab.
func (g *Gitlab) DeepCopy() *Gitlab {
	if g == nil {
		return nil
	}
	cpy := *g // Start with a shallow copy.

	// Deep copy for pointer to string fields.
	if g.Repository != nil {
		repoCopy := *g.Repository
		cpy.Repository = &repoCopy
	}

	return &cpy
}

// DeepCopy creates a deep copy of SCMPlatform.
func (s *SCMPlatform) DeepCopy() *SCMPlatform {
	if s == nil {
		return nil
	}

	cpy := *s // Start with a shallow copy for the basic type fields.

	// Deep copy for pointer to string fields.
	if s.Repository != nil {
		repoCopy := *s.Repository
		cpy.Repository = &repoCopy
	}

	return &cpy
}

// DeepCopy creates a deep copy of State.
func (s *State) DeepCopy() *State {
	if s == nil {
		return nil
	}

	cpy := *s // Start with a shallow copy for the basic type fields.

	// Deep copy for Logs.
	if s.Logs != nil {
		cpy.Logs = new(Logs)
		*cpy.Logs = *s.Logs // Assuming Logs contains only basic types or are properly handling their own deep copy.
	}

	// Deep copy for Runners map.
	if s.Runners != nil {
		cpy.Runners = make(map[string]*Runner, len(s.Runners))
		for key, runner := range s.Runners {
			cpy.Runners[key] = runner.DeepCopy() // Assuming Runner has a DeepCopy method.
		}
	}

	// Deep copy for Environments.
	if s.Environments != nil {
		cpy.Environments = s.Environments.DeepCopy() // Assuming Environments has a DeepCopy method.
	}

	// Deep copy for Github, Gitlab, HCloud, and Hetzner.
	if s.Github != nil {
		cpy.Github = s.Github.DeepCopy() // Assuming Github has a DeepCopy method.
	}
	if s.Gitlab != nil {
		cpy.Gitlab = s.Gitlab.DeepCopy() // Assuming Gitlab has a DeepCopy method.
	}
	if s.HCloud != nil {
		cpy.HCloud = s.HCloud.DeepCopy() // Assuming HCloud has a DeepCopy method.
	}
	if s.Hetzner != nil {
		cpy.Hetzner = s.Hetzner.DeepCopy() // Assuming Hetzner has a DeepCopy method.
	}

	return &cpy
}
