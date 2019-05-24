// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.
// +build !windows

package systemd

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/coreos/go-systemd/dbus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Conn struct {
}

func dbusNewFake() (*dbus.Conn, error) {
	fmt.Println("my dbusNewFake")
	return nil, nil
}

func connListUnitsFake(c *dbus.Conn) ([]dbus.UnitStatus, error) {
	fmt.Println("my connListUnitsFake")
	return []dbus.UnitStatus{
		{Name: "unit1", ActiveState: "active"},
		{Name: "unit2", ActiveState: "active"},
		{Name: "unit3", ActiveState: "inactive"},
	}, nil
}

func connGetUnitPropertiesFake(c *dbus.Conn, unitName string) (map[string]interface{}, error) {
	props := map[string]interface{}{
		"CPUShares": uint64(10),
	}
	return props, nil
}

func connCloseFake(c *dbus.Conn) {
}

func TestDefaultConfiguration(t *testing.T) {
	check := Check{}
	check.Configure([]byte(``), []byte(``))

	assert.Equal(t, []string(nil), check.config.instance.UnitNames)
	assert.Equal(t, []string(nil), check.config.instance.UnitRegexStrings)
	assert.Equal(t, []*regexp.Regexp(nil), check.config.instance.UnitRegexPatterns)
}

func TestConfiguration(t *testing.T) {
	check := Check{}
	rawInstanceConfig := []byte(`
unit_names:
 - ssh.service
 - syslog.socket
unit_regex:
 - lvm2-.*
 - cloud-.*
`)
	err := check.Configure(rawInstanceConfig, []byte(``))

	assert.Nil(t, err)
	// assert.Equal(t, true, check.config.instance.UnitNames)
	assert.ElementsMatch(t, []string{"ssh.service", "syslog.socket"}, check.config.instance.UnitNames)
	regexes := []*regexp.Regexp{
		regexp.MustCompile("lvm2-.*"),
		regexp.MustCompile("cloud-.*"),
	}
	assert.Equal(t, regexes, check.config.instance.UnitRegexPatterns)
}
func TestSystemdCheck(t *testing.T) {
	dbusNew = dbusNewFake
	connListUnits = connListUnitsFake
	connClose = connCloseFake
	connGetUnitProperties = connGetUnitPropertiesFake

	// create an instance of our test object
	check := new(Check)
	check.Configure(nil, nil)

	// setup expectations
	mockSender := mocksender.NewMockSender(check.ID())
	mockSender.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockSender.On("Commit").Return()

	check.Run()

	mockSender.AssertCalled(t, "Gauge", "systemd.unit.active.count", float64(2), "", []string(nil))
	mockSender.AssertCalled(t, "Gauge", "systemd.unit.count", float64(1), "", []string{"unit:unit1"})
	mockSender.AssertCalled(t, "Gauge", "systemd.unit.count", float64(1), "", []string{"unit:unit2"})
	mockSender.AssertCalled(t, "Gauge", "systemd.unit.count", float64(1), "", []string{"unit:unit3"})
	mockSender.AssertNumberOfCalls(t, "Gauge", 4)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)
}
