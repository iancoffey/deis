package fleet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/coreos/fleet/unit"
)

// Units returns a list of units filtered by target
func (c *FleetClient) Units(target string) (units []string, err error) {
	allUnits, err := c.Fleet.Units()
	if err != nil {
		return
	}
	for _, u := range allUnits {
		if strings.Contains(u.Name, target) {
			units = append(units, u.Name)
		}
	}
	if len(units) == 0 {
		err = fmt.Errorf("could not find unit: %s", target)
	}
	return
}

// nextUnit returns the next unit number for a given component
func (c *FleetClient) nextUnit(component string) (num int, err error) {
	units, err := c.Units(component)
	if err != nil {
		return
	}
	num, err = nextUnitNum(units)
	if err != nil {
		return
	}
	return
}

// lastUnit returns the last unit number for a given component
func (c *FleetClient) lastUnit(component string) (num int, err error) {
	units, err := c.Units(component)
	if err != nil {
		return
	}
	num, err = lastUnitNum(units)
	if err != nil {
		return
	}
	return
}

// NewUnit takes a component type and returns a Fleet unit
// that includes the relevant systemd service template
func NewUnit(component string, templatePaths []string) (uf *unit.UnitFile, err error) {
	template, err := readTemplate(component, templatePaths)
	if err != nil {
		return
	}
	uf, err = unit.NewUnitFile(string(template))
	if err != nil {
		return
	}
	return
}

// formatUnitName returns a properly formatted systemd service name
// using the given component type and number
func formatUnitName(component string, num int) (unitName string, err error) {
	component = strings.TrimPrefix(component, "deis-")
	if num == 0 {
		return "deis-" + component + ".service", nil
	}
	return "deis-" + component + "@" + strconv.Itoa(num) + ".service", nil
}

// readTemplate returns the contents of a systemd template for the given component
func readTemplate(component string, templatePaths []string) (out []byte, err error) {
	templateName := "deis-" + component + ".service"
	var templateFile string

	// look in $DEISCTL_UNITS env var, then the local and global root paths
	for _, p := range templatePaths {
		if p == "" {
			continue
		}
		filename := path.Join(p, templateName)
		if _, err := os.Stat(filename); err == nil {
			templateFile = filename
			break
		}
	}

	if templateFile == "" {
		return nil, fmt.Errorf("Could not find unit template for %v", component)
	}
	out, err = ioutil.ReadFile(templateFile)
	if err != nil {
		return
	}
	return
}
