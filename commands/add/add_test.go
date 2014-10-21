package add

import (
	. "github.com/franela/goblin"
	"github.com/franela/vault/vault"
	"github.com/franela/vault/vault/testutils"
	"github.com/mitchellh/cli"
	"testing"
)

func TestAdd(t *testing.T) {
	g := Goblin(t)

	g.Describe("Add", func() {
		g.BeforeEach(func() {
			vault.SetHomeDir(testutils.GetTemporaryHomeDir())
		})

		g.AfterEach(func() {
			testutils.RemoveTemporaryHomeDir(vault.UnsetHomeDir())
		})

		g.Describe("#Run", func() {
			g.It("Should not fail if recipient already exist")
			g.It("Should allow to add multiple recipients")

			g.It("Should add new recipients", func() {
				v := vault.Vaultfile{}
				v.Recipients = []string{"bob@example.com"}
				v.Save()

				c, _ := Factory()

				repairCommand := cli.MockCommand{}
				addCmd, _ := c.(addCommand)
				addCmd.Repair = &repairCommand

				addCmd.Run([]string{"alice@example.com"})

				newVaultfile, _ := vault.LoadVaultfile()
				g.Assert(newVaultfile.Recipients).Equal([]string{"bob@example.com", "alice@example.com"})
				g.Assert(repairCommand.RunCalled).IsTrue()
			})
		})
	})
}