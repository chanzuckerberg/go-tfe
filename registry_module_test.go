package tfe

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModulesCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	testOrg, testOrgCleanup := createOrganization(t, client)
	defer testOrgCleanup()

	optionsModule := RegistryModuleCreateOptions{
		Name:     *String(randomString(t)),
		Provider: "random",
	}

	t.Run("creating a module", func(t *testing.T) {
		m, err := client.RegistryModules.Create(ctx, testOrg.Name, optionsModule)
		require.NoError(t, err)
		assert.Equal(t, optionsModule.Name, m.Name)
		assert.Equal(t, optionsModule.Provider, m.Provider)
	})

	t.Run("creating a module version", func(t *testing.T) {
		optionsModuleVersion := RegistryModuleCreateVersionOptions{
			Version: "1.2.3",
		}

		mv, err := client.RegistryModules.CreateVersion(ctx, testOrg.Name, optionsModule.Name, optionsModule.Provider, optionsModuleVersion)
		require.NoError(t, err)
		assert.Equal(t, optionsModuleVersion.Version, mv.Version)
	})
}

func TestModulesDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	testOrg, testOrgCleanup := createOrganization(t, client)
	defer testOrgCleanup()

	optionsModule := ModuleCreateOptions{
		Name:     *String(randomString(t)),
		Provider: "random",
	}

	m, err := client.Registry.CreateModule(ctx, testOrg.Name, optionsModule)

	t.Run("creating a module", func(t *testing.T) {
		require.NoError(t, err)
		assert.Equal(t, optionsModule.Name, m.Name)
		assert.Equal(t, optionsModule.Provider, m.Provider)
	})

	t.Run("deleting a module", func(t *testing.T) {
		//testModule, _ := client.Registry.CreateModule(ctx, testOrg.Name, optionsModule)
		//fmt.Printf("%+v\n", testModule)
		//fmt.Printf("%+v\n", testModule.Organization.Name)
		fmt.Print("foo")
		deleteErr := client.Registry.DeleteModule(ctx, testOrg.Name, m.Name)
		require.NoError(t, deleteErr)

		//try to get the module

	})
}
