package set

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/franela/vault/gpg"
	"github.com/franela/vault/ui"
	"github.com/franela/vault/vault"
	"github.com/mitchellh/cli"
)

const setHelpText = `
Usage: vault set [options] vaultpath [text]

   Sets something into the vault. Either plain text or files using the -f 
   parameter are allowed. You can also read from stdin using "-".

Options: 
    
    -f	Input file which will be added to the vault in the specified vaultpath

`

func Factory() (cli.Command, error) {
	return setCommand{}, nil
}

type setCommand struct {
}

func (setCommand) Help() string {
	return setHelpText
}

func (setCommand) Run(args []string) int {
	var fileName string

	// parse flags
	cmdFlags := flag.NewFlagSet("set", flag.ContinueOnError)
	cmdFlags.StringVar(&fileName, "f", "", "specify the file to encrypt")

	if err := cmdFlags.Parse(args); err != nil {
		ui.Printf(setHelpText)
		return 1
	}

	args = cmdFlags.Args()

	// make sure no excess arguments were given
	if len(args) > 2 {
		ui.Printf(setHelpText)
		return 1
	}

	// check that at least a file, or plain text input was given
	if len(args) != 2 && fileName == "" {
		ui.Printf(setHelpText)
		return 1
	}

	// in case a file was given, add the fileName to the arguments to populate the key
	if len(args) == 0 {
		args = append(args, path.Base(fileName))
	}

	// load the Vaultfile
	vaultFile, err := vault.LoadVaultfile()

	if err != nil {
		ui.Printf("%s", err)
		return 1
	}

	// make sure the current Vaultfile has recipients
	if len(vaultFile.Recipients) == 0 {
		ui.Printf("Cannot set in vault if Vaultfile has no recipients. Use `vault add` to add one or more recipients.\n")
		return 3
	}

	// make sure the path to encrypt is under the current path
	vaultPath := args[0]
	if ok, err := isUnderCurrentPath(vaultPath); err != nil || !ok {
		if err != nil {
			ui.Printf("%s\n", err)
		} else {
			ui.Printf("Destination should be under current path.\n")
		}
		return 1
	}

	// make sure the vault path has the ".asc" extension
	if filepath.Ext(vaultPath) != ".asc" {
		vaultPath = vaultPath + ".asc"
	}

	if len(fileName) > 0 {
		// encrypt a file
		err = gpg.EncryptFile(path.Join(vault.GetHomeDir(), vaultPath), fileName, vaultFile.Recipients)
	} else if args[1] == "-" {
		// encrypt from stdin
		ui.Printf("Enter the text below, end with EOF (Ctrl + D):\n")
		err = gpg.EncryptFile(path.Join(vault.GetHomeDir(), vaultPath), "-", vaultFile.Recipients)
	} else {
		// encrypt plain text
		err = gpg.Encrypt(path.Join(vault.GetHomeDir(), vaultPath), args[1], vaultFile.Recipients)
	}

	if err != nil {
		ui.Printf("%s", err)
		return 1
	}

	return 0
}

func isUnderCurrentPath(p string) (bool, error) {
	abs := ""
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	if path.IsAbs(p) {
		abs = p
	} else {
		abs = path.Join(wd, path.Dir(p))
	}

	return strings.HasPrefix(abs, wd), nil
}

func (setCommand) Synopsis() string {
	return "Store files or plain-text intro your vault."
}
